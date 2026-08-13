// models/department.go
package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Department struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"name"`
	Code        string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"code"`
	Description string         `gorm:"type:text" json:"description"`
	HeadOfDept  *uuid.UUID     `gorm:"type:uuid;index" json:"head_of_dept,omitempty"`
	Status      string         `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedBy   uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Head     User      `gorm:"foreignKey:HeadOfDept;references:ID" json:"head,omitempty"`
	Subjects []Subject `gorm:"foreignKey:DepartmentID" json:"subjects,omitempty"`
}

func (Department) TableName() string {
	return "departments"
}