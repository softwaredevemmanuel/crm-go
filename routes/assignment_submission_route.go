package routes

import (
	"github.com/gin-gonic/gin"
)

func AssignmentSubmissionRoutes(r *gin.Engine) {
	r.POST("/assignments_submission")
	r.GET("/assignments_submission")
	r.GET("/assignments_submission/:id")
	r.PUT("/assignments_submission/:id")
	r.DELETE("/assignments_submission/:id")
}
