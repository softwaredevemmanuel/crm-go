package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"crm-go/dto"

)

// Send Password Reset Email
// @Summary Send password reset email
// @Description Send a password reset link to the user's email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.SendPasswordResetRequest true "Email"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/auth/send-password-reset [post]
func (c *AuthController) SendPasswordReset(ctx *gin.Context) {
	var req dto.SendPasswordResetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	if err := c.authService.SendPasswordReset(&req); err != nil {
		if err.Error() == "account is deactivated" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Password reset email sent successfully",
	})
}