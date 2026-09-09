// models/objective_question.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ObjectiveQuestion struct {
	ID                uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	QuestionText      string         `gorm:"type:text;not null" json:"question_text"`
	QuestionType      string         `gorm:"type:varchar(50);default:'multiple_choice';check:question_type IN ('multiple_choice', 'true_false', 'multiple_response', 'matching', 'ordering')" json:"question_type"`
	DifficultyLevel   string         `gorm:"type:varchar(20);default:'medium';check:difficulty_level IN ('easy', 'medium', 'hard', 'expert')" json:"difficulty_level"`
	Points            int            `gorm:"default:1;check:points >= 0" json:"points"`
	ImageURL          string         `gorm:"type:varchar(500)" json:"image_url,omitempty"`
	VideoURL          string         `gorm:"type:varchar(500)" json:"video_url,omitempty"`
	SubjectID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"subject_id"`
	SchemeOfWorkID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"scheme_of_work_id"`
	ModuleID          uuid.UUID      `gorm:"type:uuid;index" json:"module_id"`
	TopicID           uuid.UUID      `gorm:"type:uuid;index" json:"topic_id"`
	LessonID          uuid.UUID      `gorm:"type:uuid;index" json:"lesson_id"`
	AnswerExplanation string         `gorm:"type:text" json:"answer_explanation,omitempty"`
	SolutionSteps     string         `gorm:"type:text" json:"solution_steps,omitempty"`
	Hint              string         `gorm:"type:text" json:"hint,omitempty"`
	Status            string         `gorm:"type:varchar(20);default:'active';check:status IN ('active','inactive','draft','archived')" json:"status"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Subject      Subject          `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	SchemeOfWork SchemeOfWork     `gorm:"foreignKey:SchemeOfWorkID" json:"scheme_of_work,omitempty"`
	Module       Module           `gorm:"foreignKey:ModuleID" json:"module,omitempty"`
	Topic        Topic            `gorm:"foreignKey:TopicID" json:"topic,omitempty"`
	Lesson       Lesson           `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
	Options      []QuestionOption `gorm:"foreignKey:QuestionID" json:"options,omitempty"`
}

func (ObjectiveQuestion) TableName() string {
	return "objective_questions"
}

type QuestionOption struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	QuestionID uuid.UUID      `gorm:"type:uuid;not null;index" json:"question_id"`
	OptionText string         `gorm:"type:text;not null" json:"option_text"`
	IsCorrect  bool           `gorm:"default:false" json:"is_correct"`
	Order      int            `gorm:"default:0" json:"order"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationship
	Question ObjectiveQuestion `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
}

func (QuestionOption) TableName() string {
	return "question_options"
}