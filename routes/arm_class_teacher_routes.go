// routes/arm_class_teacher_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/arm_class_teacher"
	"crm-go/services/arm_class_teacher"
	"crm-go/middleware"
)

func ArmClassTeacherRoutes(router *gin.RouterGroup, db *gorm.DB) {
	assignmentService := services.NewArmClassTeacherService(db)
	assignmentHandler := handlers.NewArmClassTeacherHandler(assignmentService)

	assignmentGroup := router.Group("/api")
	assignmentGroup.Use(middleware.AuthMiddleware())
	{
		// ============================================================
		// CREATE - Assignment creation endpoints
		// ============================================================
		
		// Create single assignment
		assignmentGroup.POST("/arm-class-teachers", assignmentHandler.CreateAssignment)
		
		// Bulk create assignments
		assignmentGroup.POST("/arm-class-teachers/bulk", assignmentHandler.BulkAssignClassTeachers)

		// ============================================================
		// READ - Assignment retrieval endpoints
		// ============================================================
		
		// Get all assignments with pagination and filters
		assignmentGroup.GET("/arm-class-teachers", assignmentHandler.GetAllAssignments)
		
		// Get assignment by ID
		assignmentGroup.GET("/arm-class-teachers/:id", assignmentHandler.GetAssignmentByID)
		
		// Get assignments by teacher
		assignmentGroup.GET("/arm-class-teachers/teacher/:teacher_id", assignmentHandler.GetAssignmentsByTeacher)
		
		// Get assignments by arm
		assignmentGroup.GET("/arm-class-teachers/arm/:arm_id", assignmentHandler.GetAssignmentsByArm)
		
		// Get all arms with their class teachers
		assignmentGroup.GET("/arm-class-teachers/arms-with-teachers", assignmentHandler.GetArmsWithClassTeachers)

		// ============================================================
		// UPDATE - Assignment update endpoints
		// ============================================================
		
		// Update assignment
		assignmentGroup.PUT("/arm-class-teachers/:id", assignmentHandler.UpdateAssignment)

		// ============================================================
		// DELETE - Assignment deletion endpoints
		// ============================================================
		
		// Delete assignment (soft delete)
		assignmentGroup.DELETE("/arm-class-teachers/:id", assignmentHandler.DeleteAssignment)
	}
}