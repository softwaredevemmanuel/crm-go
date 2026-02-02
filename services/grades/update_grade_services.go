// services/grade_service.go
package services

import (
    "errors"
    "strings"
    "time"
    
    "crm-go/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)


// UpdateGradeWithTx - for use with transactions
func (s *GradeService) UpdateGradeWithTx(tx *gorm.DB, gradeID uuid.UUID, req models.GradeUpdateInput) (*models.GradeResponse, error) {
    // Fetch existing grade
    var grade models.Grade
    if err := tx.Preload("Student").Preload("Course").First(&grade, "id = ?", gradeID).Error; err != nil {
        return nil, errors.New("grade not found")
    }
    
    // Track changes
    changes := make(map[string]interface{})
    var updatedFields []string
    
      // Update score and grade letter together
    if req.Score != nil {
        if *req.Score < 0 || *req.Score > 100 {
            return nil, errors.New("score must be between 0 and 100")
        }
        
        if grade.Score != *req.Score {
            grade.Score = *req.Score
            grade.Grade = s.calculateGradeLetter(*req.Score) // This will save!
            updatedFields = append(updatedFields, "score", "grade")
        }
    }
    
    // Update remarks if provided
    if req.Remarks != nil {
        newRemarks := strings.TrimSpace(*req.Remarks)
        if grade.Remarks != newRemarks {
            changes["remarks"] = true
            updatedFields = append(updatedFields, "remarks")
            grade.Remarks = newRemarks
        }
    }
    
    // Update assignment if provided
    if req.AssignmentID != nil {
        if *req.AssignmentID != uuid.Nil {
            var assignment models.Assignment
            if err := tx.First(&assignment, "id = ?", req.AssignmentID).Error; err != nil {
                return nil, errors.New("assignment not found")
            }
            
            if assignment.CourseID != grade.CourseID {
                return nil, errors.New("assignment does not belong to this course")
            }
            
            // Check for duplicate (excluding current grade)
            var duplicateCount int64
            query := tx.Model(&models.Grade{}).
                Where("student_id = ? AND course_id = ? AND assignment_id = ? AND id != ?",
                    grade.StudentID, grade.CourseID, req.AssignmentID, gradeID)
            
            if err := query.Count(&duplicateCount).Error; err != nil {
                return nil, errors.New("failed to check for duplicate grade")
            }
            
            if duplicateCount > 0 {
                return nil, errors.New("grade already exists for this student, course, and assignment")
            }
            
            if (grade.AssignmentID == nil && req.AssignmentID != nil) ||
               (grade.AssignmentID != nil && req.AssignmentID == nil) ||
               (grade.AssignmentID != nil && req.AssignmentID != nil && *grade.AssignmentID != *req.AssignmentID) {
                changes["assignment_id"] = true
                updatedFields = append(updatedFields, "assignment_id")
                grade.AssignmentID = req.AssignmentID
            }
        } else {
            // Setting to nil
            if grade.AssignmentID != nil {
                changes["assignment_id"] = "removed"
                updatedFields = append(updatedFields, "assignment_id")
                grade.AssignmentID = nil
            }
        }
    }
    
    // Check if any changes
    if len(updatedFields) == 0 {
        return nil, errors.New("no changes provided")
    }
    
    // Update timestamp
    grade.UpdatedAt = time.Now()
    
    // Save with selective update (only changed fields)
    if err := tx.Model(&grade).Select(updatedFields).Updates(grade).Error; err != nil {
        return nil, errors.New("failed to update grade: " + err.Error())
    }
    
    // Return response
    return s.gradeToResponse(&grade), nil
}


// GetGradeHistory - get update history for a grade
func (s *GradeService) GetGradeHistory(gradeID uuid.UUID) ([]models.GradeHistory, error) {
    // Assuming you have a grade_history table
    var history []models.GradeHistory
    
    if err := s.db.Where("grade_id = ?", gradeID).
        Order("created_at DESC").
        Find(&history).Error; err != nil {
        return nil, errors.New("failed to fetch grade history: " + err.Error())
    }
    
    return history, nil
}