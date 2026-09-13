// routes/inventory_section_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/inventory_section"
	"crm-go/middleware"
	"crm-go/services/inventory_section"
)

func InventorySectionRoutes(router *gin.RouterGroup, db *gorm.DB) {
	sectionService := services.NewInventorySectionService(db)
	sectionController := controllers.NewInventorySectionController(sectionService)

	sectionGroup := router.Group("/api")
	sectionGroup.Use(middleware.AuthMiddleware())
	{
		// CREATE
		sectionGroup.POST("/inventory-sections", sectionController.CreateInventorySection)

		// READ
		sectionGroup.GET("/inventory-sections", sectionController.GetAllInventorySections)
		sectionGroup.GET("/inventory-sections/:id", sectionController.GetInventorySectionByID)

		// UPDATE
		sectionGroup.PUT("/inventory-sections/:id", sectionController.UpdateInventorySection)

		// DELETE
		sectionGroup.DELETE("/inventory-sections/:id", sectionController.DeleteInventorySection)
	}
}