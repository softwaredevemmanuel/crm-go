// routes/arm_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/arm"
	"crm-go/services/arm"
	"crm-go/middleware"
)

func ArmRoutes(router *gin.RouterGroup, db *gorm.DB) {
	armService := services.NewArmService(db)
	armHandler := controllers.NewArmHandler(armService)

	armGroup := router.Group("/api")
	armGroup.Use(middleware.AuthMiddleware())
	{
		// Create arm
		armGroup.POST("/arms", armHandler.CreateArm)

		// Get all arms with pagination and filters
		armGroup.GET("/arms", armHandler.GetAllArms)

		// Get arms by grade
		armGroup.GET("/arms/grade/:grade_id", armHandler.GetArmsByGrade)

		// Get arm by ID
		armGroup.GET("/arms/:id", armHandler.GetArmByID)

		// Update arm
		armGroup.PUT("/arms/:id", armHandler.UpdateArm)

		// Soft Delete arm
		armGroup.DELETE("/arms/:id", armHandler.DeleteArm)
		
		// Permanent Delete arm
		armGroup.DELETE("/arms/permanent/:id", armHandler.DeleteArmPermanently)
	}
}