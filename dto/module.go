// dto/module_dto.go
package dto

import (
	"time"
)

// CreateModuleRequest represents the request to create a module
type CreateModuleRequest struct {
	SubjectID   string `json:"subject_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Sequence    int    `json:"sequence"`
}

// UpdateModuleRequest represents the request to update a module
type UpdateModuleRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Sequence    int    `json:"sequence"`
}

// ModuleResponse represents the response for a module
type ModuleResponse struct {
	ID          string    `json:"id"`
	SubjectID   string    `json:"subject_id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Sequence    int       `json:"sequence"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Nested relationships
	Subject *SubjectResponse `json:"subject,omitempty"`
}

// ModuleListResponse represents a paginated list of modules
type ModuleListResponse struct {
	Modules    []ModuleResponse `json:"modules"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}

// ModuleQueryParams represents query parameters for filtering modules
type ModuleQueryParams struct {
	SubjectID string `form:"subject_id"`
	Search    string `form:"search"`
	Page      int    `form:"page" default:"1"`
	Limit     int    `form:"limit" default:"20"`
	SortBy    string `form:"sort_by" default:"sequence"`
	SortOrder string `form:"sort_order" default:"asc"`
}

// BulkCreateModulesRequest represents the request to bulk create modules
type BulkCreateModulesRequest struct {
	SubjectID string `json:"subject_id" binding:"required"`
	Modules   []struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code"`
		Description string `json:"description"`
		Sequence    int    `json:"sequence"`
	} `json:"modules" binding:"required,min=1"`
}

// BulkModuleResult represents the result of a bulk operation
type BulkModuleResult struct {
	SuccessCount int                 `json:"success_count"`
	FailedCount  int                 `json:"failed_count"`
	Created      []ModuleResponse    `json:"created"`
	Errors       []BulkModuleError   `json:"errors"`
}

// BulkModuleError represents an error in bulk module creation
type BulkModuleError struct {
	Name  string `json:"name"`
	Code  string `json:"code"`
	Error string `json:"error"`
}

// ModuleStats represents statistics for modules
type ModuleStats struct {
	TotalModules   int64 `json:"total_modules"`
	ModulesPerSubject int64 `json:"modules_per_subject"`
	TotalSubjects  int64 `json:"total_subjects"`
}