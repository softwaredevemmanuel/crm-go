// dto/exam_scheme_item_dto.go
package dto

import (
	"time"
)

// CreateExamSchemeItemRequest represents the request to create an exam scheme item
type CreateExamSchemeItemRequest struct {
	ExamID             string `json:"exam_id" binding:"required"`
	SchemeOfWorkItemID string `json:"scheme_of_work_item_id" binding:"required"`
}

// BulkCreateExamSchemeItemsRequest represents the request to bulk create exam scheme items
type BulkCreateExamSchemeItemsRequest struct {
	ExamID             string   `json:"exam_id" binding:"required"`
	SchemeOfWorkItemIDs []string `json:"scheme_of_work_item_ids" binding:"required,min=1"`
}

// ExamSchemeItemResponse represents the response for an exam scheme item
type ExamSchemeItemResponse struct {
	ExamID             string    `json:"exam_id"`
	SchemeOfWorkItemID string    `json:"scheme_of_work_item_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`

	// Nested relationships
	Exam             *ExamResponse             `json:"exam,omitempty"`
	SchemeOfWorkItem *SchemeOfWorkItemResponse `json:"scheme_of_work_item,omitempty"`
}

// ExamSchemeItemListResponse represents a paginated list of exam scheme items
type ExamSchemeItemListResponse struct {
	Items      []ExamSchemeItemResponse `json:"items"`
	Total      int64                    `json:"total"`
	Page       int                      `json:"page"`
	Limit      int                      `json:"limit"`
	TotalPages int                      `json:"total_pages"`
}

// ExamSchemeItemQueryParams represents query parameters for filtering exam scheme items
type ExamSchemeItemQueryParams struct {
	ExamID             string `form:"exam_id"`
	SchemeOfWorkItemID string `form:"scheme_of_work_item_id"`
	Page               int    `form:"page" default:"1"`
	Limit              int    `form:"limit" default:"20"`
}

// BulkExamSchemeItemResult represents the result of a bulk operation
type BulkExamSchemeItemResult struct {
	SuccessCount int                         `json:"success_count"`
	FailedCount  int                         `json:"failed_count"`
	Created      []ExamSchemeItemResponse    `json:"created"`
	Errors       []BulkExamSchemeItemError   `json:"errors"`
}

// BulkExamSchemeItemError represents an error in bulk exam scheme item creation
type BulkExamSchemeItemError struct {
	SchemeOfWorkItemID string `json:"scheme_of_work_item_id"`
	Error              string `json:"error"`
}

// ExamSchemeItemStats represents statistics for exam scheme items
type ExamSchemeItemStats struct {
	TotalItems          int64 `json:"total_items"`
	TotalExams          int64 `json:"total_exams"`
	TotalSchemeItems    int64 `json:"total_scheme_items"`
	ItemsPerExam        int64 `json:"items_per_exam"`
}