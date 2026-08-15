// models/guardian.go
package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Guardian struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Occupation   string         `gorm:"type:varchar(255)" json:"occupation"`
	Relationship string         `gorm:"type:varchar(50);not null" json:"relationship"`
	Address      string         `gorm:"type:text" json:"address"`
	StudentID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"student_id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Status       string         `gorm:"type:varchar(20);default:'active';check:status IN ('active', 'inactive')" json:"status"`
	IsPrimary    bool           `gorm:"default:false" json:"is_primary"`
	IsEmergency  bool           `gorm:"default:false" json:"is_emergency"`
	CreatedBy    uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Student User `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	User    User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name
func (Guardian) TableName() string {
	return "guardians"
}