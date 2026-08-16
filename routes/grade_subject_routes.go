// routes/grade_subject_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/grade_subject"
	"crm-go/middleware"
	"crm-go/services/grade_subject"
)

func GradeSubjectRoutes(router *gin.RouterGroup, db *gorm.DB) {
	gradeSubjectService := services.NewGradeSubjectService(db)
	gradeSubjectHandler := controllers.NewGradeSubjectHandler(gradeSubjectService)

	gradeSubjectGroup := router.Group("/api")
	gradeSubjectGroup.Use(middleware.AuthMiddleware())
	{
		// Create single grade-subject mapping
		gradeSubjectGroup.POST("/grade-subjects", gradeSubjectHandler.CreateGradeSubject)

		// Bulk create grade-subject mappings
		gradeSubjectGroup.POST("/grade-subjects/bulk", gradeSubjectHandler.BulkCreateGradeSubjects)

		// Get all grade-subject mappings with pagination and filters
		gradeSubjectGroup.GET("/grade-subjects", gradeSubjectHandler.GetAllGradeSubjects)

		// Get subjects by grade
		gradeSubjectGroup.GET("/grade-subjects/grade/:grade_id", gradeSubjectHandler.GetSubjectsByGrade)

		// Get grades by subject
		gradeSubjectGroup.GET("/grade-subjects/subject/:subject_id", gradeSubjectHandler.GetGradesBySubject)

		// Get grade-subject mapping by ID
		gradeSubjectGroup.GET("/grade-subjects/:id", gradeSubjectHandler.GetGradeSubjectByID)

		// Update grade-subject mapping
		gradeSubjectGroup.PUT("/grade-subjects/:id", gradeSubjectHandler.UpdateGradeSubject)

		// Delete grade-subject mapping
		gradeSubjectGroup.DELETE("/grade-subjects/:id", gradeSubjectHandler.DeleteGradeSubject)
	}
}