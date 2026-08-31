// routes/lesson_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/lesson"
	"crm-go/services/lesson"
	"crm-go/middleware"
)

func LessonRoutes(router *gin.RouterGroup, db *gorm.DB) {
	lessonService := services.NewLessonService(db)
	lessonHandler := handlers.NewLessonHandler(lessonService)

	lessonGroup := router.Group("/api")
	lessonGroup.Use(middleware.AuthMiddleware())
	{
		// ============================================================
		// CREATE - Lesson creation endpoints
		// ============================================================
		
		// Create single lesson
		lessonGroup.POST("/lessons", lessonHandler.CreateLesson)
		
		// Bulk create lessons
		lessonGroup.POST("/lessons/bulk", lessonHandler.BulkCreateLessons)

		// ============================================================
		// READ - Lesson retrieval endpoints
		// ============================================================
		
		// Get all lessons with pagination and filters
		lessonGroup.GET("/lessons", lessonHandler.GetAllLessons)
		
		// Get lesson by ID
		lessonGroup.GET("/lessons/:id", lessonHandler.GetLessonByID)
		
	
		// Get lessons by class
		lessonGroup.GET("/lessons/class/:class_id", lessonHandler.GetLessonsByClass)

		// ============================================================
		// UPDATE - Lesson update endpoints
		// ============================================================
		
		// Update lesson
		lessonGroup.PUT("/lessons/:id", lessonHandler.UpdateLesson)

		// ============================================================
		// DELETE - Lesson deletion endpoints
		// ============================================================
		
		// Delete lesson (soft delete)
		lessonGroup.DELETE("/lessons/:id", lessonHandler.DeleteLesson)
	}
}