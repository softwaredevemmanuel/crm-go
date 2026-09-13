// routes/inventory_request_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/inventory_request"
	"crm-go/middleware"
	"crm-go/services/inventory_request"
)

func InventoryRequestRoutes(router *gin.RouterGroup, db *gorm.DB) {
	requestService := services.NewInventoryRequestService(db)
	requestController := controllers.NewInventoryRequestController(requestService)

	requestGroup := router.Group("/api/inventory-requests")
	requestGroup.Use(middleware.AuthMiddleware())
	{
		// CREATE
		requestGroup.POST("/", requestController.CreateInventoryRequest)

		// READ
		requestGroup.GET("/", requestController.GetAllInventoryRequests)
		requestGroup.GET("/stats", requestController.GetInventoryRequestStats)
		requestGroup.GET("/my-requests", requestController.GetMyInventoryRequests)
		requestGroup.GET("/:id", requestController.GetInventoryRequestByID)

		// UPDATE
		requestGroup.PUT("/:id", requestController.UpdateInventoryRequest)

		// WORKFLOW ACTIONS
		requestGroup.POST("/:id/approve", requestController.ApproveInventoryRequest)
		requestGroup.POST("/:id/reject", requestController.RejectInventoryRequest)
		requestGroup.POST("/:id/fulfill", requestController.FulfillInventoryRequest)
		requestGroup.POST("/:id/cancel", requestController.CancelInventoryRequest)

		// DELETE
		requestGroup.DELETE("/:id", requestController.DeleteInventoryRequest)
	}
}