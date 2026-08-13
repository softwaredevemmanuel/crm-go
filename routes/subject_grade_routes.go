package routes

import (
	"crm-go/controllers/subject_grade"
	"crm-go/middleware"
	"crm-go/services/subject_grade"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SubjectGradeRoutes(r *gin.Engine, db *gorm.DB) {

	// Initialize service and handler
	subjectGradeService := services.NewSubjectGradeService(db)
	subjectGradeController := controllers.NewSubjectGradeController(db, subjectGradeService)

	protected := r.Group("/api/subject-grades")
	protected.Use(middleware.AuthMiddleware())

	// Create single subject-grade relationship
	protected.POST("", subjectGradeController.CreateSubjectGrade)

	// Bulk create subject-grade relationships
	protected.POST("/bulk", subjectGradeController.BulkCreateSubjectGrades)

	// Get all subject-grade relationships
	protected.GET("/relationship", subjectGradeController.GetAllSubjectGrades)

	// Get subjects by grade
	protected.GET("/grade/:grade_id", subjectGradeController.GetSubjectsByGrade)

	// Get grades by subject
	protected.GET("/subject/:subject_id", subjectGradeController.GetGradesBySubject)

	// Get subject-grade by ID
	protected.GET("/:id", subjectGradeController.GetSubjectGradeByID)

	// Delete subject-grade relationship
	protected.DELETE("/:id", subjectGradeController.DeleteSubjectGrade)

}
