// routes/test_scheme_item_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/test_scheme_item"
	"crm-go/services/test_scheme_item"
	"crm-go/middleware"
)

func TestSchemeItemRoutes(router *gin.RouterGroup, db *gorm.DB) {
	testSchemeItemService := services.NewTestSchemeItemService(db)
	testSchemeItemHandler := handlers.NewTestSchemeItemHandler(testSchemeItemService)

	testSchemeItemGroup := router.Group("/api")
	testSchemeItemGroup.Use(middleware.AuthMiddleware())
	{
		// ============================================================
		// CREATE - Test scheme item creation endpoints
		// ============================================================
		
		// Create single test scheme item
		testSchemeItemGroup.POST("/test-scheme-items", testSchemeItemHandler.CreateTestSchemeItem)
		
		// Bulk create test scheme items
		testSchemeItemGroup.POST("/test-scheme-items/bulk", testSchemeItemHandler.BulkCreateTestSchemeItems)

		// ============================================================
		// READ - Test scheme item retrieval endpoints
		// ============================================================
		
		// Get all test scheme items with pagination and filters
		testSchemeItemGroup.GET("/test-scheme-items", testSchemeItemHandler.GetAllTestSchemeItems)
		
		// Get test scheme items by test
		testSchemeItemGroup.GET("/test-scheme-items/test/:test_id", testSchemeItemHandler.GetTestSchemeItemsByTest)
		
		// Get test scheme items by scheme of work item
		testSchemeItemGroup.GET("/test-scheme-items/scheme-item/:scheme_of_work_item_id", testSchemeItemHandler.GetTestSchemeItemsBySchemeItem)

		// ============================================================
		// DELETE - Test scheme item deletion endpoints
		// ============================================================
		
		// Delete single test scheme item
		testSchemeItemGroup.DELETE("/test-scheme-items/:test_id/:scheme_of_work_item_id", testSchemeItemHandler.DeleteTestSchemeItem)
		
		// Delete all test scheme items for a test
		testSchemeItemGroup.DELETE("/test-scheme-items/test/:test_id", testSchemeItemHandler.DeleteAllTestSchemeItemsByTest)
	}
}