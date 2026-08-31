// routes/subject_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/subject"
	"crm-go/services/subjects"
	"crm-go/middleware"
)

func SubjectRoutes(router *gin.RouterGroup, db *gorm.DB) {
	subjectService := services.NewSubjectService(db)
	subjectHandler := handlers.NewSubjectHandler(subjectService)

	subjectGroup := router.Group("/api")
	subjectGroup.Use(middleware.AuthMiddleware())
	{
		// ============================================================
		// CREATE - Subject creation endpoints
		// ============================================================
		
		// Create single subject
		subjectGroup.POST("/subjects", subjectHandler.CreateSubject)
		
		// Bulk create subjects
		subjectGroup.POST("/subjects/bulk", subjectHandler.BulkCreateSubjects)

		// ============================================================
		// READ - Subject retrieval endpoints
		// ============================================================
		
		// Get all subjects with pagination and filters
		subjectGroup.GET("/subjects", subjectHandler.GetAllSubjects)
		
		// Get subject by ID
		subjectGroup.GET("/subjects/:id", subjectHandler.GetSubjectByID)
		
		// Get subjects by department
		subjectGroup.GET("/subjects/department/:department_id", subjectHandler.GetSubjectsByDepartment)
		
		// Get active subjects
		subjectGroup.GET("/subjects/active", subjectHandler.GetActiveSubjects)
		
		// Get subject statistics
		subjectGroup.GET("/subjects/stats", subjectHandler.GetSubjectStats)

		// ============================================================
		// UPDATE - Subject update endpoints
		// ============================================================
		
		// Update subject
		subjectGroup.PUT("/subjects/:id", subjectHandler.UpdateSubject)

		// ============================================================
		// DELETE - Subject deletion endpoints
		// ============================================================
		
		// Delete subject (soft delete)
		subjectGroup.DELETE("/subjects/:id", subjectHandler.DeleteSubject)
	}
}