package routes

import (
	lessonController "crm-go/controllers/lessons"
	"crm-go/middleware"
	"crm-go/services/activity"
	services "crm-go/services/topics"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LessonRoutes(r *gin.Engine, db *gorm.DB) {
	activitySvc := activity.NewService(db)

	lessons := r.Group("/lessons")

	{
		createLessonSvc := services.NewCreateLessonService(db)
		getLessonSvc := services.NewGetLessonService(db)
		updateLessonSvc := services.NewUpdateLessonService(db)
		lessonCtrl := lessonController.NewCreateLessonController(
			db,
			createLessonSvc,
			getLessonSvc,
			updateLessonSvc,
			activitySvc,
		)

		lessons.GET("/", lessonCtrl.GetAllLessons)
		// lessons.GET("/:id", lessonController.GetLessonByID)

		// Protected routes

		protected := r.Group("/api")
		protected.Use(middleware.AuthMiddleware())

		protected.POST("/lessons", middleware.RoleMiddleware("admin"), lessonCtrl.CreateLesson)
		protected.PUT("/lessons/:id", middleware.RoleMiddleware("admin"), lessonCtrl.UpdateLesson)
		// protected.DELETE("/lessons/:id", middleware.RoleMiddleware("admin"), lessonController.DeleteLesson)

	}
}
