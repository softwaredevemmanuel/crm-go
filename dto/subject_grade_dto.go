// dto/subject_grade_dto.go
package dto

import (
	"time"
)

// CreateSubjectGradeRequest represents the request body for creating a subject-grade relationship
type CreateSubjectGradeRequest struct {
	SubjectID    string `json:"subject_id" binding:"required"`
	GradeID      string `json:"grade_id" binding:"required"`
	AcademicYear string `json:"academic_year" binding:"required"`
	Status       string `json:"status" binding:"omitempty,oneof=active inactive archived"`
	IsRequired   bool   `json:"is_required"`
	Credits      int    `json:"credits" binding:"min=0,max=10"`
	Description  string `json:"description"`
}

// UpdateSubjectGradeRequest represents the request body for updating a subject-grade relationship
type UpdateSubjectGradeRequest struct {
	Status       string `json:"status" binding:"omitempty,oneof=active inactive archived"`
	IsRequired   *bool  `json:"is_required"`
	Credits      int    `json:"credits" binding:"min=0,max=10"`
	Description  string `json:"description"`
	AcademicYear string `json:"academic_year"`
}

// SubjectGradeResponse represents the subject-grade relationship response
type SubjectGradeResponse struct {
	ID           string    `json:"id"`
	SubjectID    string    `json:"subject_id"`
	GradeID      string    `json:"grade_id"`
	AcademicYear string    `json:"academic_year"`
	Status       string    `json:"status"`
	IsRequired   bool      `json:"is_required"`
	Credits      int       `json:"credits"`
	Description  string    `json:"description"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Subject      *SubjectResponse `json:"subject,omitempty"`
	Grade        *ClassGradeResponse `json:"grade,omitempty"`
}

// SubjectGradeListResponse represents paginated subject-grade list response
type SubjectGradeListResponse struct {
	SubjectGrades []SubjectGradeResponse `json:"subject_grades"`
	Total         int64                  `json:"total"`
	Page          int                    `json:"page"`
	Limit         int                    `json:"limit"`
	TotalPages    int                    `json:"total_pages"`
}

// SubjectGradeQueryParams represents query parameters for filtering subject-grade relationships
type SubjectGradeQueryParams struct {
	SubjectID    string `form:"subject_id"`
	GradeID      string `form:"grade_id"`
	AcademicYear string `form:"academic_year"`
	Status       string `form:"status"`
	IsRequired   *bool  `form:"is_required"`
	Page         int    `form:"page" binding:"min=1"`
	Limit        int    `form:"limit" binding:"min=1,max=100"`
	SortBy       string `form:"sort_by" binding:"omitempty,oneof=subject_id grade_id academic_year status created_at"`
	SortOrder    string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// BulkCreateSubjectGradeRequest represents the request body for bulk creating subject-grade relationships
type BulkCreateSubjectGradeRequest struct {
	SubjectIDs   []string `json:"subject_ids" binding:"required,min=1"`
	GradeID      string   `json:"grade_id" binding:"required"`
	AcademicYear string   `json:"academic_year" binding:"required"`
	Status       string   `json:"status" binding:"omitempty,oneof=active inactive archived"`
	IsRequired   bool     `json:"is_required"`
	Credits      int      `json:"credits" binding:"min=0,max=10"`
	Description  string   `json:"description"`
}

// BulkCreateSubjectGradeResponse represents the response for bulk creation
type BulkCreateSubjectGradeResponse struct {
	SuccessCount int                     `json:"success_count"`
	FailedCount  int                     `json:"failed_count"`
	Created      []SubjectGradeResponse  `json:"created"`
	Errors       []BulkCreateError       `json:"errors,omitempty"`
}

type BulkCreateError struct {
	SubjectID string `json:"subject_id"`
	Error     string `json:"error"`
}