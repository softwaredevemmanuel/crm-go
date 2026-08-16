// dto/grade_subject_dto.go
package dto

import (
	"time"
)

// CreateGradeSubjectRequest represents the request body for creating a grade-subject mapping
type CreateGradeSubjectRequest struct {
	GradeID      string `json:"grade_id" binding:"required"`
	SubjectID    string `json:"subject_id" binding:"required"`
	Status       string `json:"status" binding:"omitempty,oneof=active inactive"`
	IsCompulsory bool   `json:"is_compulsory"`
}

// BulkCreateGradeSubjectRequest represents the request body for bulk creating grade-subject mappings
type BulkCreateGradeSubjectRequest struct {
	GradeID      string   `json:"grade_id" binding:"required"`
	SubjectIDs   []string `json:"subject_ids" binding:"required,min=1"`
	Status       string   `json:"status" binding:"omitempty,oneof=active inactive"`
	IsCompulsory bool     `json:"is_compulsory"`
}

// UpdateGradeSubjectRequest represents the request body for updating a grade-subject mapping
type UpdateGradeSubjectRequest struct {
	Status       string `json:"status" binding:"omitempty,oneof=active inactive"`
	IsCompulsory *bool  `json:"is_compulsory"`
}

// GradeSubjectResponse represents the grade-subject response
type GradeSubjectResponse struct {
	ID           string    `json:"id"`
	GradeID      string    `json:"grade_id"`
	SubjectID    string    `json:"subject_id"`
	Status       string    `json:"status"`
	IsCompulsory bool      `json:"is_compulsory"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	// Nested relationships
	Grade   *ClassGradeResponse `json:"grade,omitempty"`
	Subject *SubjectResponse    `json:"subject,omitempty"`
}

// GradeSubjectListResponse represents paginated grade-subject list response
type GradeSubjectListResponse struct {
	GradeSubjects []GradeSubjectResponse `json:"grade_subjects"`
	Total         int64                  `json:"total"`
	Page          int                    `json:"page"`
	Limit         int                    `json:"limit"`
	TotalPages    int                    `json:"total_pages"`
}

// GradeSubjectQueryParams represents query parameters for filtering grade-subjects
type GradeSubjectQueryParams struct {
	GradeID      string `form:"grade_id"`
	SubjectID    string `form:"subject_id"`
	Status       string `form:"status" binding:"omitempty,oneof=active inactive"`
	IsCompulsory *bool  `form:"is_compulsory"`
	Page         int    `form:"page" binding:"min=1"`
	Limit        int    `form:"limit" binding:"min=1,max=100"`
	SortBy       string `form:"sort_by" binding:"omitempty,oneof=grade_id subject_id status is_compulsory created_at"`
	SortOrder    string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// BulkCreateResult represents the result of bulk create operation
type BulkCreateResult struct {
	SuccessCount int                      `json:"success_count"`
	FailedCount  int                      `json:"failed_count"`
	Created      []GradeSubjectResponse   `json:"created"`
	Errors       []BulkCreateError        `json:"errors,omitempty"`
}

type BulkCreateError struct {
	SubjectID string `json:"subject_id"`
	Error     string `json:"error"`
}