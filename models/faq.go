package models

import (
	"time"

	"github.com/google/uuid"
)

type FAQ struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Question  string    `gorm:"type:varchar(255);not null" json:"question"`
	Answer    string    `gorm:"type:text;not null" json:"answer"`
	Category  string    `gorm:"type:varchar(100)" json:"category"`
	Audience  string    `gorm:"type:varchar(50);default:'all'" json:"audience"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedBy uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (FAQ) TableName() string {
	return "faqs"
}	