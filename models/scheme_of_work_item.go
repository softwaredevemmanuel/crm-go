// models/scheme_of_work_item.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SchemeOfWorkItem struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SchemeOfWorkID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"scheme_of_work_id"`
	ModuleID        uuid.UUID      `gorm:"type:uuid;index" json:"module_id"`
	WeekStart       int            `gorm:"type:int;not null" json:"week_start"`
	WeekEnd         int            `gorm:"type:int" json:"week_end"`
	Topic           string         `gorm:"type:varchar(255);not null" json:"topic"`
	Subtopic        string         `gorm:"type:varchar(255)" json:"subtopic"`
	Content         string         `gorm:"type:text" json:"content"`
	Sequence        int            `gorm:"type:int;default:1" json:"sequence"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	SchemeOfWork        SchemeOfWork          `gorm:"foreignKey:SchemeOfWorkID" json:"scheme_of_work,omitempty"`
	Module              Module                `gorm:"foreignKey:ModuleID" json:"module,omitempty"`
	LearningObjectives  []LearningObjective   `gorm:"foreignKey:SchemeOfWorkItemID" json:"learning_objectives,omitempty"`
	Lessons             []Lesson              `gorm:"foreignKey:SchemeOfWorkItemID" json:"lessons,omitempty"`
	Assignments         []Assignment          `gorm:"foreignKey:SchemeOfWorkItemID" json:"assignments,omitempty"`
	TestSchemeItems     []TestSchemeItem      `gorm:"foreignKey:SchemeOfWorkItemID" json:"test_scheme_items,omitempty"`
	ExamSchemeItems     []ExamSchemeItem      `gorm:"foreignKey:SchemeOfWorkItemID" json:"exam_scheme_items,omitempty"`
}

func (SchemeOfWorkItem) TableName() string {
	return "scheme_of_work_items"
}