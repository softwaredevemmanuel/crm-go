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
	Grade              string         `gorm:"type:varchar(50);not null;index" json:"grade"`
	Term               string         `gorm:"type:varchar(20);not null;check:term IN ('first','second','third')" json:"term"`
	Week               int            `gorm:"type:int;not null" json:"week"`
	Topic              string         `gorm:"type:varchar(255);not null" json:"topic"`
	Subtopics          string         `gorm:"type:text" json:"subtopics"`
	Objectives         string         `gorm:"type:text;not null" json:"objectives"`
	Activities         string         `gorm:"type:text" json:"activities"`
	TeachingResources  string         `gorm:"type:text" json:"teaching_resources"`
	Assessment         string         `gorm:"type:text" json:"assessment"`
	Status             string         `gorm:"type:varchar(20);default:'draft'" json:"status"`
	CreatedBy          uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Subject Subject `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	Creator User    `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (SchemeOfWork) TableName() string {
	return "schemes_of_work"
}