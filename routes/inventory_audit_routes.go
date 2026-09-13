// routes/inventory_audit_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/inventory_audit"
	"crm-go/middleware"
	"crm-go/services/inventory_audit"
)

func InventoryAuditRoutes(router *gin.RouterGroup, db *gorm.DB) {
	auditService := services.NewInventoryAuditService(db)
	auditController := controllers.NewInventoryAuditController(auditService)

	auditGroup := router.Group("/api")
	auditGroup.Use(middleware.AuthMiddleware())
	{
		// CREATE
		auditGroup.POST("/inventory-audits", auditController.CreateInventoryAudit)

		// READ
		auditGroup.GET("/inventory-audits", auditController.GetAllInventoryAudits)
		auditGroup.GET("/inventory-audits/stats", auditController.GetInventoryAuditStats)
		auditGroup.GET("/inventory-audits/:id", auditController.GetInventoryAuditByID)

		// UPDATE
		auditGroup.PUT("/inventory-audits/:id", auditController.UpdateInventoryAudit)

		// AUDIT ITEMS
		auditGroup.POST("/inventory-audits/:id/items", auditController.AddAuditItem)
		auditGroup.PUT("/inventory-audits/:id/items/:item_id", auditController.UpdateAuditItem)
		auditGroup.DELETE("/inventory-audits/:id/items/:item_id", auditController.RemoveAuditItem)

		// WORKFLOW ACTIONS
		auditGroup.POST("/inventory-audits/:id/complete", auditController.CompleteInventoryAudit)
		auditGroup.POST("/inventory-audits/:id/approve", auditController.ApproveInventoryAudit)
		auditGroup.POST("/inventory-audits/:id/reject", auditController.RejectInventoryAudit)
		auditGroup.POST("/inventory-audits/:id/cancel", auditController.CancelInventoryAudit)

		// DELETE
		auditGroup.DELETE("/inventory-audits/:id", auditController.DeleteInventoryAudit)
	}
}