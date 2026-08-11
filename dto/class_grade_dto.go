// dto/class_grade_dto.go
package dto

import (
	"time"
)

// CreateClassGradeRequest represents the request body for creating a class grade
type CreateClassGradeRequest struct {
	Name         string `json:"name" binding:"required,min=2,max=255"`
	Code         string `json:"code" binding:"required,min=2,max=50"`
	Level        int    `json:"level" binding:"required,min=1,max=6"`
	Description  string `json:"description"`
	AcademicYear string `json:"academic_year" binding:"required"`
	Capacity     int    `json:"capacity" binding:"min=1,max=100"`
	Status       string `json:"status" binding:"omitempty,oneof=active inactive archived"`
}

// UpdateClassGradeRequest represents the request body for updating a class grade
type UpdateClassGradeRequest struct {
	Name         string `json:"name" binding:"omitempty,min=2,max=255"`
	Code         string `json:"code" binding:"omitempty,min=2,max=50"`
	Level        int    `json:"level" binding:"omitempty,min=1,max=6"`
	Description  string `json:"description"`
	AcademicYear string `json:"academic_year"`
	Capacity     int    `json:"capacity" binding:"omitempty,min=1,max=100"`
	Status       string `json:"status" binding:"omitempty,oneof=active inactive archived"`
}

// ClassGradeResponse represents the class grade response
type ClassGradeResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	Level        int       `json:"level"`
	Description  string    `json:"description"`
	AcademicYear string    `json:"academic_year"`
	Capacity     int       `json:"capacity"`
	Status       string    `json:"status"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ClassGradeListResponse represents paginated class grade list response
type ClassGradeListResponse struct {
	ClassGrades []ClassGradeResponse `json:"class_grades"`
	Total       int64                `json:"total"`
	Page        int                  `json:"page"`
	Limit       int                  `json:"limit"`
	TotalPages  int                  `json:"total_pages"`
}

// ClassGradeQueryParams represents query parameters for filtering class grades
type ClassGradeQueryParams struct {
	Search       string `form:"search"`
	Level        int    `form:"level"`
	AcademicYear string `form:"academic_year"`
	Status       string `form:"status"`
	Page         int    `form:"page" binding:"min=1"`
	Limit        int    `form:"limit" binding:"min=1,max=100"`
	SortBy       string `form:"sort_by" binding:"omitempty,oneof=name code level academic_year capacity status created_at"`
	SortOrder    string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}