// models/Lesson_report.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LessonReport represents a teacher's Lesson teaching progress report
type LessonReport struct {
	ID                uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TeacherID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"teacher_id"`
	SchemeOfWorkID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"scheme_of_work_id"`
	ModuleID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"module_id"`
	TopicID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"topic_id"`
	LessonID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"lesson_id"`
	AcademicSessionID uuid.UUID      `gorm:"type:uuid;not null;index" json:"academic_session_id"` // Fixed field name
	ReportDate        time.Time      `gorm:"not null;index" json:"report_date"`
	Status            string         `gorm:"type:varchar(20);not null;default:'not_taught';check:status IN ('not_taught','in_progress','completed')" json:"status"`
	StudentsCovered   string         `gorm:"type:text" json:"students_covered"`
	Report            string         `gorm:"type:text" json:"report"`
	Challenges        string         `gorm:"type:text" json:"challenges"`
	Remarks           string         `gorm:"type:text" json:"remarks"`
	StudentsPresent   int            `json:"students_present"`
	StudentsAbsent    int            `json:"students_absent"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Teacher         User            `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	SchemeOfWork    SchemeOfWork    `gorm:"foreignKey:SchemeOfWorkID" json:"scheme_of_work,omitempty"`
	Module          Module          `gorm:"foreignKey:ModuleID" json:"module,omitempty"`
	Topic           Topic           `gorm:"foreignKey:TopicID" json:"topic,omitempty"`
	Lesson          Lesson          `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
	AcademicSession AcademicSession `gorm:"foreignKey:AcademicSessionID" json:"academic_session,omitempty"`
}

func (LessonReport) TableName() string {
	return "lesson_reports"
}