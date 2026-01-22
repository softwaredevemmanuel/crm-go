package routes

import (
	"github.com/gin-gonic/gin"
	"crm-go/controllers/assignments"
)

func AssignmentRoutes(r *gin.Engine) {
	r.POST("/assignments", assignments.CreateAssignment)
	r.GET("/assignments", assignments.GetAssignments)
	r.GET("/assignments/:id", assignments.GetAssignmentByID)
	r.PUT("/assignments/:id", assignments.UpdateAssignment)
	r.DELETE("/assignments/:id", assignments.DeleteAssignment)
}
