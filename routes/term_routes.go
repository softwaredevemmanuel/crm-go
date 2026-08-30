// routes/term_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/term"
	"crm-go/services/term"
	"crm-go/middleware"
)

func TermRoutes(router *gin.RouterGroup, db *gorm.DB) {
	termService := services.NewTermService(db)
	termHandler := handlers.NewTermHandler(termService)

	termGroup := router.Group("/api")
	termGroup.Use(middleware.AuthMiddleware())
	{
		// ============================================================
		// CREATE - Term creation endpoints
		// ============================================================
		
		// Create single term
		termGroup.POST("/terms", termHandler.CreateTerm)
		
		// Get all terms with pagination and filters
		termGroup.GET("/terms", termHandler.GetAllTerms)
		
		// Get term by ID
		termGroup.GET("/terms/:id", termHandler.GetTermByID)
		
		// Get terms by academic session
		termGroup.GET("/terms/session/:session_id", termHandler.GetTermsByAcademicSession)
		
		// Get current term for an academic session
		termGroup.GET("/terms/session/:session_id/current", termHandler.GetCurrentTerm)
		
		// Get term statistics
		termGroup.GET("/terms/stats", termHandler.GetTermStats)

		// ============================================================
		// UPDATE - Term update endpoints
		// ============================================================
		
		// Update term
		termGroup.PUT("/terms/:id", termHandler.UpdateTerm)

		// ============================================================
		// DELETE - Term deletion endpoints
		// ============================================================
		
		// Delete term (soft delete)
		termGroup.DELETE("/terms/:id", termHandler.DeleteTerm)
	}
}