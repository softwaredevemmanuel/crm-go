package routes

import (
	"crm-go/controllers/category"
	"crm-go/middleware"
	"github.com/gin-gonic/gin"
)

func CourseProductRoutes(r *gin.Engine) {
	course_products := r.Group("/course-products")
	{
		course_products.GET("/:id/courses", controllers.GetCoursesByProduct)

		// Protected routes
		protected := r.Group("/api")
		protected.Use(middleware.AuthMiddleware())
		protected.POST("/course-products", middleware.RoleMiddleware("admin"), controllers.CreateCourseProduct)
		protected.DELETE("/course-products/:id", middleware.RoleMiddleware("admin"), controllers.DeleteCourseProduct)

	}
}
