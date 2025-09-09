package models

import (
	"time"
	"github.com/google/uuid"


)

type ObjectiveQuestion struct {
    ID                 uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
    QuestionText       string         `gorm:"type:text;not null" json:"question_text"`
    QuestionType       string         `gorm:"type:varchar(50);default:'multiple_choice';check:question_type IN ('multiple_choice', 'true_false', 'multiple_response', 'matching', 'ordering')" json:"question_type"`
    DifficultyLevel    string         `gorm:"type:varchar(20);default:'medium';check:difficulty_level IN ('easy', 'medium', 'hard', 'expert')" json:"difficulty_level"`
    Points             int            `gorm:"default:1;check:points >= 0" json:"points"`
    
    // Media Support
    ImageURL           string         `gorm:"type:varchar(500)" json:"image_url,omitempty"`
    VideoURL           string         `gorm:"type:varchar(500)" json:"video_url,omitempty"`

    // Content Relationships
    CourseID           uuid.UUID      `gorm:"type:uuid;not null" json:"course_id"`
    ChapterID          *uuid.UUID     `gorm:"type:uuid" json:"chapter_id,omitempty"`
    TopicID            *uuid.UUID     `gorm:"type:uuid" json:"topic_id,omitempty"`
    QuizID             *uuid.UUID     `gorm:"type:uuid" json:"quiz_id,omitempty"`
    QuestionBankID     *uuid.UUID     `gorm:"type:uuid" json:"question_bank_id,omitempty"`
    

    // Answer Configuration
    AnswerExplanation  string         `gorm:"type:text" json:"answer_explanation,omitempty"`
    SolutionSteps      string         `gorm:"type:text" json:"solution_steps,omitempty"`
    Hint               string         `gorm:"type:text" json:"hint,omitempty"`
    MaxAttempts        int            `gorm:"default:1;check:max_attempts >= 1" json:"max_attempts"`
    TimeLimit          int            `gorm:"default:0" json:"time_limit"` // in seconds (0 = no limit)
 

    // Metadata
    CreatedBy          uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
    ReviewedBy         *uuid.UUID     `gorm:"type:uuid" json:"reviewed_by,omitempty"`
    ApprovedBy         *uuid.UUID     `gorm:"type:uuid" json:"approved_by,omitempty"`
    CreatedAt          time.Time      `json:"created_at"`
    UpdatedAt          time.Time      `json:"updated_at"`
    ReviewedAt         *time.Time     `json:"reviewed_at,omitempty"`
    ApprovedAt         *time.Time     `json:"approved_at,omitempty"`
    
    // Relationships
    Course             Course         `gorm:"foreignKey:CourseID" json:"course,omitempty"`

}