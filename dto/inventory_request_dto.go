// dto/inventory_request_dto.go
package dto

import (
	"time"
)

// CreateInventoryRequestRequest represents the request to create an inventory request
type CreateInventoryRequestRequest struct {
	InventorySectionID string `json:"inventory_section_id" binding:"required"`
	InventoryItemID    string `json:"inventory_item_id" binding:"required"`
	QuantityRequested  int    `json:"quantity_requested" binding:"required,min=1"`
	Purpose            string `json:"purpose" binding:"required,min=3"`
	NeededByDate       string `json:"needed_by_date"` // Format: YYYY-MM-DD
	Notes              string `json:"notes"`
}

// UpdateInventoryRequestRequest represents the request to update an inventory request
type UpdateInventoryRequestRequest struct {
	QuantityRequested *int   `json:"quantity_requested"`
	Purpose           string `json:"purpose"`
	NeededByDate      string `json:"needed_by_date"`
	Notes             string `json:"notes"`
}

// ApproveInventoryRequestRequest represents the request to approve a request
type ApproveInventoryRequestRequest struct {
	QuantityApproved int    `json:"quantity_approved" binding:"required,min=1"`
	Notes            string `json:"notes"`
}

// RejectInventoryRequestRequest represents the request to reject a request
type RejectInventoryRequestRequest struct {
	RejectionReason string `json:"rejection_reason" binding:"required,min=3"`
}

// FulfillInventoryRequestRequest represents the request to fulfill a request
type FulfillInventoryRequestRequest struct {
	QuantityIssued int    `json:"quantity_issued" binding:"required,min=1"`
	Notes          string `json:"notes"`
}

// InventoryRequestResponse represents the response for an inventory request
type InventoryRequestResponse struct {
	ID                 string     `json:"id"`
	RequestCode        string     `json:"request_code"`
	RequestedByID      string     `json:"requested_by_id"`
	InventorySectionID string     `json:"inventory_section_id"`
	InventoryItemID    string     `json:"inventory_item_id"`
	QuantityRequested  int        `json:"quantity_requested"`
	QuantityApproved   int        `json:"quantity_approved"`
	QuantityIssued     int        `json:"quantity_issued"`
	Purpose            string     `json:"purpose"`
	Status             string     `json:"status"`
	ApprovedByID       string     `json:"approved_by_id,omitempty"`
	ApprovedAt         *time.Time `json:"approved_at,omitempty"`
	RejectionReason    string     `json:"rejection_reason,omitempty"`
	NeededByDate       *time.Time `json:"needed_by_date,omitempty"`
	FulfilledAt        *time.Time `json:"fulfilled_at,omitempty"`
	Notes              string     `json:"notes"`
	CreatedBy          string     `json:"created_by"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`

	// Nested relationships
	RequestedBy      *UserResponse                `json:"requested_by,omitempty"`
	InventorySection *InventorySectionResponse    `json:"inventory_section,omitempty"`
	InventoryItem    *InventoryItemResponse       `json:"inventory_item,omitempty"`
	ApprovedBy       *UserResponse                `json:"approved_by,omitempty"`
	Creator          *UserResponse                `json:"creator,omitempty"`
}

// InventoryRequestListResponse represents a paginated list of inventory requests
type InventoryRequestListResponse struct {
	Requests   []InventoryRequestResponse `json:"requests"`
	Total      int64                      `json:"total"`
	Page       int                        `json:"page"`
	Limit      int                        `json:"limit"`
	TotalPages int                        `json:"total_pages"`
}

// InventoryRequestQueryParams represents query parameters for filtering requests
type InventoryRequestQueryParams struct {
	RequestedByID      string `form:"requested_by_id"`
	InventorySectionID string `form:"inventory_section_id"`
	InventoryItemID    string `form:"inventory_item_id"`
	Status             string `form:"status" binding:"omitempty,oneof=pending approved rejected fulfilled cancelled"`
	Search             string `form:"search"`
	DateFrom           string `form:"date_from"`
	DateTo             string `form:"date_to"`
	Page               int    `form:"page" default:"1"`
	Limit              int    `form:"limit" default:"20"`
	SortBy             string `form:"sort_by" default:"created_at"`
	SortOrder          string `form:"sort_order" default:"desc" binding:"omitempty,oneof=asc desc"`
}

// InventoryRequestStats represents statistics for inventory requests
type InventoryRequestStats struct {
	TotalRequests     int64 `json:"total_requests"`
	PendingRequests   int64 `json:"pending_requests"`
	ApprovedRequests  int64 `json:"approved_requests"`
	RejectedRequests  int64 `json:"rejected_requests"`
	FulfilledRequests int64 `json:"fulfilled_requests"`
	CancelledRequests int64 `json:"cancelled_requests"`
}