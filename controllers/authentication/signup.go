// controllers/auth_controller.go
package controllers

import (
	"crm-go/config"
	"crm-go/models"
	"crm-go/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

// SignUp handles user registration with email verification
// @Summary Register a new user
// @Description Create a new user account with first name, last name, email, password, and role
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param   input body SignUpInput true "User signup credentials"
// @Success 201 {object} map[string]interface{} "User created successfully. Verification email sent."
// @Failure 400 {object} map[string]interface{} "Invalid input or email already exists"
// @Failure 500 {object} map[string]interface{} "Failed to create user or send email"
// @Router /auth/signup [post]
func SignUp(c *gin.Context) {
	var input SignUpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User with this email already exists",
		})
		return
	}

	// Validate password strength
	if len(input.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must be at least 8 characters long",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create user
	user := models.User{
		ID:         uuid.New(),
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		MiddleName: input.MiddleName,
		Email:      input.Email,
		Phone:      input.Phone,
		Password:   string(hashedPassword),
		Role:       input.Role,
		Position:   input.Position,
		IsActive:   true,
		IsVerified: false, // User starts as unverified
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Save user to database
	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	// Generate verification token
	token := uuid.New().String()

	// Delete any existing verification tokens for this user (just in case)
	config.DB.Where("user_id = ?", user.ID).Delete(&models.EmailVerification{})

	// Create verification record
	verification := models.EmailVerification{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		Email:     user.Email,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours expiry
		Used:      false,
	}

	if err := config.DB.Create(&verification).Error; err != nil {
		// Log error but don't fail the signup
		// We can still return success but inform user about email issue
		c.JSON(http.StatusCreated, gin.H{
			"message": "User created successfully, but we couldn't send the verification email. Please contact support.",
			"user": gin.H{
				"id":          user.ID,
				"first_name":  user.FirstName,
				"last_name":   user.LastName,
				"email":       user.Email,
				"role":        user.Role,
				"is_active":   user.IsActive,
				"is_verified": user.IsVerified,
			},
			"warning": "Verification email could not be sent",
		})
		return
	}

	// Get frontend URL
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	// Build verification link
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", frontendURL, token)

	// Send verification email
	emailService := services.NewEmailService()
	err = emailService.SendVerificationEmail(
		user.Email,
		user.FirstName,
		user.LastName,
		token,
		frontendURL,
	)

	if err != nil {
		// Log the error but still return success with a warning
		c.JSON(http.StatusCreated, gin.H{
			"message": "User created successfully, but we couldn't send the verification email. Please check your email settings or contact support.",
			"user": gin.H{
				"id":          user.ID,
				"first_name":  user.FirstName,
				"last_name":   user.LastName,
				"email":       user.Email,
				"role":        user.Role,
				"is_active":   user.IsActive,
				"is_verified": user.IsVerified,
			},
			"warning": "Verification email could not be sent",
		})
		return
	}

	// Success response with verification email sent
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully. A verification email has been sent to your email address. Please verify your email to continue.",
		"user": gin.H{
			"id":          user.ID,
			"first_name":  user.FirstName,
			"last_name":   user.LastName,
			"email":       user.Email,
			"role":        user.Role,
			"is_active":   user.IsActive,
			"is_verified": user.IsVerified,
		},
		"verification": gin.H{
			"token":            token, // For testing only
			"expires_in":       "24 hours",
			"verification_url": verificationURL, // For testing only
		},
	})
}
