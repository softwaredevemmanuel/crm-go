package models

import (
	"time"
	"github.com/google/uuid"
)

type Lesson struct {
    ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    ChapterID uuid.UUID `gorm:"type:uuid;not null;index"`                  // Points to module (level=1 chapter)
    CourseID  uuid.UUID `gorm:"type:uuid;not null;index"`                  // Denormalized for performance
    Title     string    `gorm:"type:varchar(255);not null"`

	ContentType string `gorm:"type:varchar(50);not null"`                  // video, pdf, text, quiz
    ContentURL string  `gorm:"type:varchar(500);not null"`                 // Actual content location

	   // Timestamps
    CreatedAt       time.Time
    UpdatedAt       time.Time
    
    // Relationships
    Chapter   Chapter   `gorm:"foreignKey:ChapterID"`                      // Belongs to module
    Course    Course    `gorm:"foreignKey:CourseID"`                       // Belongs to course
}

// TableName specifies the table name
func (Lesson) TableName() string {
    return "lessons"
}