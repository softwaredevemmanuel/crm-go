package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Topic struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ModuleID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"module_id"`
	Title      string         `gorm:"type:varchar(255);not null" json:"title"`
	Description string        `gorm:"type:text" json:"description"`
	TopicOrder int            `gorm:"not null;default:1" json:"topic_order"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Module  Module   `gorm:"foreignKey:ModuleID" json:"module,omitempty"`
	Lessons []Lesson `gorm:"foreignKey:TopicID" json:"lessons,omitempty"`
}

func (Topic) TableName() string {
	return "topics"
}