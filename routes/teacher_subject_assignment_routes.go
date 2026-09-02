// routes/teacher_subject_assignment_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/teacher_subject"
	"crm-go/services/teacher_subject"
	"crm-go/middleware"
)

func TeacherSubjectAssignmentRoutes(router *gin.RouterGroup, db *gorm.DB) {
	assignmentService := services.NewTeacherSubjectAssignmentService(db)
	assignmentHandler := handlers.NewTeacherSubjectAssignmentHandler(assignmentService)

	assignmentGroup := router.Group("/api")
	assignmentGroup.Use(middleware.AuthMiddleware())
	{
		// ============================================================
		// CREATE - Assignment creation endpoints
		// ============================================================
		
		// Create single assignment
		assignmentGroup.POST("/teacher-subject-assignments", assignmentHandler.CreateAssignment)
		
		// Bulk create assignments
		assignmentGroup.POST("/teacher-subject-assignments/bulk", assignmentHandler.BulkAssignSubjects)

		// ============================================================
		// READ - Assignment retrieval endpoints
		// ============================================================
		
		// Get all assignments with pagination and filters
		assignmentGroup.GET("/teacher-subject-assignments", assignmentHandler.GetAllAssignments)
		
		// Get assignment by ID
		assignmentGroup.GET("/teacher-subject-assignments/:id", assignmentHandler.GetAssignmentByID)
		
		// Get assignments by teacher
		assignmentGroup.GET("/teacher-subject-assignments/teacher/:teacher_id", assignmentHandler.GetAssignmentsByTeacher)
		
		// Get assignments by grade
		assignmentGroup.GET("/teacher-subject-assignments/grade/:grade_id", assignmentHandler.GetAssignmentsByGrade)
		
		// Get assignments by teacher and grade
		assignmentGroup.GET("/teacher-subject-assignments/teacher/:teacher_id/grade/:grade_id", assignmentHandler.GetAssignmentsByTeacherAndGrade)

		// ============================================================
		// UPDATE - Assignment update endpoints
		// ============================================================
		
		// Update assignment
		assignmentGroup.PUT("/teacher-subject-assignments/:id", assignmentHandler.UpdateAssignment)

		// ============================================================
		// DELETE - Assignment deletion endpoints
		// ============================================================
		
		// Delete assignment (soft delete)
		assignmentGroup.DELETE("/teacher-subject-assignments/:id", assignmentHandler.DeleteAssignment)
	}
}