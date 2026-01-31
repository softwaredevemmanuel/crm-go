package routes

import (
	topicController "crm-go/controllers/topic"
	"crm-go/middleware"
	"github.com/gin-gonic/gin"
	"crm-go/services/activity"
	"crm-go/services"
	"gorm.io/gorm"
)



func TopicRoutes(r *gin.Engine, db *gorm.DB) {
		activitySvc := activity.NewService(db)

	// topics := r.Group("/topics")

	{
		// topics.GET("/", topicController.GetAllTopics)
		// topics.GET("/:id", topicController.GetTopicByID)

		// Protected routes
	
		protected := r.Group("/api")
		protected.Use(middleware.AuthMiddleware())
		topicSvc := services.NewTopicService(db)
		topicCtrl := topicController.NewTopicController(
			db,
			topicSvc,
			activitySvc,
		)
		protected.POST("/topics", middleware.RoleMiddleware("admin"), topicCtrl.CreateTopic)
		// protected.PUT("/topics/:id", middleware.RoleMiddleware("admin"), topicController.UpdateTopic)
		// protected.DELETE("/topics/:id", middleware.RoleMiddleware("admin"), topicController.DeleteTopic)

	}
}
