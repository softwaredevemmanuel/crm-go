package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"crm-go/config"
	"crm-go/models"
	"errors"
	"time"
)

var jwtSecret = []byte("supersecretkey") // 🔥 should come from ENV in production

// AuthMiddleware verifies the JWT token AND checks session activity
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Expected format: Bearer <token>
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// ✅ NEW: Check if session is still active in database
		var session models.UserSession
		result := config.DB.Where("session_token = ? AND is_active = true AND expires_at > ?", 
			tokenString, time.Now()).
			First(&session)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Session expired or invalidated",
					"message": "Please login again",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to verify session status",
				})
			}
			c.Abort()
			return
		}

		// ✅ NEW: Update last used timestamp (optional but recommended)
		config.DB.Model(&session).Update("last_used_at", time.Now())

		// Save claims into context
		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])
		c.Set("session_id", session.ID) // ✅ Add session ID to context

		c.Next()
	}
}