// models/class_grade.go
package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ClassGrade struct {
	ID                uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name              string         `gorm:"type:varchar(255);not null" json:"name"`
	Code              string         `gorm:"type:varchar(50);not null;" json:"code"`
	Level             int            `gorm:"not null" json:"level"` // 1, 2, 3, 4, 5, 6
	Description       string         `gorm:"type:text" json:"description"`
	AcademicSessionID uuid.UUID       `gorm:"type:uuid;not null;index" json:"academic_session_id"`
	ClassTeacherID    uuid.UUID       `gorm:"type:uuid;not null;index" json:"class_teacher_id"`
	Capacity          int            `gorm:"default:30" json:"capacity"`
	Status            string         `gorm:"type:varchar(20);default:'active';check:status IN ('active', 'inactive', 'completed')" json:"status"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	AcademicSession AcademicSession `gorm:"foreignKey:AcademicSessionID" json:"academic_session,omitempty"`
	ClassTeacher User `gorm:"foreignKey:ClassTeacherID" json:"user,omitempty"`
}

// TableName specifies the table name
func (ClassGrade) TableName() string {
	return "class_grades"
}
