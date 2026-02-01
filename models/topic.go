package models

import (
	"time"
	"github.com/google/uuid"
)

type Topic struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CourseID  uuid.UUID `gorm:"type:uuid;not null;index"`
	ChapterID uuid.UUID `gorm:"type:uuid;not null;index"`
	TutorID   uuid.UUID `gorm:"type:uuid;not null;index"`

	Title       string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
	Order       int    `gorm:"not null"` // Controls topic sequence within a chapter

	CreatedAt time.Time
	UpdatedAt time.Time

	// Relationships
	Course  Course  `gorm:"foreignKey:CourseID"`
	Chapter Chapter `gorm:"foreignKey:ChapterID"`
	Lessons []Lessons `gorm:"foreignKey:TopicID"`
}

type TopicInput struct {
	CourseID  uuid.UUID `json:"course_id" binding:"required"`
	ChapterID uuid.UUID `json:"chapter_id" binding:"required"`
	TutorID   uuid.UUID `json:"tutor_id" binding:"required"`
	Title     string    `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Order     int       `json:"order" binding:"required"`
}

type TopicResponse struct {
	ID        uuid.UUID `json:"id"`
	CourseID  uuid.UUID `json:"course_id"`
	ChapterID uuid.UUID `json:"chapter_id"`
	TutorID   uuid.UUID `json:"tutor_id"`
	Title     string    `json:"title"`
	Description string  `json:"description"`
	Order     int       `json:"order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TopicMiniResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Order     int       `json:"order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TopicViewResponse struct {
	ID          uuid.UUID `json:"id"`
	CourseID    uuid.UUID `json:"course_id"`
	ChapterID   uuid.UUID `json:"chapter_id"`
	TutorID     uuid.UUID `json:"tutor_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Order       int       `json:"order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Course      CourseMiniResponse                   `json:"course"`
	Chapter     *ChapterMiniResponse                 `json:"chapter,omitempty"`
	Tutor      *UserResponse                          `json:"tutor,omitempty"`
	Lessons     []LessonMiniResponse                 `json:"lessons,omitempty"`
}

// TopicFilters for querying topics
type TopicFilters struct {
    CourseID  uuid.UUID `form:"course_id"`
    ChapterID uuid.UUID `form:"chapter_id"`
    TutorID   uuid.UUID
    Search    string    `form:"search"`
    SortBy    string    `form:"sort_by"`    // title, order, created_at, updated_at
    SortOrder string    `form:"sort_order"` // asc, desc
    Page      int       `form:"page,default=1"`
    Limit     int       `form:"limit,default=10"`
}

// PaginatedTopicsResponse for paginated results
type PaginatedTopicsResponse struct {
    Data       []TopicResponse `json:"data"`
    Total      int64           `json:"total"`
    Page       int             `json:"page"`
    Limit      int             `json:"limit"`
    TotalPages int             `json:"total_pages"`
}