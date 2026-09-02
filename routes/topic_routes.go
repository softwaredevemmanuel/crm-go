// routes/topic_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/topics"
	"crm-go/services/topics"
	"crm-go/middleware"
)

func TopicRoutes(router *gin.RouterGroup, db *gorm.DB) {
	topicService := services.NewTopicService(db)
	topicHandler := handlers.NewTopicHandler(topicService)

	topicGroup := router.Group("/api")
	topicGroup.Use(middleware.AuthMiddleware())
	{
		// ============================================================
		// CREATE - Topic creation endpoints
		// ============================================================
		
		// Create single topic
		topicGroup.POST("/topics", topicHandler.CreateTopic)
		
		// Bulk create topics
		topicGroup.POST("/topics/bulk", topicHandler.BulkCreateTopics)

		// ============================================================
		// READ - Topic retrieval endpoints
		// ============================================================
		
		// Get all topics with pagination and filters
		topicGroup.GET("/topics", topicHandler.GetAllTopics)
		
		// Get topic by ID
		topicGroup.GET("/topics/:id", topicHandler.GetTopicByID)
		
		// Get topics by module
		topicGroup.GET("/topics/module/:module_id", topicHandler.GetTopicsByModule)

		// ============================================================
		// UPDATE - Topic update endpoints
		// ============================================================
		
		// Update topic
		topicGroup.PUT("/topics/:id", topicHandler.UpdateTopic)
		
		// Reorder topics
		topicGroup.PUT("/topics/reorder", topicHandler.ReorderTopics)

		// ============================================================
		// DELETE - Topic deletion endpoints
		// ============================================================
		
		// Delete topic (soft delete)
		topicGroup.DELETE("/topics/:id", topicHandler.DeleteTopic)
	}
}