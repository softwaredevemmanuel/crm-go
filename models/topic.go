package models

import (
	"time"

	"github.com/google/uuid"
)

type Lesson struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CourseID uuid.UUID `gorm:"type:uuid;not null;index"`
	ModuleID uuid.UUID `gorm:"type:uuid;not null;index"`
	TutorID  uuid.UUID `gorm:"type:uuid;not null;index"`

	Title       string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
	Order       int    `gorm:"not null"` // Controls lesson sequence within a module

	CreatedAt time.Time
	UpdatedAt time.Time

	// Relationships
	Course Course   `gorm:"foreignKey:CourseID"`
	Module Module   `gorm:"foreignKey:ModuleID"`
	Topics []Topics `gorm:"foreignKey:LessonID"`
}

type LessonInput struct {
	CourseID    uuid.UUID `json:"course_id" binding:"required"`
	ModuleID    uuid.UUID `json:"module_id" binding:"required"`
	TutorID     uuid.UUID `json:"tutor_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Order       int       `json:"order" binding:"required"`
}

type LessonResponse struct {
	ID          uuid.UUID `json:"id"`
	CourseID    uuid.UUID `json:"course_id"`
	ModuleID    uuid.UUID `json:"module_id"`
	TutorID     uuid.UUID `json:"tutor_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Order       int       `json:"order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LessonMiniResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Order     int       `json:"order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LessonViewResponse struct {
	ID          uuid.UUID           `json:"id"`
	CourseID    uuid.UUID           `json:"course_id"`
	ModuleID    uuid.UUID           `json:"module_id"`
	TutorID     uuid.UUID           `json:"tutor_id"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Order       int                 `json:"order"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Course      CourseMiniResponse  `json:"course"`
	Module      *ModuleMiniResponse `json:"module,omitempty"`
	Tutor       *UserResponse       `json:"tutor,omitempty"`
	Topics      []TopicMiniResponse `json:"topics,omitempty"`
}

// LessonFilters for querying lessons
type LessonFilters struct {
	CourseID  uuid.UUID `form:"course_id"`
	ModuleID  uuid.UUID `form:"module_id"`
	TutorID   uuid.UUID
	Search    string `form:"search"`
	SortBy    string `form:"sort_by"`    // title, order, created_at, updated_at
	SortOrder string `form:"sort_order"` // asc, desc
	Page      int    `form:"page,default=1"`
	Limit     int    `form:"limit,default=10"`
}

// PaginatedLessonsResponse for paginated results
type PaginatedLessonsResponse struct {
	Data       []LessonResponse `json:"data"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}
