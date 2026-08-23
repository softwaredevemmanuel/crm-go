// models/grade_subject.go
package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GradeSubject struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	GradeID      uuid.UUID      `gorm:"type:uuid;not null;index:idx_grade_subject_unique,unique" json:"grade_id"`
	SubjectID    uuid.UUID      `gorm:"type:uuid;not null;index:idx_grade_subject_unique,unique" json:"subject_id"`
	Status       string         `gorm:"type:varchar(20);not null;default:'active';check:status IN ('active', 'inactive')" json:"status"`
	IsCompulsory bool           `gorm:"not null;" json:"is_compulsory"`
	CreatedBy    uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Grade   ClassGrade `gorm:"foreignKey:GradeID" json:"grade,omitempty"`
	Subject Subject    `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	User    User       `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
}

// TableName specifies the table name
func (GradeSubject) TableName() string {
	return "grade_subjects"
}