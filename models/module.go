package models

import (
	"time"

	"github.com/google/uuid"
)

type Module struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CourseID      uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_course_module_number"`
	Title         string    `gorm:"type:varchar(255);not null"`
	Slug          string    `gorm:"type:varchar(300);not null;index"` // introduction-to-python
	Description   string    `gorm:"type:text"`
	ModuleNumber  int       `gorm:"default:1;index;uniqueIndex:idx_course_module_number"`
	IsFree        bool      `gorm:"default:false"` // Free preview module
	Status        string    `gorm:"type:varchar(20);default:'draft';check:status IN ('draft', 'review', 'approved', 'published', 'archived')"`
	EstimatedTime int       `gorm:"default:0"` // Estimated minutes to complete
	TotalTopics   int       `gorm:"default:0"` // Auto-calculated topic count
	TotalDuration int       `gorm:"default:0"` // Auto-calculated total minutes

	// Relationships
	Course Course    `gorm:"foreignKey:CourseID"`
	Topics *[]Topics `gorm:"foreignKey:ModuleID"`

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ModuleInput struct {
	CourseID      uuid.UUID `json:"course_id" binding:"required"`
	Title         string    `json:"title" binding:"required"`
	Slug          string    `json:"slug" binding:"required"`
	Description   string    `json:"description"`
	ModuleNumber  int       `json:"module_number"`
	IsFree        bool      `json:"is_free"`
	Status        string    `json:"status" default:"draft"`
	EstimatedTime int       `json:"estimated_time"`
}

type ModuleResponse struct {
	ID            uuid.UUID `json:"id"`
	CourseID      uuid.UUID `json:"course_id"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	ModuleNumber  int       `json:"module_number"`
	IsFree        bool      `json:"is_free"`
	Status        string    `json:"status"`
	EstimatedTime int       `json:"estimated_time"`
	TotalTopics   int       `json:"total_topics"`
	TotalDuration int       `json:"total_duration"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ModuleMiniResponse struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	ModuleNumber int       `json:"module_number"`
}

type ModuleViewResponse struct {
	ID            uuid.UUID `json:"id"`
	CourseID      uuid.UUID `json:"course_id"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	ModuleNumber  int       `json:"module_number"`
	IsFree        bool      `json:"is_free"`
	Status        string    `json:"status"`
	EstimatedTime int       `json:"estimated_time"`
	TotalDuration int       `json:"total_duration"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Course CourseMiniResponse `json:"course"`
	Topics *TopicMiniResponse `json:"topics,omitempty"`
}

// TableName specifies the table name
func (Module) TableName() string {
	return "modules"
}
