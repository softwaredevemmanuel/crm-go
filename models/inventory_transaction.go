// models/inventory_transaction.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryTransactionType string

const (
	InventoryTxnPurchase     InventoryTransactionType = "purchase"
	InventoryTxnDonation     InventoryTransactionType = "donation"
	InventoryTxnIssue        InventoryTransactionType = "issue"
	InventoryTxnReturn       InventoryTransactionType = "return"
	InventoryTxnTransfer     InventoryTransactionType = "transfer"
	InventoryTxnBorrow       InventoryTransactionType = "borrow"
	InventoryTxnUsage        InventoryTransactionType = "usage"
	InventoryTxnExperiment   InventoryTransactionType = "experiment"
	InventoryTxnDamage       InventoryTransactionType = "damage"
	InventoryTxnLoss         InventoryTransactionType = "loss"
	InventoryTxnRepair       InventoryTransactionType = "repair"
	InventoryTxnRepairReturn InventoryTransactionType = "repair_return"
	InventoryTxnDisposal     InventoryTransactionType = "disposal"
	InventoryTxnAdjustment   InventoryTransactionType = "adjustment"
	InventoryTxnAudit        InventoryTransactionType = "audit"
	InventoryTxnExpiry       InventoryTransactionType = "expiry"
	InventoryTxnCalibration  InventoryTransactionType = "calibration"
)

type InventoryTransactionStatus string

const (
	InventoryTxnStatusPending   InventoryTransactionStatus = "pending"
	InventoryTxnStatusApproved  InventoryTransactionStatus = "approved"
	InventoryTxnStatusRejected  InventoryTransactionStatus = "rejected"
	InventoryTxnStatusCompleted InventoryTransactionStatus = "completed"
	InventoryTxnStatusCancelled InventoryTransactionStatus = "cancelled"
)

type InventoryReferenceType string

const (
	InventoryRefPurchaseOrder InventoryReferenceType = "purchase_order"
	InventoryRefInvoice       InventoryReferenceType = "invoice"
	InventoryRefRequisition   InventoryReferenceType = "requisition"
	InventoryRefReturn        InventoryReferenceType = "return"
	InventoryRefTransfer      InventoryReferenceType = "transfer"
	InventoryRefExperiment    InventoryReferenceType = "experiment"
	InventoryRefPractical     InventoryReferenceType = "practical"
)

type InventoryTransaction struct {
	ID                 uuid.UUID                  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TransactionCode    string                     `gorm:"type:varchar(100);uniqueIndex;not null" json:"transaction_code"`
	ItemID             uuid.UUID                  `gorm:"type:uuid;not null;index" json:"item_id"`
	InventorySectionID uuid.UUID                  `gorm:"type:uuid;not null;index" json:"inventory_section_id"`
	TransactionType    InventoryTransactionType   `gorm:"type:varchar(50);not null;index;check:transaction_type IN ('purchase','donation','issue','return','transfer','borrow','usage','experiment','damage','loss','repair','repair_return','disposal','adjustment','audit','expiry','calibration')" json:"transaction_type"`
	Quantity           int                        `gorm:"not null;check:quantity > 0" json:"quantity"`
	QuantityBefore     int                        `gorm:"default:0" json:"quantity_before"`
	QuantityAfter      int                        `gorm:"default:0" json:"quantity_after"`
	UnitCost           float64                    `gorm:"type:decimal(15,2);default:0" json:"unit_cost"`
	TotalCost          float64                    `gorm:"type:decimal(15,2);default:0" json:"total_cost"`

	// References
	ReferenceNumber string                  `gorm:"type:varchar(100)" json:"reference_number"`
	ReferenceType   *InventoryReferenceType `gorm:"type:varchar(50)" json:"reference_type,omitempty"`

	// Locations
	FromLocation     string     `gorm:"type:varchar(255)" json:"from_location"`
	ToLocation       string     `gorm:"type:varchar(255)" json:"to_location"`
	FromSectionID    *uuid.UUID `gorm:"type:uuid;index" json:"from_section_id,omitempty"`
	ToSectionID      *uuid.UUID `gorm:"type:uuid;index" json:"to_section_id,omitempty"`

	// People
	IssuedToID   *uuid.UUID `gorm:"type:uuid;index" json:"issued_to_id,omitempty"`
	IssuedByID   *uuid.UUID `gorm:"type:uuid;index" json:"issued_by_id,omitempty"`
	ApprovedByID *uuid.UUID `gorm:"type:uuid;index" json:"approved_by_id,omitempty"`
	ReceivedByID *uuid.UUID `gorm:"type:uuid;index" json:"received_by_id,omitempty"`
	SupervisorID *uuid.UUID `gorm:"type:uuid;index" json:"supervisor_id,omitempty"`

	// Context
	ClassID        *uuid.UUID `gorm:"type:uuid;index" json:"class_id,omitempty"`
	Purpose        string     `gorm:"type:text" json:"purpose"`
	ExperimentName string     `gorm:"type:varchar(255)" json:"experiment_name"`
	StudentCount   int        `gorm:"default:0" json:"student_count"`
	Notes          string     `gorm:"type:text" json:"notes"`

	// Status
	Status            InventoryTransactionStatus `gorm:"type:varchar(20);default:'pending';check:status IN ('pending','approved','rejected','completed','cancelled')" json:"status"`
	TransactionDate   time.Time                  `gorm:"not null;index" json:"transaction_date"`
	ExpectedReturnDate *time.Time                `gorm:"type:date" json:"expected_return_date,omitempty"`
	ActualReturnDate   *time.Time                `gorm:"type:date" json:"actual_return_date,omitempty"`

	// Conditions
	ConditionBefore *InventoryItemCondition `gorm:"type:varchar(20)" json:"condition_before,omitempty"`
	ConditionAfter  *InventoryItemCondition `gorm:"type:varchar(20)" json:"condition_after,omitempty"`

	// Attachments
	AttachmentURL string `gorm:"type:varchar(500)" json:"attachment_url"`

	// Audit
	CreatedBy uuid.UUID      `gorm:"type:uuid;not null;index" json:"created_by"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Item             InventoryItem      `gorm:"foreignKey:ItemID" json:"item,omitempty"`
	InventorySection InventorySection   `gorm:"foreignKey:InventorySectionID" json:"inventory_section,omitempty"`
	FromSection      *InventorySection  `gorm:"foreignKey:FromSectionID" json:"from_section,omitempty"`
	ToSection        *InventorySection  `gorm:"foreignKey:ToSectionID" json:"to_section,omitempty"`
	IssuedTo         *User              `gorm:"foreignKey:IssuedToID" json:"issued_to,omitempty"`
	IssuedBy         *User              `gorm:"foreignKey:IssuedByID" json:"issued_by,omitempty"`
	ApprovedBy       *User              `gorm:"foreignKey:ApprovedByID" json:"approved_by,omitempty"`
	ReceivedBy       *User              `gorm:"foreignKey:ReceivedByID" json:"received_by,omitempty"`
	Supervisor       *User              `gorm:"foreignKey:SupervisorID" json:"supervisor,omitempty"`
	Creator          User               `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (InventoryTransaction) TableName() string {
	return "inventory_transactions"
}