// dto/exercise_dto.go
package dto

import (
	"time"
)

// CreateExerciseRequest represents the request to create an exercise
type CreateExerciseRequest struct {
	LessonID     string  `json:"lesson_id" binding:"required"`
	Title        string  `json:"title" binding:"required"`
	Instructions string  `json:"instructions"`
	Content      string  `json:"content"`
	TotalMarks   float64 `json:"total_marks"`
}

// UpdateExerciseRequest represents the request to update an exercise
type UpdateExerciseRequest struct {
	Title        string  `json:"title"`
	Instructions string  `json:"instructions"`
	Content      string  `json:"content"`
	TotalMarks   float64 `json:"total_marks"`
}

// ExerciseResponse represents the response for an exercise
type ExerciseResponse struct {
	ID           string    `json:"id"`
	LessonID     string    `json:"lesson_id"`
	Title        string    `json:"title"`
	Instructions string    `json:"instructions"`
	Content      string    `json:"content"`
	TotalMarks   float64   `json:"total_marks"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Nested relationships
	Lesson  *LessonResponse `json:"lesson,omitempty"`
	Creator *UserResponse   `json:"creator,omitempty"`
}

// ExerciseListResponse represents a paginated list of exercises
type ExerciseListResponse struct {
	Exercises  []ExerciseResponse `json:"exercises"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}

// ExerciseQueryParams represents query parameters for filtering exercises
type ExerciseQueryParams struct {
	LessonID string `form:"lesson_id"`
	Search   string `form:"search"`
	Page     int    `form:"page" default:"1"`
	Limit    int    `form:"limit" default:"20"`
	SortBy   string `form:"sort_by" default:"created_at"`
	SortOrder string `form:"sort_order" default:"desc"`
}

// BulkCreateExercisesRequest represents the request to bulk create exercises
type BulkCreateExercisesRequest struct {
	LessonID string `json:"lesson_id" binding:"required"`
	Exercises []struct {
		Title        string  `json:"title" binding:"required"`
		Instructions string  `json:"instructions"`
		Content      string  `json:"content"`
		TotalMarks   float64 `json:"total_marks"`
	} `json:"exercises" binding:"required,min=1"`
}

// BulkExerciseResult represents the result of a bulk operation
type BulkExerciseResult struct {
	SuccessCount int                  `json:"success_count"`
	FailedCount  int                  `json:"failed_count"`
	Created      []ExerciseResponse   `json:"created"`
	Errors       []BulkExerciseError  `json:"errors"`
}

// BulkExerciseError represents an error in bulk exercise creation
type BulkExerciseError struct {
	Title string `json:"title"`
	Error string `json:"error"`
}

// ExerciseStats represents statistics for exercises
type ExerciseStats struct {
	TotalExercises   int64 `json:"total_exercises"`
	TotalMarks       float64 `json:"total_marks"`
	AverageMarks     float64 `json:"average_marks"`
	ExercisesPerLesson float64 `json:"exercises_per_lesson"`
}