// routes/user_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/users"
	"crm-go/services/users"
	"crm-go/middleware"
)

func UserRoutes(router *gin.RouterGroup, db *gorm.DB) {
	userService := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService)

	userGroup := router.Group("/api")
	userGroup.Use(middleware.AuthMiddleware())
	{
		// Get all users with pagination and filters
		userGroup.GET("/users", userHandler.GetAllUsers)

		// Get users by role
		userGroup.GET("/users/role/:role", userHandler.GetUsersByRole)

		// Delete user
		userGroup.DELETE("/users/:id", userHandler.DeleteUser)
	}
}