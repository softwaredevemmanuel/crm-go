package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpskillProgram struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(150);not null;uniqueIndex" json:"name"`
	Slug        string    `gorm:"type:varchar(150);not null;uniqueIndex" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	IsActive    bool      `gorm:"not null;default:true;index" json:"is_active"`

	Courses []UpskillProgramCourse `gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE" json:"courses,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (UpskillProgram) TableName() string { return "upskill_programs" }