package routes

import (
	assignmentController "crm-go/controllers/assignmentsubmission"
	"crm-go/middleware"
	"github.com/gin-gonic/gin"
	"crm-go/config"
	"crm-go/services/activity"
	"gorm.io/gorm"
)



func AssignmentSubmissionRoutes(r *gin.Engine, db *gorm.DB) {
		activitySvc := activity.NewService(db)

	// assignmentSubmissions := r.Group("/assignment_submissions")

	{
		// assignmentSubmissions.GET("/", assignmentController.GetAssignmentSubmissions)
		// assignmentSubmissions.GET("/:id", assignmentController.GetAssignmentSubmissionByID)

		// Protected routes
	
		protected := r.Group("/api")
		protected.Use(middleware.AuthMiddleware())
		assignmentController := assignmentController.NewAssignmentController(
		config.DB,
		activitySvc,
	)
		protected.POST("/assignment_submissions", middleware.RoleMiddleware("admin"), assignmentController.CreateAssignmentSubmission)
		// protected.PUT("/assignment_submissions/:id", middleware.RoleMiddleware("admin"), assignmentController.UpdateAssignmentSubmission)
		// protected.DELETE("/assignment_submissions/:id", middleware.RoleMiddleware("admin"), assignmentController.DeleteAssignmentSubmission)

	}
}
