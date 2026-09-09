// dto/objective_question_dto.go
package dto

import (
	"time"
)

// CreateObjectiveQuestionRequest represents the request to create an objective question
type CreateObjectiveQuestionRequest struct {
	QuestionText      string `json:"question_text" binding:"required"`
	QuestionType      string `json:"question_type" binding:"required,oneof=multiple_choice true_false multiple_response matching ordering"`
	DifficultyLevel   string `json:"difficulty_level" binding:"required,oneof=easy medium hard expert"`
	Points            int    `json:"points" binding:"min=1"`
	ImageURL          string `json:"image_url"`
	VideoURL          string `json:"video_url"`
	SubjectID         string `json:"subject_id" binding:"required"`
	SchemeOfWorkID    string `json:"scheme_of_work_id" binding:"required"`
	ModuleID          string `json:"module_id"`
	TopicID           string `json:"topic_id"`
	LessonID          string `json:"lesson_id"`
	AnswerExplanation string `json:"answer_explanation"`
	SolutionSteps     string `json:"solution_steps"`
	Hint              string `json:"hint"`
	Status            string `json:"status" binding:"omitempty,oneof=active inactive draft archived"`
	Options           []CreateOptionRequest `json:"options" binding:"required,min=2"`
}

// CreateOptionRequest represents a question option
type CreateOptionRequest struct {
	OptionText string `json:"option_text" binding:"required"`
	IsCorrect  bool   `json:"is_correct"`
	Order      int    `json:"order"`
}

// UpdateObjectiveQuestionRequest represents the request to update an objective question
type UpdateObjectiveQuestionRequest struct {
	QuestionText      string `json:"question_text"`
	QuestionType      string `json:"question_type" binding:"omitempty,oneof=multiple_choice true_false multiple_response matching ordering"`
	DifficultyLevel   string `json:"difficulty_level" binding:"omitempty,oneof=easy medium hard expert"`
	Points            int    `json:"points" binding:"min=1"`
	ImageURL          string `json:"image_url"`
	VideoURL          string `json:"video_url"`
	AnswerExplanation string `json:"answer_explanation"`
	SolutionSteps     string `json:"solution_steps"`
	Hint              string `json:"hint"`
	Status            string `json:"status" binding:"omitempty,oneof=active inactive draft archived"`
	Options           []UpdateOptionRequest `json:"options"`
}

// UpdateOptionRequest represents a question option update
type UpdateOptionRequest struct {
	ID         string `json:"id"`
	OptionText string `json:"option_text" binding:"required"`
	IsCorrect  bool   `json:"is_correct"`
	Order      int    `json:"order"`
}

// ObjectiveQuestionQueryParams represents query parameters for filtering questions
type ObjectiveQuestionQueryParams struct {
	SubjectID      string `form:"subject_id"`
	SchemeOfWorkID string `form:"scheme_of_work_id"`
	ModuleID       string `form:"module_id"`
	TopicID        string `form:"topic_id"`
	LessonID       string `form:"lesson_id"`
	QuestionType   string `form:"question_type" binding:"omitempty,oneof=multiple_choice true_false multiple_response matching ordering"`
	DifficultyLevel string `form:"difficulty_level" binding:"omitempty,oneof=easy medium hard expert"`
	Status         string `form:"status" binding:"omitempty,oneof=active inactive draft archived"`
	Search         string `form:"search"`
	Page           int    `form:"page" default:"1"`
	Limit          int    `form:"limit" default:"20"`
	SortBy         string `form:"sort_by" default:"created_at"`
	SortOrder      string `form:"sort_order" default:"desc"`
}

// ObjectiveQuestionResponse represents the response for an objective question
type ObjectiveQuestionResponse struct {
	ID                string    `json:"id"`
	QuestionText      string    `json:"question_text"`
	QuestionType      string    `json:"question_type"`
	DifficultyLevel   string    `json:"difficulty_level"`
	Points            int       `json:"points"`
	ImageURL          string    `json:"image_url,omitempty"`
	VideoURL          string    `json:"video_url,omitempty"`
	SubjectID         string    `json:"subject_id"`
	SchemeOfWorkID    string    `json:"scheme_of_work_id"`
	ModuleID          string    `json:"module_id"`
	TopicID           string    `json:"topic_id"`
	LessonID          string    `json:"lesson_id"`
	AnswerExplanation string    `json:"answer_explanation,omitempty"`
	SolutionSteps     string    `json:"solution_steps,omitempty"`
	Hint              string    `json:"hint,omitempty"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// Nested relationships
	Subject      *SubjectResponse      `json:"subject,omitempty"`
	SchemeOfWork *SchemeOfWorkResponse `json:"scheme_of_work,omitempty"`
	Module       *ModuleResponse       `json:"module,omitempty"`
	Topic        *TopicResponse        `json:"topic,omitempty"`
	Lesson       *LessonResponse       `json:"lesson,omitempty"`
	Options      []OptionResponse      `json:"options,omitempty"`
}

// OptionResponse represents a question option
type OptionResponse struct {
	ID         string `json:"id"`
	OptionText string `json:"option_text"`
	IsCorrect  bool   `json:"is_correct"`
	Order      int    `json:"order"`
}

// ObjectiveQuestionListResponse represents a paginated list of questions
type ObjectiveQuestionListResponse struct {
	Questions  []ObjectiveQuestionResponse `json:"questions"`
	Total      int64                       `json:"total"`
	Page       int                         `json:"page"`
	Limit      int                         `json:"limit"`
	TotalPages int                         `json:"total_pages"`
}


