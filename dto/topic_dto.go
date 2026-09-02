// dto/topic_dto.go
package dto

import (
	"time"
)

// CreateTopicRequest represents the request to create a topic
type CreateTopicRequest struct {
	ModuleID    string `json:"module_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	TopicOrder  int    `json:"topic_order"`
}

// UpdateTopicRequest represents the request to update a topic
type UpdateTopicRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TopicOrder  int    `json:"topic_order"`
}

// TopicResponse represents the response for a topic
type TopicResponse struct {
	ID          string    `json:"id"`
	ModuleID    string    `json:"module_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TopicOrder  int       `json:"topic_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Nested relationships
	Module  *ModuleResponse  `json:"module,omitempty"`
	Lessons []LessonResponse `json:"lessons,omitempty"`
}

// TopicListResponse represents a paginated list of topics
type TopicListResponse struct {
	Topics     []TopicResponse `json:"topics"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

// TopicQueryParams represents query parameters for filtering topics
type TopicQueryParams struct {
	ModuleID string `form:"module_id"`
	Search   string `form:"search"`
	Page     int    `form:"page" default:"1"`
	Limit    int    `form:"limit" default:"20"`
	SortBy   string `form:"sort_by" default:"topic_order"`
	SortOrder string `form:"sort_order" default:"asc"`
}

// BulkCreateTopicsRequest represents the request to bulk create topics
type BulkCreateTopicsRequest struct {
	ModuleID string `json:"module_id" binding:"required"`
	Topics   []struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		TopicOrder  int    `json:"topic_order"`
	} `json:"topics" binding:"required,min=1"`
}

// BulkTopicResult represents the result of a bulk operation
type BulkTopicResult struct {
	SuccessCount int               `json:"success_count"`
	FailedCount  int               `json:"failed_count"`
	Created      []TopicResponse   `json:"created"`
	Errors       []BulkTopicError  `json:"errors"`
}

// BulkTopicError represents an error in bulk topic creation
type BulkTopicError struct {
	Title string `json:"title"`
	Error string `json:"error"`
}

// ReorderTopicsRequest represents the request to reorder topics
type ReorderTopicsRequest struct {
	TopicOrders []struct {
		ID         string `json:"id" binding:"required"`
		TopicOrder int    `json:"topic_order" binding:"required"`
	} `json:"topic_orders" binding:"required,min=1"`
}

// TopicStats represents statistics for topics
type TopicStats struct {
	TotalTopics   int64 `json:"total_topics"`
	TotalModules  int64 `json:"total_modules"`
	TopicsPerModule int64 `json:"topics_per_module"`
}