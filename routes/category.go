package routes

import (
	"crm-go/controllers/category"
	"crm-go/middleware"
	"github.com/gin-gonic/gin"
)

func CategoryRoutes(r *gin.Engine) {
	categories := r.Group("/categories")
	{
		categories.GET("/", controllers.GetCategories)

		// Protected routes
		protected := r.Group("/api")
		protected.Use(middleware.AuthMiddleware())
		protected.POST("/categories", middleware.RoleMiddleware("admin"), controllers.CreateCategory)
		protected.PUT("/categories/:id", middleware.RoleMiddleware("admin"), controllers.UpdateCategory)
		protected.DELETE("/categories/:id", middleware.RoleMiddleware("admin"), controllers.DeleteCategory)

	}
}
