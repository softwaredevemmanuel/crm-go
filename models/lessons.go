// models/lesson.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Lesson struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SchemeOfWorkItemID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"scheme_of_work_item_id"`
	TeacherID            uuid.UUID      `gorm:"type:uuid;not null;index" json:"teacher_id"`
	ClassID              uuid.UUID      `gorm:"type:uuid;not null;index" json:"class_id"`
	ArmID                uuid.UUID      `gorm:"type:uuid;index" json:"arm_id"`
	Title                string         `gorm:"type:varchar(255);not null" json:"title"`
	LessonDate           *time.Time     `gorm:"type:date" json:"lesson_date,omitempty"`
	Week                 int            `gorm:"type:int" json:"week"`
	Period               int            `gorm:"type:int" json:"period"`
	Duration             int            `gorm:"type:int" json:"duration"`
	Status               string         `gorm:"type:varchar(20);default:'planned'" json:"status"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships - using foreign keys
	SchemeOfWorkItem SchemeOfWorkItem `gorm:"foreignKey:SchemeOfWorkItemID" json:"scheme_of_work_item,omitempty"`
	Teacher          User             `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Class            ClassGrade       `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Arm              Arm              `gorm:"foreignKey:ArmID" json:"arm,omitempty"`
	
	// One-to-One relationship with LessonPlan (use separate file or same package)
	// This will be resolved by having LessonPlan reference LessonID
}

func (Lesson) TableName() string {
	return "lessons"
}