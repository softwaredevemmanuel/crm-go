package models

import (
	"time"
	"github.com/google/uuid"


)

type Assignment struct {
    ID                 uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    
    // Course Hierarchy
    CourseID           uuid.UUID      `gorm:"type:uuid;not null;index"`
    ChapterID          *uuid.UUID     `gorm:"type:uuid;index"`        // Optional chapter association
    LessonID           *uuid.UUID     `gorm:"type:uuid;index"`        // Optional module association
    Publisher          User           `gorm:"foreignKey:PublishedBy"`
    
    // Assignment Identification
    Title              string         `gorm:"type:varchar(255);not null"`
    Slug               string         `gorm:"type:varchar(300);uniqueIndex;not null"` // URL-friendly
    Description        string         `gorm:"type:text"`                              // Detailed instructions
    LearningObjectives []string       `gorm:"type:text[]"`                            // What students will learn
    
    // Assignment Type & Category
    Type               string         `gorm:"type:varchar(50);default:'homework';check:type IN ('homework', 'project', 'essay', 'quiz', 'exam', 'lab', 'presentation', 'discussion', 'peer_review', 'group', 'research', 'creative')"`

    // Timing & Scheduling
    DueDate            time.Time      `gorm:"not null;index"`                         // Submission deadline

 
    
 
    // Relationships
    Course             Course         `gorm:"foreignKey:CourseID"`
    Chapter            Chapter        `gorm:"foreignKey:ChapterID"`
    Lesson             Lesson        `gorm:"foreignKey:LessonID"`
    Submissions        []AssignmentSubmission `gorm:"foreignKey:AssignmentID"`
    
    // Timestamps
    CreatedAt          time.Time
    UpdatedAt          time.Time
    ApprovedAt         *time.Time
    PublishedAt        *time.Time
    ArchivedAt         *time.Time
}

// TableName specifies the table name
func (Assignment) TableName() string {
    return "assignments"
}