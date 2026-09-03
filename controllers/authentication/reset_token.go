
package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"crm-go/dto"
)


// Verify Reset Token
// @Summary Verify password reset token
// @Description Check if a password reset token is valid
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.VerifyResetTokenRequest true "Reset token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/auth/verify-reset-token [post]
func (c *AuthController) VerifyResetToken(ctx *gin.Context) {
	var req dto.VerifyResetTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	response, err := c.authService.VerifyResetToken(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}