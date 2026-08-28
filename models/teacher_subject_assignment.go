// models/teacher_subject_assignment.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeacherSubjectAssignment struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	GradeID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"grade_id"`
	SubjectID uuid.UUID      `gorm:"type:uuid;not null;index" json:"subject_id"`
	TeacherID uuid.UUID      `gorm:"type:uuid;not null;index" json:"teacher_id"`
	Status    string         `gorm:"type:varchar(20);default:'active'" json:"status"` // active, inactive
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Grade   ClassGrade `gorm:"foreignKey:GradeID" json:"grade,omitempty"`
	Subject Subject    `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	Teacher User       `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
}

func (TeacherSubjectAssignment) TableName() string {
	return "teacher_subject_assignments"
}