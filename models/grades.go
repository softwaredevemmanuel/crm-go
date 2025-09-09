package models

import (
	"time"
	"github.com/google/uuid"


)
type Grade struct {
    ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    StudentID  uuid.UUID `gorm:"type:uuid;not null"`   // link to student
    CourseID   uuid.UUID `gorm:"type:uuid;not null"`   // link to course
	AssignmentID     *uuid.UUID     `gorm:"type:uuid;index"`      // optional link to assignment
    Score      float64   `gorm:"not null"`             // raw score (e.g., 85.5)
    Grade      string    `gorm:"type:varchar(5)"`      // A, B, C, D, F
    Remarks    string    `gorm:"type:text"`            // optional feedback

	    
    // Relationships
    Student          User           `gorm:"foreignKey:StudentID"`
    Course           Course         `gorm:"foreignKey:CourseID"`
    Assignment       Assignment     `gorm:"foreignKey:AssignmentID"`
	
    CreatedAt  time.Time
    UpdatedAt  time.Time
}
