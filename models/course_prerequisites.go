package models

import (
	"time"
	
	"github.com/google/uuid"
)

type CoursePrerequisite struct {
    ID                uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    
    // Course Relationships
    CourseID          uuid.UUID   `gorm:"type:uuid;not null;index"` // The course that requires prerequisites
    PrerequisiteID    uuid.UUID   `gorm:"type:uuid;not null;index"` // The prerequisite course
    
    // Timestamps
    CreatedAt         time.Time
    UpdatedAt         time.Time
}

// TableName specifies the table name
func (CoursePrerequisite) TableName() string {
    return "course_prerequisites"
}