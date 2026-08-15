package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/guardian"
	"crm-go/services/guardian"
	"crm-go/middleware"
)

func GuardianRoutes(router *gin.RouterGroup, db *gorm.DB) {
	guardianService := services.NewGuardianService(db)
	guardianHandler := controllers.NewGuardianHandler(guardianService)

	guardianGroup := router.Group("/api")
	guardianGroup.Use(middleware.AuthMiddleware())
	{
		// Create guardian
		guardianGroup.POST("/guardians", guardianHandler.CreateGuardian)

		// Get all guardians with pagination and filters
		guardianGroup.GET("/guardians", guardianHandler.GetAllGuardians)

		// Get guardians by student
		guardianGroup.GET("/guardians/student/:student_id", guardianHandler.GetGuardiansByStudent)

		// Get guardian by ID
		guardianGroup.GET("/guardians/:id", guardianHandler.GetGuardianByID)

		// Update guardian
		guardianGroup.PUT("/guardians/:id", guardianHandler.UpdateGuardian)

		// Delete guardian
		guardianGroup.DELETE("/guardians/:id", guardianHandler.DeleteGuardian)
	}
}