// routes/exam_scheme_item_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/exam_scheme_item"
	"crm-go/services/exam_scheme_item"
	"crm-go/middleware"
)

func ExamSchemeItemRoutes(router *gin.RouterGroup, db *gorm.DB) {
	examSchemeItemService := services.NewExamSchemeItemService(db)
	examSchemeItemHandler := handlers.NewExamSchemeItemHandler(examSchemeItemService)

	examSchemeItemGroup := router.Group("/api")
	examSchemeItemGroup.Use(middleware.AuthMiddleware())
	{
		// ============================================================
		// CREATE - Exam scheme item creation endpoints
		// ============================================================
		
		// Create single exam scheme item
		examSchemeItemGroup.POST("/exam-scheme-items", examSchemeItemHandler.CreateExamSchemeItem)
		
		// Bulk create exam scheme items
		examSchemeItemGroup.POST("/exam-scheme-items/bulk", examSchemeItemHandler.BulkCreateExamSchemeItems)

		// ============================================================
		// READ - Exam scheme item retrieval endpoints
		// ============================================================
		
		// Get all exam scheme items with pagination and filters
		examSchemeItemGroup.GET("/exam-scheme-items", examSchemeItemHandler.GetAllExamSchemeItems)
		
		// Get exam scheme items by exam
		examSchemeItemGroup.GET("/exam-scheme-items/exam/:exam_id", examSchemeItemHandler.GetExamSchemeItemsByExam)
		
		// Get exam scheme items by scheme of work item
		examSchemeItemGroup.GET("/exam-scheme-items/scheme-item/:scheme_of_work_item_id", examSchemeItemHandler.GetExamSchemeItemsBySchemeItem)

		// ============================================================
		// DELETE - Exam scheme item deletion endpoints
		// ============================================================
		
		// Delete single exam scheme item
		examSchemeItemGroup.DELETE("/exam-scheme-items/:exam_id/:scheme_of_work_item_id", examSchemeItemHandler.DeleteExamSchemeItem)
		
		// Delete all exam scheme items for an exam
		examSchemeItemGroup.DELETE("/exam-scheme-items/exam/:exam_id", examSchemeItemHandler.DeleteAllExamSchemeItemsByExam)
	}
}