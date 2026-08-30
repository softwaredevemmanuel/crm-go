// models/assignment.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Assignment struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SchemeOfWorkItemID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"scheme_of_work_item_id"`
	LessonID             uuid.UUID      `gorm:"type:uuid;index" json:"lesson_id"`
	TeacherID            uuid.UUID      `gorm:"type:uuid;not null;index" json:"teacher_id"`
	ClassID              uuid.UUID      `gorm:"type:uuid;not null;index" json:"class_id"`
	ArmID                uuid.UUID      `gorm:"type:uuid;index" json:"arm_id"`
	Title                string         `gorm:"type:varchar(255);not null" json:"title"`
	Description          string         `gorm:"type:text" json:"description"`
	Type                 string         `gorm:"type:varchar(30);not null;default:'assignment'" json:"type"`
	AssignedDate         *time.Time     `gorm:"type:date" json:"assigned_date,omitempty"`
	DueDate              *time.Time     `gorm:"type:date" json:"due_date,omitempty"`
	TotalMarks           float64        `gorm:"type:decimal(10,2)" json:"total_marks"`
	Status               string         `gorm:"type:varchar(20);default:'draft'" json:"status"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	SchemeOfWorkItem SchemeOfWorkItem `gorm:"foreignKey:SchemeOfWorkItemID" json:"scheme_of_work_item,omitempty"`
	Lesson           Lesson           `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
	Teacher          User             `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Class            ClassGrade       `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Arm              Arm              `gorm:"foreignKey:ArmID" json:"arm,omitempty"`
}

func (Assignment) TableName() string {
	return "assignments"
}