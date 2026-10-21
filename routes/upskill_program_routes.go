package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"crm-go/middleware"
	"crm-go/controllers/upskill_program"
	"crm-go/services/upskill_program"
)

func UpskillProgramRoutes(router *gin.RouterGroup, db *gorm.DB) {
	svc := services.NewUpskillProgramService(db)
	ctrl := controllers.NewUpskillProgramController(svc)

	programs := router.Group("/api")
	// Public read routes (no auth)
	{
		programs.GET("/upskill-programs", ctrl.ListPrograms)
		programs.GET("/upskill-programs/:id", ctrl.GetProgram)
		programs.GET("/upskill-programs/slug/:slug", ctrl.GetProgramBySlug)
	}

	// Protected write routes
	protected := programs.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/upskill-programs", ctrl.CreateProgram)
		protected.PATCH("/upskill-programs/:id", ctrl.UpdateProgram)
		protected.PUT("/upskill-programs/:id", ctrl.UpdateProgram)
		protected.DELETE("/upskill-programs/:id", ctrl.DeleteProgram)
		protected.PATCH("/upskill-programs/:id/active", ctrl.ToggleProgramActive)

		// Course management inside a program
		protected.POST("/upskill-programs/:id/courses", ctrl.AddCourseToProgram)
		protected.DELETE("/upskill-programs/:id/courses/:course_id", ctrl.RemoveCourseFromProgram)
		protected.POST("/upskill-programs/:id/courses/reorder", ctrl.ReorderProgramCourses)
	}
}