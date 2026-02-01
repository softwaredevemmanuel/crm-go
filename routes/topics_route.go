package routes

import (
	topicController "crm-go/controllers/topic"
	"crm-go/middleware"
	"crm-go/services"
	"crm-go/services/activity"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TopicRoutes(r *gin.Engine, db *gorm.DB) {
	activitySvc := activity.NewService(db)

	topics := r.Group("/topics")

	{
		topicSvc := services.NewTopicService(db)
		topicCtrl := topicController.NewTopicController(
			db,
			topicSvc,
			activitySvc,
		)
		
		topics.GET("/", topicCtrl.GetAllTopics)
		// topics.GET("/:id", topicController.GetTopicByID)

		// Protected routes

		protected := r.Group("/api")
		protected.Use(middleware.AuthMiddleware())

		protected.POST("/topics", middleware.RoleMiddleware("admin"), topicCtrl.CreateTopic)
		protected.PUT("/topics/:id", middleware.RoleMiddleware("admin"), topicCtrl.UpdateTopic)
		// protected.DELETE("/topics/:id", middleware.RoleMiddleware("admin"), topicController.DeleteTopic)

	}
}
