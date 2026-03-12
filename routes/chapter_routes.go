package routes

import (
	chapters "crm-go/controllers/chapters"
	"github.com/gin-gonic/gin"
	"crm-go/middleware"
)

func ChapterRoutes(r *gin.Engine) {
	r.GET("/chapters", chapters.GetAllChapters)
	r.GET("/chapters/:id", chapters.GetChapterByID)
	
	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	protected.POST("/chapters", middleware.RoleMiddleware("admin"), chapters.CreateChapter)
	protected.PUT("/chapters/:id", middleware.RoleMiddleware("admin"), chapters.UpdateChapter)
	protected.DELETE("/chapters/:id", middleware.RoleMiddleware("admin"), chapters.DeleteChapter)
}
