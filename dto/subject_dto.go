// dto/subject_dto.go
package dto

import (
	"time"
)

// CreateSubjectRequest represents the request body for creating a subject
type CreateSubjectRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=255"`
	Code        string `json:"code" binding:"required,min=2,max=50"`
	Description string `json:"description"`
	Department  string `json:"department"`
	Credits     int    `json:"credits" binding:"min=0,max=10"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// SubjectResponse represents the subject response
type SubjectResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Department  string    `json:"department"`
	Credits     int       `json:"credits"`
	Status      string    `json:"status"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SubjectListResponse represents paginated subject list response
type SubjectListResponse struct {
	Subjects []SubjectResponse `json:"subjects"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
	TotalPages    int               `json:"total_pages"`
}

// UpdateSubjectRequest represents the request body for updating a subject
type UpdateSubjectRequest struct {
	Name        string `json:"name" binding:"omitempty,min=3,max=255"`
	Code        string `json:"code" binding:"omitempty,min=2,max=50"`
	Description string `json:"description"`
	Department  string `json:"department"`
	Credits     int    `json:"credits" binding:"min=0,max=10"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
}


// SubjectQueryParams represents query parameters for filtering subjects
type SubjectQueryParams struct {
	Search     string `form:"search"`
	Department string `form:"department"`
	Status     string `form:"status"`
	Page       int    `form:"page" binding:"min=1"`
	Limit      int    `form:"limit" binding:"min=1,max=100"`
	SortBy     string `form:"sort_by" binding:"omitempty,oneof=name code department credits created_at updated_at"`
	SortOrder  string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}