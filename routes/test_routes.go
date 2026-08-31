// routes/test_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/test"
	"crm-go/services/test"
	"crm-go/middleware"
)

func TestRoutes(router *gin.RouterGroup, db *gorm.DB) {
	testService := services.NewTestService(db)
	testHandler := handlers.NewTestHandler(testService)

	testGroup := router.Group("/api")
	testGroup.Use(middleware.AuthMiddleware())
	{
		// ============================================================
		// CREATE - Test creation endpoints
		// ============================================================
		
		// Create single test
		testGroup.POST("/tests", testHandler.CreateTest)
		
		// Bulk create tests
		testGroup.POST("/tests/bulk", testHandler.BulkCreateTests)

		// ============================================================
		// READ - Test retrieval endpoints
		// ============================================================
		
		// Get all tests with pagination and filters
		testGroup.GET("/tests", testHandler.GetAllTests)
		
		// Get test by ID
		testGroup.GET("/tests/:id", testHandler.GetTestByID)
		
		// Get tests by subject
		testGroup.GET("/tests/subject/:subject_id", testHandler.GetTestsBySubject)
		
		// Get tests by class
		testGroup.GET("/tests/class/:class_id", testHandler.GetTestsByClass)

		// ============================================================
		// UPDATE - Test update endpoints
		// ============================================================
		
		// Update test
		testGroup.PUT("/tests/:id", testHandler.UpdateTest)

		// ============================================================
		// DELETE - Test deletion endpoints
		// ============================================================
		
		// Delete test (soft delete)
		testGroup.DELETE("/tests/:id", testHandler.DeleteTest)
	}
}