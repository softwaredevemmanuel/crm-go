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



// ForgotPassword handles password reset requests
// @Summary Request password reset
// @Description Initiates a password reset process by sending a reset link to the provided email (if exists)
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body ForgotPasswordInput true "Email address for password reset"
// @Success 200 {object} map[string]interface{} "Reset email sent (always returns success for security)"
// @Success 200 {object} object{message=string,link=string} "Success response with reset link (for testing)"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 500 {object} map[string]string "Failed to send email"
// @Router /auth/forgot-password [post]
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

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}