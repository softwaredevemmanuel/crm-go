package controllers

import (
	"crm-go/config"
	"crm-go/models"
	"crm-go/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var cfg = config.LoadEnv()

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID.String(), string(user.Email), string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	
	// Create user session
	session := models.UserSession{
		SessionToken: token, // Or generate a separate session token if preferred
		UserID:       user.ID,
		UserAgent:    c.Request.UserAgent(),
		UserIP:       c.ClientIP(),
		DeviceType:   getDeviceType(c.Request.UserAgent()),
		DeviceOS:     getOS(c.Request.UserAgent()),
		Browser:      getBrowser(c.Request.UserAgent()),
		IsActive:     true,
		LoginType:    "password",
		IssuedAt:     time.Now(),
		ExpiresAt:    time.Now().Add(time.Duration(cfg.JWTExpire) * time.Hour), // Match your JWT expiration
		LastUsedAt:   time.Now(),
	}

	// Save session to database
	if err := config.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// Set session cookie (optional)
	c.SetCookie("session_token", token, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"session_id": session.ID,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.FirstName + " " + user.LastName,
			"email": user.Email,
			"role":  user.Role,
		},
		"session": gin.H{
			"expires_at": session.ExpiresAt,
			"device":     session.DeviceType,
			"browser":    session.Browser,
		},
	})
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