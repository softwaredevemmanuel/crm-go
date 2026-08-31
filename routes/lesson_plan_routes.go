// routes/lesson_plan_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/lesson_plan"
	"crm-go/services/lesson_plan"
	"crm-go/middleware"
)

func LessonPlanRoutes(router *gin.RouterGroup, db *gorm.DB) {
	lessonPlanService := services.NewLessonPlanService(db)
	lessonPlanHandler := handlers.NewLessonPlanHandler(lessonPlanService)

	lessonPlanGroup := router.Group("/api")
	lessonPlanGroup.Use(middleware.AuthMiddleware())
	{
		// ============================================================
		// CREATE - Lesson plan creation endpoints
		// ============================================================
		
		// Create lesson plan
		lessonPlanGroup.POST("/lesson-plans", lessonPlanHandler.CreateLessonPlan)

		// ============================================================
		// READ - Lesson plan retrieval endpoints
		// ============================================================
		
		// Get all lesson plans with pagination and filters
		lessonPlanGroup.GET("/lesson-plans", lessonPlanHandler.GetAllLessonPlans)
		
		// Get lesson plan by ID
		lessonPlanGroup.GET("/lesson-plans/:id", lessonPlanHandler.GetLessonPlanByID)
		
		// Get lesson plan by lesson ID
		lessonPlanGroup.GET("/lesson-plans/lesson/:lesson_id", lessonPlanHandler.GetLessonPlanByLesson)

		// ============================================================
		// UPDATE - Lesson plan update endpoints
		// ============================================================
		
		// Update lesson plan
		lessonPlanGroup.PUT("/lesson-plans/:id", lessonPlanHandler.UpdateLessonPlan)

		// ============================================================
		// DELETE - Lesson plan deletion endpoints
		// ============================================================
		
		// Delete lesson plan (soft delete)
		lessonPlanGroup.DELETE("/lesson-plans/:id", lessonPlanHandler.DeleteLessonPlan)
	}
}