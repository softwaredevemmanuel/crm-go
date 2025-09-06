package routes

import (
	"crm-go/controllers/courses"
	"crm-go/middleware"
	"github.com/gin-gonic/gin"
)

func CourseProductRoutes(r *gin.Engine) {
	course_product_tables := r.Group("/course-products")
	{
		course_product_tables.GET("/:id/courses", controllers.GetCoursesByProduct)

		// Protected routes
		protected := r.Group("/api")
		protected.Use(middleware.AuthMiddleware())
		protected.POST("/course-products", middleware.RoleMiddleware("admin"), controllers.CreateCourseProduct)
		protected.DELETE("/course-products/:id", middleware.RoleMiddleware("admin"), controllers.DeleteCourseProduct)

	}
}
