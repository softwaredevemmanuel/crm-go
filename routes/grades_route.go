// routes/grade_routes.go
package routes

import (
	"crm-go/controllers/grades"
	"crm-go/middleware"
	"crm-go/services/activity"
	"crm-go/services/grades"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GradeRoutes(r *gin.Engine, db *gorm.DB) {
	// Initialize services
	gradeService := services.NewGradeService(db)
	activityService := activity.NewService(db) // Assuming you have this

	// Initialize controller
	gradeController := controllers.NewGradeController(db, gradeService, activityService)

	// Protected routes

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())

	protected.POST("/grades", middleware.RoleMiddleware("admin"), gradeController.CreateGrade)

	// New update routes
	protected.PUT("/grades/:id", middleware.RoleMiddleware("admin"), gradeController.UpdateGrade)
}
