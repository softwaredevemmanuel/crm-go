package models

import (
	"time"
	"github.com/google/uuid"

)

type CourseMaterial struct {
    ID                 uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`    
    CourseID           uuid.UUID      `gorm:"type:uuid;not null;index"`
    ChapterID          *uuid.UUID     `gorm:"type:uuid;index"` // Optional chapter association
    LessonID           *uuid.UUID     `gorm:"type:uuid;index"` // Optional lesson association
	Title              string         `gorm:"type:varchar(255);not null;index:idx_course_material_unique,unique"`
    Description        string         `gorm:"type:text"`
    Slug               string         `gorm:"type:varchar(300);index"` // URL-friendly
    Type               string         `gorm:"type:varchar(50);not null;check:type IN ('document', 'video', 'audio', 'image', 'code', 'presentation', 'spreadsheet', 'archive', 'link', 'external', 'exercise', 'quiz', 'template')"`    
    FileURL            string         `gorm:"type:varchar(500)"` // Direct file URL    
    Status             string         `gorm:"type:varchar(20);default:'draft';check:status IN ('draft', 'review', 'approved', 'published', 'archived', 'hidden')"`
    CreatedAt          time.Time
        
    // Relationships
    Course             Course         `gorm:"foreignKey:CourseID"`
    Chapter            Chapter        `gorm:"foreignKey:ChapterID"`
    Lesson             Lessons         `gorm:"foreignKey:LessonID"`

}
// TableName specifies the table name
func (CourseMaterial) TableName() string {
    return "course_materials"
}

type CreateCourseMaterialRequest struct {
	CourseID    uuid.UUID  `json:"course_id" binding:"required"`
	ChapterID   *uuid.UUID `json:"chapter_id"`
	LessonID    *uuid.UUID `json:"lesson_id"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Type        string     `json:"type" binding:"required"`
	FileURL     string     `json:"file_url"`
	Status      string     `json:"status"` 
}

type CourseMaterialResponse struct {
    ID                 uuid.UUID      `json:"id"`
    CourseID           uuid.UUID      `json:"course_id"`
    ChapterID          *uuid.UUID     `json:"chapter_id"`
    LessonID           *uuid.UUID     `json:"lesson_id"`
	Title              string         `json:"title"`
    Description        string         `json:"description"`
    Slug               string         `json:"slug"`
    Type               string         `json:"type"`
    FileURL            string         `json:"file_url"`
    Status             string         `json:"status"`
    CreatedAt          time.Time
}


type CourseMaterialViewResponse struct {
    ID                 uuid.UUID      `json:"id"`    
    CourseID           uuid.UUID      `json:"course_id"`
    ChapterID          *uuid.UUID     `json:"chapter_id"`
    LessonID           *uuid.UUID     `json:"lesson_id"`
	Title              string         `json:"title"`
    Description        string         `json:"description"`
    Slug               string         `json:"slug"`
    Type               string         `json:"type"`
    FileURL            string         `json:"file_url"`
    Status             string         `json:"status"`
    CreatedAt          time.Time     `json:"created_at"`
        
    // Relationships
    Course             CourseMiniResponse        `gorm:"foreignKey:CourseID"`
    Chapter            *ChapterMiniResponse        `gorm:"foreignKey:ChapterID"`
    Lesson             *LessonMiniResponse         `gorm:"foreignKey:LessonID"`

}

type UpdateCourseMaterialRequest struct {
	ChapterID   *uuid.UUID `json:"chapter_id"`
	LessonID    *uuid.UUID `json:"lesson_id"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Type        string     `json:"type" binding:"required"`
	FileURL     string     `json:"file_url"`
	Status      string     `json:"status"` 
}