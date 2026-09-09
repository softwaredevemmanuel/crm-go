// models/student_answer.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ObjectiveQuestionAnswer struct {
	ID                uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	StudentID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"student_id"`
	QuestionID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"question_id"`
	SelectedOptionIDs string         `gorm:"type:text" json:"selected_option_ids"` // Comma-separated UUIDs for multiple selection
	SelectedOptionID  uuid.UUID      `gorm:"type:uuid" json:"selected_option_id"`  // Single selection
	IsCorrect         bool           `gorm:"default:false" json:"is_correct"`
	Score             int            `gorm:"default:0" json:"score"`
	TimeSpent         int            `gorm:"default:0" json:"time_spent"` // Time spent in seconds
	AttemptNumber     int            `gorm:"default:1" json:"attempt_number"`
	Status            string         `gorm:"type:varchar(20);default:'active';check:status IN ('active','inactive')" json:"status"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Student  User              `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Question ObjectiveQuestion `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
}

func (ObjectiveQuestionAnswer) TableName() string {
	return "objective_question_answers"
}
