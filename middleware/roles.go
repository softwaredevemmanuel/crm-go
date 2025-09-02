package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware checks if user has the required role
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role not found"})
			c.Abort()
			return
		}

		userRole := role.(string)

		// Check if user's role is in allowedRoles
		for _, r := range allowedRoles {
			if userRole == r {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		c.Abort()
	}
}
