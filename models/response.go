package models

import (
	"time"
	
	"github.com/google/uuid"
)

// LoginInput represents login request
// @Description Login request payload
type LoginInput struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// LoginResponse represents login response
// @Description Login response with token and session information
type LoginResponse struct {
	Message   string     `json:"message" example:"Login successful"`
	Token     string     `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	SessionID uuid.UUID  `json:"session_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	User      UserInfo   `json:"user"`
	Session   SessionInfo `json:"session"`
}

// UserInfo represents user information in login response
// @Description User information
type UserInfo struct {
	ID    uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name  string    `json:"name" example:"John Doe"`
	Email string    `json:"email" example:"user@example.com"`
	Role  string    `json:"role" example:"student"`
}

// SessionInfo represents session information in login response
// @Description Session information
type SessionInfo struct {
	ExpiresAt time.Time `json:"expires_at" example:"2024-01-20T15:04:05Z"`
	Device    string    `json:"device" example:"desktop"`
	Browser   string    `json:"browser" example:"Chrome"`
	IPAddress string    `json:"ip_address" example:"192.168.1.1"`
}

// ErrorResponse represents an error response
// @Description Error response structure
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid request"`
	Message string `json:"message,omitempty" example:"Please check your input"`
}

// LogoutResponse represents successful logout response
// @Description Logout successful response
type LogoutResponse struct {
	Message string `json:"message" example:"Logout successful"`
	Status  string `json:"status" example:"logged_out"`
	Details string `json:"details" example:"Session terminated and cookies cleared"`
}

// AlreadyLoggedOutResponse represents response when session was already invalidated
// @Description Already logged out response
type AlreadyLoggedOutResponse struct {
	Message string `json:"message" example:"Session already invalidated"`
	Warning string `json:"warning" example:"No active session found for this token"`
	Status  string `json:"status" example:"already_logged_out"`
}


// CategoryResponse represents the JSON returned after creating a category
type CategoryResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}