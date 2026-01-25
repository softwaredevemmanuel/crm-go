package models

import (
	"time"
	"github.com/google/uuid"
)

type Lessons struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ChapterID   uuid.UUID `gorm:"type:uuid;not null;index"`
	CourseID    uuid.UUID `gorm:"type:uuid;not null;index"`

	Title       string    `gorm:"type:varchar(255);not null"`
	ContentType string    `gorm:"type:varchar(50);not null"`
	ContentURL  string    `gorm:"type:varchar(500);not null"`

	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Relationships
	Chapter Chapter `gorm:"foreignKey:ChapterID"`
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

type LessonMiniResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	ContentType string    `json:"content_type"`
	ContentURL  string    `json:"content_url"`

}
type LessonViewResponse struct {
	ID          uuid.UUID `json:"id"`
	ChapterID   uuid.UUID `json:"chapter_id"`
	CourseID    uuid.UUID `json:"course_id"`
	Title       string    `json:"title"`
	ContentType string    `json:"content_type"`
	ContentURL  string    `json:"content_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Course      CourseMiniResponse                   `json:"course"`
	Chapter     *ChapterMiniResponse                 `json:"chapter,omitempty"`
}

type LessonUpdateInput struct {
	Title       string    `json:"title" binding:"required"`
	ContentType string    `json:"content_type" binding:"required"` // video, pdf, text, quiz
	ContentURL  string    `json:"content_url" binding:"required"`
}

// TableName specifies the table name
func (Lessons) TableName() string {
    return "lessons"
}