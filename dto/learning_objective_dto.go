// dto/learning_objective_dto.go
package dto

import (
	"time"
)

// CreateLearningObjectiveRequest represents the request to create a learning objective
type CreateLearningObjectiveRequest struct {
	SchemeOfWorkItemID string `json:"scheme_of_work_item_id" binding:"required"`
	Objective          string `json:"objective" binding:"required"`
	Sequence           int    `json:"sequence"`
}

// UpdateLearningObjectiveRequest represents the request to update a learning objective
type UpdateLearningObjectiveRequest struct {
	Objective string `json:"objective"`
	Sequence  int    `json:"sequence"`
}

// LearningObjectiveResponse represents the response for a learning objective
type LearningObjectiveResponse struct {
	ID                 string    `json:"id"`
	SchemeOfWorkItemID string    `json:"scheme_of_work_item_id"`
	Objective          string    `json:"objective"`
	Sequence           int       `json:"sequence"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`

	// Nested relationships
	SchemeOfWorkItem *SchemeOfWorkItemResponse `json:"scheme_of_work_item,omitempty"`
}

// LearningObjectiveListResponse represents a paginated list of learning objectives
type LearningObjectiveListResponse struct {
	Objectives []LearningObjectiveResponse `json:"objectives"`
	Total      int64                       `json:"total"`
	Page       int                         `json:"page"`
	Limit      int                         `json:"limit"`
	TotalPages int                         `json:"total_pages"`
}

// LearningObjectiveQueryParams represents query parameters for filtering learning objectives
type LearningObjectiveQueryParams struct {
	SchemeOfWorkItemID string `form:"scheme_of_work_item_id"`
	Search             string `form:"search"`
	Page               int    `form:"page" default:"1"`
	Limit              int    `form:"limit" default:"20"`
	SortBy             string `form:"sort_by" default:"sequence"`
	SortOrder          string `form:"sort_order" default:"asc"`
}

// BulkCreateLearningObjectivesRequest represents the request to bulk create learning objectives
type BulkCreateLearningObjectivesRequest struct {
	SchemeOfWorkItemID string `json:"scheme_of_work_item_id" binding:"required"`
	Objectives         []struct {
		Objective string `json:"objective" binding:"required"`
		Sequence  int    `json:"sequence"`
	} `json:"objectives" binding:"required,min=1"`
}

// BulkLearningObjectiveResult represents the result of a bulk operation
type BulkLearningObjectiveResult struct {
	SuccessCount int                            `json:"success_count"`
	FailedCount  int                            `json:"failed_count"`
	Created      []LearningObjectiveResponse    `json:"created"`
	Errors       []BulkLearningObjectiveError   `json:"errors"`
}

// BulkLearningObjectiveError represents an error in bulk learning objective creation
type BulkLearningObjectiveError struct {
	Objective string `json:"objective"`
	Error     string `json:"error"`
}