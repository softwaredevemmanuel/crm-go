// routes/push_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers"
	"crm-go/services"
	"crm-go/middleware"
)

func PushRoutes(router *gin.RouterGroup, db *gorm.DB) {
	pushService := services.NewPushService(db)
	pushHandler := controllers.NewPushHandler(pushService)

	pushGroup := router.Group("/api")
	pushGroup.Use(middleware.AuthMiddleware())
	{
		pushGroup.POST("/push/subscribe", pushHandler.Subscribe)
		pushGroup.DELETE("/push/subscribe", pushHandler.Unsubscribe)
		pushGroup.POST("/push/test", pushHandler.SendTestNotification)
	}
}