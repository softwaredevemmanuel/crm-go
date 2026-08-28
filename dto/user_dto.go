// dto/user_dto.go
package dto

import (
)



// UserListResponse represents paginated user list response
type UserListResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

// UserQueryParams represents query parameters for filtering users
type UserQueryParams struct {
	Search    string `form:"search"`
	Role      string `form:"role" binding:"omitempty,oneof=admin staff student parent user"`
	Position  string `form:"position" binding:"omitempty,oneof=admin staff student parent user teacher"`
	Status    string `form:"status" binding:"omitempty,oneof=active inactive"`
	IsActive  *bool  `form:"is_active"`
	IsVerified *bool `form:"is_verified"`
	Page      int    `form:"page" binding:"min=1"`
	Limit     int    `form:"limit" binding:"min=1,max=100"`
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=first_name last_name email role status created_at"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	FirstName  string `json:"first_name" binding:"omitempty,min=2,max=100"`
	LastName   string `json:"last_name" binding:"omitempty,min=2,max=100"`
	MiddleName string `json:"middle_name"`
	Phone      string `json:"phone"`
	Role       string `json:"role" binding:"omitempty,oneof=admin staff student parent user"`
	Position   string `json:"position" binding:"omitempty,oneof=admin staff student parent user teacher"`
	Status     string `json:"status" binding:"omitempty,oneof=active inactive"`
	IsActive   *bool  `json:"is_active"`
	IsVerified *bool  `json:"is_verified"`
	Location   string `json:"location"`
}