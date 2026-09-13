// routes/inventory_item_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/inventory_item"
	"crm-go/middleware"
	"crm-go/services/inventory_item"
)

func InventoryItemRoutes(router *gin.RouterGroup, db *gorm.DB) {
	itemService := services.NewInventoryItemService(db)
	itemController := controllers.NewInventoryItemController(itemService)

	itemGroup := router.Group("/api")
	itemGroup.Use(middleware.AuthMiddleware())
	{
		// CREATE
		itemGroup.POST("/inventory-items", itemController.CreateInventoryItem)

		// READ
		itemGroup.GET("/inventory-items", itemController.GetAllInventoryItems)
		itemGroup.GET("/inventory-items/stats", itemController.GetInventoryItemStats)
		itemGroup.GET("/inventory-items/section/:section_id", itemController.GetInventoryItemsBySection)
		itemGroup.GET("/inventory-items/:id", itemController.GetInventoryItemByID)

		// UPDATE
		itemGroup.PUT("/inventory-items/:id", itemController.UpdateInventoryItem)

		// DELETE
		itemGroup.DELETE("/inventory-items/:id", itemController.DeleteInventoryItem)
	}
}