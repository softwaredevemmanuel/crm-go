package controllers

import (
	"crm-go/config"
	"crm-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type SignUpInput struct {
	FirstName     string       `json:"first_name" binding:"required"`
	LastName     string       `json:"last_name" binding:"required"`
	Email    string       `json:"email" binding:"required,email"`
	Password string       `json:"password" binding:"required,min=6"`
	Role     models.Role  `json:"role" binding:"required"`
}

func SignUp(c *gin.Context) {
	var input SignUpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)

	user := models.User{
		FirstName:     input.FirstName,
		LastName:     input.LastName,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": user})
}

// controllers/auth.go
func Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
		return
	}

	// Invalidate session
	result := config.DB.Model(&models.UserSession{}).
		Where("session_token = ? AND is_active = true", token).
		Updates(map[string]interface{}{
			"is_active":    false,
			"logged_out_at": time.Now(),
		})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	 if result.RowsAffected == 0 {
        // Token was already invalid or doesn't exist
        c.JSON(http.StatusOK, gin.H{
            "message": "Session already invalidated",
            "warning": "No active session found for this token",
        })
        return
    }
	   // Clear all auth cookies
    c.SetCookie("session_token", "", -1, "/", "", false, true)
    c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}
