// dto/module_dto.go
package dto

import (
	"time"
)

// CreateModuleRequest represents the request to create a module
type CreateModuleRequest struct {
	SchemeOfWorkID string `json:"scheme_of_work_id" binding:"required"`
	Title          string `json:"title" binding:"required"`
	Description    string `json:"description"`
	ModuleOrder    int    `json:"module_order"`
}

// UpdateModuleRequest represents the request to update a module
type UpdateModuleRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ModuleOrder int    `json:"module_order"`
}

// ModuleResponse represents the response for a module
type ModuleResponse struct {
	ID             string    `json:"id"`
	SchemeOfWorkID string    `json:"scheme_of_work_id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	ModuleOrder    int       `json:"module_order"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Nested relationships
	SchemeOfWork *SchemeOfWorkResponse `json:"scheme_of_work,omitempty"`
	Lessons      []LessonResponse      `json:"lessons,omitempty"`
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
	SchemeOfWorkID string `form:"scheme_of_work_id"`
	Search         string `form:"search"`
	Page           int    `form:"page" default:"1"`
	Limit          int    `form:"limit" default:"20"`
	SortBy         string `form:"sort_by" default:"module_order"`
	SortOrder      string `form:"sort_order" default:"asc"`
}

// BulkCreateModulesRequest represents the request to bulk create modules
type BulkCreateModulesRequest struct {
	SchemeOfWorkID string `json:"scheme_of_work_id" binding:"required"`
	Modules        []struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		ModuleOrder int    `json:"module_order"`
	} `json:"modules" binding:"required,min=1"`
}

// BulkModuleResult represents the result of a bulk operation
type BulkModuleResult struct {
	SuccessCount int                `json:"success_count"`
	FailedCount  int                `json:"failed_count"`
	Created      []ModuleResponse   `json:"created"`
	Errors       []BulkModuleError  `json:"errors"`
}

// BulkModuleError represents an error in bulk module creation
type BulkModuleError struct {
	Title string `json:"title"`
	Error string `json:"error"`
}

// ReorderModulesRequest represents the request to reorder modules
type ReorderModulesRequest struct {
	ModuleOrders []struct {
		ID          string `json:"id" binding:"required"`
		ModuleOrder int    `json:"module_order" binding:"required"`
	} `json:"module_orders" binding:"required,min=1"`
}

// ModuleStats represents statistics for modules
type ModuleStats struct {
	TotalModules   int64 `json:"total_modules"`
	TotalSchemes   int64 `json:"total_schemes"`
	ModulesPerScheme int64 `json:"modules_per_scheme"`
}