// dto/subject_dto.go
package dto

import (
	"time"
)

// ============ REQUEST DTOs ============

type CreateSubjectRequest struct {
	Name         string `json:"name" binding:"required,min=3,max=255"`
	Code         string `json:"code" binding:"required,min=2,max=50"`
	Description  string `json:"description"`
	DepartmentID string `json:"department_id" binding:"required,uuid"`
	Credits      int    `json:"credits" binding:"min=0,max=10"`
	Status       string `json:"status"`
}

type UpdateSubjectRequest struct {
	Name         string `json:"name,omitempty"`
	Code         string `json:"code,omitempty"`
	Description  string `json:"description,omitempty"`
	DepartmentID string `json:"department_id,omitempty"`
	Credits      int    `json:"credits,omitempty"`
	Status       string `json:"status,omitempty"`
}

// dto/subject_dto.go
type SubjectQueryParams struct {
	Search       string `form:"search"`
	DepartmentID string `form:"department_id"` // Keep for UUID
	Department   string `form:"department"`    // ✅ Add for name filtering
	Status       string `form:"status"`
	Page         int    `form:"page" binding:"min=1"`
	Limit        int    `form:"limit" binding:"min=1,max=100"`
	SortBy       string `form:"sort_by" binding:"omitempty,oneof=name code department credits status created_at updated_at"`
	SortOrder    string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// ============ RESPONSE DTOs ============

type DepartmentBriefResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type DepartmentWithHeadBriefResponse struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Code        string      `json:"code"`
	Description string      `json:"description"`
	HeadID      *string     `json:"head_id,omitempty"`
	Head        *UserBrief  `json:"head,omitempty"`
}

type SubjectResponse struct {
	ID           string                  `json:"id"`
	Name         string                  `json:"name"`
	Code         string                  `json:"code"`
	Description  string                  `json:"description"`
	DepartmentID string                  `json:"department_id"`
	Department   *DepartmentBriefResponse `json:"department,omitempty"`
	Credits      int                     `json:"credits"`
	Status       string                  `json:"status"`
	CreatedBy    string                  `json:"created_by"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
}

type SubjectWithDepartmentResponse struct {
	ID           string                         `json:"id"`
	Name         string                         `json:"name"`
	Code         string                         `json:"code"`
	Description  string                         `json:"description"`
	DepartmentID string                         `json:"department_id"`
	Department   *DepartmentWithHeadBriefResponse `json:"department,omitempty"`
	Credits      int                            `json:"credits"`
	Status       string                         `json:"status"`
	CreatedBy    string                         `json:"created_by"`
	CreatedAt    time.Time                      `json:"created_at"`
	UpdatedAt    time.Time                      `json:"updated_at"`
}

type SubjectListResponse struct {
	Subjects   []SubjectWithDepartmentResponse `json:"subjects"`
	Total      int64                          `json:"total"`
	Page       int                            `json:"page"`
	Limit      int                            `json:"limit"`
	TotalPages int                            `json:"total_pages"`
}

