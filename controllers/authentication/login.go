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



// Login handles user login
// @Summary User login
// @Description Authenticate user and return JWT token with session information
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body models.LoginInput true "Login credentials"
// @Success 200 {object} models.LoginResponse "Login successful"
// @Failure 400 {object} models.ErrorResponse "Invalid input"
// @Failure 401 {object} models.ErrorResponse "Invalid credentials"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var input models.LoginInput
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid input",
			Message: err.Error(),
		})
		return
	}

	// Find user by email
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Invalid credentials",
			Message: "Invalid email or password",
		})
		return
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Invalid credentials",
			Message: "Invalid email or password",
		})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID.String(), user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal server error",
			Message: "Failed to generate token",
		})
		return
	}

	// Create user session
	session := models.UserSession{
		SessionToken: token,
		UserID:       user.ID,
		UserAgent:    c.Request.UserAgent(),
		UserIP:       c.ClientIP(),
		DeviceType:   getDeviceType(c.Request.UserAgent()),
		DeviceOS:     getOS(c.Request.UserAgent()),
		Browser:      getBrowser(c.Request.UserAgent()),
		IsActive:     true,
		LoginType:    "password",
		IssuedAt:     time.Now(),
		ExpiresAt:    time.Now().Add(time.Duration(cfg.JWTExpire) * time.Hour),
		LastUsedAt:   time.Now(),
	}

	// Save session to database
	if err := config.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal server error",
			Message: "Failed to create session",
		})
		return
	}

	// Set session cookie
	c.SetCookie("session_token", token, 3600, "/", "", false, true)

	// Return structured response
	response := models.LoginResponse{
		Message:    "Login successful",
		Token:      token,
		SessionID:  session.ID,
		User: models.UserInfo{
			ID:    user.ID,
			Name:  user.FirstName + " " + user.LastName,
			Email: user.Email,
			Role:  string(user.Role),
		},
		Session: models.SessionInfo{
			ExpiresAt: session.ExpiresAt,
			Device:    session.DeviceType,
			Browser:   session.Browser,
			IPAddress: session.UserIP,
		},
	}

	c.JSON(http.StatusOK, response)
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