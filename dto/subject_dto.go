// dto/subject_dto.go
package dto

import (
	"time"
)

// CreateSubjectRequest represents the request to create a subject
type CreateSubjectRequest struct {
	Name         string `json:"name" binding:"required"`
	Code         string `json:"code" binding:"required"`
	Description  string `json:"description"`
	DepartmentID string `json:"department_id" binding:"required"`
	Credits      int    `json:"credits"`
	Status       string `json:"status"` // active, inactive
}

// UpdateSubjectRequest represents the request to update a subject
type UpdateSubjectRequest struct {
	Name         string `json:"name"`
	Code         string `json:"code"`
	Description  string `json:"description"`
	DepartmentID string `json:"department_id"`
	Credits      int    `json:"credits"`
	Status       string `json:"status"`
}

// SubjectResponse represents the response for a subject
type SubjectResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	Description  string    `json:"description"`
	DepartmentID string    `json:"department_id"`
	Credits      int       `json:"credits"`
	Status       string    `json:"status"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Nested relationships
	Department *DepartmentResponse `json:"department,omitempty"`
	Creator    *UserResponse       `json:"creator,omitempty"`
}

// SubjectListResponse represents a paginated list of subjects
type SubjectListResponse struct {
	Subjects   []SubjectResponse `json:"subjects"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// SubjectQueryParams represents query parameters for filtering subjects
type SubjectQueryParams struct {
	DepartmentID string `form:"department_id"`
	Status       string `form:"status"`
	Search       string `form:"search"`
	Page         int    `form:"page" default:"1"`
	Limit        int    `form:"limit" default:"20"`
	SortBy       string `form:"sort_by" default:"name"`
	SortOrder    string `form:"sort_order" default:"asc"`
}

// BulkCreateSubjectsRequest represents the request to bulk create subjects
type BulkCreateSubjectsRequest struct {
	Subjects []struct {
		Name         string `json:"name" binding:"required"`
		Code         string `json:"code" binding:"required"`
		Description  string `json:"description"`
		DepartmentID string `json:"department_id" binding:"required"`
		Credits      int    `json:"credits"`
	} `json:"subjects" binding:"required,min=1"`
	Status string `json:"status"`
}

// BulkSubjectResult represents the result of a bulk operation
type BulkSubjectResult struct {
	SuccessCount int                `json:"success_count"`
	FailedCount  int                `json:"failed_count"`
	Created      []SubjectResponse  `json:"created"`
	Errors       []BulkSubjectError `json:"errors"`
}

// BulkSubjectError represents an error in bulk subject creation
type BulkSubjectError struct {
	Name  string `json:"name"`
	Code  string `json:"code"`
	Error string `json:"error"`
}

// SubjectStats represents statistics for subjects
type SubjectStats struct {
	TotalSubjects   int64 `json:"total_subjects"`
	ActiveSubjects  int64 `json:"active_subjects"`
	InactiveSubjects int64 `json:"inactive_subjects"`
	TotalDepartments int64 `json:"total_departments"`
}