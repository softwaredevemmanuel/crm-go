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
    

    // Answer Configuration
    AnswerExplanation  string         `gorm:"type:text" json:"answer_explanation,omitempty"`
    SolutionSteps      string         `gorm:"type:text" json:"solution_steps,omitempty"`
    Hint               string         `gorm:"type:text" json:"hint,omitempty"`

 
    // Metadata
    TutorID            uuid.UUID      `gorm:"type:uuid;not null" json:"tutor_id"`
    CreatedAt          time.Time      `json:"created_at"`
    UpdatedAt          time.Time      `json:"updated_at"`
    IsApproved         bool           `json:"is_approved"`
    // Relationships
    Course             Course         `gorm:"foreignKey:CourseID" json:"course,omitempty"`

}


// ObjectiveQuestionInput - for creating objective questions
type ObjectiveQuestionInput struct {
    // Required fields
    QuestionText    string    `json:"question_text" binding:"required,min=10,max=1000"`
    CourseID        uuid.UUID `json:"course_id" binding:"required"`
    TutorID        uuid.UUID `json:"tutor_id" binding:"required"`

    // Optional fields with defaults
    QuestionType    string    `json:"question_type" binding:"omitempty,oneof=multiple_choice true_false multiple_response matching ordering"`
    DifficultyLevel string    `json:"difficulty_level" binding:"omitempty,oneof=easy medium hard expert"`
    Points          int       `json:"points" binding:"omitempty,min=0,max=100"`
    
    // Media support
    ImageURL        string    `json:"image_url" binding:"omitempty,url,max=500"`
    VideoURL        string    `json:"video_url" binding:"omitempty,url,max=500"`
    
    // Content relationships
    ChapterID       *uuid.UUID `json:"chapter_id"`
    TopicID         *uuid.UUID `json:"topic_id"`
    
    // Answer configuration
    AnswerExplanation string   `json:"answer_explanation" binding:"max=2000"`
    SolutionSteps     string   `json:"solution_steps" binding:"max=2000"`
    Hint              string   `json:"hint" binding:"max=500"`
    
    // Options for multiple choice questions
    Options         []QuestionOptionInput `json:"options"`
    
    // Approval
    IsApproved      bool      `json:"is_approved"`
}

// QuestionOptionInput - for question options
type QuestionOptionInput struct {
    OptionText      string    `json:"option_text" binding:"required,min=1,max=500"`
    IsCorrect       bool      `json:"is_correct"`
    Explanation     string    `json:"explanation" binding:"max=500"`
    SortOrder       int       `json:"sort_order" binding:"min=0"`
}

// ObjectiveQuestionResponse - for API responses
type ObjectiveQuestionResponse struct {
    ID                 uuid.UUID               `json:"id"`
    QuestionText       string                  `json:"question_text"`
    QuestionType       string                  `json:"question_type"`
    DifficultyLevel    string                  `json:"difficulty_level"`
    Points             int                     `json:"points"`
    
    ImageURL           string                  `json:"image_url,omitempty"`
    VideoURL           string                  `json:"video_url,omitempty"`
    
    CourseID           uuid.UUID               `json:"course_id"`
    CourseName         string                  `json:"course_name,omitempty"`
    ChapterID          *uuid.UUID              `json:"chapter_id,omitempty"`
    ChapterName        string                  `json:"chapter_name,omitempty"`
    TopicID            *uuid.UUID              `json:"topic_id,omitempty"`
    TopicName          string                  `json:"topic_name,omitempty"`
    
    AnswerExplanation  string                  `json:"answer_explanation,omitempty"`
    SolutionSteps      string                  `json:"solution_steps,omitempty"`
    Hint               string                  `json:"hint,omitempty"`

    TutorID            uuid.UUID               `json:"tutor_id"`
    CreatorName        string                  `json:"creator_name,omitempty"`
    CreatedAt          time.Time               `json:"created_at"`
    UpdatedAt          time.Time               `json:"updated_at"`
    IsApproved         bool                    `json:"is_approved"`
    
    // Options (for multiple choice questions)
    Options            []QuestionOptionResponse `json:"options,omitempty"`
    
    // Statistics (optional)
    TotalAttempts      int                     `json:"total_attempts,omitempty"`
    CorrectAttempts    int                     `json:"correct_attempts,omitempty"`
    SuccessRate        float64                 `json:"success_rate,omitempty"`
}

// QuestionOptionResponse - option response
type QuestionOptionResponse struct {
    ID           uuid.UUID `json:"id"`
    QuestionID   uuid.UUID `json:"question_id"`
    OptionText   string    `json:"option_text"`
    IsCorrect    bool      `json:"is_correct"`
    Explanation  string    `json:"explanation,omitempty"`
    SortOrder    int       `json:"sort_order"`
}

// QuestionOption model for database
type QuestionOption struct {
    ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    QuestionID  uuid.UUID `gorm:"type:uuid;not null;index"`
    OptionText  string    `gorm:"type:varchar(500);not null"`
    IsCorrect   bool      `gorm:"default:false"`
    Explanation string    `gorm:"type:text"`
    SortOrder   int       `gorm:"default:0;check:sort_order >= 0"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    
    Question    ObjectiveQuestion `gorm:"foreignKey:QuestionID"`
}

func (QuestionOption) TableName() string {
    return "question_options"
}