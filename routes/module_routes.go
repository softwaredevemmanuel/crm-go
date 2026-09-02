// routes/module_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/modules"
	"crm-go/services/modules"
	"crm-go/middleware"
)

func ModuleRoutes(router *gin.RouterGroup, db *gorm.DB) {
	moduleService := services.NewModuleService(db)
	moduleHandler := controllers.NewModuleHandler(moduleService)

	moduleGroup := router.Group("/api")
	moduleGroup.Use(middleware.AuthMiddleware())
	{
		// ============================================================
		// CREATE - Module creation endpoints
		// ============================================================
		
		// Create single module
		moduleGroup.POST("/modules", moduleHandler.CreateModule)
		
		// Bulk create modules
		moduleGroup.POST("/modules/bulk", moduleHandler.BulkCreateModules)

		// ============================================================
		// READ - Module retrieval endpoints
		// ============================================================
		
		// Get all modules with pagination and filters
		moduleGroup.GET("/modules", moduleHandler.GetAllModules)
		
		// Get module by ID
		moduleGroup.GET("/modules/:id", moduleHandler.GetModuleByID)
		
		// Get modules by scheme of work
		moduleGroup.GET("/modules/scheme/:scheme_of_work_id", moduleHandler.GetModulesBySchemeOfWork)

		// ============================================================
		// UPDATE - Module update endpoints
		// ============================================================
		
		// Update module
		moduleGroup.PUT("/modules/:id", moduleHandler.UpdateModule)
		
		// Reorder modules
		moduleGroup.PUT("/modules/reorder", moduleHandler.ReorderModules)

		// ============================================================
		// DELETE - Module deletion endpoints
		// ============================================================
		
		// Delete module (soft delete)
		moduleGroup.DELETE("/modules/:id", moduleHandler.DeleteModule)
	}
}