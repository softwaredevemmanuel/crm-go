// models/scheme_of_work.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SchemeOfWork struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SubjectID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"subject_id"`
	GradeID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"grade_id"`
	Term               string         `gorm:"type:varchar(20);not null;check:term IN ('first','second','third')" json:"term"`
	Title             string         `gorm:"type:varchar(255);not null" json:"title"`
	Description       string         `gorm:"type:text" json:"description"`
	Status 				string `gorm:"type:varchar(20);default:'draft';check:status IN ('draft','published','archived')" json:"status"`
	CreatedBy          uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Subject Subject `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	Creator User    `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Lessons           []Lesson        `gorm:"foreignKey:SchemeOfWorkID" json:"lessons,omitempty"`
	Grade   ClassGrade `gorm:"foreignKey:GradeID" json:"grade,omitempty"`


}

func (SchemeOfWork) TableName() string {
	return "schemes_of_work"
}