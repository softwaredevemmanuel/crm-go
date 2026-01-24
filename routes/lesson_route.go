package routes

import (
	lessonController "crm-go/controllers/lessons"   // alias for lessons
	"crm-go/middleware"
	"github.com/gin-gonic/gin"
)

func LessonRoutes(r *gin.Engine) {
	{
		lessons := r.Group("/lessons")
		lessons.GET("/", lessonController.GetLessons)
		lessons.GET("/:id", lessonController.GetLessonByID)

		// Protected routes
		protected := r.Group("/api")
		protected.Use(middleware.AuthMiddleware())
		protected.POST("/lessons", middleware.RoleMiddleware("admin"), lessonController.CreateLesson)
		protected.PUT("/lessons/:id", middleware.RoleMiddleware("admin"), lessonController.UpdateLesson)
		protected.DELETE("/lessons/:id", middleware.RoleMiddleware("admin"), lessonController.DeleteLesson)

	}
}
