// routes/learning_objective_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/learning_objective"
	"crm-go/services/learning_objective"
	"crm-go/middleware"
)

func LearningObjectiveRoutes(router *gin.RouterGroup, db *gorm.DB) {
	objectiveService := services.NewLearningObjectiveService(db)
	objectiveHandler := handlers.NewLearningObjectiveHandler(objectiveService)

	objectiveGroup := router.Group("/api")
	objectiveGroup.Use(middleware.AuthMiddleware())
	{
		// ============================================================
		// CREATE - Objective creation endpoints
		// ============================================================
		
		// Create single objective
		objectiveGroup.POST("/learning-objectives", objectiveHandler.CreateLearningObjective)
		
		// Bulk create objectives
		objectiveGroup.POST("/learning-objectives/bulk", objectiveHandler.BulkCreateLearningObjectives)

		// ============================================================
		// READ - Objective retrieval endpoints
		// ============================================================
		
		// Get all objectives with pagination and filters
		objectiveGroup.GET("/learning-objectives", objectiveHandler.GetAllObjectives)
		
		// Get objective by ID
		objectiveGroup.GET("/learning-objectives/:id", objectiveHandler.GetObjectiveByID)
		
		// Get objectives by scheme of work item
		objectiveGroup.GET("/learning-objectives/scheme-item/:scheme_of_work_item_id", objectiveHandler.GetObjectivesBySchemeItem)

		// ============================================================
		// UPDATE - Objective update endpoints
		// ============================================================
		
		// Update objective
		objectiveGroup.PUT("/learning-objectives/:id", objectiveHandler.UpdateLearningObjective)

		// ============================================================
		// DELETE - Objective deletion endpoints
		// ============================================================
		
		// Delete objective (soft delete)
		objectiveGroup.DELETE("/learning-objectives/:id", objectiveHandler.DeleteLearningObjective)
	}
}