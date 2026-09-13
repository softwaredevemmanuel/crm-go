// routes/inventory_transaction_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/inventory_transaction"
	"crm-go/middleware"
	"crm-go/services/inventory_transaction"
)

func InventoryTransactionRoutes(router *gin.RouterGroup, db *gorm.DB) {
	transactionService := services.NewInventoryTransactionService(db)
	transactionController := controllers.NewInventoryTransactionController(transactionService)

	transactionGroup := router.Group("/api")
	transactionGroup.Use(middleware.AuthMiddleware())
	{
		// CREATE
		transactionGroup.POST("/inventory-transactions", transactionController.CreateInventoryTransaction)

		// READ
		transactionGroup.GET("/inventory-transactions", transactionController.GetAllInventoryTransactions)
		transactionGroup.GET("/inventory-transactions/stats", transactionController.GetInventoryTransactionStats)
		transactionGroup.GET("/inventory-transactions/item/:item_id", transactionController.GetTransactionsByItem)
		transactionGroup.GET("/inventory-transactions/:id", transactionController.GetInventoryTransactionByID)

		// UPDATE
		transactionGroup.PUT("/inventory-transactions/:id", transactionController.UpdateInventoryTransaction)

		// WORKFLOW ACTIONS
		transactionGroup.POST("/inventory-transactions/:id/approve", transactionController.ApproveInventoryTransaction)
		transactionGroup.POST("/inventory-transactions/:id/reject", transactionController.RejectInventoryTransaction)
		transactionGroup.POST("/inventory-transactions/:id/complete", transactionController.CompleteInventoryTransaction)
		transactionGroup.POST("/inventory-transactions/:id/cancel", transactionController.CancelInventoryTransaction)

		// DELETE
		transactionGroup.DELETE("/inventory-transactions/:id", transactionController.DeleteInventoryTransaction)
	}
}