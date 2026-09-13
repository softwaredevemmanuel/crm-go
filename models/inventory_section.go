// models/inventory_section.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryStatus string

const (
	InventoryStatusActive     InventoryStatus = "active"
	InventoryStatusInactive   InventoryStatus = "inactive"
	InventoryStatusClosed     InventoryStatus = "closed"
	InventoryStatusRenovation InventoryStatus = "renovation"
)

type InventorySection struct {
	ID           uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name         string          `gorm:"type:varchar(255);not null" json:"name"`
	Code         string          `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Description  string          `gorm:"type:text" json:"description"`
	Location     string          `gorm:"type:varchar(255)" json:"location"`
	HeadOfDeptID *uuid.UUID      `gorm:"type:uuid;index" json:"head_of_dept_id,omitempty"`
	Status       InventoryStatus `gorm:"type:varchar(20);default:'active';check:status IN ('active','inactive','closed','renovation')" json:"status"`
	Notes        string          `gorm:"type:text" json:"notes"`
	CreatedAt    time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt  `gorm:"index" json:"-"`

	// Relationships — FIX foreign key names to match actual fields
	HeadOfDept   *User                  `gorm:"foreignKey:HeadOfDeptID" json:"head_of_dept,omitempty"`
	Items        []InventoryItem        `gorm:"foreignKey:InventorySectionID" json:"items,omitempty"`         // FIXED
	Transactions []InventoryTransaction `gorm:"foreignKey:InventorySectionID" json:"transactions,omitempty"`  // FIXED
	Requests     []InventoryRequest     `gorm:"foreignKey:InventorySectionID" json:"requests,omitempty"`      // Added
	Audits       []InventoryAudit       `gorm:"foreignKey:InventorySectionID" json:"audits,omitempty"`        // Added
}

func (InventorySection) TableName() string {
	return "inventory_sections"
}