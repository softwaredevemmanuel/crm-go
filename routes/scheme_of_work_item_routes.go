// routes/scheme_of_work_item_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/scheme_of_work_item"
	"crm-go/services/scheme_of_work_item"
	"crm-go/middleware"
)

func SchemeOfWorkItemRoutes(router *gin.RouterGroup, db *gorm.DB) {
	itemService := services.NewSchemeOfWorkItemService(db)
	itemHandler := handlers.NewSchemeOfWorkItemHandler(itemService)

	itemGroup := router.Group("/api")
	itemGroup.Use(middleware.AuthMiddleware())
	{
		// ============================================================
		// CREATE - Item creation endpoints
		// ============================================================
		
		// Create single item
		itemGroup.POST("/scheme-items", itemHandler.CreateSchemeOfWorkItem)
		
		// Bulk create items
		itemGroup.POST("/scheme-items/bulk", itemHandler.BulkCreateSchemeItems)

		// ============================================================
		// READ - Item retrieval endpoints
		// ============================================================
		
		// Get all items with pagination and filters
		itemGroup.GET("/scheme-items", itemHandler.GetAllItems)
		
		// Get item by ID
		itemGroup.GET("/scheme-items/:id", itemHandler.GetItemByID)
		
		// Get items by scheme of work
		itemGroup.GET("/scheme-items/scheme/:scheme_of_work_id", itemHandler.GetItemsBySchemeOfWork)
		
		// Get items by module
		itemGroup.GET("/scheme-items/module/:module_id", itemHandler.GetItemsByModule)

		// ============================================================
		// UPDATE - Item update endpoints
		// ============================================================
		
		// Update item
		itemGroup.PUT("/scheme-items/:id", itemHandler.UpdateSchemeOfWorkItem)

		// ============================================================
		// DELETE - Item deletion endpoints
		// ============================================================
		
		// Delete item (soft delete)
		itemGroup.DELETE("/scheme-items/:id", itemHandler.DeleteSchemeOfWorkItem)
	}
}