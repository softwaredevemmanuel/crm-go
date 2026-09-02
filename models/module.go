package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Module struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	SchemeOfWorkID uuid.UUID      `gorm:"type:uuid;not null;index" json:"scheme_of_work_id"`
	Title          string         `gorm:"type:varchar(255);not null" json:"title"`
	Description    string         `gorm:"type:text" json:"description"`
	ModuleOrder    int            `gorm:"not null;default:1" json:"module_order"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	SchemeOfWork SchemeOfWork `gorm:"foreignKey:SchemeOfWorkID" json:"scheme_of_work,omitempty"`
	Lessons      []Lesson     `gorm:"foreignKey:ModuleID" json:"lessons,omitempty"`
}

func (Module) TableName() string {
	return "modules"
}