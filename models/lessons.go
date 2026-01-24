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
    CreatedAt       time.Time
    UpdatedAt       time.Time
    
    // Relationships
    Chapter   Chapter   `gorm:"foreignKey:ChapterID"`                      // Belongs to module
}

type LessonInput struct {
	ChapterID   uuid.UUID `json:"chapter_id" binding:"required"`
	CourseID    uuid.UUID `json:"course_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	ContentType string    `json:"content_type" binding:"required"` // video, pdf, text, quiz
	ContentURL  string    `json:"content_url" binding:"required"`
}

type LessonResponse struct {
	ID          uuid.UUID `json:"id"`
	ChapterID   uuid.UUID `json:"chapter_id"`
	CourseID    uuid.UUID `json:"course_id"`
	Title       string    `json:"title"`
	ContentType string    `json:"content_type"`
	ContentURL  string    `json:"content_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName specifies the table name
func (Lesson) TableName() string {
    return "lessons"
}