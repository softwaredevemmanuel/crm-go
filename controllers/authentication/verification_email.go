package controllers

import (
	"crm-go/dto"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Send Verification Email
// @Summary Send email verification
// @Description Send a verification email to the user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.SendVerificationEmailRequest true "Email"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/auth/send-verification [post]
func (c *AuthController) SendVerificationEmail(ctx *gin.Context) {
	var req dto.SendVerificationEmailRequest
	log.Fatalln("drop here")
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	if err := c.authService.SendVerificationEmail(&req); err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		if err.Error() == "email already verified" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Email already verified",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Verification email sent successfully",
	})
}