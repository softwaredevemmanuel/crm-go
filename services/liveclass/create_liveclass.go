// services/live_class_service.go
package services

import (
	"crm-go/models"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type LiveClassService struct {
	db *gorm.DB
}

func NewLiveClassService(db *gorm.DB) *LiveClassService {
	return &LiveClassService{db: db}
}

// Helper to generate slug from title
func (s *LiveClassService) generateSlug(title string) (string, error) {
	// Convert to lowercase and replace spaces with hyphens
	slug := strings.ToLower(strings.TrimSpace(title))
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	slug = result.String()

	// Ensure uniqueness
	baseSlug := slug
	counter := 1

	for {
		var count int64
		s.db.Model(&models.LiveClass{}).Where("slug = ?", slug).Count(&count)

		if count == 0 {
			break
		}

		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++

		if counter > 100 {
			return "", errors.New("failed to generate unique slug")
		}
	}

	return slug, nil
}

// Validate live class input
func (s *LiveClassService) validateLiveClass(req models.LiveClassInput) error {
	// Validate course exists
	var course models.Course
	if err := s.db.First(&course, "id = ?", req.CourseID).Error; err != nil {
		return errors.New("course not found")
	}

	// Validate tutor exists and is a tutor
	var tutor models.User
	if err := s.db.First(&tutor, "id = ? AND role = ?", req.TutorID, "tutor").Error; err != nil {
		return errors.New("tutor not found or user is not a tutor")
	}

	// Validate tutor teaches this course
	if tutor.ID != course.TutorID {
		return errors.New("tutor does not teach this course")
	}

	// Validate chapter if provided
	if req.ChapterID != nil && *req.ChapterID != uuid.Nil {
		var chapter models.Chapter
		if err := s.db.First(&chapter, "id = ?", req.ChapterID).Error; err != nil {
			return errors.New("chapter not found")
		}
		if chapter.CourseID != req.CourseID {
			return errors.New("chapter does not belong to this course")
		}
	}

	// Validate topic if provided
	if req.TopicID != nil && *req.TopicID != uuid.Nil {
		var topic models.Topic
		if err := s.db.First(&topic, "id = ?", req.TopicID).Error; err != nil {
			return errors.New("topic not found")
		}
		if topic.CourseID != req.CourseID {
			return errors.New("topic does not belong to this course")
		}
	}

	// Validate lesson if provided
	if req.LessonID != nil && *req.LessonID != uuid.Nil {
		var lesson models.Lessons
		if err := s.db.First(&lesson, "id = ?", req.LessonID).Error; err != nil {
			return errors.New("lesson not found")
		}
		if lesson.CourseID != req.CourseID {
			return errors.New("lesson does not belong to this course")
		}
	}

	

	// Time validation
	if !req.StartTime.IsZero() || !req.EndTime.IsZero() {
		// both times must be provided together
		if req.StartTime.IsZero() || req.EndTime.IsZero() {
			return errors.New("both start time and end time must be provided")
		}

		now := time.Now()
		// cannot schedule in the past
		if req.StartTime.Before(now) {
			return errors.New("cannot schedule class in the past")
		}

		// end must be after start
		if !req.EndTime.After(req.StartTime) {
			return errors.New("end time must be after start time")
		}

		// enforce reasonable duration bounds
		newDuration := int(req.EndTime.Sub(req.StartTime).Minutes())
		if newDuration < 5 {
			return errors.New("duration must be at least 5 minutes")
		}
		if req.Duration != 0 && newDuration > req.Duration {
			return errors.New("resulting duration must not exceed the provided duration")
		}
	}

	// Validate capacity
	if req.MaxAttendees < req.MinAttendees {
		return errors.New("max attendees cannot be less than min attendees")
	}

	// Validate waitlist
	if req.WaitlistEnabled && req.WaitlistCapacity < 0 {
		return errors.New("waitlist capacity cannot be negative")
	}

	// Check for overlapping classes for same tutor
	var overlappingCount int64
	err := s.db.Model(&models.LiveClass{}).
		Where("tutor_id = ?", req.TutorID).
		Where("(start_time, end_time) OVERLAPS (?, ?)", req.StartTime, req.EndTime).
		Count(&overlappingCount).Error

	if err == nil && overlappingCount > 0 {
		return errors.New("tutor has another class scheduled at this time")
	}

	// Check for overlapping classes for same course
	if req.ChapterID != nil {
		err = s.db.Model(&models.LiveClass{}).
			Where("course_id = ? AND chapter_id = ?", req.CourseID, req.ChapterID).
			Where("(start_time, end_time) OVERLAPS (?, ?)", req.StartTime, req.EndTime).
			Count(&overlappingCount).Error

		if err == nil && overlappingCount > 0 {
			return errors.New("this chapter already has a class scheduled at this time")
		}
	}

	return nil
}

// CreateMeeting - abstract meeting creation (to be implemented per platform)
func (s *LiveClassService) createMeeting(platform string, liveClass *models.LiveClass) error {
	switch platform {
	case "zoom":
		return s.createZoomMeeting(liveClass)
	case "google_meet":
		return s.createGoogleMeet(liveClass)
	case "teams":
		return s.createTeamsMeeting(liveClass)
	case "jitsi":
		return s.createJitsiMeeting(liveClass)
	case "bigbluebutton":
		return s.createBigBlueButtonMeeting(liveClass)
	default:
		// For custom/platform, generate meeting details manually
		liveClass.MeetingID = uuid.New().String()
		liveClass.MeetingURL = fmt.Sprintf("https://meet.example.com/%s", liveClass.Slug)
		liveClass.MeetingPassword = generateRandomPassword(8)
		return nil
	}
}

// CreateLiveClass - main create function
func (s *LiveClassService) CreateLiveClass(req models.LiveClassInput) (*models.LiveClassResponse, error) {
	// Validate input
	if err := s.validateLiveClass(req); err != nil {
		return nil, err
	}

	// Generate slug
	slug, err := s.generateSlug(req.Title)
	if err != nil {
		return nil, err
	}

	// Calculate duration if not provided
	duration := req.Duration
	if duration == 0 {
		duration = int(req.EndTime.Sub(req.StartTime).Minutes())
		if duration <= 0 {
			return nil, errors.New("invalid duration calculated from start and end times")
		}
	}

	// Set defaults for optional fields
	if req.Platform == "" {
		req.Platform = "zoom" // Default platform
	}

	if req.AccessLevel == "" {
		req.AccessLevel = "enrolled"
	}

	if req.Timezone == "" {
		req.Timezone = "UTC"
	}

	if req.RecordingStorage == "" {
		req.RecordingStorage = "platform"
	}

	// Create live class model
	liveClass := models.LiveClass{
		ID:                    uuid.New(),
		CourseID:              req.CourseID,
		ChapterID:             req.ChapterID,
		TopicID:               req.TopicID,
		LessonID:              req.LessonID,
		Title:                 strings.TrimSpace(req.Title),
		Description:           strings.TrimSpace(req.Description),
		Slug:                  slug,
		StartTime:             req.StartTime,
		EndTime:               req.EndTime,
		Duration:              duration,
		Timezone:              req.Timezone,
		TutorID:               req.TutorID,
		HostNotes:             strings.TrimSpace(req.HostNotes),
		MaxAttendees:          req.MaxAttendees,
		MinAttendees:          req.MinAttendees,
		WaitlistEnabled:       req.WaitlistEnabled,
		WaitlistCapacity:      req.WaitlistCapacity,
		AccessLevel:           req.AccessLevel,
		Platform:              req.Platform,
		Agenda:                strings.TrimSpace(req.Agenda),
		RecommendedSetup:      strings.TrimSpace(req.RecommendedSetup),
		RecordAutomatically:   req.RecordAutomatically,
		RecordingStorage:      req.RecordingStorage,
		AutoPublishRecordings: req.AutoPublishRecordings,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		IsCancelled:           func() *bool { b := false; return &b }(),
	}

	// Create meeting on selected platform
	if err := s.createMeeting(req.Platform, &liveClass); err != nil {
		return nil, fmt.Errorf("failed to create meeting: %v", err)
	}

	// Save to database
	if err := s.db.Create(&liveClass).Error; err != nil {
		// Check for duplicate slug
		if strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "slug") {
			return nil, errors.New("a class with similar title already exists")
		}
		return nil, errors.New("failed to save live class: " + err.Error())
	}

	// Convert to response
	return s.liveClassToResponse(&liveClass, true), nil
}

// CreateLiveClassWithTx - for use with transactions
func (s *LiveClassService) CreateLiveClassWithTx(tx *gorm.DB, req models.LiveClassInput) (*models.LiveClassResponse, error) {
	if err := s.validateLiveClass(req); err != nil {
		return nil, err
	}

	slug, err := s.generateSlug(req.Title)
	if err != nil {
		return nil, err
	}

	duration := req.Duration
	if duration == 0 {
		duration = int(req.EndTime.Sub(req.StartTime).Minutes())
	}

	// Set defaults
	if req.Platform == "" {
		req.Platform = "zoom"
	}
	if req.AccessLevel == "" {
		req.AccessLevel = "enrolled"
	}
	if req.Timezone == "" {
		req.Timezone = "UTC"
	}
	if req.RecordingStorage == "" {
		req.RecordingStorage = "platform"
	}

	liveClass := models.LiveClass{
		ID:                    uuid.New(),
		CourseID:              req.CourseID,
		ChapterID:             req.ChapterID,
		TopicID:               req.TopicID,
		LessonID:              req.LessonID,
		Title:                 strings.TrimSpace(req.Title),
		Description:           strings.TrimSpace(req.Description),
		Slug:                  slug,
		StartTime:             req.StartTime,
		EndTime:               req.EndTime,
		Duration:              duration,
		Timezone:              req.Timezone,
		TutorID:               req.TutorID,
		HostNotes:             strings.TrimSpace(req.HostNotes),
		MaxAttendees:          req.MaxAttendees,
		MinAttendees:          req.MinAttendees,
		WaitlistEnabled:       req.WaitlistEnabled,
		WaitlistCapacity:      req.WaitlistCapacity,
		AccessLevel:           req.AccessLevel,
		Platform:              req.Platform,
		Agenda:                strings.TrimSpace(req.Agenda),
		RecommendedSetup:      strings.TrimSpace(req.RecommendedSetup),
		RecordAutomatically:   req.RecordAutomatically,
		RecordingStorage:      req.RecordingStorage,
		AutoPublishRecordings: req.AutoPublishRecordings,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		IsCancelled:           func() *bool { b := false; return &b }(),
	}

	// Create meeting
	if err := s.createMeeting(req.Platform, &liveClass); err != nil {
		return nil, fmt.Errorf("failed to create meeting: %v", err)
	}

	if err := tx.Create(&liveClass).Error; err != nil {
		return nil, errors.New("failed to save live class: " + err.Error())
	}

	return s.liveClassToResponse(&liveClass, false), nil
}

// Helper to convert LiveClass to LiveClassResponse
func (s *LiveClassService) liveClassToResponse(liveClass *models.LiveClass, withDetails bool) *models.LiveClassResponse {
	response := &models.LiveClassResponse{
		ID:                    liveClass.ID,
		CourseID:              liveClass.CourseID,
		ChapterID:             liveClass.ChapterID,
		TopicID:               liveClass.TopicID,
		LessonID:              liveClass.LessonID,
		Title:                 liveClass.Title,
		Description:           liveClass.Description,
		Slug:                  liveClass.Slug,
		StartTime:             liveClass.StartTime,
		EndTime:               liveClass.EndTime,
		Duration:              liveClass.Duration,
		Timezone:              liveClass.Timezone,
		TutorID:               liveClass.TutorID,
		MaxAttendees:          liveClass.MaxAttendees,
		MinAttendees:          liveClass.MinAttendees,
		WaitlistEnabled:       liveClass.WaitlistEnabled,
		WaitlistCapacity:      liveClass.WaitlistCapacity,
		AccessLevel:           liveClass.AccessLevel,
		Platform:              liveClass.Platform,
		MeetingID:             liveClass.MeetingID,
		MeetingURL:            liveClass.MeetingURL,
		MeetingPassword:       liveClass.MeetingPassword,
		Agenda:                liveClass.Agenda,
		RecommendedSetup:      liveClass.RecommendedSetup,
		TestURL:               liveClass.TestURL,
		RecordAutomatically:   liveClass.RecordAutomatically,
		RecordingStorage:      liveClass.RecordingStorage,
		AutoPublishRecordings: liveClass.AutoPublishRecordings,
		HostNotes:             liveClass.HostNotes,
		CreatedAt:             liveClass.CreatedAt,
		UpdatedAt:             liveClass.UpdatedAt,
	}

	// Calculate status
	now := time.Now()
	if now.After(liveClass.StartTime) && now.Before(liveClass.EndTime) {
		response.Status = "ongoing"
		response.IsLiveNow = true
	} else if now.Before(liveClass.StartTime) {
		response.Status = "scheduled"
		response.IsUpcoming = true
	} else {
		response.Status = "completed"
	}

	// Add details if requested
	if withDetails {
		// Get course name
		var course models.Course
		if err := s.db.First(&course, "id = ?", liveClass.CourseID).Error; err == nil {
			response.CourseName = course.Title
		}

		// Get tutor name
		var tutor models.User
		if err := s.db.First(&tutor, "id = ?", liveClass.TutorID).Error; err == nil {
			response.TutorName = fmt.Sprintf("%s %s", tutor.FirstName, tutor.LastName)
		}

		// Get enrollment count
		var enrollmentCount int64
		s.db.Model(&models.LiveClassEnrollment{}).
			Where("live_class_id = ? AND status = ?", liveClass.ID, "confirmed").
			Count(&enrollmentCount)

		response.TotalEnrolled = int(enrollmentCount)
		response.AvailableSeats = liveClass.MaxAttendees - int(enrollmentCount)
	}

	return response
}

// Helper functions
func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// Placeholder meeting creation functions
func (s *LiveClassService) createZoomMeeting(liveClass *models.LiveClass) error {
	// Implement Zoom API integration
	liveClass.MeetingID = fmt.Sprintf("zoom-%s", uuid.New().String())
	liveClass.MeetingURL = fmt.Sprintf("https://zoom.us/j/%s", liveClass.MeetingID)
	liveClass.MeetingPassword = generateRandomPassword(6)
	return nil
}

func (s *LiveClassService) createGoogleMeet(liveClass *models.LiveClass) error {
	liveClass.MeetingID = fmt.Sprintf("meet-%s", uuid.New().String())
	liveClass.MeetingURL = fmt.Sprintf("https://meet.google.com/%s", liveClass.MeetingID)
	return nil
}

func (s *LiveClassService) createTeamsMeeting(liveClass *models.LiveClass) error {
	liveClass.MeetingID = fmt.Sprintf("teams-%s", uuid.New().String())
	liveClass.MeetingURL = fmt.Sprintf("https://teams.microsoft.com/l/meetup-join/%s", liveClass.MeetingID)
	return nil
}

func (s *LiveClassService) createJitsiMeeting(liveClass *models.LiveClass) error {
	liveClass.MeetingID = liveClass.Slug
	liveClass.MeetingURL = fmt.Sprintf("https://meet.jit.si/%s", liveClass.Slug)
	return nil
}

func (s *LiveClassService) createBigBlueButtonMeeting(liveClass *models.LiveClass) error {
	liveClass.MeetingID = liveClass.Slug
	liveClass.MeetingURL = fmt.Sprintf("https://bbb.example.com/b/%s", liveClass.Slug)
	liveClass.MeetingPassword = generateRandomPassword(8)
	return nil
}
