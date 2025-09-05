package routes

import (
	controllers "crm-go/controllers/products" // alias for courses
	"crm-go/middleware"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(r *gin.Engine) {
	products := r.Group("/products")
	{
		products.GET("/", controllers.GetProducts)
		products.GET("/:id", controllers.GetProductByID)
		products.GET("/:id/products", controllers.GetProductWithCourseMates)

		// Protected routes
		protected := r.Group("/api")
		protected.Use(middleware.AuthMiddleware())
		protected.POST("/products", middleware.RoleMiddleware("admin"), controllers.CreateProduct)
		protected.PUT("/products/:id", middleware.RoleMiddleware("admin"), controllers.UpdateProduct)
		protected.DELETE("/products/:id", middleware.RoleMiddleware("admin"), controllers.DeleteProduct)

	}
}
