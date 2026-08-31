// dto/test_dto.go
package dto

import (
	"time"
)

// CreateTestRequest represents the request to create a test
type CreateTestRequest struct {
	AcademicSessionID string  `json:"academic_session_id" binding:"required"`
	TermID            string  `json:"term_id" binding:"required"`
	SubjectID         string  `json:"subject_id" binding:"required"`
	ClassID           string  `json:"class_id" binding:"required"`
	Title             string  `json:"title" binding:"required"`
	TestType          string  `json:"test_type"`
	TestDate          string  `json:"test_date"`
	Duration          int     `json:"duration"`
	TotalMarks        float64 `json:"total_marks"`
	Status            string  `json:"status"` // draft, published, completed
}

// UpdateTestRequest represents the request to update a test
type UpdateTestRequest struct {
	AcademicSessionID string  `json:"academic_session_id"`
	TermID            string  `json:"term_id"`
	SubjectID         string  `json:"subject_id"`
	ClassID           string  `json:"class_id"`
	Title             string  `json:"title"`
	TestType          string  `json:"test_type"`
	TestDate          string  `json:"test_date"`
	Duration          int     `json:"duration"`
	TotalMarks        float64 `json:"total_marks"`
	Status            string  `json:"status"`
}

// TestResponse represents the response for a test
type TestResponse struct {
	ID                string     `json:"id"`
	AcademicSessionID string     `json:"academic_session_id"`
	TermID            string     `json:"term_id"`
	SubjectID         string     `json:"subject_id"`
	ClassID           string     `json:"class_id"`
	Title             string     `json:"title"`
	TestType          string     `json:"test_type"`
	TestDate          *time.Time `json:"test_date,omitempty"`
	Duration          int        `json:"duration"`
	TotalMarks        float64    `json:"total_marks"`
	Status            string     `json:"status"`
	CreatedBy         string     `json:"created_by"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// Nested relationships
	AcademicSession *AcademicSessionResponse `json:"academic_session,omitempty"`
	Term            *TermResponse            `json:"term,omitempty"`
	Subject         *SubjectResponse         `json:"subject,omitempty"`
	Class           *ClassGradeResponse      `json:"class,omitempty"`
	Creator         *UserResponse            `json:"creator,omitempty"`
}

// TestListResponse represents a paginated list of tests
type TestListResponse struct {
	Tests      []TestResponse `json:"tests"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

// TestQueryParams represents query parameters for filtering tests
type TestQueryParams struct {
	AcademicSessionID string `form:"academic_session_id"`
	TermID            string `form:"term_id"`
	SubjectID         string `form:"subject_id"`
	ClassID           string `form:"class_id"`
	Status            string `form:"status"`
	Search            string `form:"search"`
	Page              int    `form:"page" default:"1"`
	Limit             int    `form:"limit" default:"20"`
	SortBy            string `form:"sort_by" default:"test_date"`
	SortOrder         string `form:"sort_order" default:"desc"`
}

// BulkCreateTestsRequest represents the request to bulk create tests
type BulkCreateTestsRequest struct {
	Tests []struct {
		AcademicSessionID string  `json:"academic_session_id" binding:"required"`
		TermID            string  `json:"term_id" binding:"required"`
		SubjectID         string  `json:"subject_id" binding:"required"`
		ClassID           string  `json:"class_id" binding:"required"`
		Title             string  `json:"title" binding:"required"`
		TestType          string  `json:"test_type"`
		TestDate          string  `json:"test_date"`
		Duration          int     `json:"duration"`
		TotalMarks        float64 `json:"total_marks"`
	} `json:"tests" binding:"required,min=1"`
	Status string `json:"status"`
}

// BulkTestResult represents the result of a bulk operation
type BulkTestResult struct {
	SuccessCount int               `json:"success_count"`
	FailedCount  int               `json:"failed_count"`
	Created      []TestResponse    `json:"created"`
	Errors       []BulkTestError   `json:"errors"`
}

// BulkTestError represents an error in bulk test creation
type BulkTestError struct {
	Title string `json:"title"`
	Error string `json:"error"`
}

// TestStats represents statistics for tests
type TestStats struct {
	TotalTests     int64   `json:"total_tests"`
	DraftTests     int64   `json:"draft_tests"`
	PublishedTests int64   `json:"published_tests"`
	CompletedTests int64   `json:"completed_tests"`
	TotalMarks     float64 `json:"total_marks"`
	AverageMarks   float64 `json:"average_marks"`
	TotalSubjects  int64   `json:"total_subjects"`
	TotalClasses   int64   `json:"total_classes"`
}