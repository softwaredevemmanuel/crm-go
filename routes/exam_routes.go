// routes/exam_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/exam"
	"crm-go/services/exam"
	"crm-go/middleware"
)

func ExamRoutes(router *gin.RouterGroup, db *gorm.DB) {
	examService := services.NewExamService(db)
	examHandler := handlers.NewExamHandler(examService)

	examGroup := router.Group("/api")
	examGroup.Use(middleware.AuthMiddleware())
	{
		// ============================================================
		// CREATE - Exam creation endpoints
		// ============================================================
		
		// Create single exam
		examGroup.POST("/exams", examHandler.CreateExam)
		
		// Bulk create exams
		examGroup.POST("/exams/bulk", examHandler.BulkCreateExams)

		// ============================================================
		// READ - Exam retrieval endpoints
		// ============================================================
		
		// Get all exams with pagination and filters
		examGroup.GET("/exams", examHandler.GetAllExams)
		
		// Get exam by ID
		examGroup.GET("/exams/:id", examHandler.GetExamByID)
		
		// Get exams by subject
		examGroup.GET("/exams/subject/:subject_id", examHandler.GetExamsBySubject)
		
		// Get exams by class
		examGroup.GET("/exams/class/:class_id", examHandler.GetExamsByClass)

		// ============================================================
		// UPDATE - Exam update endpoints
		// ============================================================
		
		// Update exam
		examGroup.PUT("/exams/:id", examHandler.UpdateExam)

		// ============================================================
		// DELETE - Exam deletion endpoints
		// ============================================================
		
		// Delete exam (soft delete)
		examGroup.DELETE("/exams/:id", examHandler.DeleteExam)
	}
}