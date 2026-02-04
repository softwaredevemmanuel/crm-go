// services/live_class_service.go
package services

import (
    "errors"
    "fmt"
    "strings"
    "time"
    "crm-go/models"
    "crm-go/utils"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

// UpdateLiveClass - main update function
func (s *LiveClassService) UpdateLiveClass(liveClassID uuid.UUID, req models.LiveClassUpdateInput, updatedBy uuid.UUID) (*models.LiveClassResponse, error) {
    // Fetch existing live class
    var liveClass models.LiveClass
    if err := s.db.First(&liveClass, "id = ?", liveClassID).Error; err != nil {
        return nil, errors.New("live class not found")
    }
    
    // Check current status
    now := time.Now()
    hasStarted := now.After(liveClass.StartTime)
    hasEnded := now.After(liveClass.EndTime)
    isCancelled := liveClass.IsCancelled != nil && *liveClass.IsCancelled
    
    // Track changes for audit
    changes := []string{}
    updates := make(map[string]interface{})
    
    // 1. Handle cancellation/uncancellation
    if req.IsCancelled != nil {
        if *req.IsCancelled {
            // Trying to cancel the class
            if isCancelled {
                return nil, errors.New("class is already cancelled")
            }
            if hasStarted && !hasEnded {
                return nil, errors.New("cannot cancel an ongoing class")
            }
            if hasEnded {
                return nil, errors.New("cannot cancel a completed class")
            }
            
            // Cancel the class
            liveClass.IsCancelled = req.IsCancelled
            updates["is_cancelled"] = true
            changes = append(changes, "cancelled")
            
            fmt.Printf("Live class %s cancelled by user %s\n", liveClassID, updatedBy)
        } else {
            // Trying to uncancel the class (set IsCancelled = false)
            if !isCancelled {
                return nil, errors.New("class is not cancelled")
            }
            
            // Can only uncancel if class hasn't started yet
            if hasStarted {
                return nil, errors.New("cannot uncancel a class that has already started")
            }
            
            // Uncancel the class
            liveClass.IsCancelled = req.IsCancelled // false
            updates["is_cancelled"] = false
            changes = append(changes, "uncancelled")
            
            fmt.Printf("Live class %s uncancelled by user %s\n", liveClassID, updatedBy)
        }
    }
    
    // If class is cancelled, restrict what can be updated
    if isCancelled && !(req.IsCancelled != nil && !*req.IsCancelled) {
        // When cancelled, only allow:
        // 1. Uncancelling (handled above)
        // 2. Host notes
        
        // Check if trying to update anything other than allowed fields
        if (req.Title != nil || req.Description != nil || req.StartTime != nil || 
            req.EndTime != nil || req.Agenda != nil || req.RecordAutomatically != nil ||
            req.MaxAttendees != nil || req.Platform != nil) {
            return nil, errors.New("cannot update cancelled class. Uncancel it first or only update host notes.")
        }
    }
    
    // 2. Fields that can ONLY be updated BEFORE class starts AND when not cancelled
    if !hasStarted && !isCancelled {
        // Title
        if req.Title != nil {
            newTitle := strings.TrimSpace(*req.Title)
            if liveClass.Title != newTitle {
                liveClass.Title = newTitle
                updates["title"] = newTitle
                changes = append(changes, "title")
                
                // Update slug if title changed
                if newSlug, err := s.generateSlug(newTitle); err == nil {
                    liveClass.Slug = newSlug
                    updates["slug"] = newSlug
                }
            }
        }
        
        // Description
        if req.Description != nil {
            newDesc := strings.TrimSpace(*req.Description)
            if liveClass.Description != newDesc {
                liveClass.Description = newDesc
                updates["description"] = newDesc
                changes = append(changes, "description")
            }
        }
        
   
        
        // Time rescheduling - only if not cancelled
        if req.StartTime != nil || req.EndTime != nil {
            var newStartTime, newEndTime time.Time
            var err error
            
            // Parse new start time
            if req.StartTime != nil {
                newStartTime, err = utils.ParseTime(*req.StartTime)
                if err != nil {
                    return nil, fmt.Errorf("invalid start time: %v. Use format like '2024-03-15T14:00:00Z'", err)
                }
            } else {
                newStartTime = liveClass.StartTime
            }
            
            // Parse new end time
            if req.EndTime != nil {
                newEndTime, err = utils.ParseTime(*req.EndTime)
                if err != nil {
                    return nil, fmt.Errorf("invalid end time: %v. Use format like '2024-03-15T15:30:00Z'", err)
                }
            } else {
                newEndTime = liveClass.EndTime
            }
            
            // Validate new times
            if newStartTime.Before(now) {
                return nil, errors.New("cannot reschedule to a past time")
            }
            
            if newEndTime.Before(newStartTime) {
                return nil, errors.New("end time must be after start time")
            }
            
            // Calculate new duration
            newDuration := int(newEndTime.Sub(newStartTime).Minutes())
            if newDuration < 5 || newDuration > *req.Duration {
                return nil, errors.New("resulting duration must be between 5 and " + fmt.Sprintf("%d", *req.Duration) + " minutes")
            }
            
            // Check if times actually changed
            if !liveClass.StartTime.Equal(newStartTime) || !liveClass.EndTime.Equal(newEndTime) {
                // Check for tutor schedule conflicts
                var conflictCount int64
                err := s.db.Model(&models.LiveClass{}).
                    Where("tutor_id = ? AND id != ? AND (is_cancelled IS NULL OR is_cancelled = false)", 
                        liveClass.TutorID, liveClassID).
                    Where("(start_time, end_time) OVERLAPS (?, ?)", newStartTime, newEndTime).
                    Count(&conflictCount).Error
                
                if err == nil && conflictCount > 0 {
                    return nil, errors.New("tutor has another class scheduled at this time")
                }
                
                liveClass.StartTime = newStartTime
                liveClass.EndTime = newEndTime
                liveClass.Duration = newDuration
                
                updates["start_time"] = newStartTime
                updates["end_time"] = newEndTime
                updates["duration"] = newDuration
                
                changes = append(changes, "schedule")
            }
        }
        
        // Duration (independent update)
        if req.Duration != nil {
            newDuration := *req.Duration
            if newDuration < 5 || newDuration > *req.Duration {
                return nil, errors.New("duration must be between 5 and " + fmt.Sprintf("%d", *req.Duration) + " minutes")
            }
            liveClass.Duration = newDuration
            updates["duration"] = newDuration
            changes = append(changes, "duration")
        }
        
        // Timezone
        if req.Timezone != nil {
            newTimezone := strings.TrimSpace(*req.Timezone)
            liveClass.Timezone = newTimezone
            updates["timezone"] = newTimezone
            changes = append(changes, "timezone")
        }
        
        // Chapter ID
        if req.ChapterID != nil {
            if *req.ChapterID != uuid.Nil {
                // Validate chapter exists and belongs to course
                var chapter models.Chapter
                if err := s.db.First(&chapter, "id = ?", *req.ChapterID).Error; err != nil {
                    return nil, errors.New("chapter not found")
                }
                if chapter.CourseID != liveClass.CourseID {
                    return nil, errors.New("chapter does not belong to this course")
                }
            }
            liveClass.ChapterID = req.ChapterID
            updates["chapter_id"] = req.ChapterID
            changes = append(changes, "chapter")
        }
        
        // Topic ID
        if req.TopicID != nil {
            if *req.TopicID != uuid.Nil {
                var topic models.Topic
                if err := s.db.First(&topic, "id = ?", *req.TopicID).Error; err != nil {
                    return nil, errors.New("topic not found")
                }
                if topic.CourseID != liveClass.CourseID {
                    return nil, errors.New("topic does not belong to this course")
                }
            }
            liveClass.TopicID = req.TopicID
            updates["topic_id"] = req.TopicID
            changes = append(changes, "topic")
        }
        
        // Lesson ID
        if req.LessonID != nil {
            if *req.LessonID != uuid.Nil {
                var lesson models.Lessons
                if err := s.db.First(&lesson, "id = ?", *req.LessonID).Error; err != nil {
                    return nil, errors.New("lesson not found")
                }
                if lesson.CourseID != liveClass.CourseID {
                    return nil, errors.New("lesson does not belong to this course")
                }
            }
            liveClass.LessonID = req.LessonID
            updates["lesson_id"] = req.LessonID
            changes = append(changes, "lesson")
        }
        
        // Capacity - Max Attendees
        if req.MaxAttendees != nil {
            newMax := *req.MaxAttendees
            if newMax < 1 || newMax > 1000 {
                return nil, errors.New("max attendees must be between 1 and 1000")
            }
            
            // Validate against current enrollments
            var confirmedCount int64
            s.db.Model(&models.LiveClassEnrollment{}).
                Where("live_class_id = ? AND status = 'confirmed'", liveClassID).
                Count(&confirmedCount)
            
            if int(confirmedCount) > newMax {
                return nil, fmt.Errorf("cannot reduce capacity below %d confirmed enrollments", confirmedCount)
            }
            
            // Ensure min attendees doesn't exceed new max
            if req.MinAttendees == nil && newMax < liveClass.MinAttendees {
                return nil, fmt.Errorf("max attendees cannot be less than current min attendees (%d)", liveClass.MinAttendees)
            }
            
            liveClass.MaxAttendees = newMax
            updates["max_attendees"] = newMax
            changes = append(changes, "max_attendees")
        }
        
        // Capacity - Min Attendees
        if req.MinAttendees != nil {
            newMin := *req.MinAttendees
            if newMin < 1 || newMin > 1000 {
                return nil, errors.New("min attendees must be between 1 and 1000")
            }
            
            // Ensure min doesn't exceed max
            currentMax := liveClass.MaxAttendees
            if req.MaxAttendees != nil {
                currentMax = *req.MaxAttendees
            }
            
            if newMin > currentMax {
                return nil, fmt.Errorf("min attendees cannot exceed max attendees (%d)", currentMax)
            }
            
            liveClass.MinAttendees = newMin
            updates["min_attendees"] = newMin
            changes = append(changes, "min_attendees")
        }
        
        // Waitlist Enabled
        if req.WaitlistEnabled != nil {
            liveClass.WaitlistEnabled = *req.WaitlistEnabled
            updates["waitlist_enabled"] = *req.WaitlistEnabled
            changes = append(changes, "waitlist_enabled")
        }
        
        // Waitlist Capacity
        if req.WaitlistCapacity != nil {
            newWaitlistCap := *req.WaitlistCapacity
            if newWaitlistCap < 0 || newWaitlistCap > 100 {
                return nil, errors.New("waitlist capacity must be between 0 and 100")
            }
            liveClass.WaitlistCapacity = newWaitlistCap
            updates["waitlist_capacity"] = newWaitlistCap
            changes = append(changes, "waitlist_capacity")
        }
        
        // Access Level
        if req.AccessLevel != nil {
            liveClass.AccessLevel = *req.AccessLevel
            updates["access_level"] = *req.AccessLevel
            changes = append(changes, "access_level")
        }
        
        
        // Platform
        if req.Platform != nil {
            newPlatform := *req.Platform
            if liveClass.Platform != newPlatform {
                // Platform change requires meeting recreation
                return nil, errors.New("platform cannot be changed after creation")
            }
        }
    } else if hasStarted && !isCancelled {
        // If class has started, check for restricted fields
        restrictedFields := []struct{
            Field interface{}
            Name  string
        }{
            {req.Title, "title"},
            {req.Description, "description"},
            {req.StartTime, "start_time"},
            {req.EndTime, "end_time"},
            {req.MaxAttendees, "max_attendees"},
            {req.Platform, "platform"},
            {req.ChapterID, "chapter_id"},
            {req.TopicID, "topic_id"},
            {req.LessonID, "lesson_id"},
        }
        
        for _, field := range restrictedFields {
            if field.Field != nil {
                switch v := field.Field.(type) {
                case *string:
                    if v != nil {
                        return nil, fmt.Errorf("cannot update %s after class has started", field.Name)
                    }
                case *uuid.UUID:
                    if v != nil {
                        return nil, fmt.Errorf("cannot update %s after class has started", field.Name)
                    }
                default:
                    if v != nil {
                        return nil, fmt.Errorf("cannot update %s after class has started", field.Name)
                    }
                }
            }
        }
    }
    
    // 3. Fields that can be updated BEFORE class ENDS (and when not cancelled)
    if (!hasEnded || isCancelled) && !isCancelled {
        // Agenda
        if req.Agenda != nil {
            newAgenda := strings.TrimSpace(*req.Agenda)
            liveClass.Agenda = newAgenda
            updates["agenda"] = newAgenda
            changes = append(changes, "agenda")
        }
        
        // Recommended Setup
        if req.RecommendedSetup != nil {
            newSetup := strings.TrimSpace(*req.RecommendedSetup)
            liveClass.RecommendedSetup = newSetup
            updates["recommended_setup"] = newSetup
            changes = append(changes, "recommended_setup")
        }
        
        // Recording settings
        if req.RecordAutomatically != nil {
            liveClass.RecordAutomatically = *req.RecordAutomatically
            updates["record_automatically"] = *req.RecordAutomatically
            changes = append(changes, "record_automatically")
        }
        
        if req.RecordingStorage != nil {
            liveClass.RecordingStorage = *req.RecordingStorage
            updates["recording_storage"] = *req.RecordingStorage
            changes = append(changes, "recording_storage")
        }
        
        if req.AutoPublishRecordings != nil {
            liveClass.AutoPublishRecordings = *req.AutoPublishRecordings
            updates["auto_publish_recordings"] = *req.AutoPublishRecordings
            changes = append(changes, "auto_publish_recordings")
        }
        
      
    }
    
    // 4. Fields that can be updated ANYTIME (even if cancelled)
    // Host Notes
    if req.HostNotes != nil {
        newHostNotes := strings.TrimSpace(*req.HostNotes)
        liveClass.HostNotes = newHostNotes
        updates["host_notes"] = newHostNotes
        changes = append(changes, "host_notes")
    }
    
    // Check if any updates were made
    if len(updates) == 0 {
        return nil, errors.New("no changes provided")
    }
    
    // Add updated timestamp
    liveClass.UpdatedAt = now
    updates["updated_at"] = now
    
    // Save changes to database
    if err := s.db.Model(&liveClass).Updates(updates).Error; err != nil {
        return nil, errors.New("failed to update live class: " + err.Error())
    }
    
    // Send notifications for important changes
    if len(changes) > 0 && !isCancelled {
        go s.notifyAboutUpdates(liveClass, changes)
    }
    
    return s.liveClassToResponse(&liveClass, true), nil
}

// UpdateLiveClassWithTx - transaction version
func (s *LiveClassService) UpdateLiveClassWithTx(tx *gorm.DB, liveClassID uuid.UUID, req models.LiveClassUpdateInput, updatedBy uuid.UUID, withDetails bool) (*models.LiveClassResponse, error) {
    var liveClass models.LiveClass
    if err := tx.First(&liveClass, "id = ?", liveClassID).Error; err != nil {
        return nil, errors.New("live class not found")
    }
    
    now := time.Now()
    hasStarted := now.After(liveClass.StartTime)
    hasEnded := now.After(liveClass.EndTime)
    isCancelled := liveClass.IsCancelled != nil && *liveClass.IsCancelled
    
    // Prepare update map
    updates := make(map[string]interface{})
 
    // Handle cancellation/uncancellation
    if req.IsCancelled != nil {
        if *req.IsCancelled {
            // Cancel
        
            if hasStarted && !hasEnded {
                return nil, errors.New("cannot cancel an ongoing class")
            }
            if hasEnded {
                return nil, errors.New("cannot cancel a completed class")
            }
            updates["is_cancelled"] = withDetails
        } else {
            // Uncancel
            if !isCancelled {
                return nil, errors.New("class is not cancelled")
            }
            if hasStarted {
                return nil, errors.New("cannot uncancel a class that has already started")
            }
            updates["is_cancelled"] = false
        }
    }
    
    // If class is cancelled, restrict other updates
    if isCancelled && !(req.IsCancelled != nil && !*req.IsCancelled) {
        // Only allow host notes when cancelled
        if (req.Title != nil || req.Description != nil || req.Agenda != nil || 
            req.RecordAutomatically != nil || req.StartTime != nil || req.EndTime != nil ||
            req.MaxAttendees != nil) {
            return nil, errors.New("cannot update cancelled class. Uncancel it first.")
        }
    }
    
    
    // Fields that can be updated anytime
    if req.HostNotes != nil {
        updates["host_notes"] = strings.TrimSpace(*req.HostNotes)
    }
    
    // Check if any updates
    if len(updates) == 0 {
        return nil, errors.New("no changes provided")
    }
    
    // Add timestamp
    updates["updated_at"] = now
    
    // Perform update
    result := tx.Model(&liveClass).Updates(updates)
    if result.Error != nil {
        return nil, result.Error
    }
    
    if result.RowsAffected == 0 {
        return nil, errors.New("live class was not updated")
    }
    
    // Refresh object
    if err := tx.First(&liveClass, "id = ?", liveClassID).Error; err != nil {
        return nil, err
    }
    
    return s.liveClassToResponse(&liveClass, false), nil
}



// Helper for notifications
func (s *LiveClassService) notifyAboutUpdates(liveClass models.LiveClass, changes []string) {
    fmt.Printf("Live class '%s' (ID: %s) updated. Changes: %v\n", 
        liveClass.Title, liveClass.ID, changes)
    
    // Notify enrolled students of important changes
    importantChanges := map[string]bool{
        "cancelled": true,
        "uncancelled": true,
        "schedule": true,
        "max_attendees": true,
    }
    
    for _, change := range changes {
        if importantChanges[change] {
            // Get enrolled students
            var enrollments []models.LiveClassEnrollment
            s.db.Where("live_class_id = ? AND status IN ('confirmed', 'waitlisted')", 
                liveClass.ID).Find(&enrollments)
            
            fmt.Printf("Notifying %d students about change: %s\n", len(enrollments), change)
            break
        }
    }
}

