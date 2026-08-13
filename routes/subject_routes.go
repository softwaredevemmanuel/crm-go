// routes/subject_routes.go
package routes

import (
	"crm-go/controllers/subject"
	"crm-go/middleware"
	"crm-go/services/subjects"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SubjectRoutes(router *gin.RouterGroup, db *gorm.DB) {
	service := services.NewSubjectService(db)
	handler := controllers.NewSubjectController(service)

	subjectGroup := router.Group("/api/subjects")
	subjectGroup.Use(middleware.AuthMiddleware()) // Your auth middleware
	{
		// 1. Create subject
		subjectGroup.POST("", handler.CreateSubject)

		// 2. Get all subjects with pagination
		subjectGroup.GET("", handler.GetAllSubjects)

		// 3. Get subject by ID with department and head
		subjectGroup.GET("/department/head-of-department/:id", handler.GetSubjectWithDepartmentAndHead)

		// 4. Update subject
		subjectGroup.PUT("/:id", handler.UpdateSubject)

		// 5. Delete subject
		subjectGroup.DELETE("/:id", handler.DeleteSubject)
	}
}

