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
		// ============================================================
		// CREATE - Enrollment creation endpoints
		// ============================================================
		
		// Create single enrollment
		enrollmentGroup.POST("/student-enrollments", enrollmentHandler.CreateStudentEnrollment)
		
		// Bulk create enrollments
		enrollmentGroup.POST("/student-enrollments/bulk", enrollmentHandler.BulkCreateStudentEnrollments)

		// ============================================================
		// READ - Enrollment retrieval endpoints
		// ============================================================
		
		// Get all enrollments with pagination and filters (admin/super admin)
		enrollmentGroup.GET("/fetch/student-enrollments", enrollmentHandler.GetAllStudentEnrollments)
		
		// Get enrollment by ID
		enrollmentGroup.GET("/student-enrollments/:id", enrollmentHandler.GetStudentEnrollmentByID)
		
		// Get all verified enrollments for a specific student
		enrollmentGroup.GET("/student-enrollments/student/:student_id", enrollmentHandler.GetEnrollmentsByStudent)
		
		// Get current enrollment by student
		enrollmentGroup.GET("/student-enrollments/student/:student_id/current", enrollmentHandler.GetCurrentEnrollmentByStudent)
		
		// Get enrollments by grade
		enrollmentGroup.GET("/student-enrollments/grade/:grade_id", enrollmentHandler.GetEnrollmentsByGrade)

		// ============================================================
		// CLASS TEACHER ENDPOINTS - New endpoints for class teachers
		// ============================================================
		
		// Get all enrollments for a class teacher (all grades taught by teacher)
		enrollmentGroup.GET("/student-enrollments/teacher/:teacher_id", enrollmentHandler.GetEnrollmentsByClassTeacher)
		
		// Get paginated enrollments for a class teacher with filters
		enrollmentGroup.GET("/student-enrollments/teacher/:teacher_id/paginated", enrollmentHandler.GetEnrollmentsByClassTeacherPaginated)
		
		// Get all grades assigned to a class teacher
		enrollmentGroup.GET("/student-enrollments/teacher/:teacher_id/grades", enrollmentHandler.GetGradesByClassTeacher)
		
		// Get dashboard statistics for a class teacher
		enrollmentGroup.GET("/student-enrollments/teacher/:teacher_id/dashboard", enrollmentHandler.GetClassTeacherDashboardStats)

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