// models/inventory_request.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryRequestStatus string

const (
	InventoryRequestStatusPending   InventoryRequestStatus = "pending"
	InventoryRequestStatusApproved  InventoryRequestStatus = "approved"
	InventoryRequestStatusRejected  InventoryRequestStatus = "rejected"
	InventoryRequestStatusFulfilled InventoryRequestStatus = "fulfilled"
	InventoryRequestStatusCancelled InventoryRequestStatus = "cancelled"
)

type InventoryRequest struct {
	ID                 uuid.UUID              `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	RequestCode        string                 `gorm:"type:varchar(100);uniqueIndex;not null" json:"request_code"`
	RequestedByID      uuid.UUID              `gorm:"type:uuid;not null;index" json:"requested_by_id"`
	InventorySectionID uuid.UUID              `gorm:"type:uuid;not null;index" json:"inventory_section_id"`
	InventoryItemID    uuid.UUID              `gorm:"type:uuid;not null;index" json:"inventory_item_id"`
	QuantityRequested  int                    `gorm:"not null;check:quantity_requested > 0" json:"quantity_requested"`
	QuantityApproved   int                    `gorm:"default:0" json:"quantity_approved"`
	QuantityIssued     int                    `gorm:"default:0" json:"quantity_issued"`
	Purpose            string                 `gorm:"type:text" json:"purpose"`
	Status             InventoryRequestStatus `gorm:"type:varchar(20);default:'pending';check:status IN ('pending','approved','rejected','fulfilled','cancelled')" json:"status"`
	ApprovedByID       *uuid.UUID             `gorm:"type:uuid;index" json:"approved_by_id,omitempty"`
	ApprovedAt         *time.Time             `json:"approved_at,omitempty"`
	RejectionReason    string                 `gorm:"type:text" json:"rejection_reason"`
	NeededByDate       *time.Time             `gorm:"type:date" json:"needed_by_date,omitempty"`
	FulfilledAt        *time.Time             `json:"fulfilled_at,omitempty"`
	Notes              string                 `gorm:"type:text" json:"notes"`
	CreatedBy          uuid.UUID              `gorm:"type:uuid;not null;index" json:"created_by"`
	CreatedAt          time.Time              `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time              `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt         `gorm:"index" json:"-"`

	// Relationships
	RequestedBy      User              `gorm:"foreignKey:RequestedByID" json:"requested_by,omitempty"`
	InventorySection InventorySection  `gorm:"foreignKey:InventorySectionID" json:"inventory_section,omitempty"`
	InventoryItem    InventoryItem     `gorm:"foreignKey:InventoryItemID" json:"inventory_item,omitempty"`
	ApprovedBy       *User             `gorm:"foreignKey:ApprovedByID" json:"approved_by,omitempty"`
	Creator          User              `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (InventoryRequest) TableName() string {
	return "inventory_requests"
}