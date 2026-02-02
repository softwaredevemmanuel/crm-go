// services/grade_service.go
package services

import (
    "errors"
    "fmt"
    "strings"
    "time"
    
    "crm-go/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type GradeService struct {
    db *gorm.DB
}

func NewGradeService(db *gorm.DB) *GradeService {
    return &GradeService{db: db}
}

// calculateGradeLetter - converts score to letter grade
func (s *GradeService) calculateGradeLetter(score float64) string {
    if score >= 90 {
        return "A"
    } else if score >= 80 {
        return "B"
    } else if score >= 70 {
        return "C"
    } else if score >= 60 {
        return "D"
    } else {
        return "F"
    }
}

// Check if grade already exists (for same student, course, assignment)
func (s *GradeService) gradeExists(studentID, courseID uuid.UUID, assignmentID *uuid.UUID) (bool, error) {
    var count int64
    query := s.db.Model(&models.Grade{}).
        Where("student_id = ? AND course_id = ?", studentID, courseID)
    
    if assignmentID != nil {
        query = query.Where("assignment_id = ?", assignmentID)
    } else {
        query = query.Where("assignment_id IS NULL")
    }
    
    if err := query.Count(&count).Error; err != nil {
        return false, err
    }
    
    return count > 0, nil
}

// Validate grade prerequisites
func (s *GradeService) validateGrade(req models.GradeInput) error {
    // 1. Check if student exists
    var student models.User
    if err := s.db.First(&student, "id = ?", req.StudentID).Error; err != nil {
        return errors.New("student not found")
    }
    
    // 2. Check if course exists
    var course models.Course
    if err := s.db.First(&course, "id = ?", req.CourseID).Error; err != nil {
        return errors.New("course not found")
    }
    
    // 3. Check if student is enrolled in the course
    var enrollmentCount int64
    if err := s.db.Model(&models.Enrollment{}).
        Where("student_id = ? AND course_id = ?", req.StudentID, req.CourseID).
        Count(&enrollmentCount).Error; err != nil {
        return errors.New("failed to check enrollment")
    }
    
    if enrollmentCount == 0 {
        return errors.New("student is not enrolled in this course")
    }
    
    // 4. If assignment ID is provided, check if it exists and belongs to course
    if req.AssignmentID != nil {
        var assignment models.Assignment
        if err := s.db.First(&assignment, "id = ?", req.AssignmentID).Error; err != nil {
            return errors.New("assignment not found")
        }
        
        if assignment.CourseID != req.CourseID {
            return errors.New("assignment does not belong to this course")
        }
        
    }
    
    // 5. Check for duplicate grade
    exists, err := s.gradeExists(req.StudentID, req.CourseID, req.AssignmentID)
    if err != nil {
        return errors.New("failed to check for duplicate grade")
    }
    if exists {
        return errors.New("grade already exists for this student, course, and assignment")
    }
    
    return nil
}

// CreateGrade - main create function
func (s *GradeService) CreateGrade(req models.GradeInput) (*models.GradeResponse, error) {
    // Validate inputs
    if req.Score < 0 || req.Score > 100 {
        return nil, errors.New("score must be between 0 and 100")
    }
    
    // Validate prerequisites
    if err := s.validateGrade(req); err != nil {
        return nil, err
    }
    
    // Calculate grade letter
    gradeLetter := s.calculateGradeLetter(req.Score)
    
    // Create grade record
    grade := models.Grade{
        ID:           uuid.New(),
        StudentID:    req.StudentID,
        CourseID:     req.CourseID,
        TutorID:      req.TutorID,
        AssignmentID: req.AssignmentID,
        Score:        req.Score,
        Grade:        gradeLetter,
        Remarks:      strings.TrimSpace(req.Remarks),
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
    
    // Save to database
    if err := s.db.Create(&grade).Error; err != nil {
        return nil, errors.New("failed to save grade: " + err.Error())
    }
    
    // Recalculate course average if needed
    go s.recalculateCourseAverage(req.CourseID, req.StudentID)
    
    // Convert to response
    return s.gradeToResponse(&grade), nil
}

// CreateGradeWithTx - for use with transactions
func (s *GradeService) CreateGradeWithTx(tx *gorm.DB, req models.GradeInput) (*models.GradeResponse, error) {
    // Same logic but using transaction
    if err := s.validateGrade(req); err != nil {
        return nil, err
    }
    
    gradeLetter := s.calculateGradeLetter(req.Score)
    
    grade := models.Grade{
        ID:           uuid.New(),
        StudentID:    req.StudentID,
        CourseID:     req.CourseID,
        TutorID:      req.TutorID,
        AssignmentID: req.AssignmentID,
        Score:        req.Score,
        Grade:        gradeLetter,
        Remarks:      strings.TrimSpace(req.Remarks),
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
    
    if err := tx.Create(&grade).Error; err != nil {
        return nil, errors.New("failed to save grade: " + err.Error())
    }
    
    return s.gradeToResponse(&grade), nil
}

// Helper to convert Grade to GradeResponse
func (s *GradeService) gradeToResponse(grade *models.Grade) *models.GradeResponse {
    // Optionally fetch related data
    var studentName, courseName string
    
    // You can preload these or fetch separately
    var student models.User
    if err := s.db.First(&student, "id = ?", grade.StudentID).Error; err == nil {
        studentName = fmt.Sprintf("%s %s", student.FirstName, student.LastName)
    }
    
    var course models.Course
    if err := s.db.First(&course, "id = ?", grade.CourseID).Error; err == nil {
        courseName = course.Title
    }
    
    return &models.GradeResponse{
        ID:           grade.ID,
        StudentID:    grade.StudentID,
        StudentName:  studentName,
        CourseID:     grade.CourseID,
        TutorID:      grade.TutorID,
        CourseName:   courseName,
        AssignmentID: grade.AssignmentID,
        Score:        grade.Score,
        Grade:        grade.Grade,
        Percentage:   grade.Score, // Already 0-100
        Remarks:      grade.Remarks,
        CreatedAt:    grade.CreatedAt,
        UpdatedAt:    grade.UpdatedAt,
    }
}

// Recalculate course average (can be called asynchronously)
func (s *GradeService) recalculateCourseAverage(courseID, studentID uuid.UUID) {
    var average float64
    
    // Calculate average of all grades for this student in this course
    err := s.db.Model(&models.Grade{}).
        Where("student_id = ? AND course_id = ?", studentID, courseID).
        Select("COALESCE(AVG(score), 0)").Scan(&average).Error
    
    if err != nil {
        // Log error but don't fail
        fmt.Printf("Failed to recalculate average: %v\n", err)
    }
    
    // You could update a student_course_average table here
    fmt.Printf("Average for student %s in course %s: %.2f\n", 
        studentID, courseID, average)
}

// BulkCreateGrades - for creating multiple grades at once
func (s *GradeService) BulkCreateGrades(requests []models.GradeInput) ([]models.GradeResponse, error) {
    tx := s.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    responses := make([]models.GradeResponse, 0, len(requests))
    
    for _, req := range requests {
        response, err := s.CreateGradeWithTx(tx, req)
        if err != nil {
            tx.Rollback()
            return nil, fmt.Errorf("failed to create grade: %v", err)
        }
        responses = append(responses, *response)
    }
    
    if err := tx.Commit().Error; err != nil {
        return nil, errors.New("failed to commit grades: " + err.Error())
    }
    
    return responses, nil
}