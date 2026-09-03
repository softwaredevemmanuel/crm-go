package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"crm-go/models"
	"crm-go/config"

)

// LoginId handles login with email from header and password from body
// @Summary Login with ID
// @Description Authenticates user using email from header and password from body
// @Tags Authentication
// @Accept json
// @Produce json
// @Param email header string true "User Email"
// @Param request body LoginIdInput true "Login credentials"
// @Success 200 {object} LoginIdResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/login/id [post]
func LoginId(c *gin.Context) {
	// Get email from header
	email := c.GetHeader("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Missing header",
			Message: "Email header is required",
		})
		return
	}

	// Bind password from JSON body
	var input LoginIdInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid input",
			Message: err.Error(),
		})
		return
	}

	// Validate password is provided
	if input.Password == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid input",
			Message: "Password is required",
		})
		return
	}

	// Find user by email
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Invalid credentials",
			Message: "Invalid email or password",
		})
		return
	}

	// Compare login IDs (password)
	if user.LoginID != input.Password {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Invalid credentials",
			Message: "Invalid email or password",
		})
		return
	}

	// Return structured response
	response := LoginIdResponse{
		Message:    "Login successful",
	}

	c.JSON(http.StatusOK, response)
}