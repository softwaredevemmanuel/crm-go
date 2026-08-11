// routes/class_grade_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"crm-go/controllers/class_grade"
	"crm-go/services/class_grade"
	"crm-go/middleware"
)

func ClassGradeRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize service and handler
	classGradeService := services.NewClassGradeService(db)
	classGradeController := controllers.NewClassGradeController(db, classGradeService)

	{
		// Create class grade
	protected := router.Group("/api/class-grades")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("", middleware.RoleMiddleware("admin"), classGradeController.GetAllClassGrades)
	protected.GET("/academic-years", middleware.RoleMiddleware("admin"), classGradeController.GetAcademicYears)
	protected.GET("/levels", middleware.RoleMiddleware("admin"), classGradeController.GetLevels)
	protected.GET("/:id", middleware.RoleMiddleware("admin"), classGradeController.GetClassGradeByID)
	protected.POST("", middleware.RoleMiddleware("admin"), classGradeController.CreateClassGrade)
	protected.PUT("/:id", middleware.RoleMiddleware("admin"), classGradeController.UpdateClassGrade)
	protected.DELETE("/:id", middleware.RoleMiddleware("admin"), classGradeController.DeleteClassGrade)
		
		// classGradeGroup.GET("/", classGradeController.GetAllClassGrades)
		// classGradeGroup.GET("/:id", classGradeController.GetClassGradeByID)
		// classGradeGroup.PUT("/:id", classGradeController.UpdateClassGrade)
		// classGradeGroup.DELETE("/:id", classGradeCOntroller.DeleteClassGrade)
	}
}