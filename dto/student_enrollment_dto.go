// dto/student_enrollment_dto.go
package dto

import (
	"time"
)

// CreateStudentEnrollmentRequest represents the request body for creating a student enrollment
type CreateStudentEnrollmentRequest struct {
	StudentID      string `json:"student_id" binding:"required"`
	GradeID        string `json:"grade_id" binding:"required"`
	Status         string `json:"status" binding:"omitempty,oneof=active inactive transferred graduated withdrawn"`
	GraduationDate string `json:"graduation_date"` // Format: "2006-01-02"
	Notes          string `json:"notes"`
	IsVerified		bool	`json:"is_verified"`
}

// UpdateStudentEnrollmentRequest represents the request body for updating a student enrollment
type UpdateStudentEnrollmentRequest struct {
	GradeID        string `json:"grade_id"`
	Status         string `json:"status" binding:"omitempty,oneof=active inactive transferred graduated withdrawn"`
	GraduationDate string `json:"graduation_date"` // Format: "2006-01-02"
	Notes          string `json:"notes"`
	IsVerified      bool `json:"is_verified"`
}

// StudentEnrollmentResponse represents the student enrollment response
type StudentEnrollmentResponse struct {
	ID             string     `json:"id"`
	StudentID      string     `json:"student_id"`
	GradeID        string     `json:"grade_id"`
	Status         string     `json:"status"`
	GraduationDate *time.Time `json:"graduation_date,omitempty"`
	Notes          string     `json:"notes"`
	IsVerified		bool	`json:"is_verified"`
	CreatedBy      string     `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Nested relationships
	Student *UserResponse       `json:"student,omitempty"`
	Grade   *ClassGradeResponse `json:"grade,omitempty"`
}

// StudentEnrollmentListResponse represents paginated student enrollment list response
type StudentEnrollmentListResponse struct {
	Enrollments []StudentEnrollmentResponse `json:"enrollments"`
	Total       int64                       `json:"total"`
	Page        int                         `json:"page"`
	Limit       int                         `json:"limit"`
	TotalPages  int                         `json:"total_pages"`
}

// StudentEnrollmentQueryParams represents query parameters for filtering student enrollments
type StudentEnrollmentQueryParams struct {
	StudentID string `form:"student_id"`
	GradeID   string `form:"grade_id"`
	Status    string `form:"status" binding:"omitempty,oneof=active inactive transferred graduated withdrawn"`
	Page      int    `form:"page" binding:"min=1"`
	Limit     int    `form:"limit" binding:"min=1,max=100"`
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=student_id grade_id status enrollment_date created_at"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// BulkCreateStudentEnrollmentRequest represents the request body for bulk creating student enrollments
type BulkCreateStudentEnrollmentRequest struct {
	StudentIDs     []string `json:"student_ids" binding:"required,min=1"`
	GradeID        string   `json:"grade_id" binding:"required"`
	Status         string   `json:"status" binding:"omitempty,oneof=active inactive transferred graduated withdrawn"`
	Notes          string   `json:"notes"`
	IsVerified		bool	`json:"is_verified"`

}

// BulkEnrollmentResult represents the result of bulk enrollment operation
type BulkEnrollmentResult struct {
	SuccessCount int                           `json:"success_count"`
	FailedCount  int                           `json:"failed_count"`
	Created      []StudentEnrollmentResponse   `json:"created"`
	Errors       []BulkEnrollmentError         `json:"errors,omitempty"`
}

type BulkEnrollmentError struct {
	StudentID string `json:"student_id"`
	Error     string `json:"error"`
}


// dto/student_enrollment_dto.go

type ClassTeacherDashboardStats struct {
	TotalGrades      int          `json:"total_grades"`
	TotalStudents    int          `json:"total_students"`
	ActiveStudents   int          `json:"active_students"`
	InactiveStudents int          `json:"inactive_students"`
	Grades           []GradeStats `json:"grades"`
}

type GradeStats struct {
	GradeID        string `json:"grade_id"`
	GradeName      string `json:"grade_name"`
	GradeCode      string `json:"grade_code"`
	Level          int    `json:"level"`
	TotalStudents  int    `json:"total_students"`
	ActiveStudents int    `json:"active_students"`
	Capacity       int    `json:"capacity"`
}