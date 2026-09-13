// models/inventory_audit.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryAuditStatus string

const (
	InventoryAuditStatusInProgress InventoryAuditStatus = "in_progress"
	InventoryAuditStatusCompleted  InventoryAuditStatus = "completed"
	InventoryAuditStatusApproved   InventoryAuditStatus = "approved"
	InventoryAuditStatusRejected   InventoryAuditStatus = "rejected"
	InventoryAuditStatusCancelled  InventoryAuditStatus = "cancelled"
)

type InventoryAudit struct {
	ID                 uuid.UUID            `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	AuditCode          string               `gorm:"type:varchar(100);uniqueIndex;not null" json:"audit_code"`
	InventorySectionID uuid.UUID            `gorm:"type:uuid;not null;index" json:"inventory_section_id"`
	AuditDate          time.Time            `gorm:"type:date;not null" json:"audit_date"`
	ConductedByID      uuid.UUID            `gorm:"type:uuid;not null;index" json:"conducted_by_id"`
	Status             InventoryAuditStatus `gorm:"type:varchar(20);default:'in_progress';check:status IN ('in_progress','completed','approved','rejected','cancelled')" json:"status"`
	TotalItemsChecked  int                  `gorm:"default:0" json:"total_items_checked"`
	TotalDiscrepancies int                  `gorm:"default:0" json:"total_discrepancies"`
	TotalItemsMissing  int                  `gorm:"default:0" json:"total_items_missing"`
	TotalItemsExcess   int                  `gorm:"default:0" json:"total_items_excess"`
	TotalItemsDamaged  int                  `gorm:"default:0" json:"total_items_damaged"`
	Notes              string               `gorm:"type:text" json:"notes"`
	RejectionReason    string               `gorm:"type:text" json:"rejection_reason"`
	ApprovedByID       *uuid.UUID           `gorm:"type:uuid;index" json:"approved_by_id,omitempty"`
	ApprovedAt         *time.Time           `json:"approved_at,omitempty"`
	CompletedAt        *time.Time           `json:"completed_at,omitempty"`
	CreatedBy          uuid.UUID            `gorm:"type:uuid;not null;index" json:"created_by"`
	CreatedAt          time.Time            `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time            `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt       `gorm:"index" json:"-"`

	// Relationships
	InventorySection InventorySection    `gorm:"foreignKey:InventorySectionID" json:"inventory_section,omitempty"`
	ConductedBy      User                `gorm:"foreignKey:ConductedByID" json:"conducted_by,omitempty"`
	ApprovedBy       *User               `gorm:"foreignKey:ApprovedByID" json:"approved_by,omitempty"`
	Creator          User                `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	AuditItems       []InventoryAuditItem `gorm:"foreignKey:AuditID" json:"audit_items,omitempty"`
}

func (InventoryAudit) TableName() string {
	return "inventory_audits"
}

// InventoryAuditItem represents individual items checked in an audit
type InventoryAuditItem struct {
	ID                uuid.UUID              `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	AuditID           uuid.UUID              `gorm:"type:uuid;not null;index" json:"audit_id"`
	ItemID            uuid.UUID              `gorm:"type:uuid;not null;index" json:"item_id"`
	SystemQuantity    int                    `gorm:"default:0" json:"system_quantity"`
	PhysicalQuantity  int                    `gorm:"default:0" json:"physical_quantity"`
	Difference        int                    `gorm:"default:0" json:"difference"`
	Condition         InventoryItemCondition `gorm:"type:varchar(20);default:'good';check:condition IN ('new','good','fair','poor','damaged','obsolete','expired')" json:"condition"`
	Location          string                 `gorm:"type:varchar(255)" json:"location"`
	IsMissing         bool                   `gorm:"default:false" json:"is_missing"`
	IsExcess          bool                   `gorm:"default:false" json:"is_excess"`
	IsDamaged         bool                   `gorm:"default:false" json:"is_damaged"`
	Notes             string                 `gorm:"type:text" json:"notes"`
	ActionTaken       string                 `gorm:"type:text" json:"action_taken"`
	ActionTakenByID   *uuid.UUID             `gorm:"type:uuid;index" json:"action_taken_by_id,omitempty"`
	CreatedBy         uuid.UUID              `gorm:"type:uuid;not null;index" json:"created_by"`
	CreatedAt         time.Time              `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time              `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt         `gorm:"index" json:"-"`

	// Relationships
	Audit         InventoryAudit `gorm:"foreignKey:AuditID" json:"audit,omitempty"`
	Item          InventoryItem  `gorm:"foreignKey:ItemID" json:"item,omitempty"`
	ActionTakenBy *User          `gorm:"foreignKey:ActionTakenByID" json:"action_taken_by,omitempty"`
	Creator       User           `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (InventoryAuditItem) TableName() string {
	return "inventory_audit_items"
}