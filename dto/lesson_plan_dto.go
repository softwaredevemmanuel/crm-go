// dto/lesson_plan_dto.go
package dto

import (
	"time"
)

// CreateLessonPlanRequest represents the request to create a lesson plan
type CreateLessonPlanRequest struct {
	LessonID               string `json:"lesson_id" binding:"required"`
	PreviousKnowledge      string `json:"previous_knowledge"`
	BehaviouralObjectives  string `json:"behavioural_objectives"`
	TeachingAids           string `json:"teaching_aids"`
	Introduction           string `json:"introduction"`
	LessonContent          string `json:"lesson_content"`
	TeacherActivities      string `json:"teacher_activities"`
	StudentActivities      string `json:"student_activities"`
	Conclusion             string `json:"conclusion"`
	Evaluation             string `json:"evaluation"`
}

// UpdateLessonPlanRequest represents the request to update a lesson plan
type UpdateLessonPlanRequest struct {
	PreviousKnowledge      string `json:"previous_knowledge"`
	BehaviouralObjectives  string `json:"behavioural_objectives"`
	TeachingAids           string `json:"teaching_aids"`
	Introduction           string `json:"introduction"`
	LessonContent          string `json:"lesson_content"`
	TeacherActivities      string `json:"teacher_activities"`
	StudentActivities      string `json:"student_activities"`
	Conclusion             string `json:"conclusion"`
	Evaluation             string `json:"evaluation"`
}

// LessonPlanResponse represents the response for a lesson plan
type LessonPlanResponse struct {
	ID                     string    `json:"id"`
	LessonID               string    `json:"lesson_id"`
	PreviousKnowledge      string    `json:"previous_knowledge"`
	BehaviouralObjectives  string    `json:"behavioural_objectives"`
	TeachingAids           string    `json:"teaching_aids"`
	Introduction           string    `json:"introduction"`
	LessonContent          string    `json:"lesson_content"`
	TeacherActivities      string    `json:"teacher_activities"`
	StudentActivities      string    `json:"student_activities"`
	Conclusion             string    `json:"conclusion"`
	Evaluation             string    `json:"evaluation"`
	CreatedBy              string    `json:"created_by"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`

	// Nested relationships
	Lesson  *LessonResponse `json:"lesson,omitempty"`
	Creator *UserResponse   `json:"creator,omitempty"`
}

// LessonPlanListResponse represents a paginated list of lesson plans
type LessonPlanListResponse struct {
	LessonPlans []LessonPlanResponse `json:"lesson_plans"`
	Total       int64                `json:"total"`
	Page        int                  `json:"page"`
	Limit       int                  `json:"limit"`
	TotalPages  int                  `json:"total_pages"`
}

// LessonPlanQueryParams represents query parameters for filtering lesson plans
type LessonPlanQueryParams struct {
	LessonID string `form:"lesson_id"`
	Search   string `form:"search"`
	Page     int    `form:"page" default:"1"`
	Limit    int    `form:"limit" default:"20"`
	SortBy   string `form:"sort_by" default:"created_at"`
	SortOrder string `form:"sort_order" default:"desc"`
}

// LessonPlanStats represents statistics for lesson plans
type LessonPlanStats struct {
	TotalLessonPlans     int64 `json:"total_lesson_plans"`
	WithObjectives       int64 `json:"with_objectives"`
	WithTeachingAids     int64 `json:"with_teaching_aids"`
	WithEvaluation       int64 `json:"with_evaluation"`
}