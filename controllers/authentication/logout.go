package controllers

import (
	"crm-go/config"
	"crm-go/models"
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
	
)

// Logout handles user logout
// @Summary User logout
// @Description Invalidates the user's session and clears authentication cookies
// @Tags Authentication
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.LogoutResponse "Logout successful"
// @Success 200 {object} models.AlreadyLoggedOutResponse "Session already invalidated"
// @Failure 400 {object} models.ErrorResponse "No token provided"
// @Failure 500 {object} models.ErrorResponse "Failed to logout"
// @Router /auth/logout [post]
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