package routes

import (
	assignmentController "crm-go/controllers/assignmentsubmission"
	"crm-go/middleware"
	"github.com/gin-gonic/gin"
)

func AssignmentSubmissionRoutes(r *gin.Engine) {
	// assignmentSubmissions := r.Group("/assignment_submissions")

	{
		// assignmentSubmissions.GET("/", assignmentController.GetAssignmentSubmissions)
		// assignmentSubmissions.GET("/:id", assignmentController.GetAssignmentSubmissionByID)

		// Protected routes
		protected := r.Group("/api")
		protected.Use(middleware.AuthMiddleware())
		protected.POST("/assignment_submissions", middleware.RoleMiddleware("admin"), assignmentController.CreateAssignmentSubmission)
		// protected.PUT("/assignment_submissions/:id", middleware.RoleMiddleware("admin"), assignmentController.UpdateAssignmentSubmission)
		// protected.DELETE("/assignment_submissions/:id", middleware.RoleMiddleware("admin"), assignmentController.DeleteAssignmentSubmission)

	}
}
