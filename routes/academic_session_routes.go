// routes/academic_session_routes.go
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
	sessionHandler := handlers.NewAcademicSessionHandler(sessionService)

	sessionGroup := router.Group("/api")
	sessionGroup.Use(middleware.AuthMiddleware())
	{
		// ============================================================
		// CREATE - Session creation endpoints
		// ============================================================
		
		// Create single session
		sessionGroup.POST("/academic-sessions", sessionHandler.CreateAcademicSession)

		// ============================================================
		// READ - Session retrieval endpoints
		// ============================================================
		
		// Get all sessions with pagination and filters
		sessionGroup.GET("/academic-sessions", sessionHandler.GetAllSessions)
		
		// Get session by ID
		sessionGroup.GET("/academic-sessions/:id", sessionHandler.GetSessionByID)
		
		// Get current session
		sessionGroup.GET("/academic-sessions/current", sessionHandler.GetCurrentSession)
		
		// Get active sessions
		sessionGroup.GET("/academic-sessions/active", sessionHandler.GetActiveSessions)
		
		// Get session statistics
		sessionGroup.GET("/academic-sessions/stats", sessionHandler.GetSessionStats)

		// ============================================================
		// UPDATE - Session update endpoints
		// ============================================================
		
		// Update session
		sessionGroup.PUT("/academic-sessions/:id", sessionHandler.UpdateAcademicSession)

		// ============================================================
		// DELETE - Session deletion endpoints
		// ============================================================
		
		// Delete session (soft delete)
		sessionGroup.DELETE("/academic-sessions/:id", sessionHandler.DeleteAcademicSession)
	}
}