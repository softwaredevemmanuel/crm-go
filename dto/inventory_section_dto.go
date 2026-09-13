// dto/inventory_section_dto.go
package dto

import (
	"time"
)

// CreateInventorySectionRequest represents the request to create an inventory section
type CreateInventorySectionRequest struct {
	Name         string `json:"name" binding:"required,min=2,max=255"`
	Code         string `json:"code" binding:"required,min=2,max=50"`
	Description  string `json:"description"`
	Location     string `json:"location" binding:"max=255"`
	HeadOfDeptID string `json:"head_of_dept_id"`
	Status       string `json:"status" binding:"omitempty,oneof=active inactive closed renovation"`
	Notes        string `json:"notes"`
}

// UpdateInventorySectionRequest represents the request to update an inventory section
type UpdateInventorySectionRequest struct {
	Name         string `json:"name" binding:"omitempty,min=2,max=255"`
	Code         string `json:"code" binding:"omitempty,min=2,max=50"`
	Description  string `json:"description"`
	Location     string `json:"location" binding:"omitempty,max=255"`
	HeadOfDeptID string `json:"head_of_dept_id"`
	Status       string `json:"status" binding:"omitempty,oneof=active inactive closed renovation"`
	Notes        string `json:"notes"`
}

// InventorySectionResponse represents the response for an inventory section
type InventorySectionResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	HeadOfDeptID string   `json:"head_of_dept_id,omitempty"`
	Status      string    `json:"status"`
	Notes       string    `json:"notes"`
	ItemCount   int64     `json:"item_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Nested relationships
	HeadOfDept *UserResponse `json:"head_of_dept,omitempty"`
}

// InventorySectionListResponse represents a paginated list of inventory sections
type InventorySectionListResponse struct {
	Sections   []InventorySectionResponse `json:"sections"`
	Total      int64                      `json:"total"`
	Page       int                        `json:"page"`
	Limit      int                        `json:"limit"`
	TotalPages int                        `json:"total_pages"`
}

// InventorySectionQueryParams represents query parameters for filtering sections
type InventorySectionQueryParams struct {
	Search    string `form:"search"`
	Status    string `form:"status" binding:"omitempty,oneof=active inactive closed renovation"`
	Location  string `form:"location"`
	Page      int    `form:"page" default:"1"`
	Limit     int    `form:"limit" default:"20"`
	SortBy    string `form:"sort_by" default:"created_at"`
	SortOrder string `form:"sort_order" default:"desc" binding:"omitempty,oneof=asc desc"`
}

