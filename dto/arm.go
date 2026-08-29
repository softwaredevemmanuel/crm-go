package dto

import (
	"time"
)

// CreateArmRequest represents the request body for creating an arm
type CreateArmRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=255"`
	Description string `json:"description"`
	Code        string `json:"code" binding:"omitempty,min=1,max=5"`
	GradeID     string `json:"grade_id" binding:"required"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive archived"`
	Capacity    int    `json:"capacity" binding:"min=1,max=100"`
}

// UpdateArmRequest represents the request body for updating an arm
type UpdateArmRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=255"`
	Description string `json:"description"`
	Code        string `json:"code" binding:"omitempty,min=1,max=5"`
	GradeID     string `json:"grade_id"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive archived"`
	Capacity    int    `json:"capacity" binding:"min=1,max=100"`
}

// ArmResponse represents the arm response
type ArmResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Code        string    `json:"code"`
	GradeID     string    `json:"grade_id"`
	Grade       *ClassGradeResponse `json:"grade,omitempty"`
	Status      string    `json:"status"`
	Capacity    int       `json:"capacity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ArmListResponse represents paginated arm list response
type ArmListResponse struct {
	Arms       []ArmResponse `json:"arms"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"total_pages"`
}

// ArmQueryParams represents query parameters for filtering arms
type ArmQueryParams struct {
	Search     string `form:"search"`
	GradeID    string `form:"grade_id"`
	Status     string `form:"status"`
	Page       int    `form:"page" binding:"min=1"`
	Limit      int    `form:"limit" binding:"min=1,max=100"`
	SortBy     string `form:"sort_by" binding:"omitempty,oneof=name code grade_id capacity status created_at"`
	SortOrder  string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}