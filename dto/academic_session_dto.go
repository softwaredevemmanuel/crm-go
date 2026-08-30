// dto/academic_session_dto.go
package dto

import (
	"time"
)

// CreateAcademicSessionRequest represents the request to create an academic session
type CreateAcademicSessionRequest struct {
	AcademicYear string `json:"academic_year" binding:"required"`
	Code         string `json:"code" binding:"required"`
	StartDate    string `json:"start_date" binding:"required"`
	EndDate      string `json:"end_date" binding:"required"`
	Status       string `json:"status"` // active, inactive, completed
	IsCurrent    bool   `json:"is_current"`
	Description  string `json:"description"`
}

// UpdateAcademicSessionRequest represents the request to update an academic session
type UpdateAcademicSessionRequest struct {
	AcademicYear string `json:"academic_year"`
	Code         string `json:"code"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Status       string `json:"status"`
	IsCurrent    *bool  `json:"is_current"`
	Description  string `json:"description"`
}

// AcademicSessionResponse represents the response for an academic session
type AcademicSessionResponse struct {
	ID           string    `json:"id"`
	AcademicYear string    `json:"academic_year"`
	Code         string    `json:"code"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Status       string    `json:"status"`
	IsCurrent    bool      `json:"is_current"`
	Description  string    `json:"description"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Nested relationships
	Creator *UserResponse `json:"creator,omitempty"`
	Terms   []TermResponse `json:"terms,omitempty"`
}

// AcademicSessionListResponse represents a paginated list of academic sessions
type AcademicSessionListResponse struct {
	Sessions   []AcademicSessionResponse `json:"sessions"`
	Total      int64                     `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"total_pages"`
}

// AcademicSessionQueryParams represents query parameters for filtering academic sessions
type AcademicSessionQueryParams struct {
	Status    string `form:"status"`
	IsCurrent *bool  `form:"is_current"`
	Search    string `form:"search"`
	Page      int    `form:"page" default:"1"`
	Limit     int    `form:"limit" default:"20"`
	SortBy    string `form:"sort_by" default:"created_at"`
	SortOrder string `form:"sort_order" default:"desc"`
}

// AcademicSessionStats represents statistics for academic sessions
type AcademicSessionStats struct {
	TotalSessions     int64 `json:"total_sessions"`
	ActiveSessions    int64 `json:"active_sessions"`
	InactiveSessions  int64 `json:"inactive_sessions"`
	CompletedSessions int64 `json:"completed_sessions"`
	CurrentSessions   int64 `json:"current_sessions"`
	TotalTerms        int64 `json:"total_terms"`
}