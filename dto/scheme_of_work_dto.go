package dto

import (
	"time"
)

// CreateSchemeOfWorkRequest represents the request to create a scheme of work
type CreateSchemeOfWorkRequest struct {
	SubjectID   string `json:"subject_id" binding:"required"`
	GradeID     string `json:"grade_id" binding:"required"`
	Term        string `json:"term" binding:"required,oneof=first second third"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status"` // draft, published, archived
}

// UpdateSchemeOfWorkRequest represents the request to update a scheme of work
type UpdateSchemeOfWorkRequest struct {
	SubjectID   string `json:"subject_id"`
	GradeID     string `json:"grade_id"`
	Term        string `json:"term" binding:"omitempty,oneof=first second third"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// SchemeOfWorkResponse represents the response for a scheme of work
type SchemeOfWorkResponse struct {
	ID          string    `json:"id"`
	SubjectID   string    `json:"subject_id"`
	GradeID     string    `json:"grade_id"`
	Term        string    `json:"term"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Nested relationships
	Subject *SubjectResponse `json:"subject,omitempty"`
	Grade   *ClassGradeResponse `json:"grade,omitempty"`
	Creator *UserResponse    `json:"creator,omitempty"`
	Lessons []LessonResponse `json:"lessons,omitempty"`
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
	GradeID   string `form:"grade_id"`
	Term      string `form:"term" binding:"omitempty,oneof=first second third"`
	Status    string `form:"status" binding:"omitempty,oneof=draft published archived"`
	Search    string `form:"search"`
	Page      int    `form:"page" default:"1"`
	Limit     int    `form:"limit" default:"20"`
	SortBy    string `form:"sort_by" default:"created_at"`
	SortOrder string `form:"sort_order" default:"desc"`
}

// BulkCreateSchemesRequest represents the request to bulk create schemes of work
type BulkCreateSchemesRequest struct {
	SubjectID string `json:"subject_id" binding:"required"`
	GradeID   string `json:"grade_id" binding:"required"`
	Term      string `json:"term" binding:"required,oneof=first second third"`
	Schemes   []struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	} `json:"schemes" binding:"required,min=1"`
	Status string `json:"status"`
}

// BulkSchemeResult represents the result of a bulk operation
type BulkSchemeResult struct {
	SuccessCount int                      `json:"success_count"`
	FailedCount  int                      `json:"failed_count"`
	Created      []SchemeOfWorkResponse   `json:"created"`
	Errors       []BulkSchemeError        `json:"errors"`
}

// BulkSchemeError represents an error in bulk scheme creation
type BulkSchemeError struct {
	Title string `json:"title"`
	Error string `json:"error"`
}

// SchemeOfWorkStats represents statistics for schemes of work
type SchemeOfWorkStats struct {
	TotalSchemes    int64 `json:"total_schemes"`
	DraftSchemes    int64 `json:"draft_schemes"`
	PublishedSchemes int64 `json:"published_schemes"`
	ArchivedSchemes  int64 `json:"archived_schemes"`
	TotalSubjects   int64 `json:"total_subjects"`
	TotalGrades     int64 `json:"total_grades"`
}