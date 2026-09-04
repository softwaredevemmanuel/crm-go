// dto/lesson_dto.go
package dto

import (
	"time"
)

// CreateLessonRequest represents the request to create a lesson
type CreateLessonRequest struct {
	SchemeOfWorkID string `json:"scheme_of_work_id" binding:"required"`
	ModuleID       string `json:"module_id" binding:"required"`
	TopicID        string `json:"topic_id" binding:"required"`
	LessonOrder    int    `json:"lesson_order"`
	Title          string `json:"title" binding:"required"`
	Description    string `json:"description"`
	LessonDate     string `json:"lesson_date"`
	Week           int    `json:"week"`
	Duration       int    `json:"duration"`
	Content        string `json:"content"` // JSON string
	Objectives     string `json:"objectives"`
	Activities     string `json:"activities"`
	Resources      string `json:"resources"`
	Assessment     string `json:"assessment"`
	Status         string `json:"status"` // planned, ongoing, completed, cancelled
}

// UpdateLessonRequest represents the request to update a lesson
type UpdateLessonRequest struct {
	SchemeOfWorkID string `json:"scheme_of_work_id"`
	ModuleID       string `json:"module_id"`
	TopicID        string `json:"topic_id"`
	LessonOrder    int    `json:"lesson_order"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	LessonDate     string `json:"lesson_date"`
	Week           int    `json:"week"`
	Duration       int    `json:"duration"`
	Content        string `json:"content"` // JSON string
	Objectives     string `json:"objectives"`
	Activities     string `json:"activities"`
	Resources      string `json:"resources"`
	Assessment     string `json:"assessment"`
	Status         string `json:"status"`
}

// LessonResponse represents the response for a lesson
type LessonResponse struct {
	ID             string     `json:"id"`
	SchemeOfWorkID string     `json:"scheme_of_work_id"`
	ModuleID       string     `json:"module_id"`
	TopicID        string     `json:"topic_id"`
	LessonOrder    int        `json:"lesson_order"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	LessonDate     *time.Time `json:"lesson_date,omitempty"`
	Week           int        `json:"week"`
	Duration       int        `json:"duration"`
	Content        string     `json:"content"` 
	Objectives     string     `json:"objectives"`
	Activities     string     `json:"activities"`
	Resources      string     `json:"resources"`
	Assessment     string     `json:"assessment"`
	Status         string     `json:"status"`
	CreatedBy      string     `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Nested relationships
	SchemeOfWork *SchemeOfWorkResponse `json:"scheme_of_work,omitempty"`
	Module       *ModuleResponse       `json:"module,omitempty"`
	Topic        *TopicResponse        `json:"topic,omitempty"`
	Creator      *UserResponse         `json:"creator,omitempty"`
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
	SchemeOfWorkID string `form:"scheme_of_work_id"`
	ModuleID       string `form:"module_id"`
	TopicID        string `form:"topic_id"`
	Status         string `form:"status"`
	Week           int    `form:"week"`
	Search         string `form:"search"`
	Page           int    `form:"page" default:"1"`
	Limit          int    `form:"limit" default:"20"`
	SortBy         string `form:"sort_by" default:"lesson_order"`
	SortOrder      string `form:"sort_order" default:"asc"`
}

// BulkCreateLessonsRequest represents the request to bulk create lessons
type BulkCreateLessonsRequest struct {
	SchemeOfWorkID string `json:"scheme_of_work_id" binding:"required"`
	ModuleID       string `json:"module_id" binding:"required"`
	TopicID        string `json:"topic_id" binding:"required"`
	Lessons        []struct {
		LessonOrder int    `json:"lesson_order"`
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		LessonDate  string `json:"lesson_date"`
		Week        int    `json:"week"`
		Duration    int    `json:"duration"`
		Content     string `json:"content"` 
		Objectives  string `json:"objectives"`
		Activities  string `json:"activities"`
		Resources   string `json:"resources"`
		Assessment  string `json:"assessment"`
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

// ReorderLessonsRequest represents the request to reorder lessons
type ReorderLessonsRequest struct {
	LessonOrders []struct {
		ID          string `json:"id" binding:"required"`
		LessonOrder int    `json:"lesson_order" binding:"required"`
	} `json:"lesson_orders" binding:"required,min=1"`
}

// LessonStats represents statistics for lessons
type LessonStats struct {
	TotalLessons    int64 `json:"total_lessons"`
	PlannedLessons  int64 `json:"planned_lessons"`
	OngoingLessons  int64 `json:"ongoing_lessons"`
	CompletedLessons int64 `json:"completed_lessons"`
	CancelledLessons int64 `json:"cancelled_lessons"`
	TotalSchemes    int64 `json:"total_schemes"`
	TotalModules    int64 `json:"total_modules"`
	TotalTopics     int64 `json:"total_topics"`
}