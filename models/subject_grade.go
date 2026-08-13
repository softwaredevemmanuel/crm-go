// models/subject_grade.go
package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubjectGrade struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	SubjectID   uuid.UUID      `gorm:"type:uuid;not null;index:idx_subject_grade_unique,unique" json:"subject_id"`
	GradeID     uuid.UUID      `gorm:"type:uuid;not null;index:idx_subject_grade_unique,unique" json:"grade_id"`
	AcademicYear string        `gorm:"type:varchar(20);not null" json:"academic_year"`
	Status      string         `gorm:"type:varchar(20);default:'active';check:status IN ('active', 'inactive', 'archived')" json:"status"`
	IsRequired  bool           `gorm:"default:false" json:"is_required"`
	Credits     int            `gorm:"default:0" json:"credits"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedBy   uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Subject Subject    `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	Grade   ClassGrade `gorm:"foreignKey:GradeID" json:"grade,omitempty"`
}

// TableName specifies the table name
func (SubjectGrade) TableName() string {
	return "subject_grades"
}