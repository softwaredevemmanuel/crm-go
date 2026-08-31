// dto/exam_dto.go
package dto

import (
	"time"
)

// CreateExamRequest represents the request to create an exam
type CreateExamRequest struct {
	AcademicSessionID string  `json:"academic_session_id" binding:"required"`
	TermID            string  `json:"term_id" binding:"required"`
	SubjectID         string  `json:"subject_id" binding:"required"`
	ClassID           string  `json:"class_id" binding:"required"`
	Title             string  `json:"title" binding:"required"`
	ExamType          string  `json:"exam_type"`
	ExamDate          string  `json:"exam_date"`
	Duration          int     `json:"duration"`
	TotalMarks        float64 `json:"total_marks"`
	Status            string  `json:"status"` // draft, published, completed
}

// UpdateExamRequest represents the request to update an exam
type UpdateExamRequest struct {
	AcademicSessionID string  `json:"academic_session_id"`
	TermID            string  `json:"term_id"`
	SubjectID         string  `json:"subject_id"`
	ClassID           string  `json:"class_id"`
	Title             string  `json:"title"`
	ExamType          string  `json:"exam_type"`
	ExamDate          string  `json:"exam_date"`
	Duration          int     `json:"duration"`
	TotalMarks        float64 `json:"total_marks"`
	Status            string  `json:"status"`
}

// ExamResponse represents the response for an exam
type ExamResponse struct {
	ID                string     `json:"id"`
	AcademicSessionID string     `json:"academic_session_id"`
	TermID            string     `json:"term_id"`
	SubjectID         string     `json:"subject_id"`
	ClassID           string     `json:"class_id"`
	Title             string     `json:"title"`
	ExamType          string     `json:"exam_type"`
	ExamDate          *time.Time `json:"exam_date,omitempty"`
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

// ExamListResponse represents a paginated list of exams
type ExamListResponse struct {
	Exams      []ExamResponse `json:"exams"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

// ExamQueryParams represents query parameters for filtering exams
type ExamQueryParams struct {
	AcademicSessionID string `form:"academic_session_id"`
	TermID            string `form:"term_id"`
	SubjectID         string `form:"subject_id"`
	ClassID           string `form:"class_id"`
	Status            string `form:"status"`
	Search            string `form:"search"`
	Page              int    `form:"page" default:"1"`
	Limit             int    `form:"limit" default:"20"`
	SortBy            string `form:"sort_by" default:"exam_date"`
	SortOrder         string `form:"sort_order" default:"desc"`
}

// BulkCreateExamsRequest represents the request to bulk create exams
type BulkCreateExamsRequest struct {
	Exams []struct {
		AcademicSessionID string  `json:"academic_session_id" binding:"required"`
		TermID            string  `json:"term_id" binding:"required"`
		SubjectID         string  `json:"subject_id" binding:"required"`
		ClassID           string  `json:"class_id" binding:"required"`
		Title             string  `json:"title" binding:"required"`
		ExamType          string  `json:"exam_type"`
		ExamDate          string  `json:"exam_date"`
		Duration          int     `json:"duration"`
		TotalMarks        float64 `json:"total_marks"`
	} `json:"exams" binding:"required,min=1"`
	Status string `json:"status"`
}

// BulkExamResult represents the result of a bulk operation
type BulkExamResult struct {
	SuccessCount int              `json:"success_count"`
	FailedCount  int              `json:"failed_count"`
	Created      []ExamResponse   `json:"created"`
	Errors       []BulkExamError  `json:"errors"`
}

// BulkExamError represents an error in bulk exam creation
type BulkExamError struct {
	Title string `json:"title"`
	Error string `json:"error"`
}

// ExamStats represents statistics for exams
type ExamStats struct {
	TotalExams     int64   `json:"total_exams"`
	DraftExams     int64   `json:"draft_exams"`
	PublishedExams int64   `json:"published_exams"`
	CompletedExams int64   `json:"completed_exams"`
	TotalMarks     float64 `json:"total_marks"`
	AverageMarks   float64 `json:"average_marks"`
	TotalSubjects  int64   `json:"total_subjects"`
	TotalClasses   int64   `json:"total_classes"`
}