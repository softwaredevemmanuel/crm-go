package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/newsletter_campaign"
	"crm-go/services/newsletter_campaign"
	"crm-go/utils"
)

func NewsletterCampaignRoutes(
	api *gin.RouterGroup,
	db *gorm.DB,
	emailService utils.EmailService,
	authMiddleware gin.HandlerFunc,
) {
	svc := services.NewNewsletterCampaignService(db, emailService)
	ctrl := controllers.NewNewsletterCampaignController(svc)

	campaigns := api.Group("/api")
	campaigns.Use(authMiddleware)
	{
		campaigns.POST("/newsletter/campaigns", ctrl.CreateCampaign)
		campaigns.GET("/fetch/newsletter/campaigns", ctrl.ListCampaigns)
		campaigns.GET("/newsletter/campaigns/preview", ctrl.PreviewRecipients)
		campaigns.GET("/newsletter/campaigns/:id", ctrl.GetCampaign)
		campaigns.PATCH("/newsletter/campaigns/:id", ctrl.UpdateCampaign)
		campaigns.PUT("/newsletter/campaigns/:id", ctrl.UpdateCampaign)
		campaigns.DELETE("/newsletter/campaigns/:id", ctrl.DeleteCampaign)
		campaigns.POST("/newsletter/campaigns/:id/send", ctrl.SendCampaign)
		campaigns.POST("/newsletter/campaigns/:id/cancel", ctrl.CancelCampaign)
	}
}