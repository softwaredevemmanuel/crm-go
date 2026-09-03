package controllers

import (
	"crm-go/services"
	"crm-go/config"

)
var cfg = config.LoadEnv()


type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}
type LoginIdInput struct {
	Password string `json:"password" binding:"required"`
}
type LoginIdResponse struct {
	Message   string     `json:"message" example:"Login successful"`
}

// SignUpInput represents the signup request body
type SignUpInput struct {
	Email      string `json:"email" binding:"required,email"`
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	MiddleName string `json:"middle_name"`
	Phone      string `json:"phone" binding:"required"`
	DOB        string `json:"dob"`
	Password   string `json:"password" binding:"required,min=8"`
	Position   string `json:"position"`
	Role       string `json:"role" binding:"required"`
}



// Helper functions to parse user agent
func getDeviceType(userAgent string) string {

	// Simple device detection - you can use a proper library like github.com/mssola/user_agent
	if contains(userAgent, "Mobile") || contains(userAgent, "Android") || contains(userAgent, "iPhone") {
		return "mobile"
	} else if contains(userAgent, "Tablet") || contains(userAgent, "iPad") {
		return "tablet"
	}
	return "desktop"
}

func getOS(userAgent string) string {
	switch {
	case contains(userAgent, "Windows"):
		return "Windows"
	case contains(userAgent, "Macintosh") || contains(userAgent, "Mac OS"):
		return "macOS"
	case contains(userAgent, "Linux"):
		return "Linux"
	case contains(userAgent, "Android"):
		return "Android"
	case contains(userAgent, "iPhone") || contains(userAgent, "iPad"):
		return "iOS"
	default:
		return "Unknown"
	}
}

func getBrowser(userAgent string) string {
	switch {
	case contains(userAgent, "Chrome"):
		return "Chrome"
	case contains(userAgent, "Firefox"):
		return "Firefox"
	case contains(userAgent, "Safari"):
		return "Safari"
	case contains(userAgent, "Edge"):
		return "Edge"
	case contains(userAgent, "Opera"):
		return "Opera"
	default:
		return "Unknown"
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}