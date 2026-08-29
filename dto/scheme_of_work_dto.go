// dto/scheme_of_work_dto.go
package dto

import (
	"time"
)

// CreateSchemeOfWorkRequest represents the request to create a scheme of work
type CreateSchemeOfWorkRequest struct {
	SubjectID         string `json:"subject_id" binding:"required"`
	Grade             string `json:"grade" binding:"required"`
	Term              string `json:"term" binding:"required,oneof=first second third"`
	Week              int    `json:"week" binding:"required,min=1,max=52"`
	Topic             string `json:"topic" binding:"required"`
	Subtopics         string `json:"subtopics"`
	Objectives        string `json:"objectives" binding:"required"`
	Activities        string `json:"activities"`
	TeachingResources string `json:"teaching_resources"`
	Assessment        string `json:"assessment"`
	Status            string `json:"status"` // draft, published, archived
}

// UpdateSchemeOfWorkRequest represents the request to update a scheme of work
type UpdateSchemeOfWorkRequest struct {
	SubjectID         string `json:"subject_id"`
	Grade             string `json:"grade"`
	Term              string `json:"term" binding:"omitempty,oneof=first second third"`
	Week              int    `json:"week"`
	Topic             string `json:"topic"`
	Subtopics         string `json:"subtopics"`
	Objectives        string `json:"objectives"`
	Activities        string `json:"activities"`
	TeachingResources string `json:"teaching_resources"`
	Assessment        string `json:"assessment"`
	Status            string `json:"status"` // draft, published, archived
}

// SchemeOfWorkResponse represents the response for a scheme of work
type SchemeOfWorkResponse struct {
	ID                 string    `json:"id"`
	SubjectID          string    `json:"subject_id"`
	Grade              string    `json:"grade"`
	Term               string    `json:"term"`
	Week               int       `json:"week"`
	Topic              string    `json:"topic"`
	Subtopics          string    `json:"subtopics"`
	Objectives         string    `json:"objectives"`
	Activities         string    `json:"activities"`
	TeachingResources  string    `json:"teaching_resources"`
	Assessment         string    `json:"assessment"`
	Status             string    `json:"status"`
	CreatedBy          string    `json:"created_by"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`

	// Nested relationships
	Subject *SubjectResponse `json:"subject,omitempty"`
	Creator *UserResponse    `json:"creator,omitempty"`
}

// SchemeOfWorkListResponse represents a paginated list of schemes of work
type SchemeOfWorkListResponse struct {
	Schemes    []SchemeOfWorkResponse `json:"schemes"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	TotalPages int                    `json:"total_pages"`
}

// SchemeOfWorkQueryParams represents query parameters for filtering schemes of work
type SchemeOfWorkQueryParams struct {
	SubjectID string `form:"subject_id"`
	Grade     string `form:"grade"`
	Term      string `form:"term" binding:"omitempty,oneof=first second third"`
	Status    string `form:"status" binding:"omitempty,oneof=draft published archived"`
	Week      int    `form:"week"`
	Search    string `form:"search"`
	Page      int    `form:"page" default:"1"`
	Limit     int    `form:"limit" default:"20"`
	SortBy    string `form:"sort_by" default:"week"`
	SortOrder string `form:"sort_order" default:"asc"`
}

// BulkCreateSchemeOfWorkRequest represents the request to bulk create schemes of work
type BulkCreateSchemeOfWorkRequest struct {
	SubjectID string `json:"subject_id" binding:"required"`
	Grade     string `json:"grade" binding:"required"`
	Term      string `json:"term" binding:"required,oneof=first second third"`
	Schemes   []struct {
		Week              int    `json:"week" binding:"required"`
		Topic             string `json:"topic" binding:"required"`
		Subtopics         string `json:"subtopics"`
		Objectives        string `json:"objectives" binding:"required"`
		Activities        string `json:"activities"`
		TeachingResources string `json:"teaching_resources"`
		Assessment        string `json:"assessment"`
	} `json:"schemes" binding:"required,min=1"`
	Status string `json:"status"` // draft, published, archived
}

// BulkSchemeResult represents the result of a bulk operation
type BulkSchemeResult struct {
	SuccessCount int                    `json:"success_count"`
	FailedCount  int                    `json:"failed_count"`
	Created      []SchemeOfWorkResponse `json:"created"`
	Errors       []BulkSchemeError      `json:"errors"`
}

// BulkSchemeError represents an error in bulk scheme creation
type BulkSchemeError struct {
	Week  int    `json:"week"`
	Topic string `json:"topic"`
	Error string `json:"error"`
}

// SchemeOfWorkStats represents statistics for schemes of work
type SchemeOfWorkStats struct {
	TotalSchemes     int64 `json:"total_schemes"`
	DraftSchemes     int64 `json:"draft_schemes"`
	PublishedSchemes int64 `json:"published_schemes"`
	ArchivedSchemes  int64 `json:"archived_schemes"`
	TotalWeeks       int64 `json:"total_weeks"`
}

// SchemeOverview represents a summary of schemes for a subject/grade
type SchemeOverview struct {
	SubjectID   string `json:"subject_id"`
	SubjectName string `json:"subject_name"`
	Grade       string `json:"grade"`
	Term        string `json:"term"`
	TotalWeeks  int    `json:"total_weeks"`
	WeeksCovered int   `json:"weeks_covered"`
	Progress    string `json:"progress"` // percentage
	Status      string `json:"status"`
}