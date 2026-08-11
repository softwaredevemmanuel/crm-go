// routes/subject_routes.go
package routes

import (
	"crm-go/controllers/subject"
	"crm-go/middleware"
	"crm-go/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SubjectRoutes(r *gin.Engine, db *gorm.DB) {

	// Initialize services
	subjectService := services.NewSubjectService(db)

	// Initialize controller
	subjectController := controllers.NewSubjectController(db, subjectService)

	subjects := r.Group("/subjects")
	subjects.GET("", subjectController.GetAllSubjects)
	subjects.GET("/:id", subjectController.GetSubjectByID)
	subjects.GET("/department", subjectController.GetSubjectDepartments)

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	protected.POST("/subjects", middleware.RoleMiddleware("admin"), subjectController.CreateSubject)
	protected.PUT("/subjects/:id", middleware.RoleMiddleware("admin"), subjectController.UpdateSubject)
	protected.DELETE("/subjects/:id", middleware.RoleMiddleware("admin"), subjectController.DeleteSubject)

}
