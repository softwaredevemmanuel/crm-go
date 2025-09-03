package routes

import (
	"crm-go/controllers"
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
	protected.POST("/courses", middleware.RoleMiddleware("admin"), controllers.CreateCourse)
	protected.PUT("/courses/:id", middleware.RoleMiddleware("admin"), controllers.UpdateCourse)
	protected.DELETE("/courses/:id", middleware.RoleMiddleware("admin"),controllers.DeleteCourse)
	}
}
