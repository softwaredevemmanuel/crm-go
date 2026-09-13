// models/Inventory_item_image.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryItemImage struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	InventoryItemID uuid.UUID      `gorm:"type:uuid;not null;index" json:"inventory_item_id"`
	ImageURL        string         `gorm:"type:varchar(500);not null" json:"image_url"`
	Caption         string         `gorm:"type:varchar(255)" json:"caption"`
	IsPrimary       bool           `gorm:"default:false" json:"is_primary"`
	DisplayOrder    int            `gorm:"default:0" json:"display_order"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (InventoryItemImage) TableName() string {
	return "inventory_item_images"
}