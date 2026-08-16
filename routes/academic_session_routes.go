package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/academic_session"
	"crm-go/services/academic_session"
	"crm-go/middleware"
)

func AcademicSessionRoutes(router *gin.RouterGroup, db *gorm.DB) {
	sessionService := services.NewAcademicSessionService(db)
	sessionHandler := controllers.NewAcademicSessionHandler(sessionService)

	sessionGroup := router.Group("/api")
	sessionGroup.Use(middleware.AuthMiddleware())
	{
		// Create academic session
		sessionGroup.POST("/academic-sessions", sessionHandler.CreateAcademicSession)

		// Get all academic sessions with pagination and filters
		sessionGroup.GET("/academic-sessions", sessionHandler.GetAllAcademicSessions)

		// Get current academic session
		sessionGroup.GET("/academic-sessions/current", sessionHandler.GetCurrentAcademicSession)

		// Get academic session by ID
		sessionGroup.GET("/academic-sessions/:id", sessionHandler.GetAcademicSessionByID)

		// Update academic session
		sessionGroup.PUT("/academic-sessions/:id", sessionHandler.UpdateAcademicSession)

		// Delete academic session
		sessionGroup.DELETE("/academic-sessions/:id", sessionHandler.DeleteAcademicSession)
	}
}