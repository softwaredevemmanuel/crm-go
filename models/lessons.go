// models/scheme_of_work.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Lesson struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	SchemeOfWorkID uuid.UUID `gorm:"type:uuid;not null;index" json:"scheme_of_work_id"`
	ModuleID       uuid.UUID `gorm:"type:uuid;not null;index" json:"module_id"`
	TopicID        uuid.UUID `gorm:"type:uuid;not null;index" json:"topic_id"`

	LessonOrder int `gorm:"not null;default:1" json:"lesson_order"`

	Title       string `gorm:"type:varchar(255);not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`

	LessonDate *time.Time     `gorm:"type:date" json:"lesson_date,omitempty"`
	Week       int            `gorm:"default:1" json:"week"`
	Duration   int            `gorm:"default:40" json:"duration"`
	Content    datatypes.JSON `gorm:"type:jsonb" json:"content"`

	Objectives string `gorm:"type:text" json:"objectives"`
	Activities string `gorm:"type:text" json:"activities"`
	Resources  string `gorm:"type:text" json:"resources"`
	Assessment string `gorm:"type:text" json:"assessment"`
	Status     string `gorm:"type:varchar(20);default:'planned'" json:"status"`

	CreatedBy uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	SchemeOfWork SchemeOfWork `gorm:"foreignKey:SchemeOfWorkID" json:"scheme_of_work,omitempty"`
	Creator      User         `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Module       Module       `gorm:"foreignKey:ModuleID" json:"module,omitempty"`
	Topic        Topic        `gorm:"foreignKey:TopicID" json:"topic,omitempty"`
}

func (Lesson) TableName() string {
	return "lessons"
}
