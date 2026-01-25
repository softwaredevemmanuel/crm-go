package models

import (
	"time"
	"github.com/google/uuid"

)

type CourseMaterial struct {
    ID                 uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    
    // Course Hierarchy
    CourseID           uuid.UUID      `gorm:"type:uuid;not null;index"`
    ChapterID          *uuid.UUID     `gorm:"type:uuid;index"` // Optional chapter association
    LessonID           *uuid.UUID     `gorm:"type:uuid;index"` // Optional lesson association

    
    // Material Identification
    Title              string         `gorm:"type:varchar(255);not null"`
    Description        string         `gorm:"type:text"`
    Slug               string         `gorm:"type:varchar(300);index"` // URL-friendly
    Type               string         `gorm:"type:varchar(50);not null;check:type IN ('document', 'video', 'audio', 'image', 'code', 'presentation', 'spreadsheet', 'archive', 'link', 'external', 'exercise', 'quiz', 'template')"`
    
    // Content Storage
    FileURL            string         `gorm:"type:varchar(500)"` // Direct file URL


    
    // Status & Workflow
    Status             string         `gorm:"type:varchar(20);default:'draft';check:status IN ('draft', 'review', 'approved', 'published', 'archived', 'hidden')"`

 
    
    // Relationships
    Course             Course         `gorm:"foreignKey:CourseID"`
    Chapter            Chapter        `gorm:"foreignKey:ChapterID"`
    Lesson             Lessons         `gorm:"foreignKey:LessonID"`

    
    // Timestamps
    CreatedAt          time.Time

}

// TableName specifies the table name
func (CourseMaterial) TableName() string {
    return "course_materials"
}