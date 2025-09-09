package models

import (
	"time"
	"github.com/google/uuid"
)

type Chapter struct {
    ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    
    // Course Relationship
    CourseID        uuid.UUID      `gorm:"type:uuid;not null;index"`
    
    // Chapter Identification
    Title           string         `gorm:"type:varchar(255);not null"`
	Slug            string         `gorm:"type:varchar(300);not null;index"` // introduction-to-python
    Description     string         `gorm:"type:text"`
    
    // Chapter Organization
    ChapterNumber   int            `gorm:"default:0;index"` // Display number (Chapter 1, Chapter 2)
    
    // Access Control
    IsFree          bool           `gorm:"default:false"` // Free preview chapter
	 // Status & Workflow
    Status          string         `gorm:"type:varchar(20);default:'draft';check:status IN ('draft', 'review', 'approved', 'published', 'archived')"`
    
    
   // Content Details
    EstimatedTime   int            `gorm:"default:0"` // Estimated minutes to complete
    TotalLessons    int            `gorm:"default:0"` // Auto-calculated lesson count
    TotalDuration   int            `gorm:"default:0"` // Auto-calculated total minutes
    
    // Relationships
    Course          Course         `gorm:"foreignKey:CourseID"`
    Lessons         []Lesson       `gorm:"foreignKey:ChapterID"`
    
    // Timestamps
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

// TableName specifies the table name
func (Chapter) TableName() string {
    return "chapters"
}