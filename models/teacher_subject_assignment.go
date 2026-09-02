// models/teacher_subject_assignment.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeacherSubjectAssignment struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SubjectID uuid.UUID      `gorm:"type:uuid;not null;index" json:"subject_id"`
	TeacherID uuid.UUID      `gorm:"type:uuid;not null;index" json:"teacher_id"`
	GradeID uuid.UUID      `gorm:"type:uuid;not null;index" json:"grade_id"`
	Status    string         `gorm:"type:varchar(20);default:'active'" json:"status"` // active, inactive
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Subject Subject    `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	Teacher User       `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Grade   ClassGrade `gorm:"foreignKey:GradeID" json:"grade,omitempty"`
}

func (TeacherSubjectAssignment) TableName() string {
	return "teacher_subject_assignments"
}