package controllers

import (
	"crm-go/config"
	"crm-go/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)


type ResetPasswordInput struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ResetPassword handles password reset confirmation
// @Summary Reset user password
// @Description Confirms password reset using a valid token and sets new password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body ResetPasswordInput true "Password reset data"
// @Success 200 {object} object{message=string} "Password reset successful"
// @Failure 400 {object} object{error=string} "Invalid token, token expired, or invalid input"
// @Router /auth/reset-password [post]
func ResetPassword(c *gin.Context) {
	var input ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var reset models.PasswordReset
	if err := config.DB.Where("token = ?", input.Token).First(&reset).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	// Check expiry
	if time.Now().After(reset.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token expired"})
		return
	}

	// Get user
	var user models.User
	config.DB.First(&user, "id = ?", reset.UserID)

	// Update password
	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), 14)
	user.Password = string(hashed)
	config.DB.Save(&user)

	// Delete used token
	config.DB.Delete(&reset)

	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully"})
}
