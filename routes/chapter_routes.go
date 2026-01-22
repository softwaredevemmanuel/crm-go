package routes

import (
	chapters "crm-go/controllers/chapters"
	"github.com/gin-gonic/gin"
)

func ChapterRoutes(r *gin.Engine) {
	r.GET("/chapters", chapters.GetAllChapters)
	r.GET("/chapters/:id", chapters.GetChapterByID)
	
	r.POST("/api/chapters", chapters.CreateChapter)
	r.PUT("/api/chapters/:id", chapters.UpdateChapter)
	r.DELETE("/api/chapters/:id", chapters.DeleteChapter)
}
