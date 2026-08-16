// dto/academic_session_dto.go
package dto

import (
	"time"
)

// CreateAcademicSessionRequest represents the request body for creating an academic session
type CreateAcademicSessionRequest struct {
	AcademicYear        string `json:"academic_year" binding:"required,min=2,max=20"`
	Code        string `json:"code" binding:"required,min=2,max=20"`
	StartDate   string `json:"start_date" binding:"required"` // Format: "2006-01-02"
	EndDate     string `json:"end_date" binding:"required"`   // Format: "2006-01-02"
	Status      string `json:"status" binding:"omitempty,oneof=active inactive completed"`
	IsCurrent   bool   `json:"is_current"`
	Description string `json:"description"`
}

// UpdateAcademicSessionRequest represents the request body for updating an academic session
type UpdateAcademicSessionRequest struct {
	AcademicYear        string `json:"academic_year" binding:"omitempty,min=2,max=20"`
	Code        string `json:"code" binding:"omitempty,min=2,max=20"`
	StartDate   string `json:"start_date"` // Format: "2006-01-02"
	EndDate     string `json:"end_date"`   // Format: "2006-01-02"
	Status      string `json:"status" binding:"omitempty,oneof=active inactive completed"`
	IsCurrent   *bool  `json:"is_current"`
	Description string `json:"description"`
}

// AcademicSessionResponse represents the academic session response
type AcademicSessionResponse struct {
	ID          string    `json:"id"`
	AcademicYear        string    `json:"academic_year"`
	Code        string    `json:"code"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Status      string    `json:"status"`
	IsCurrent   bool      `json:"is_current"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DaysRemaining int     `json:"days_remaining,omitempty"`
	IsActive     bool     `json:"is_active,omitempty"`
}

// AcademicSessionListResponse represents paginated academic session list response
type AcademicSessionListResponse struct {
	Sessions   []AcademicSessionResponse `json:"sessions"`
	Total      int64                     `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"total_pages"`
}

// AcademicSessionQueryParams represents query parameters for filtering academic sessions
type AcademicSessionQueryParams struct {
	Search    string `form:"search"`
	Status    string `form:"status" binding:"omitempty,oneof=active inactive completed"`
	IsCurrent *bool  `form:"is_current"`
	Page      int    `form:"page" binding:"min=1"`
	Limit     int    `form:"limit" binding:"min=1,max=100"`
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=name code start_date end_date status created_at"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}