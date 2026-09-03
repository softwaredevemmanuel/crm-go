// dto/auth_dto.go
package dto

import "time"

// Email Verification DTOs
type SendVerificationEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type VerifyEmailResponse struct {
	UserID         string    `json:"user_id"`
	Email          string    `json:"email"`
	IsVerified     bool      `json:"is_verified"`
	VerifiedAt     time.Time `json:"verified_at"`
}

// Password Reset DTOs
type SendPasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyResetTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

type VerifyResetTokenResponse struct {
	Valid   bool   `json:"valid"`
	Email   string `json:"email"`
	UserID  string `json:"user_id"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// Response DTOs
type AuthResponse struct {
	Message string `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}