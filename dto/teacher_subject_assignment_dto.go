// dto/teacher_subject_assignment_dto.go
package dto

import (
	"time"
)

// CreateTeacherSubjectAssignmentRequest represents the request to assign a subject to a teacher
type CreateTeacherSubjectAssignmentRequest struct {
	GradeID   string `json:"grade_id" binding:"required"`
	SubjectID string `json:"subject_id" binding:"required"`
	TeacherID string `json:"teacher_id" binding:"required"`
	Status    string `json:"status"` // active, inactive (default: active)
}

// UpdateTeacherSubjectAssignmentRequest represents the request to update a subject assignment
type UpdateTeacherSubjectAssignmentRequest struct {
	GradeID   string `json:"grade_id"`
	SubjectID string `json:"subject_id"`
	TeacherID string `json:"teacher_id"`
	Status    string `json:"status"`
}

// BulkAssignSubjectsRequest represents the request to assign multiple subjects to a teacher
type BulkAssignSubjectsRequest struct {
	TeacherID string   `json:"teacher_id" binding:"required"`
	GradeID   string   `json:"grade_id" binding:"required"`
	SubjectIDs []string `json:"subject_ids" binding:"required,min=1"`
	Status    string   `json:"status"`
}

// TeacherSubjectAssignmentResponse represents the response for a subject assignment
type TeacherSubjectAssignmentResponse struct {
	ID        string    `json:"id"`
	GradeID   string    `json:"grade_id"`
	SubjectID string    `json:"subject_id"`
	TeacherID string    `json:"teacher_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Nested relationships
	Grade   *ClassGradeResponse `json:"grade,omitempty"`
	Subject *SubjectResponse    `json:"subject,omitempty"`
	Teacher *UserResponse       `json:"teacher,omitempty"`
}


// TeacherSubjectAssignmentListResponse represents a paginated list of assignments
type TeacherSubjectAssignmentListResponse struct {
	Assignments []TeacherSubjectAssignmentResponse `json:"assignments"`
	Total       int64                              `json:"total"`
	Page        int                                `json:"page"`
	Limit       int                                `json:"limit"`
	TotalPages  int                                `json:"total_pages"`
}

// TeacherSubjectAssignmentQueryParams represents query parameters for filtering assignments
type TeacherSubjectAssignmentQueryParams struct {
	GradeID   string `form:"grade_id"`
	SubjectID string `form:"subject_id"`
	TeacherID string `form:"teacher_id"`
	Status    string `form:"status"`
	Page      int    `form:"page" default:"1"`
	Limit     int    `form:"limit" default:"20"`
	SortBy    string `form:"sort_by" default:"created_at"`
	SortOrder string `form:"sort_order" default:"desc"`
}

// BulkAssignmentResult represents the result of a bulk assignment operation
type BulkAssignmentResult struct {
	SuccessCount int                              `json:"success_count"`
	FailedCount  int                              `json:"failed_count"`
	Created      []TeacherSubjectAssignmentResponse `json:"created"`
	Errors       []BulkAssignmentError            `json:"errors"`
}

// BulkAssignmentError represents an error in bulk assignment
type BulkAssignmentError struct {
	SubjectID string `json:"subject_id"`
	Error     string `json:"error"`
}