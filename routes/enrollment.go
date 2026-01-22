package routes

import (
	enrollmentsController "crm-go/controllers/enrollments"   // alias for courses
	"github.com/gin-gonic/gin"
)

func EnrollmentRoutes(r *gin.Engine) {
	enrollments := r.Group("/enrollments")
	{
		enrollments.GET("/", enrollmentsController.GetEnrollments)
		enrollments.GET("/:id", enrollmentsController.GetEnrollmentByID)
		enrollments.POST("/", enrollmentsController.CreateEnrollment)

	}
}


	