// dto/department_dto.go
package dto

import (
	"time"
)

// ============ REQUEST DTOs ============

type CreateDepartmentRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Code        string `json:"code" binding:"required,min=2,max=20"`
	Description string `json:"description"`
	HeadID      string `json:"head_id"`
	Status      string `json:"status"`
}

type UpdateDepartmentRequest struct {
	Name        string `json:"name,omitempty"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	HeadID      string `json:"head_id,omitempty"`
	Status      string `json:"status,omitempty"`
}

type DepartmentQueryParams struct {
	Page      int    `form:"page" binding:"min=1"`
	Limit     int    `form:"limit" binding:"min=1,max=100"`
	Search    string `form:"search"`
	Status    string `form:"status"`
	SortBy    string `form:"sort_by"`
	SortOrder string `form:"sort_order"`
}

// ============ RESPONSE DTOs ============

type DepartmentResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	HeadID      *string   `json:"head_id,omitempty"`
	Status      string    `json:"status"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DepartmentWithHeadResponse struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Code        string      `json:"code"`
	Description string      `json:"description"`
	HeadID      *string     `json:"head_id,omitempty"`
	Head        *UserBrief  `json:"head,omitempty"`
	Status      string      `json:"status"`
	CreatedBy   string      `json:"created_by"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type DepartmentWithSubjectsResponse struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Code        string                   `json:"code"`
	Description string                   `json:"description"`
	HeadID      *string                  `json:"head_id,omitempty"`
	Head        *UserBrief               `json:"head,omitempty"`
	Subjects    []SubjectBriefResponse   `json:"subjects,omitempty"`
	SubjectCount int                     `json:"subject_count"`
	Status      string                   `json:"status"`
	CreatedBy   string                   `json:"created_by"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

type SubjectBriefResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Credits     int    `json:"credits"`
	Status      string `json:"status"`
}

type UserBrief struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type DepartmentListResponse struct {
	Departments []DepartmentWithSubjectsResponse `json:"departments"`
	Total       int64                           `json:"total"`
	Page        int                             `json:"page"`
	Limit       int                             `json:"limit"`
	TotalPages  int                             `json:"total_pages"`
}