// dto/student_answer_dto.go
package dto

import (
	"time"

)

// SubmitAnswerRequest represents the request to submit an answer
type SubmitAnswerRequest struct {
	QuestionID        string   `json:"question_id" binding:"required"`
	SelectedOptionIDs []string `json:"selected_option_ids"` // For multiple selection
	SelectedOptionID  string   `json:"selected_option_id"`  // For single selection
	TimeSpent         int      `json:"time_spent"`          // Time spent in seconds
}

// SubmitAnswerResponse represents the response after submitting an answer
type SubmitAnswerResponse struct {
	IsCorrect       bool   `json:"is_correct"`
	Score           int    `json:"score"`
	CorrectAnswer   string `json:"correct_answer,omitempty"`
	Explanation     string `json:"explanation,omitempty"`
	TotalAttempts   int    `json:"total_attempts"`
	PreviousAttempts int   `json:"previous_attempts"`
	TimeSpent       int    `json:"time_spent"`
}

// StudentAnswerResponse represents the response for a student answer
type StudentAnswerResponse struct {
	ID                string    `json:"id"`
	StudentID         string    `json:"student_id"`
	QuestionID        string    `json:"question_id"`
	SelectedOptionIDs string    `json:"selected_option_ids,omitempty"`
	SelectedOptionID  string    `json:"selected_option_id,omitempty"`
	IsCorrect         bool      `json:"is_correct"`
	Score             int       `json:"score"`
	TimeSpent         int       `json:"time_spent"`
	AttemptNumber     int       `json:"attempt_number"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// StudentAnswerStats represents statistics for a student's answers
type StudentAnswerStats struct {
	TotalQuestions   int     `json:"total_questions"`
	CorrectAnswers   int     `json:"correct_answers"`
	IncorrectAnswers int     `json:"incorrect_answers"`
	ScorePercentage  float64 `json:"score_percentage"`
	AverageTime      float64 `json:"average_time"`
	TotalScore       int     `json:"total_score"`
}

// StudentAnswerQueryParams represents query parameters for filtering answers
type StudentAnswerQueryParams struct {
	StudentID    string `form:"student_id"`
	QuestionID   string `form:"question_id"`
	LessonID     string `form:"lesson_id"`
	TopicID      string `form:"topic_id"`
	IsCorrect    *bool  `form:"is_correct"`
	Status       string `form:"status" binding:"omitempty,oneof=active inactive"`
	Page         int    `form:"page" default:"1"`
	Limit        int    `form:"limit" default:"20"`
	SortBy       string `form:"sort_by" default:"created_at"`
	SortOrder    string `form:"sort_order" default:"desc"`
}