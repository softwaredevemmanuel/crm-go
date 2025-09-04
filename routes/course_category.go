package routes

import (
	"crm-go/controllers/category"
	"crm-go/middleware"
	"github.com/gin-gonic/gin"
)

func CourseCategoryRoutes(r *gin.Engine) {
	course_categories := r.Group("/course-categories")
	{
		course_categories.GET("/:id/courses", controllers.GetCoursesByCategory)

		// Protected routes
		protected := r.Group("/api")
		protected.Use(middleware.AuthMiddleware())
		protected.POST("/course-categories", middleware.RoleMiddleware("admin"), controllers.CreateCourseCategory)
		protected.DELETE("/course-categories/:id", middleware.RoleMiddleware("admin"), controllers.DeleteCourseCategory)

	}
}
