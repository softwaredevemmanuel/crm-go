package models

import (
	"time"

	"github.com/google/uuid"
)

type Topics struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CourseID uuid.UUID `gorm:"type:uuid;not null;index"`
	ModuleID uuid.UUID `gorm:"type:uuid;not null;index"`
	LessonID uuid.UUID `gorm:"type:uuid;index"`

	Title       string `gorm:"type:varchar(255);not null"`
	ContentType string `gorm:"type:varchar(50);not null"`
	ContentURL  string `gorm:"type:varchar(500);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	// Relationships
	Module Module `gorm:"foreignKey:ModuleID"`
	Lesson Lesson `gorm:"foreignKey:LessonID"`
}

type TopicInput struct {
	CourseID    uuid.UUID `json:"course_id" binding:"required"`
	ModuleID    uuid.UUID `json:"module_id" binding:"required"`
	LessonID    uuid.UUID `json:"lesson_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	ContentType string    `json:"content_type" binding:"required"` // video, pdf, text, quiz
	ContentURL  string    `json:"content_url" binding:"required"`
}

type TopicResponse struct {
	ID          uuid.UUID `json:"id"`
	CourseID    uuid.UUID `json:"course_id"`
	ModuleID    uuid.UUID `json:"module_id"`
	LessonID    uuid.UUID `json:"lesson_id"`
	Title       string    `json:"title"`
	ContentType string    `json:"content_type"`
	ContentURL  string    `json:"content_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TopicMiniResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	ContentType string    `json:"content_type"`
	ContentURL  string    `json:"content_url"`
}
type TopicViewResponse struct {
	ID          uuid.UUID           `json:"id"`
	CourseID    uuid.UUID           `json:"course_id"`
	ModuleID    uuid.UUID           `json:"module_id"`
	LessonID    uuid.UUID           `json:"lesson_id"`
	Title       string              `json:"title"`
	ContentType string              `json:"content_type"`
	ContentURL  string              `json:"content_url"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Course      CourseMiniResponse  `json:"course"`
	Module      *ModuleMiniResponse `json:"module,omitempty"`
}

type TopicUpdateInput struct {
	Title       string `json:"title" binding:"required"`
	ContentType string `json:"content_type" binding:"required"` // video, pdf, text, quiz
	ContentURL  string `json:"content_url" binding:"required"`
}

// TableName specifies the table name
func (Topics) TableName() string {
	return "topics"
}
