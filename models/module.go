// models/module.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Module struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SubjectID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"subject_id"`
	Name        string         `gorm:"type:varchar(200);not null" json:"name"`
	Code        string         `gorm:"type:varchar(50)" json:"code"`
	Description string         `gorm:"type:text" json:"description"`
	Sequence    int            `gorm:"type:int;default:1" json:"sequence"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Subject Subject `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	SchemeItems []SchemeOfWorkItem `gorm:"foreignKey:ModuleID" json:"scheme_items,omitempty"`
}

func (Module) TableName() string {
	return "modules"
}