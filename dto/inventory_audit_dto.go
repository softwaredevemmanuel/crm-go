// dto/inventory_audit_dto.go
package dto

import (
	"time"
)

// AuditItemInput represents an item being checked during an audit
type AuditItemInput struct {
	ItemID           string `json:"item_id" binding:"required"`
	PhysicalQuantity int    `json:"physical_quantity" binding:"min=0"`
	Condition        string `json:"condition" binding:"omitempty,oneof=new good fair poor damaged obsolete expired"`
	Location         string `json:"location"`
	IsMissing        bool   `json:"is_missing"`
	IsExcess         bool   `json:"is_excess"`
	IsDamaged        bool   `json:"is_damaged"`
	Notes            string `json:"notes"`
	ActionTaken      string `json:"action_taken"`
}

// CreateInventoryAuditRequest represents the request to create an audit
type CreateInventoryAuditRequest struct {
	InventorySectionID string           `json:"inventory_section_id" binding:"required"`
	AuditDate          string           `json:"audit_date" binding:"required"`
	Notes              string           `json:"notes"`
	Items              []AuditItemInput `json:"items"`
}

// UpdateInventoryAuditRequest represents the request to update an audit
type UpdateInventoryAuditRequest struct {
	AuditDate string `json:"audit_date"`
	Notes     string `json:"notes"`
}

// UpdateAuditItemRequest represents the request to update an audit item
type UpdateAuditItemRequest struct {
	PhysicalQuantity int    `json:"physical_quantity" binding:"min=0"`
	Condition        string `json:"condition" binding:"omitempty,oneof=new good fair poor damaged obsolete expired"`
	Location         string `json:"location"`
	IsMissing        *bool  `json:"is_missing"`
	IsExcess         *bool  `json:"is_excess"`
	IsDamaged        *bool  `json:"is_damaged"`
	Notes            string `json:"notes"`
	ActionTaken      string `json:"action_taken"`
}

// ApproveInventoryAuditRequest represents the request to approve an audit
type ApproveInventoryAuditRequest struct {
	Notes string `json:"notes"`
}

// RejectInventoryAuditRequest represents the request to reject an audit
type RejectInventoryAuditRequest struct {
	RejectionReason string `json:"rejection_reason" binding:"required,min=3"`
}

// AuditItemResponse represents the response for an audit item
type AuditItemResponse struct {
	ID               string    `json:"id"`
	AuditID          string    `json:"audit_id"`
	ItemID           string    `json:"item_id"`
	SystemQuantity   int       `json:"system_quantity"`
	PhysicalQuantity int       `json:"physical_quantity"`
	Difference       int       `json:"difference"`
	Condition        string    `json:"condition"`
	Location         string    `json:"location"`
	IsMissing        bool      `json:"is_missing"`
	IsExcess         bool      `json:"is_excess"`
	IsDamaged        bool      `json:"is_damaged"`
	Notes            string    `json:"notes"`
	ActionTaken      string    `json:"action_taken"`
	ActionTakenByID  string    `json:"action_taken_by_id,omitempty"`
	CreatedBy        string    `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Nested relationships
	Item          *InventoryItemResponse `json:"item,omitempty"`
	ActionTakenBy *UserResponse          `json:"action_taken_by,omitempty"`
	Creator       *UserResponse          `json:"creator,omitempty"`
}

// InventoryAuditResponse represents the response for an audit
type InventoryAuditResponse struct {
	ID                 string     `json:"id"`
	AuditCode          string     `json:"audit_code"`
	InventorySectionID string     `json:"inventory_section_id"`
	AuditDate          time.Time  `json:"audit_date"`
	ConductedByID      string     `json:"conducted_by_id"`
	Status             string     `json:"status"`
	TotalItemsChecked  int        `json:"total_items_checked"`
	TotalDiscrepancies int        `json:"total_discrepancies"`
	TotalItemsMissing  int        `json:"total_items_missing"`
	TotalItemsExcess   int        `json:"total_items_excess"`
	TotalItemsDamaged  int        `json:"total_items_damaged"`
	Notes              string     `json:"notes"`
	RejectionReason    string     `json:"rejection_reason,omitempty"`
	ApprovedByID       string     `json:"approved_by_id,omitempty"`
	ApprovedAt         *time.Time `json:"approved_at,omitempty"`
	CompletedAt        *time.Time `json:"completed_at,omitempty"`
	CreatedBy          string     `json:"created_by"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`

	// Nested relationships
	InventorySection *InventorySectionResponse `json:"inventory_section,omitempty"`
	ConductedBy      *UserResponse             `json:"conducted_by,omitempty"`
	ApprovedBy       *UserResponse             `json:"approved_by,omitempty"`
	Creator          *UserResponse             `json:"creator,omitempty"`
	AuditItems       []AuditItemResponse       `json:"audit_items,omitempty"`
}

// InventoryAuditListResponse represents a paginated list of audits
type InventoryAuditListResponse struct {
	Audits     []InventoryAuditResponse `json:"audits"`
	Total      int64                    `json:"total"`
	Page       int                      `json:"page"`
	Limit      int                      `json:"limit"`
	TotalPages int                      `json:"total_pages"`
}

// InventoryAuditQueryParams represents query parameters for filtering audits
type InventoryAuditQueryParams struct {
	InventorySectionID string `form:"inventory_section_id"`
	ConductedByID      string `form:"conducted_by_id"`
	Status             string `form:"status" binding:"omitempty,oneof=in_progress completed approved rejected cancelled"`
	Search             string `form:"search"`
	DateFrom           string `form:"date_from"`
	DateTo             string `form:"date_to"`
	Page               int    `form:"page" default:"1"`
	Limit              int    `form:"limit" default:"20"`
	SortBy             string `form:"sort_by" default:"created_at"`
	SortOrder          string `form:"sort_order" default:"desc" binding:"omitempty,oneof=asc desc"`
}

// InventoryAuditStats represents statistics for audits
type InventoryAuditStats struct {
	TotalAudits        int64   `json:"total_audits"`
	InProgressAudits   int64   `json:"in_progress_audits"`
	CompletedAudits    int64   `json:"completed_audits"`
	ApprovedAudits     int64   `json:"approved_audits"`
	RejectedAudits     int64   `json:"rejected_audits"`
	TotalDiscrepancies int64   `json:"total_discrepancies"`
	AverageDiscrepancy float64 `json:"average_discrepancy"`
}