// dto/test_scheme_item_dto.go
package dto

import (
	"time"
)

// CreateTestSchemeItemRequest represents the request to create a test scheme item
type CreateTestSchemeItemRequest struct {
	TestID             string `json:"test_id" binding:"required"`
	SchemeOfWorkItemID string `json:"scheme_of_work_item_id" binding:"required"`
}

// BulkCreateTestSchemeItemsRequest represents the request to bulk create test scheme items
type BulkCreateTestSchemeItemsRequest struct {
	TestID             string   `json:"test_id" binding:"required"`
	SchemeOfWorkItemIDs []string `json:"scheme_of_work_item_ids" binding:"required,min=1"`
}

// TestSchemeItemResponse represents the response for a test scheme item
type TestSchemeItemResponse struct {
	TestID             string    `json:"test_id"`
	SchemeOfWorkItemID string    `json:"scheme_of_work_item_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`

	// Nested relationships
	Test             *TestResponse             `json:"test,omitempty"`
	SchemeOfWorkItem *SchemeOfWorkItemResponse `json:"scheme_of_work_item,omitempty"`
}

// TestSchemeItemListResponse represents a paginated list of test scheme items
type TestSchemeItemListResponse struct {
	Items      []TestSchemeItemResponse `json:"items"`
	Total      int64                    `json:"total"`
	Page       int                      `json:"page"`
	Limit      int                      `json:"limit"`
	TotalPages int                      `json:"total_pages"`
}

// TestSchemeItemQueryParams represents query parameters for filtering test scheme items
type TestSchemeItemQueryParams struct {
	TestID             string `form:"test_id"`
	SchemeOfWorkItemID string `form:"scheme_of_work_item_id"`
	Page               int    `form:"page" default:"1"`
	Limit              int    `form:"limit" default:"20"`
}

// BulkTestSchemeItemResult represents the result of a bulk operation
type BulkTestSchemeItemResult struct {
	SuccessCount int                         `json:"success_count"`
	FailedCount  int                         `json:"failed_count"`
	Created      []TestSchemeItemResponse    `json:"created"`
	Errors       []BulkTestSchemeItemError   `json:"errors"`
}

// BulkTestSchemeItemError represents an error in bulk test scheme item creation
type BulkTestSchemeItemError struct {
	SchemeOfWorkItemID string `json:"scheme_of_work_item_id"`
	Error              string `json:"error"`
}

// TestSchemeItemStats represents statistics for test scheme items
type TestSchemeItemStats struct {
	TotalItems          int64 `json:"total_items"`
	TotalTests          int64 `json:"total_tests"`
	TotalSchemeItems    int64 `json:"total_scheme_items"`
	ItemsPerTest        int64 `json:"items_per_test"`
}