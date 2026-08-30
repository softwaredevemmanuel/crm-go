// models/academic_session.go
package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type AcademicSession struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	AcademicYear string         `gorm:"type:varchar(20);not null;" json:"academic_year"`
	Code         string         `gorm:"type:varchar(20);not null;" json:"code"`
	StartDate    time.Time      `gorm:"type:date;not null" json:"start_date"`
	EndDate      time.Time      `gorm:"type:date;not null" json:"end_date"`
	Status       string         `gorm:"type:varchar(20);not null;default:'active';check:status IN ('active', 'inactive', 'completed')" json:"status"`
	IsCurrent    bool           `gorm:"not null;default:false" json:"is_current"`
	Description  string         `gorm:"type:text" json:"description"`
	CreatedBy    uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Terms    []Term    `gorm:"foreignKey:AcademicSessionID" json:"terms,omitempty"`
	Schemes  []SchemeOfWork `gorm:"foreignKey:AcademicSessionID" json:"schemes,omitempty"`
	Tests    []Test    `gorm:"foreignKey:AcademicSessionID" json:"tests,omitempty"`
	Exams    []Exam    `gorm:"foreignKey:AcademicSessionID" json:"exams,omitempty"`
	Creator  User            `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`

}

// TableName specifies the table name
func (AcademicSession) TableName() string {
	return "academic_sessions"
}
