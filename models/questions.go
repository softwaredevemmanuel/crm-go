package models

import (
	"time"
	"github.com/google/uuid"


)

type Question struct {
    ID               uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    
    // Basic Info
    QuestionText     string         `gorm:"type:text;not null"`
    QuestionType     string         `gorm:"type:varchar(20);default:'mcq';check:question_type IN ('mcq', 'true_false', 'multiple_response')"`
    Points           int            `gorm:"default:1"`
    
    // Course Relationship
    CourseID         uuid.UUID      `gorm:"type:uuid;not null;index"`
    ChapterID        *uuid.UUID     `gorm:"type:uuid;index"`
    
    // Options & Answers
    Options          map[string]interface{} `gorm:"type:jsonb;not null"` // Multiple choice options
    CorrectAnswers   map[string]interface{} `gorm:"type:jsonb;not null"` // Correct answer(s)
    Explanation      string         `gorm:"type:text"`           // Answer explanation
    
    // Media Support
    ImageURL         string         `gorm:"type:varchar(500)"`   // Question image
    
    
    // Timestamps
    CreatedAt        time.Time
    UpdatedAt        time.Time
}

// TableName specifies the table name
func (Question) TableName() string {
    return "questions"
}