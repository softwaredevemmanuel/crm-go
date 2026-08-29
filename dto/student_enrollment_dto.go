// dto/student_enrollment_dto.go
package dto

import (
	"time"
)

// CreateStudentEnrollmentRequest represents the request to create a student enrollment
type CreateStudentEnrollmentRequest struct {
	StudentID      string `json:"student_id" binding:"required"`
	ArmID          string `json:"arm_id" binding:"required"`
	Status         string `json:"status"` // active, inactive, transferred, graduated, withdrawn
	GraduationDate string `json:"graduation_date"`
	Notes          string `json:"notes"`
	IsVerified     bool   `json:"is_verified"`
}

// UpdateStudentEnrollmentRequest represents the request to update a student enrollment
type UpdateStudentEnrollmentRequest struct {
	ArmID          string `json:"arm_id"`
	Status         string `json:"status"`
	GraduationDate string `json:"graduation_date"`
	Notes          string `json:"notes"`
	IsVerified     *bool  `json:"is_verified"`
}

// StudentEnrollmentResponse represents the response for a student enrollment
type StudentEnrollmentResponse struct {
	ID             string     `json:"id"`
	StudentID      string     `json:"student_id"`
	ArmID          string     `json:"arm_id"`
	Status         string     `json:"status"`
	GraduationDate *time.Time `json:"graduation_date,omitempty"`
	Notes          string     `json:"notes"`
	IsVerified     bool       `json:"is_verified"`
	CreatedBy      string     `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Nested relationships
	Student     *UserResponse     `json:"student,omitempty"`
	Arm         *ArmResponse      `json:"arm,omitempty"`
	ClassTeacher *UserResponse    `json:"class_teacher,omitempty"`
}



// StudentEnrollmentListResponse represents a paginated list of student enrollments
type StudentEnrollmentListResponse struct {
	Enrollments []StudentEnrollmentResponse `json:"enrollments"`
	Total       int64                       `json:"total"`
	Page        int                         `json:"page"`
	Limit       int                         `json:"limit"`
	TotalPages  int                         `json:"total_pages"`
}

// StudentEnrollmentQueryParams represents query parameters for filtering enrollments
type StudentEnrollmentQueryParams struct {
	StudentID string `form:"student_id"`
	ArmID     string `form:"arm_id"`
	Status    string `form:"status"`
	IsVerified *bool  `form:"is_verified"`
	Page      int    `form:"page" default:"1"`
	Limit     int    `form:"limit" default:"20"`
	SortBy    string `form:"sort_by" default:"created_at"`
	SortOrder string `form:"sort_order" default:"desc"`
}

// BulkCreateStudentEnrollmentsRequest represents the request to bulk create enrollments
type BulkCreateStudentEnrollmentsRequest struct {
	ArmID          string   `json:"arm_id" binding:"required"`
	StudentIDs     []string `json:"student_ids" binding:"required,min=1"`
	Status         string   `json:"status"`
	GraduationDate string   `json:"graduation_date"`
	Notes          string   `json:"notes"`
	IsVerified     bool     `json:"is_verified"`
}

// BulkEnrollmentResult represents the result of a bulk enrollment operation
type BulkEnrollmentResult struct {
	SuccessCount int                           `json:"success_count"`
	FailedCount  int                           `json:"failed_count"`
	Created      []StudentEnrollmentResponse   `json:"created"`
	Errors       []BulkEnrollmentError         `json:"errors"`
}

// BulkEnrollmentError represents an error in bulk enrollment
type BulkEnrollmentError struct {
	StudentID string `json:"student_id"`
	Error     string `json:"error"`
}

// StudentEnrollmentStats represents statistics for student enrollments
type StudentEnrollmentStats struct {
	TotalEnrollments   int64 `json:"total_enrollments"`
	ActiveEnrollments  int64 `json:"active_enrollments"`
	InactiveEnrollments int64 `json:"inactive_enrollments"`
	GraduatedEnrollments int64 `json:"graduated_enrollments"`
	TransferredEnrollments int64 `json:"transferred_enrollments"`
	WithdrawnEnrollments int64 `json:"withdrawn_enrollments"`
	VerifiedEnrollments int64 `json:"verified_enrollments"`
	UnverifiedEnrollments int64 `json:"unverified_enrollments"`
}