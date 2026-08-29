// routes/student_enrollment_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/student_enrollment"
	"crm-go/services/student_enrollment"
	"crm-go/middleware"
)

func StudentEnrollmentRoutes(router *gin.RouterGroup, db *gorm.DB) {
	enrollmentService := services.NewStudentEnrollmentService(db)
	enrollmentHandler := handlers.NewStudentEnrollmentHandler(enrollmentService)

	enrollmentGroup := router.Group("/api")
	enrollmentGroup.Use(middleware.AuthMiddleware())
	{
		
		// Create single enrollment
		enrollmentGroup.POST("/student-enrollments", enrollmentHandler.CreateStudentEnrollment)
		
		// Bulk create enrollments
		enrollmentGroup.POST("/student-enrollments/bulk", enrollmentHandler.BulkCreateStudentEnrollments)

		// ============================================================
		// READ - Enrollment retrieval endpoints
		// ============================================================
		
		// Get all enrollments with pagination and filters
		enrollmentGroup.GET("/fetch-student-enrollments", enrollmentHandler.GetAllStudentEnrollments)
		
		// Get enrollment by ID
		enrollmentGroup.GET("/student-enrollments/:id", enrollmentHandler.GetStudentEnrollmentByID)
		
		// Get enrollments by student
		enrollmentGroup.GET("/student-enrollments/student/:student_id", enrollmentHandler.GetEnrollmentsByStudent)
		
		// Get current enrollment by student
		enrollmentGroup.GET("/student-enrollments/student/:student_id/current", enrollmentHandler.GetCurrentEnrollmentByStudent)
		
		// Get enrollments by arm
		enrollmentGroup.GET("/student-enrollments/arm/:arm_id", enrollmentHandler.GetEnrollmentsByArm)
		
		// Get enrollments by arm and status
		enrollmentGroup.GET("/student-enrollments/arm/:arm_id/status/:status", enrollmentHandler.GetEnrollmentsByArmAndStatus)
		
		// Get enrollment statistics
		enrollmentGroup.GET("/student-enrollments/stats", enrollmentHandler.GetEnrollmentStats)

		// ============================================================
		// UPDATE - Enrollment update endpoints
		// ============================================================
		
		// Update enrollment
		enrollmentGroup.PUT("/student-enrollments/:id", enrollmentHandler.UpdateStudentEnrollment)

		// ============================================================
		// DELETE - Enrollment deletion endpoints
		// ============================================================
		
		// Delete enrollment (soft delete)
		enrollmentGroup.DELETE("/student-enrollments/:id", enrollmentHandler.DeleteStudentEnrollment)
	}
}