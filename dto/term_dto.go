// dto/term_dto.go
package dto

import (
	"time"
)

// CreateTermRequest represents the request to create a term
type CreateTermRequest struct {
	AcademicSessionID string `json:"academic_session_id" binding:"required"`
	Name              string `json:"name" binding:"required"`
	Code              string `json:"code"`
	TermNumber        int    `json:"term_number" binding:"required,min=1,max=3"`
	StartDate         string `json:"start_date"`
	EndDate           string `json:"end_date"`
	IsCurrent         bool   `json:"is_current"`
	Status            string `json:"status"` // active, inactive, completed
	Description       string `json:"description"`
}

// UpdateTermRequest represents the request to update a term
type UpdateTermRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	TermNumber  int    `json:"term_number"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	IsCurrent   *bool  `json:"is_current"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// TermResponse represents the response for a term
type TermResponse struct {
	ID                string    `json:"id"`
	AcademicSessionID string    `json:"academic_session_id"`
	Name              string    `json:"name"`
	Code              string    `json:"code"`
	TermNumber        int       `json:"term_number"`
	StartDate         *time.Time `json:"start_date,omitempty"`
	EndDate           *time.Time `json:"end_date,omitempty"`
	IsCurrent         bool      `json:"is_current"`
	Status            string    `json:"status"`
	Description       string    `json:"description"`
	CreatedBy         string    `json:"created_by"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// Nested relationships
	AcademicSession *AcademicSessionResponse `json:"academic_session,omitempty"`
	Creator         *UserResponse            `json:"creator,omitempty"`
}

// TermListResponse represents a paginated list of terms
type TermListResponse struct {
	Terms      []TermResponse `json:"terms"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

// TermQueryParams represents query parameters for filtering terms
type TermQueryParams struct {
	AcademicSessionID string `form:"academic_session_id"`
	Status            string `form:"status"`
	IsCurrent         *bool  `form:"is_current"`
	TermNumber        int    `form:"term_number"`
	Search            string `form:"search"`
	Page              int    `form:"page" default:"1"`
	Limit             int    `form:"limit" default:"20"`
	SortBy            string `form:"sort_by" default:"term_number"`
	SortOrder         string `form:"sort_order" default:"asc"`
}




// TermStats represents statistics for terms
type TermStats struct {
	TotalTerms        int64 `json:"total_terms"`
	ActiveTerms       int64 `json:"active_terms"`
	InactiveTerms     int64 `json:"inactive_terms"`
	CompletedTerms    int64 `json:"completed_terms"`
	CurrentTerms      int64 `json:"current_terms"`
	TermsPerSession   int64 `json:"terms_per_session"`
}