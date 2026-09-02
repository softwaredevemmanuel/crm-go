// routes/lesson_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/lessons"
	"crm-go/services/lessons"
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
		
		// Get lessons by scheme of work
		lessonGroup.GET("/lessons/scheme/:scheme_of_work_id", lessonHandler.GetLessonsBySchemeOfWork)
		
		// Get lessons by module
		lessonGroup.GET("/lessons/module/:module_id", lessonHandler.GetLessonsByModule)
		
		// Get lessons by topic
		lessonGroup.GET("/lessons/topic/:topic_id", lessonHandler.GetLessonsByTopic)

		// ============================================================
		// UPDATE - Lesson update endpoints
		// ============================================================
		
		// Update lesson
		lessonGroup.PUT("/lessons/:id", lessonHandler.UpdateLesson)
		
		// Reorder lessons
		lessonGroup.PUT("/lessons/reorder", lessonHandler.ReorderLessons)

		// ============================================================
		// DELETE - Lesson deletion endpoints
		// ============================================================
		
		// Delete lesson (soft delete)
		lessonGroup.DELETE("/lessons/:id", lessonHandler.DeleteLesson)
	}
}