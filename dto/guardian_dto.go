// dto/guardian_dto.go
package dto

import (
	"time"
)

// CreateGuardianRequest represents the request body for creating a guardian
type CreateGuardianRequest struct {
	FirstName    string `json:"first_name" binding:"required,min=2,max=100"`
	LastName     string `json:"last_name" binding:"required,min=2,max=100"`
	MiddleName   string `json:"middle_name"`
	Email        string `json:"email" binding:"required,email"`
	Phone        string `json:"phone" binding:"required"`
	Occupation   string `json:"occupation"`
	Relationship string `json:"relationship" binding:"required,oneof=Father Mother Guardian Uncle Aunt Grandfather Grandmother Brother Sister Other"`
	Address      string `json:"address"`
	StudentID    string `json:"student_id" binding:"required"`
	IsPrimary    bool   `json:"is_primary"`
	IsEmergency  bool   `json:"is_emergency"`
	Status       string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// UpdateGuardianRequest represents the request body for updating a guardian
type UpdateGuardianRequest struct {
	
	Occupation   string `json:"occupation"`
	Relationship string `json:"relationship" binding:"omitempty,oneof=Father Mother Guardian Uncle Aunt Grandfather Grandmother Brother Sister Other"`
	Address      string `json:"address"`
	StudentID    string `json:"student_id"`
	IsPrimary    *bool  `json:"is_primary"`
	IsEmergency  *bool  `json:"is_emergency"`
	Status       string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// GuardianResponse represents the guardian response with linked user details
type GuardianResponse struct {
	ID           string    `json:"id"`
	Occupation   string    `json:"occupation"`
	Relationship string    `json:"relationship"`
	Address      string    `json:"address"`
	StudentID    string    `json:"student_id"`
	UserID       string    `json:"user_id"`
	IsPrimary    bool      `json:"is_primary"`
	IsEmergency  bool      `json:"is_emergency"`
	Status       string    `json:"status"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	// Nested objects for complete details
	Student *UserResponse `json:"student,omitempty"`
	User    *UserResponse `json:"user,omitempty"`
}

// UserResponse represents the user response
type UserResponse struct {
	ID         string    `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	MiddleName string    `json:"middle_name"`
	FullName   string    `json:"full_name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Role       string    `json:"role"`
	Picture    string    `json:"picture"`
	IsVerified bool      `json:"is_verified"`
	IsActive   bool      `json:"is_active"`
	Location   string    `json:"location"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// GuardianListResponse represents paginated guardian list response
type GuardianListResponse struct {
	Guardians  []GuardianResponse `json:"guardians"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}

// GuardianQueryParams represents query parameters for filtering guardians
type GuardianQueryParams struct {
	Search       string `form:"search"`
	StudentID    string `form:"student_id"`
	Relationship string `form:"relationship"`
	Status       string `form:"status"`
	IsPrimary    *bool  `form:"is_primary"`
	IsEmergency  *bool  `form:"is_emergency"`
	Page         int    `form:"page" binding:"min=1"`
	Limit        int    `form:"limit" binding:"min=1,max=100"`
	SortBy       string `form:"sort_by" binding:"omitempty,oneof=relationship status created_at"`
	SortOrder    string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}