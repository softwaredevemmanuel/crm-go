// middleware/session.go
package middleware

import (
	"crm-go/config"
	"crm-go/models"
	"time"

	"github.com/gin-gonic/gin"
)

func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip session tracking for auth endpoints
		if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/register" {
			c.Next()
			return
		}

		token := c.GetHeader("Authorization")
		if token == "" {
			token, _ = c.Cookie("session_token")
		}

		if token != "" {
			// Remove "Bearer " prefix if present
			if len(token) > 7 && token[:7] == "Bearer " {
				token = token[7:]
			}

			// Update session last used time
			var session models.UserSession
			if err := config.DB.Where("session_token = ? AND is_active = true AND expires_at > ?", 
				token, time.Now()).First(&session).Error; err == nil {
				
				// Update last used time
				config.DB.Model(&session).Update("last_used_at", time.Now())
			}
		}

		c.Next()
	}
}