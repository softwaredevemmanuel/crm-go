package models

import (
	"time"

	"github.com/google/uuid"
)

type UpskillProgramCourse struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	ProgramID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_program_course;index" json:"program_id"`
	CourseID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_program_course;index" json:"course_id"`

	CourseOrder int  `gorm:"not null;default:1" json:"course_order"`
	IsRequired  bool `gorm:"not null;default:true" json:"is_required"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Course Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

func (UpskillProgramCourse) TableName() string { return "upskill_program_courses" }