package models

import (
	"time"
	"github.com/google/uuid"
)

type Chapter struct {
    ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    CourseID     uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_course_chapter_number"`    
    Title           string         `gorm:"type:varchar(255);not null"`
	Slug            string         `gorm:"type:varchar(300);not null;index"` // introduction-to-python
    Description     string         `gorm:"type:text"`
    ChapterNumber   int            `gorm:"default:1;index;uniqueIndex:idx_course_chapter_number"`
    IsFree          bool           `gorm:"default:false"` // Free preview chapter
    Status          string         `gorm:"type:varchar(20);default:'draft';check:status IN ('draft', 'review', 'approved', 'published', 'archived')"`    
    EstimatedTime   int            `gorm:"default:0"` // Estimated minutes to complete
    TotalLessons    int            `gorm:"default:0"` // Auto-calculated lesson count
    TotalDuration   int            `gorm:"default:0"` // Auto-calculated total minutes
    
    // Relationships
    Course          Course         `gorm:"foreignKey:CourseID"`
    Lessons         *[]Lessons      `gorm:"foreignKey:ChapterID"`

    // Timestamps
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

type ChapterInput struct {
	CourseID      uuid.UUID `json:"course_id" binding:"required"`
	Title         string    `json:"title" binding:"required"`
	Slug          string    `json:"slug" binding:"required"`
	Description   string    `json:"description"`
	ChapterNumber int       `json:"chapter_number"`
	IsFree        bool      `json:"is_free"`
	Status        string    `json:"status" default:"draft"`
	EstimatedTime int       `json:"estimated_time"`
}


type ChapterResponse struct {
	ID            uuid.UUID `json:"id"`
	CourseID      uuid.UUID `json:"course_id"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	ChapterNumber int       `json:"chapter_number"`
	IsFree        bool      `json:"is_free"`
	Status        string    `json:"status"`
	EstimatedTime int       `json:"estimated_time"`
	TotalLessons  int       `json:"total_lessons"`
	TotalDuration int       `json:"total_duration"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ChapterMiniResponse struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	ChapterNumber int       `json:"chapter_number"`

}

type ChapterViewResponse struct {
	ID            uuid.UUID `json:"id"`
	CourseID      uuid.UUID `json:"course_id"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	ChapterNumber int       `json:"chapter_number"`
	IsFree        bool      `json:"is_free"`
	Status        string    `json:"status"`
	EstimatedTime int       `json:"estimated_time"`
	TotalDuration int       `json:"total_duration"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Course      CourseMiniResponse                   `json:"course"`
	Lessons     *LessonMiniResponse                 `json:"lessons,omitempty"`
}

// TableName specifies the table name
func (Chapter) TableName() string {
    return "chapters"
}