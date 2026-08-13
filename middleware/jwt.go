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
	"fmt"
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

		fmt.Println("========== AUTH USER ==========")
		// fmt.Printf("User ID: %v\n", claims["user_id"])
		// fmt.Printf("Email: %v\n", claims["email"])
		// fmt.Printf("Role: %v\n", claims["role"])
		// fmt.Printf("Issuer: %v\n", claims["issuer"])
		// fmt.Printf("Issued At: %v\n", claims["iat"])
		// fmt.Printf("Expires At: %v\n", claims["exp"])
		// fmt.Println("===============================")

		// ✅ NEW: Check if session is still active in database
		var session models.UserSession
		result := config.DB.Where("session_token = ? AND is_active = true AND expires_at > ?",
			tokenString, time.Now()).
			First(&session)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":   "Session expired or invalidated",
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
		// fmt.Println("========== SESSION ==========")
		// fmt.Printf("Session ID: %v\n", session.ID)
		// fmt.Printf("Session User ID: %v\n", session.UserID)
		// fmt.Printf("Session Active: %v\n", session.IsActive)
		// fmt.Printf("Expires At: %v\n", session.ExpiresAt)
		// fmt.Println("=============================")

		// ✅ NEW: Update last used timestamp (optional but recommended)
		config.DB.Model(&session).Update("last_used_at", time.Now())

		// Save claims into context
		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])
		c.Set("session_id", session.ID) // ✅ Add session ID to context

		c.Next()
	}
}

// verifyTokenWithEmail validates the JWT token and checks if the email matches
func VerifyTokenWithEmail(tokenString string, expectedEmail string) (bool, uint, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		// Use the same secret as your middleware
		return []byte("supersecretkey"), nil // Should come from env
	})

	if err != nil {
		return false, 0, err
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check if token is expired
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				return false, 0, jwt.ErrTokenExpired
			}
		}

		// Get email from claims (assuming it's stored as "email" or "localEmail")
		email, emailOk := claims["email"].(string)
		if !emailOk {
			// Try alternative claim name
			email, emailOk = claims["localEmail"].(string)
			if !emailOk {
				return false, 0, jwt.ErrInvalidType
			}
		}

		// Verify email matches
		if email != expectedEmail {
			return false, 0, errors.New("email mismatch")
		}

		// Get user_id from claims
		userIDFloat, userIDOk := claims["user_id"].(float64)
		if !userIDOk {
			// Try int format
			userIDInt, ok := claims["user_id"].(int)
			if !ok {
				return false, 0, jwt.ErrInvalidType
			}
			return true, uint(userIDInt), nil
		}
		return true, uint(userIDFloat), nil
	}

	return false, 0, jwt.ErrTokenInvalidClaims
}

// Optional: Helper function to generate tokens (for testing)
func GenerateTestToken(userID uint, email string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(48 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"issuer":  "crm-go",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("supersecretkey"))
}
