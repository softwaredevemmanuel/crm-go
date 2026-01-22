package routes

import (
	announcementController "crm-go/controllers/announcements"
	"crm-go/middleware"
	"github.com/gin-gonic/gin"
)

func AnnouncementRoutes(r *gin.Engine) {
	announcements := r.Group("/announcements")

	{
		announcements.GET("/", announcementController.GetAnnouncements)
		announcements.GET("/:id", announcementController.GetAnnouncementByID)
		
		// Protected routes
		protected := r.Group("/api")
		protected.Use(middleware.AuthMiddleware())
		protected.POST("/announcements", middleware.RoleMiddleware("admin"), announcementController.CreateAnnouncement)
		protected.PUT("/announcements/:id", middleware.RoleMiddleware("admin"), announcementController.UpdateAnnouncement)
		protected.DELETE("/announcements/:id", middleware.RoleMiddleware("admin"), announcementController.DeleteAnnouncement)

	}
}
