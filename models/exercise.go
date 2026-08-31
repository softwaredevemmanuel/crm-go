// models/exercise.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Exercise struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	LessonID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"lesson_id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Instructions string        `gorm:"type:text" json:"instructions"`
	Content     string         `gorm:"type:text" json:"content"`
	TotalMarks  float64        `gorm:"type:decimal(10,2)" json:"total_marks"`
	CreatedBy            uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Lesson Lesson `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
	Creator User   `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`

}

func (Exercise) TableName() string {
	return "exercises"
}