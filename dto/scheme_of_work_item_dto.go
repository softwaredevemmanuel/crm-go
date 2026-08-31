// dto/scheme_of_work_item_dto.go
package dto

import (
	"time"
)

// CreateSchemeOfWorkItemRequest represents the request to create a scheme of work item
type CreateSchemeOfWorkItemRequest struct {
	SchemeOfWorkID string `json:"scheme_of_work_id" binding:"required"`
	ModuleID       string `json:"module_id"`
	WeekStart      int    `json:"week_start" binding:"required,min=1"`
	WeekEnd        int    `json:"week_end"`
	Topic          string `json:"topic" binding:"required"`
	Subtopic       string `json:"subtopic"`
	Content        string `json:"content"`
	Sequence       int    `json:"sequence"`
}

// UpdateSchemeOfWorkItemRequest represents the request to update a scheme of work item
type UpdateSchemeOfWorkItemRequest struct {
	ModuleID   string `json:"module_id"`
	WeekStart  int    `json:"week_start"`
	WeekEnd    int    `json:"week_end"`
	Topic      string `json:"topic"`
	Subtopic   string `json:"subtopic"`
	Content    string `json:"content"`
	Sequence   int    `json:"sequence"`
}

// SchemeOfWorkItemResponse represents the response for a scheme of work item
type SchemeOfWorkItemResponse struct {
	ID             string    `json:"id"`
	SchemeOfWorkID string    `json:"scheme_of_work_id"`
	ModuleID       string    `json:"module_id"`
	WeekStart      int       `json:"week_start"`
	WeekEnd        int       `json:"week_end"`
	Topic          string    `json:"topic"`
	Subtopic       string    `json:"subtopic"`
	Content        string    `json:"content"`
	Sequence       int       `json:"sequence"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Nested relationships
	SchemeOfWork *SchemeOfWorkResponse `json:"scheme_of_work,omitempty"`
	Module       *ModuleResponse       `json:"module,omitempty"`
}

// SchemeOfWorkItemListResponse represents a paginated list of scheme of work items
type SchemeOfWorkItemListResponse struct {
	Items      []SchemeOfWorkItemResponse `json:"items"`
	Total      int64                      `json:"total"`
	Page       int                        `json:"page"`
	Limit      int                        `json:"limit"`
	TotalPages int                        `json:"total_pages"`
}

// SchemeOfWorkItemQueryParams represents query parameters for filtering scheme of work items
type SchemeOfWorkItemQueryParams struct {
	SchemeOfWorkID string `form:"scheme_of_work_id"`
	ModuleID       string `form:"module_id"`
	Search         string `form:"search"`
	Page           int    `form:"page" default:"1"`
	Limit          int    `form:"limit" default:"20"`
	SortBy         string `form:"sort_by" default:"sequence"`
	SortOrder      string `form:"sort_order" default:"asc"`
}

// BulkCreateSchemeItemsRequest represents the request to bulk create scheme of work items
type BulkCreateSchemeItemsRequest struct {
	SchemeOfWorkID string `json:"scheme_of_work_id" binding:"required"`
	Items          []struct {
		ModuleID  string `json:"module_id"`
		WeekStart int    `json:"week_start" binding:"required,min=1"`
		WeekEnd   int    `json:"week_end"`
		Topic     string `json:"topic" binding:"required"`
		Subtopic  string `json:"subtopic"`
		Content   string `json:"content"`
		Sequence  int    `json:"sequence"`
	} `json:"items" binding:"required,min=1"`
}

// BulkSchemeItemResult represents the result of a bulk operation
type BulkSchemeItemResult struct {
	SuccessCount int                         `json:"success_count"`
	FailedCount  int                         `json:"failed_count"`
	Created      []SchemeOfWorkItemResponse  `json:"created"`
	Errors       []BulkSchemeItemError       `json:"errors"`
}

// BulkSchemeItemError represents an error in bulk scheme item creation
type BulkSchemeItemError struct {
	Topic string `json:"topic"`
	Error string `json:"error"`
}