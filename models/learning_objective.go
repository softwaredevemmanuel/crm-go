// models/learning_objective.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LearningObjective struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SchemeOfWorkItemID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"scheme_of_work_item_id"`
	Objective            string         `gorm:"type:text;not null" json:"objective"`
	Sequence             int            `gorm:"type:int;default:1" json:"sequence"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	SchemeOfWorkItem SchemeOfWorkItem `gorm:"foreignKey:SchemeOfWorkItemID" json:"scheme_of_work_item,omitempty"`
}

func (LearningObjective) TableName() string {
	return "learning_objectives"
}