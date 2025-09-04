package routes

import (
	"github.com/gin-gonic/gin"
	"crm-go/middleware"
)

func AdminRoutes(r *gin.Engine) {
	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
	// Course routes (admin only ideally)
	protected.GET("/admin", middleware.RoleMiddleware("admin"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome Admin!"})
	})

	}
}
