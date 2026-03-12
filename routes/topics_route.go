package routes

import (
	topicController "crm-go/controllers/topics" // alias for topics
	"crm-go/middleware"

	"github.com/gin-gonic/gin"
)

func TopicRoutes(r *gin.Engine) {
	{
		topics := r.Group("/topics")
		topics.GET("/", topicController.GetAllTopics)
		topics.GET("/:id", topicController.GetTopicByID)

		// Protected routes
		protected := r.Group("/api")
		protected.Use(middleware.AuthMiddleware())
		protected.POST("/topics", middleware.RoleMiddleware("admin"), topicController.CreateTopic)
		protected.PUT("/topics/:id", middleware.RoleMiddleware("admin"), topicController.UpdateTopic)
		protected.DELETE("/topics/:id", middleware.RoleMiddleware("admin"), topicController.DeleteTopic)

	}
}
