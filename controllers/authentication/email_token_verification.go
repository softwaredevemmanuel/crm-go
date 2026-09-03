package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"crm-go/dto"

)

// Verify Email
// @Summary Verify email
// @Description Verify user's email using token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.VerifyEmailRequest true "Verification token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/auth/verify-email [post]
func (c *AuthController) VerifyEmail(ctx *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	response, err := c.authService.VerifyEmail(&req)
	if err != nil {
		if err.Error() == "invalid or expired verification token" ||
			err.Error() == "verification token has expired" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
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
		"message": "Email verified successfully",
		"data":    response,
	})
}