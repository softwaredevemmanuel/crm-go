package controllers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"crm-go/config"
	"crm-go/models"
	"crm-go/utils"

	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GoogleLoginHandler initiates Google OAuth2 login
// @Summary Initiate Google OAuth2 login
// @Description Redirects to Google OAuth2 consent screen with optional role parameter
// @Tags Authentication
// @Produce json
// @Param role query string false "User role (student/tutor)" Enums(student, tutor) default(student)
// @Success 302 {string} string "Redirect to Google OAuth2"
// @Router /auth/google/login [get]
func GoogleLoginHandler(c *gin.Context) {
    role := c.Query("role")
    if role == "" {
        role = "student"
    }
    url := config.GoogleOauthConfig.AuthCodeURL(role)
    c.Redirect(http.StatusFound, url)
}

// GoogleCallbackHandler handles Google OAuth2 callback
// @Summary Handle Google OAuth2 callback
// @Description Processes Google OAuth2 callback, creates/updates user, and returns JWT token
// @Tags Authentication
// @Produce json
// @Param code query string true "OAuth2 authorization code from Google"
// @Param state query string true "State parameter containing user role"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/google/callback [get]
func GoogleCallbackHandler(c *gin.Context) {
	
	// Get code from query
	code := c.Query("code")
	role := c.Query("state") // state carries "student" or "tutor"

	log.Printf("✅ Value of role: %v", role)

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code in request"})
		return
	}

	// 1. Exchange code for token
	token, err := config.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("❌ Token exchange failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// 2. Get user info from Google
	client := config.GoogleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Printf("❌ Failed to get user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("❌ Google API error: %s", string(body))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Google API returned non-200"})
		return
	}

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Printf("❌ Failed to parse user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	log.Printf("✅ Google user info: %+v", userInfo)

	// Safely extract fields
	email, _ := userInfo["email"].(string)
	first_name, _ := userInfo["name"].(string)
	last_name, _ := userInfo["family_name"].(string)
	picture, _ := userInfo["picture"].(string)

	log.Printf("✅ Google user info: 1")

	// Get DB instance
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not available"})
		return
	}

	log.Printf("✅ Google user info: 2")
	// Check if user exists in DB
	var user models.User
	var userID int

	result := db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// User not found → create
			user = models.User{
				FirstName: first_name,
				LastName:  last_name,
				Role:      models.Role(role),
				Email:     email,
				Picture:   picture,
				Provider:  "google",
			}
			if err := db.Create(&user).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}
		} else {
			// Some other DB error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	}

	log.Printf("✅ Google user info: 3")

	// Generate JWT
	tokenString, err := utils.GenerateJWT(email, first_name, string(user.Role)) // default role as student
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}
	log.Printf("✅ Google user info: 4")

	// Return user info + token
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":         userID,
			"first_name": first_name,
			"last_name":  last_name,
			"role":       role, // default role
			"email":      email,
			"picture":    picture,
		},
	})
}
