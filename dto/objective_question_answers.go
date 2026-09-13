// dto/student_answer_dto.go
package dto

import (
	"time"
)

// SubmitAnswerRequest represents the request to submit an answer
type SubmitAnswerRequest struct {
	QuestionID        string   `json:"question_id" binding:"required"`
	GradeID           string   `json:"grade_id" binding:"required"` // Added
	SelectedOptionIDs []string `json:"selected_option_ids"`          // For multiple selection
	SelectedOptionID  string   `json:"selected_option_id"`           // For single selection
	TimeSpent         int      `json:"time_spent"`                   // Time spent in seconds
}

// SubmitAnswerResponse represents the response after submitting an answer
type SubmitAnswerResponse struct {
	IsCorrect        bool   `json:"is_correct"`
	Score            int    `json:"score"`
	CorrectAnswer    string `json:"correct_answer,omitempty"`
	Explanation      string `json:"explanation,omitempty"`
	TotalAttempts    int    `json:"total_attempts"`
	PreviousAttempts int    `json:"previous_attempts"`
	TimeSpent        int    `json:"time_spent"`
}

// StudentAnswerResponse represents the response for a student answer
type StudentAnswerResponse struct {
	ID                string    `json:"id"`
	GradeID           string    `json:"grade_id"`
	StudentID         string    `json:"student_id"`
	QuestionID        string    `json:"question_id"`
	SelectedOptionIDs string    `json:"selected_option_ids,omitempty"`
	SelectedOptionID  string    `json:"selected_option_id,omitempty"`
	IsCorrect         bool      `json:"is_correct"`
	Score             int       `json:"score"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// Nested relationships
	Student  *UserResponse              `json:"student,omitempty"`
	Question *ObjectiveQuestionResponse `json:"question,omitempty"`
	Grade    *ClassGradeResponse        `json:"class_grade,omitempty"`
}

// StudentAnswerQueryParams represents query parameters for filtering answers
type StudentAnswerQueryParams struct {
	StudentID  string `form:"student_id"`
	LessonID   string `form:"lesson_id"`
	GradeID    string `form:"grade_id"`
	QuestionID string `form:"question_id"`
	IsCorrect  *bool  `form:"is_correct"`
	Status     string `form:"status" binding:"omitempty,oneof=active inactive"`
	Page       int    `form:"page" default:"1"`
	Limit      int    `form:"limit" default:"20"`
	SortBy     string `form:"sort_by" default:"created_at"`
	SortOrder  string `form:"sort_order" default:"desc"`
}

// StudentAnswerStats represents statistics for a student's answers
type StudentAnswerStats struct {
	TotalQuestions   int     `json:"total_questions"`
	CorrectAnswers   int     `json:"correct_answers"`
	IncorrectAnswers int     `json:"incorrect_answers"`
	ScorePercentage  float64 `json:"score_percentage"`
	TotalScore       int     `json:"total_score"`
}

// LessonTestResult represents aggregated results for a student in a lesson
type LessonTestResult struct {
	StudentID        string  `json:"student_id"`
	StudentFirstName string  `json:"student_first_name"`
	StudentLastName  string  `json:"student_last_name"`
	StudentEmail     string  `json:"student_email"`
	GradeID          string  `json:"grade_id"`
	GradeName        string  `json:"grade_name"`
	TotalQuestions   int     `json:"total_questions"`
	CorrectAnswers   int     `json:"correct_answers"`
	IncorrectAnswers int     `json:"incorrect_answers"`
	TotalScore       int     `json:"total_score"`
	MaxScore         int     `json:"max_score"`
	ScorePercentage  float64 `json:"score_percentage"`
	Status           string  `json:"status"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

// LessonTestStats represents statistics for a lesson test
type LessonTestStats struct {
	TotalStudents int     `json:"total_students"`
	AverageScore  float64 `json:"average_score"`
	HighestScore  float64 `json:"highest_score"`
	LowestScore   float64 `json:"lowest_score"`
	PassedCount   int     `json:"passed_count"`
	FailedCount   int     `json:"failed_count"`
}

// LessonTestResultsResponse represents the response for lesson test results
type LessonTestResultsResponse struct {
	Results []LessonTestResult `json:"results"`
	Total   int64              `json:"total"`
	Stats   *LessonTestStats   `json:"stats"`
}

// GradeTestResult represents aggregated results for a student in a grade
type GradeTestResult struct {
	StudentID        string  `json:"student_id"`
	StudentFirstName string  `json:"student_first_name"`
	StudentLastName  string  `json:"student_last_name"`
	StudentEmail     string  `json:"student_email"`
	GradeID          string  `json:"grade_id"`
	GradeName        string  `json:"grade_name"`
	TotalQuestions   int     `json:"total_questions"`
	CorrectAnswers   int     `json:"correct_answers"`
	IncorrectAnswers int     `json:"incorrect_answers"`
	TotalScore       int     `json:"total_score"`
	MaxScore         int     `json:"max_score"`
	ScorePercentage  float64 `json:"score_percentage"`
	Status           string  `json:"status"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

// GradeTestStats represents statistics for a grade test
type GradeTestStats struct {
	TotalStudents int     `json:"total_students"`
	AverageScore  float64 `json:"average_score"`
	HighestScore  float64 `json:"highest_score"`
	LowestScore   float64 `json:"lowest_score"`
	PassedCount   int     `json:"passed_count"`
	FailedCount   int     `json:"failed_count"`
}

// GradeTestResultsResponse represents the response for grade test results
type GradeTestResultsResponse struct {
	Results []GradeTestResult `json:"results"`
	Total   int64             `json:"total"`
	Stats   *GradeTestStats   `json:"stats"`
}





