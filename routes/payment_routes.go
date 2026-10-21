package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"crm-go/middleware"
	"crm-go/controllers/payment"
	"crm-go/services/payment"
)

func PaymentRoutes(router *gin.RouterGroup, db *gorm.DB) {
	svc := services.NewPaymentService(db)
	ctrl := controllers.NewPaymentController(svc)

	// Nested under invoices
	invoices := router.Group("/api")
	invoices.Use(middleware.AuthMiddleware())
	{
		invoices.POST("/invoices/:id/payments", ctrl.CreatePayment)
		invoices.GET("/invoices/:id/payments", ctrl.ListInvoicePayments)
	}

	// Global payments group
	payments := router.Group("/api")
	payments.Use(middleware.AuthMiddleware())
	{
		payments.GET("/payments", ctrl.ListPayments)
		payments.GET("/payments/:id", ctrl.GetPayment)
		payments.PATCH("/payments/:id", ctrl.UpdatePayment)
		payments.PUT("/payments/:id", ctrl.UpdatePayment)
		payments.DELETE("/payments/:id", ctrl.DeletePayment)
	}
}