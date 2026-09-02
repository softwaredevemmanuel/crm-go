// models/term.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Term struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AcademicSessionID uuid.UUID      `gorm:"type:uuid;not null;index" json:"academic_session_id"`
	Name              string         `gorm:"type:varchar(50);not null" json:"name"`
	Code              string         `gorm:"type:varchar(30)" json:"code"`
	TermNumber        int            `gorm:"type:int;not null;check:term_number BETWEEN 1 AND 3" json:"term_number"`
	StartDate         *time.Time     `gorm:"type:date" json:"start_date,omitempty"`
	EndDate           *time.Time     `gorm:"type:date" json:"end_date,omitempty"`
	IsCurrent         bool           `gorm:"not null;default:false" json:"is_current"`
	Status            string         `gorm:"type:varchar(20);default:'active';check:status IN ('active', 'inactive', 'completed')" json:"status"`
	Description       string         `gorm:"type:text" json:"description"`
	CreatedBy         uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	AcademicSession AcademicSession `gorm:"foreignKey:AcademicSessionID" json:"academic_session,omitempty"`
	Creator         User            `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	
}

func (Term) TableName() string {
	return "terms"
}