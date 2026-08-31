// dto/lesson_dto.go
package dto

import (
	"time"
)

// CreateLessonRequest represents the request to create a lesson
type CreateLessonRequest struct {
	SchemeOfWorkItemID string `json:"scheme_of_work_item_id" binding:"required"`
	ClassID            string `json:"class_id" binding:"required"`
	ArmID              string `json:"arm_id"`
	Title              string `json:"title" binding:"required"`
	LessonDate         string `json:"lesson_date"`
	Week               int    `json:"week"`
	Period             int    `json:"period"`
	Duration           int    `json:"duration"`
	Status             string `json:"status"` // planned, ongoing, completed, cancelled
}

// UpdateLessonRequest represents the request to update a lesson
type UpdateLessonRequest struct {
	SchemeOfWorkItemID string `json:"scheme_of_work_item_id"`
	ClassID            string `json:"class_id"`
	ArmID              string `json:"arm_id"`
	Title              string `json:"title"`
	LessonDate         string `json:"lesson_date"`
	Week               int    `json:"week"`
	Period             int    `json:"period"`
	Duration           int    `json:"duration"`
	Status             string `json:"status"`
}

// LessonResponse represents the response for a lesson
type LessonResponse struct {
	ID                 string     `json:"id"`
	SchemeOfWorkItemID string     `json:"scheme_of_work_item_id"`
	ClassID            string     `json:"class_id"`
	ArmID              string     `json:"arm_id"`
	Title              string     `json:"title"`
	LessonDate         *time.Time `json:"lesson_date,omitempty"`
	Week               int        `json:"week"`
	Period             int        `json:"period"`
	Duration           int        `json:"duration"`
	Status             string     `json:"status"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`

	// Nested relationships
	SchemeOfWorkItem *SchemeOfWorkItemResponse `json:"scheme_of_work_item,omitempty"`
	Class            *ClassGradeResponse      `json:"class,omitempty"`
	Arm              *ArmResponse             `json:"arm,omitempty"`
}

// LessonListResponse represents a paginated list of lessons
type LessonListResponse struct {
	Lessons    []LessonResponse `json:"lessons"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}

// LessonQueryParams represents query parameters for filtering lessons
type LessonQueryParams struct {
	SchemeOfWorkItemID string `form:"scheme_of_work_item_id"`
	ClassID            string `form:"class_id"`
	ArmID              string `form:"arm_id"`
	Status             string `form:"status"`
	Week               int    `form:"week"`
	Search             string `form:"search"`
	Page               int    `form:"page" default:"1"`
	Limit              int    `form:"limit" default:"20"`
	SortBy             string `form:"sort_by" default:"lesson_date"`
	SortOrder          string `form:"sort_order" default:"desc"`
}

// BulkCreateLessonsRequest represents the request to bulk create lessons
type BulkCreateLessonsRequest struct {
	ClassID   string `json:"class_id" binding:"required"`
	ArmID     string `json:"arm_id"`
	Lessons   []struct {
		SchemeOfWorkItemID string `json:"scheme_of_work_item_id" binding:"required"`
		Title              string `json:"title" binding:"required"`
		LessonDate         string `json:"lesson_date"`
		Week               int    `json:"week"`
		Period             int    `json:"period"`
		Duration           int    `json:"duration"`
	} `json:"lessons" binding:"required,min=1"`
	Status string `json:"status"`
}

// BulkLessonResult represents the result of a bulk operation
type BulkLessonResult struct {
	SuccessCount int                `json:"success_count"`
	FailedCount  int                `json:"failed_count"`
	Created      []LessonResponse   `json:"created"`
	Errors       []BulkLessonError  `json:"errors"`
}

// BulkLessonError represents an error in bulk lesson creation
type BulkLessonError struct {
	Title string `json:"title"`
	Error string `json:"error"`
}

// LessonStats represents statistics for lessons
type LessonStats struct {
	TotalLessons    int64 `json:"total_lessons"`
	PlannedLessons  int64 `json:"planned_lessons"`
	OngoingLessons  int64 `json:"ongoing_lessons"`
	CompletedLessons int64 `json:"completed_lessons"`
	CancelledLessons int64 `json:"cancelled_lessons"`
	TotalClasses    int64 `json:"total_classes"`
}