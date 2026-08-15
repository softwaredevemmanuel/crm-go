// models/address.go
package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Address struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Address     string         `gorm:"type:text;not null" json:"address"`
	City        string         `gorm:"type:varchar(100)" json:"city"`
	State       string         `gorm:"type:varchar(100)" json:"state"`
	Country     string         `gorm:"type:varchar(100);default:'Nigeria'" json:"country"`
	PostalCode  string         `gorm:"type:varchar(20)" json:"postal_code"`
	AddressType string         `gorm:"type:varchar(50);default:'home';check:address_type IN ('home', 'school', 'office', 'other')" json:"address_type"`
	IsPrimary   bool           `gorm:"default:false" json:"is_primary"`
	Status      string         `gorm:"type:varchar(20);default:'active';check:status IN ('active', 'inactive')" json:"status"`
	CreatedBy   uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name
func (Address) TableName() string {
	return "addresses"
}