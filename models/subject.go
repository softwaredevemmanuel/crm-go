// models/subject.go
package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subject struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"name"`
	Code         string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"code"`
	Description  string         `gorm:"type:text" json:"description"`
	DepartmentID uuid.UUID      `gorm:"type:uuid;not null;index" json:"department_id"`
	Credits      int            `gorm:"default:3" json:"credits"`
	Status       string         `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedBy    uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Department Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Creator         User            `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	// Relationships
	Modules     []Module     `gorm:"foreignKey:SubjectID" json:"modules,omitempty"`
	Schemes     []SchemeOfWork `gorm:"foreignKey:SubjectID" json:"schemes,omitempty"`
	Tests       []Test       `gorm:"foreignKey:SubjectID" json:"tests,omitempty"`
	Exams       []Exam       `gorm:"foreignKey:SubjectID" json:"exams,omitempty"`
}

func (Subject) TableName() string {
	return "subjects"
}