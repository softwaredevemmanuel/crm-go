// routes/notification_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/notifications"
	"crm-go/services/notifications"
	"crm-go/middleware"
)

func NotificationRoutes(router *gin.RouterGroup, db *gorm.DB) {
	notificationService := services.NewNotificationService(db)
	notificationHandler := controllers.NewNotificationHandler(notificationService)

	notificationGroup := router.Group("/api")
	notificationGroup.Use(middleware.AuthMiddleware())
	{
		// Create notification
		notificationGroup.POST("/notifications", notificationHandler.CreateNotification)

		// Bulk create notifications
		notificationGroup.POST("/notifications/bulk", notificationHandler.BulkCreateNotifications)

		// Get all notifications with pagination and filters
		notificationGroup.GET("/notifications", notificationHandler.GetAllNotifications)

		// Get notification statistics
		notificationGroup.GET("/notifications/stats", notificationHandler.GetNotificationStats)

		// Mark all as read
		notificationGroup.POST("/notifications/read-all", notificationHandler.MarkAllAsRead)

		// Get notification by ID
		notificationGroup.GET("/notifications/:id", notificationHandler.GetNotificationByID)

		// Update notification
		notificationGroup.PUT("/notifications/:id", notificationHandler.UpdateNotification)

		// Mark notification as read
		notificationGroup.POST("/notifications/:id/read", notificationHandler.MarkAsRead)

		// Dismiss notification
		notificationGroup.POST("/notifications/:id/dismiss", notificationHandler.DismissNotification)

		// Delete notification
		notificationGroup.DELETE("/notifications/:id", notificationHandler.DeleteNotification)
	}
}