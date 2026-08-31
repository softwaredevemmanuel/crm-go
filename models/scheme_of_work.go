// models/scheme_of_work.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SchemeOfWork struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AcademicSessionID uuid.UUID      `gorm:"type:uuid;not null;index" json:"academic_session_id"`
	TermID            uuid.UUID      `gorm:"type:uuid;not null;index" json:"term_id"`
	SubjectID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"subject_id"`
	ClassID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"class_id"`
	Title             string         `gorm:"type:varchar(255);not null" json:"title"`
	Description       string         `gorm:"type:text" json:"description"`
	Status            string         `gorm:"type:varchar(20);default:'draft'" json:"status"`
	CreatedBy         uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	AcademicSession AcademicSession       `gorm:"foreignKey:AcademicSessionID" json:"academic_session,omitempty"`
	Term            Term                  `gorm:"foreignKey:TermID" json:"term,omitempty"`
	Subject         Subject               `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	Class           ClassGrade            `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Creator         User                  `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Items           []SchemeOfWorkItem    `gorm:"foreignKey:SchemeOfWorkID" json:"items,omitempty"`
	// ✅ REMOVED: Lessons []Lesson `gorm:"foreignKey:SchemeOfWorkItemID" json:"lessons,omitempty"`
	// ✅ REMOVED: Assignments []Assignment `gorm:"foreignKey:SchemeOfWorkItemID" json:"assignments,omitempty"`
}

func (SchemeOfWork) TableName() string {
	return "schemes_of_work"
}