// models/exam.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Exam struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AcademicSessionID uuid.UUID      `gorm:"type:uuid;not null;index" json:"academic_session_id"`
	TermID            uuid.UUID      `gorm:"type:uuid;not null;index" json:"term_id"`
	SubjectID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"subject_id"`
	ClassID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"class_id"`
	ArmID             uuid.UUID      `gorm:"type:uuid;index" json:"arm_id"`
	Title             string         `gorm:"type:varchar(255);not null" json:"title"`
	ExamType          string         `gorm:"type:varchar(50)" json:"exam_type"`
	ExamDate          *time.Time     `gorm:"type:date" json:"exam_date,omitempty"`
	Duration          int            `gorm:"type:int" json:"duration"`
	TotalMarks        float64        `gorm:"type:decimal(10,2)" json:"total_marks"`
	Status            string         `gorm:"type:varchar(20);default:'draft'" json:"status"`
	CreatedBy         uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	AcademicSession AcademicSession    `gorm:"foreignKey:AcademicSessionID" json:"academic_session,omitempty"`
	Term            Term               `gorm:"foreignKey:TermID" json:"term,omitempty"`
	Subject         Subject            `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	Class           ClassGrade         `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Arm             Arm                `gorm:"foreignKey:ArmID" json:"arm,omitempty"`
	Creator         User               `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	ExamSchemeItems []ExamSchemeItem   `gorm:"foreignKey:ExamID" json:"exam_scheme_items,omitempty"`
}

func (Exam) TableName() string {
	return "exams"
}