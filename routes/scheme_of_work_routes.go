// routes/scheme_of_work_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/scheme_of_work"
	"crm-go/services/scheme_of_work"
	"crm-go/middleware"
)

func SchemeOfWorkRoutes(router *gin.RouterGroup, db *gorm.DB) {
	schemeService := services.NewSchemeOfWorkService(db)
	schemeHandler := handlers.NewSchemeOfWorkHandler(schemeService)

	schemeGroup := router.Group("/api")
	schemeGroup.Use(middleware.AuthMiddleware())
	{
		// ============================================================
		// CREATE - Scheme creation endpoints
		// ============================================================
		
		// Create single scheme
		schemeGroup.POST("/schemes", schemeHandler.CreateSchemeOfWork)
		
		// Bulk create schemes
		schemeGroup.POST("/schemes/bulk", schemeHandler.BulkCreateSchemes)

		// ============================================================
		// READ - Scheme retrieval endpoints
		// ============================================================
		
		// Get all schemes with pagination and filters
		schemeGroup.GET("/schemes", schemeHandler.GetAllSchemes)
		
		// Get scheme by ID
		schemeGroup.GET("/schemes/:id", schemeHandler.GetSchemeByID)
		
		// Get schemes by subject
		schemeGroup.GET("/schemes/subject/:subject_id", schemeHandler.GetSchemesBySubject)
		
		// Get schemes by class
		schemeGroup.GET("/schemes/class/:class_id", schemeHandler.GetSchemesByClass)

		// ============================================================
		// UPDATE - Scheme update endpoints
		// ============================================================
		
		// Update scheme
		schemeGroup.PUT("/schemes/:id", schemeHandler.UpdateSchemeOfWork)

		// ============================================================
		// DELETE - Scheme deletion endpoints
		// ============================================================
		
		// Delete scheme (soft delete)
		schemeGroup.DELETE("/schemes/:id", schemeHandler.DeleteSchemeOfWork)
	}
}