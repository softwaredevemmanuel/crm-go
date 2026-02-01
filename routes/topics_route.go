package routes

import (
	topicController "crm-go/controllers/topics"
	"crm-go/middleware"
	"crm-go/services/topics"
	"crm-go/services/activity"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TopicRoutes(r *gin.Engine, db *gorm.DB) {
	activitySvc := activity.NewService(db)

	topics := r.Group("/topics")

	{
		createTopicSvc := services.NewCreateTopicService(db)
		getTopicSvc := services.NewGetTopicService(db)
		updateTopicSvc := services.NewUpdateTopicService(db)
		topicCtrl := topicController.NewCreateTopicController(
			db,
			createTopicSvc,
			getTopicSvc,
			updateTopicSvc,
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
