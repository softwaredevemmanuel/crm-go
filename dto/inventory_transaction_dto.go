// dto/inventory_transaction_dto.go
package dto

import (
	"time"
)

// CreateInventoryTransactionRequest represents the request to create a transaction
type CreateInventoryTransactionRequest struct {
	ItemID             string  `json:"item_id" binding:"required"`
	InventorySectionID string  `json:"inventory_section_id" binding:"required"`
	TransactionType    string  `json:"transaction_type" binding:"required,oneof=purchase donation issue return transfer borrow usage experiment damage loss repair repair_return disposal adjustment audit expiry calibration"`
	Quantity           int     `json:"quantity" binding:"required,min=1"`
	UnitCost           float64 `json:"unit_cost"`
	ReferenceNumber    string  `json:"reference_number"`
	ReferenceType      string  `json:"reference_type" binding:"omitempty,oneof=purchase_order invoice requisition return transfer experiment practical"`
	FromLocation       string  `json:"from_location"`
	ToLocation         string  `json:"to_location"`
	FromSectionID      string  `json:"from_section_id"`
	ToSectionID        string  `json:"to_section_id"`
	IssuedToID         string  `json:"issued_to_id"`
	IssuedByID         string  `json:"issued_by_id"`
	ApprovedByID       string  `json:"approved_by_id"`
	ReceivedByID       string  `json:"received_by_id"`
	SupervisorID       string  `json:"supervisor_id"`
	ClassID            string  `json:"class_id"`
	Purpose            string  `json:"purpose"`
	ExperimentName     string  `json:"experiment_name"`
	StudentCount       int     `json:"student_count"`
	Notes              string  `json:"notes"`
	Status             string  `json:"status" binding:"omitempty,oneof=pending approved rejected completed cancelled"`
	TransactionDate    string  `json:"transaction_date"`
	ExpectedReturnDate string  `json:"expected_return_date"`
	ConditionBefore    string  `json:"condition_before" binding:"omitempty,oneof=new good fair poor damaged obsolete expired"`
	ConditionAfter     string  `json:"condition_after" binding:"omitempty,oneof=new good fair poor damaged obsolete expired"`
	AttachmentURL      string  `json:"attachment_url"`
}

// UpdateInventoryTransactionRequest represents the request to update a transaction
type UpdateInventoryTransactionRequest struct {
	Quantity           *int    `json:"quantity"`
	UnitCost           *float64 `json:"unit_cost"`
	ReferenceNumber    string  `json:"reference_number"`
	ReferenceType      string  `json:"reference_type" binding:"omitempty,oneof=purchase_order invoice requisition return transfer experiment practical"`
	FromLocation       string  `json:"from_location"`
	ToLocation         string  `json:"to_location"`
	Purpose            string  `json:"purpose"`
	ExperimentName     string  `json:"experiment_name"`
	StudentCount       *int    `json:"student_count"`
	Notes              string  `json:"notes"`
	Status             string  `json:"status" binding:"omitempty,oneof=pending approved rejected completed cancelled"`
	TransactionDate    string  `json:"transaction_date"`
	ExpectedReturnDate string  `json:"expected_return_date"`
	ConditionAfter     string  `json:"condition_after" binding:"omitempty,oneof=new good fair poor damaged obsolete expired"`
	AttachmentURL      string  `json:"attachment_url"`
}

// CompleteTransactionRequest represents the request to complete a transaction
type CompleteTransactionRequest struct {
	ConditionAfter   string `json:"condition_after" binding:"omitempty,oneof=new good fair poor damaged obsolete expired"`
	ActualReturnDate string `json:"actual_return_date"`
	Notes            string `json:"notes"`
	AttachmentURL    string `json:"attachment_url"`
}

// InventoryTransactionResponse represents the response for a transaction
type InventoryTransactionResponse struct {
	ID                 string     `json:"id"`
	TransactionCode    string     `json:"transaction_code"`
	ItemID             string     `json:"item_id"`
	InventorySectionID string     `json:"inventory_section_id"`
	TransactionType    string     `json:"transaction_type"`
	Quantity           int        `json:"quantity"`
	QuantityBefore     int        `json:"quantity_before"`
	QuantityAfter      int        `json:"quantity_after"`
	UnitCost           float64    `json:"unit_cost"`
	TotalCost          float64    `json:"total_cost"`
	ReferenceNumber    string     `json:"reference_number"`
	ReferenceType      string     `json:"reference_type,omitempty"`
	FromLocation       string     `json:"from_location"`
	ToLocation         string     `json:"to_location"`
	FromSectionID      string     `json:"from_section_id,omitempty"`
	ToSectionID        string     `json:"to_section_id,omitempty"`
	IssuedToID         string     `json:"issued_to_id,omitempty"`
	IssuedByID         string     `json:"issued_by_id,omitempty"`
	ApprovedByID       string     `json:"approved_by_id,omitempty"`
	ReceivedByID       string     `json:"received_by_id,omitempty"`
	SupervisorID       string     `json:"supervisor_id,omitempty"`
	ClassID            string     `json:"class_id,omitempty"`
	Purpose            string     `json:"purpose"`
	ExperimentName     string     `json:"experiment_name"`
	StudentCount       int        `json:"student_count"`
	Notes              string     `json:"notes"`
	Status             string     `json:"status"`
	TransactionDate    time.Time  `json:"transaction_date"`
	ExpectedReturnDate *time.Time `json:"expected_return_date,omitempty"`
	ActualReturnDate   *time.Time `json:"actual_return_date,omitempty"`
	ConditionBefore    string     `json:"condition_before,omitempty"`
	ConditionAfter     string     `json:"condition_after,omitempty"`
	AttachmentURL      string     `json:"attachment_url"`
	CreatedBy          string     `json:"created_by"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`

	// Nested relationships
	Item             *InventoryItemResponse    `json:"item,omitempty"`
	InventorySection *InventorySectionResponse `json:"inventory_section,omitempty"`
	FromSection      *InventorySectionResponse `json:"from_section,omitempty"`
	ToSection        *InventorySectionResponse `json:"to_section,omitempty"`
	IssuedTo         *UserResponse             `json:"issued_to,omitempty"`
	IssuedBy         *UserResponse             `json:"issued_by,omitempty"`
	ApprovedBy       *UserResponse             `json:"approved_by,omitempty"`
	ReceivedBy       *UserResponse             `json:"received_by,omitempty"`
	Supervisor       *UserResponse             `json:"supervisor,omitempty"`
	Creator          *UserResponse             `json:"creator,omitempty"`
}

// InventoryTransactionListResponse represents a paginated list of transactions
type InventoryTransactionListResponse struct {
	Transactions []InventoryTransactionResponse `json:"transactions"`
	Total        int64                          `json:"total"`
	Page         int                            `json:"page"`
	Limit        int                            `json:"limit"`
	TotalPages   int                            `json:"total_pages"`
}

// InventoryTransactionQueryParams represents query parameters for filtering transactions
type InventoryTransactionQueryParams struct {
	ItemID             string `form:"item_id"`
	InventorySectionID string `form:"inventory_section_id"`
	TransactionType    string `form:"transaction_type"`
	Status             string `form:"status" binding:"omitempty,oneof=pending approved rejected completed cancelled"`
	IssuedToID         string `form:"issued_to_id"`
	IssuedByID         string `form:"issued_by_id"`
	ApprovedByID       string `form:"approved_by_id"`
	ClassID            string `form:"class_id"`
	Search             string `form:"search"`
	DateFrom           string `form:"date_from"`
	DateTo             string `form:"date_to"`
	Page               int    `form:"page" default:"1"`
	Limit              int    `form:"limit" default:"20"`
	SortBy             string `form:"sort_by" default:"transaction_date"`
	SortOrder          string `form:"sort_order" default:"desc" binding:"omitempty,oneof=asc desc"`
}

// InventoryTransactionStats represents statistics for inventory transactions
type InventoryTransactionStats struct {
	TotalTransactions   int64   `json:"total_transactions"`
	PendingTransactions int64   `json:"pending_transactions"`
	CompletedTransactions int64 `json:"completed_transactions"`
	CancelledTransactions int64 `json:"cancelled_transactions"`
	TotalQuantityIn     int64   `json:"total_quantity_in"`
	TotalQuantityOut    int64   `json:"total_quantity_out"`
	TotalValueIn        float64 `json:"total_value_in"`
	TotalValueOut       float64 `json:"total_value_out"`
}