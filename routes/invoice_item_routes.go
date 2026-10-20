package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/invoice_item"
	"crm-go/services/invoice_item"
	"crm-go/middleware"

)

func InvoiceItemRoutes(router *gin.RouterGroup, db *gorm.DB) {
	svc := services.NewInvoiceItemService(db)
	ctrl := controllers.NewInvoiceItemController(svc)

	// Standalone invoice-items group
	items := router.Group("/api")
	items.Use(middleware.AuthMiddleware())
	{
		items.POST("/invoice-items", ctrl.CreateInvoiceItem)
		items.GET("/invoice-items", ctrl.ListInvoiceItems)
		items.GET("/invoice-items/:id", ctrl.GetInvoiceItem)
		items.PATCH("/invoice-items/:id", ctrl.UpdateInvoiceItem)
		items.PUT("/invoice-items/:id", ctrl.UpdateInvoiceItem) // alias
		items.DELETE("/invoice-items/:id", ctrl.DeleteInvoiceItem)
	}



		// Nested under terms
	terms := router.Group("/api")
	terms.Use(middleware.AuthMiddleware())
	{
		terms.GET("/term/:term_id/invoice-items", ctrl.ListInvoiceItemsByTerm)
	}
}