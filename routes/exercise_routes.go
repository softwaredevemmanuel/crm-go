// routes/exercise_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/exercise"
	"crm-go/services/exercise"
	"crm-go/middleware"
)

func ExerciseRoutes(router *gin.RouterGroup, db *gorm.DB) {
	exerciseService := services.NewExerciseService(db)
	exerciseHandler := handlers.NewExerciseHandler(exerciseService)

	exerciseGroup := router.Group("/api")
	exerciseGroup.Use(middleware.AuthMiddleware())
	{
		// ============================================================
		// CREATE - Exercise creation endpoints
		// ============================================================
		
		// Create single exercise
		exerciseGroup.POST("/exercises", exerciseHandler.CreateExercise)
		
		// Bulk create exercises
		exerciseGroup.POST("/exercises/bulk", exerciseHandler.BulkCreateExercises)

		// ============================================================
		// READ - Exercise retrieval endpoints
		// ============================================================
		
		// Get all exercises with pagination and filters
		exerciseGroup.GET("/exercises", exerciseHandler.GetAllExercises)
		
		// Get exercise by ID
		exerciseGroup.GET("/exercises/:id", exerciseHandler.GetExerciseByID)
		
		// Get exercises by lesson
		exerciseGroup.GET("/exercises/lesson/:lesson_id", exerciseHandler.GetExercisesByLesson)

		// ============================================================
		// UPDATE - Exercise update endpoints
		// ============================================================
		
		// Update exercise
		exerciseGroup.PUT("/exercises/:id", exerciseHandler.UpdateExercise)

		// ============================================================
		// DELETE - Exercise deletion endpoints
		// ============================================================
		
		// Delete exercise (soft delete)
		exerciseGroup.DELETE("/exercises/:id", exerciseHandler.DeleteExercise)
	}
}