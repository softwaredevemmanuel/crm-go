package models

import (
	"time"
	"github.com/google/uuid"


)
type Grade struct {
    ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    StudentID  uuid.UUID `gorm:"type:uuid;not null"`   // link to student
    CourseID   uuid.UUID `gorm:"type:uuid;not null"`   // link to course
    TutorID    uuid.UUID `gorm:"type:uuid;not null"`   // link to tutor
	AssignmentID     *uuid.UUID     `gorm:"type:uuid;index"`      // optional link to assignment
    Score      float64   `gorm:"not null"`             // raw score (e.g., 85.5)
    Grade      string    `gorm:"type:varchar(5)"`      // A, B, C, D, F
    Remarks    string    `gorm:"type:text"`            // optional feedback

	    
    // Relationships
    Student          User           `gorm:"foreignKey:StudentID"`
    Tutor           User           `gorm:"foreignKey:TutorID"`
    Course           Course         `gorm:"foreignKey:CourseID"`
    Assignment       Assignment     `gorm:"foreignKey:AssignmentID"`
	
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

// GradeInput - for creating/updating grades
type GradeInput struct {
    StudentID    uuid.UUID  `json:"student_id" binding:"required"`
    CourseID     uuid.UUID  `json:"course_id" binding:"required"`
    TutorID      uuid.UUID  `json:"tutor_id" binding:"required"`
    AssignmentID *uuid.UUID `json:"assignment_id"` // Optional
    Score        float64    `json:"score" binding:"required,min=0,max=100"`
    Remarks      string     `json:"remarks" binding:"max=500"`
}

// GradeResponse - for API responses
type GradeResponse struct {
    ID           uuid.UUID  `json:"id"`
    StudentID    uuid.UUID  `json:"student_id"`
    TutorID      uuid.UUID  `json:"tutor_id"`
    StudentName  string     `json:"student_name,omitempty"` // Optional, can be populated
    CourseID     uuid.UUID  `json:"course_id"`
    CourseName   string     `json:"course_name,omitempty"`  // Optional
    AssignmentID *uuid.UUID `json:"assignment_id,omitempty"`
    Score        float64    `json:"score"`
    Grade        string     `json:"grade"`     // A, B, C, etc.
    Percentage   float64    `json:"percentage"` // Score as percentage
    Remarks      string     `json:"remarks"`
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
}

// GradeUpdateInput - for updating grades
type GradeUpdateInput struct {
    TutorID      uuid.UUID  `json:"tutor_id"`
    Score        *float64   `json:"score" binding:"omitempty,min=0,max=100"` // Pointer to distinguish between 0 and not provided
    Remarks      *string    `json:"remarks" binding:"omitempty,max=500"`     // Pointer for optional update
    AssignmentID *uuid.UUID `json:"assignment_id"`                           // Can change assignment link
}

// models/grade.go - add these
type BulkGradeUpdate struct {
    GradeID uuid.UUID  `json:"grade_id" binding:"required"`
    Score   *float64   `json:"score" binding:"omitempty,min=0,max=100"`
    Remarks *string    `json:"remarks" binding:"omitempty,max=500"`
}

type GradeHistory struct {
    ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    GradeID   uuid.UUID `gorm:"type:uuid;not null"`
    Field     string    `gorm:"type:varchar(50);not null"` // "score", "remarks", etc.
    OldValue  string    `gorm:"type:text"`
    NewValue  string    `gorm:"type:text;not null"`
    ChangedBy uuid.UUID `gorm:"type:uuid;not null"` // User who made change
    CreatedAt time.Time
}