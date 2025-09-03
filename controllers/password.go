package controllers

import (
	"crm-go/config"
	"crm-go/models"
	"net/http"
	"time"

	"crm-go/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

func ForgotPassword(c *gin.Context) {
	var input ForgotPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "If that email exists, a reset link has been sent"}) 
		return // security: don't reveal whether email exists
	}

	// Create reset token
	token := uuid.New().String()
	reset := models.PasswordReset{
		UserID:    user.ID.String(),
		Token:     token,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	config.DB.Create(&reset)

	// TODO: send email with link
	resetLink := "http://localhost:8080/reset-password?token=" + token

	// Send email via SMTP
	emailBody := "<p>Hello,</p><p>Click the link below to reset your password:</p>" +
	"<a href='" + resetLink + "'>Reset Password</a>"

	if err := utils.SendEmail(user.Email, "Password Reset - Go CRM", emailBody); err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset email"})
	return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "If that email exists, a reset link has been sent",
		"link":    resetLink, // just for testing; in prod, send via email
	})
}
