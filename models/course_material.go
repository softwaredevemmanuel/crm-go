package models

import (
	"time"

	"github.com/google/uuid"
)

type CourseMaterial struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CourseID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	ModuleID    *uuid.UUID `gorm:"type:uuid;index"` // Optional module association
	TopicID     *uuid.UUID `gorm:"type:uuid;index"` // Optional topic association
	Title       string     `gorm:"type:varchar(255);not null;index:idx_course_material_unique,unique"`
	Description string     `gorm:"type:text"`
	Slug        string     `gorm:"type:varchar(300);index"` // URL-friendly
	Type        string     `gorm:"type:varchar(50);not null;check:type IN ('document', 'video', 'audio', 'image', 'code', 'presentation', 'spreadsheet', 'archive', 'link', 'external', 'exercise', 'quiz', 'template')"`
	FileURL     string     `gorm:"type:varchar(500)"` // Direct file URL
	Status      string     `gorm:"type:varchar(20);default:'draft';check:status IN ('draft', 'review', 'approved', 'published', 'archived', 'hidden')"`
	CreatedAt   time.Time

	// Relationships
	Course Course `gorm:"foreignKey:CourseID"`
	Module Module `gorm:"foreignKey:ModuleID"`
	Topic  Topics `gorm:"foreignKey:TopicID"`
}

// TableName specifies the table name
func (CourseMaterial) TableName() string {
	return "course_materials"
}

type CreateCourseMaterialRequest struct {
	CourseID    uuid.UUID  `json:"course_id" binding:"required"`
	ModuleID    *uuid.UUID `json:"module_id"`
	TopicID     *uuid.UUID `json:"topic_id"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Type        string     `json:"type" binding:"required"`
	FileURL     string     `json:"file_url"`
	Status      string     `json:"status"`
}

type CourseMaterialResponse struct {
	ID          uuid.UUID  `json:"id"`
	CourseID    uuid.UUID  `json:"course_id"`
	ModuleID    *uuid.UUID `json:"module_id"`
	TopicID     *uuid.UUID `json:"topic_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Slug        string     `json:"slug"`
	Type        string     `json:"type"`
	FileURL     string     `json:"file_url"`
	Status      string     `json:"status"`
	CreatedAt   time.Time
}

type CourseMaterialViewResponse struct {
	ID          uuid.UUID  `json:"id"`
	CourseID    uuid.UUID  `json:"course_id"`
	ModuleID    *uuid.UUID `json:"module_id"`
	TopicID     *uuid.UUID `json:"topic_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Slug        string     `json:"slug"`
	Type        string     `json:"type"`
	FileURL     string     `json:"file_url"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`

	// Relationships
	Course CourseMiniResponse  `gorm:"foreignKey:CourseID"`
	Module *ModuleMiniResponse `gorm:"foreignKey:ModuleID"`
	Topic  *TopicMiniResponse  `gorm:"foreignKey:TopicID"`
}

type UpdateCourseMaterialRequest struct {
	ModuleID    *uuid.UUID `json:"module_id"`
	TopicID     *uuid.UUID `json:"topic_id"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Type        string     `json:"type" binding:"required"`
	FileURL     string     `json:"file_url"`
	Status      string     `json:"status"`
}
