package routes

import (
	courseController "crm-go/controllers/courses"   // alias for courses
	"crm-go/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterCourseRoutes(r *gin.Engine) {
	courses := r.Group("/courses")
	{
		courses.GET("/", courseController.GetCourses)
		courses.GET("/:id", courseController.GetCourseByID)
		courses.GET("/:id/products", courseController.GetCourseWithProducts)

		// Protected routes
		protected := r.Group("/api")
		protected.Use(middleware.AuthMiddleware())
		protected.POST("/courses", middleware.RoleMiddleware("admin"), courseController.CreateCourse)
		protected.PUT("/courses/:id", middleware.RoleMiddleware("admin"), courseController.UpdateCourse)
		protected.DELETE("/courses/:id", middleware.RoleMiddleware("admin"), courseController.DeleteCourse)

	}
}
