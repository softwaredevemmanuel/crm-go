// dto/scheme_of_work_dto.go
package dto

import (
	"time"
)

// CreateSchemeOfWorkRequest represents the request to create a scheme of work
type CreateSchemeOfWorkRequest struct {
	AcademicSessionID string `json:"academic_session_id" binding:"required"`
	TermID            string `json:"term_id" binding:"required"`
	SubjectID         string `json:"subject_id" binding:"required"`
	ClassID           string `json:"class_id" binding:"required"`
	Title             string `json:"title" binding:"required"`
	Description       string `json:"description"`
	Status            string `json:"status"` // draft, published, archived
}

// UpdateSchemeOfWorkRequest represents the request to update a scheme of work
type UpdateSchemeOfWorkRequest struct {
	AcademicSessionID string `json:"academic_session_id"`
	TermID            string `json:"term_id"`
	SubjectID         string `json:"subject_id"`
	ClassID           string `json:"class_id"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	Status            string `json:"status"`
}

// SchemeOfWorkResponse represents the response for a scheme of work
type SchemeOfWorkResponse struct {
	ID                string    `json:"id"`
	AcademicSessionID string    `json:"academic_session_id"`
	TermID            string    `json:"term_id"`
	SubjectID         string    `json:"subject_id"`
	ClassID           string    `json:"class_id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	Status            string    `json:"status"`
	CreatedBy         string    `json:"created_by"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// Nested relationships
	AcademicSession *AcademicSessionResponse `json:"academic_session,omitempty"`
	Term            *TermResponse            `json:"term,omitempty"`
	Subject         *SubjectResponse         `json:"subject,omitempty"`
	Class           *ClassGradeResponse      `json:"class,omitempty"`
	Creator         *UserResponse            `json:"creator,omitempty"`
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
	AcademicSessionID string `form:"academic_session_id"`
	TermID            string `form:"term_id"`
	SubjectID         string `form:"subject_id"`
	ClassID           string `form:"class_id"`
	Status            string `form:"status"`
	Search            string `form:"search"`
	Page              int    `form:"page" default:"1"`
	Limit             int    `form:"limit" default:"20"`
	SortBy            string `form:"sort_by" default:"created_at"`
	SortOrder         string `form:"sort_order" default:"desc"`
}

// BulkCreateSchemesRequest represents the request to bulk create schemes of work
type BulkCreateSchemesRequest struct {
	AcademicSessionID string `json:"academic_session_id" binding:"required"`
	TermID            string `json:"term_id" binding:"required"`
	SubjectID         string `json:"subject_id" binding:"required"`
	ClassID           string `json:"class_id" binding:"required"`
	Schemes           []struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	} `json:"schemes" binding:"required,min=1"`
	Status string `json:"status"`
}

// BulkSchemeResult represents the result of a bulk operation
type BulkSchemeResult struct {
	SuccessCount int                 `json:"success_count"`
	FailedCount  int                 `json:"failed_count"`
	Created      []SchemeOfWorkResponse `json:"created"`
	Errors       []BulkSchemeError   `json:"errors"`
}

// BulkSchemeError represents an error in bulk scheme creation
type BulkSchemeError struct {
	Title string `json:"title"`
	Error string `json:"error"`
}

// SchemeOfWorkStats represents statistics for schemes of work
type SchemeOfWorkStats struct {
	TotalSchemes      int64 `json:"total_schemes"`
	DraftSchemes      int64 `json:"draft_schemes"`
	PublishedSchemes  int64 `json:"published_schemes"`
	ArchivedSchemes   int64 `json:"archived_schemes"`
	TotalSubjects     int64 `json:"total_subjects"`
	TotalClasses      int64 `json:"total_classes"`
	TotalTerms        int64 `json:"total_terms"`
	TotalAcademicSessions int64 `json:"total_academic_sessions"`
}