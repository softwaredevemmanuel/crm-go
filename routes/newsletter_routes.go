package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/newsletter"
	"crm-go/services/newsletter"
	"crm-go/utils"
)

func NewsletterRoutes(api *gin.RouterGroup, db *gorm.DB, emailService utils.EmailService, authMiddleware gin.HandlerFunc) {
	svc := services.NewNewsletterService(db, emailService)
	ctrl := controllers.NewNewsletterController(svc)

	// Public group
	public := api.Group("/api")
	{
		public.POST("/newsletter/subscribe", ctrl.Subscribe)
		public.GET("/newsletter/confirm", ctrl.ConfirmSubscription)
		public.GET("/newsletter/unsubscribe", ctrl.Unsubscribe)
	}

	// Admin group (add your auth middleware here)
	admin := api.Group("/api")
	admin.Use(authMiddleware) // JWT / admin role check
	{
		admin.GET("/newsletter/subscribers", ctrl.ListSubscribers)
		admin.POST("/newsletter/subscribers", ctrl.CreateSubscriber)
		admin.GET("/newsletter/subscribers/:id", ctrl.GetSubscriber)
		admin.PATCH("/newsletter/subscribers/:id", ctrl.UpdateSubscriber)
		admin.PUT("/newsletter/subscribers/:id", ctrl.UpdateSubscriber) // alias
		admin.DELETE("/newsletter/subscribers/:id", ctrl.DeleteSubscriber)
		admin.GET("/newsletter/stats", ctrl.NewsletterStats)
	}
}