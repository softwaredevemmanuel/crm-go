package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"crm-go/middleware"
	"crm-go/controllers/invoice"
	"crm-go/services/invoice"
)

func InvoiceRoutes(router *gin.RouterGroup, db *gorm.DB) {
	invoiceSvc := services.NewInvoiceService(db)
	invCtrl := controllers.NewInvoiceController(invoiceSvc)

	invoices := router.Group("/api")
	invoices.Use(middleware.AuthMiddleware())
	{
		invoices.POST("/invoices", invCtrl.CreateInvoice)
		invoices.GET("/invoices", invCtrl.ListInvoices)
		invoices.GET("/invoices/summary", invCtrl.InvoiceSummary)
		invoices.GET("/invoices/:id", invCtrl.GetInvoice)
		invoices.PATCH("/invoices/:id", invCtrl.UpdateInvoice)
		invoices.PUT("/invoices/:id", invCtrl.UpdateInvoice)
		invoices.DELETE("/invoices/:id", invCtrl.DeleteInvoice)
		invoices.POST("/invoices/:id/recalculate", invCtrl.RecalculateInvoice)
	}
}