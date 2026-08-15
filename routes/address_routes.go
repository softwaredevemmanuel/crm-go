// routes/address_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/address"
	"crm-go/services/address"
	"crm-go/middleware"
)

func AddressRoutes(router *gin.RouterGroup, db *gorm.DB) {
	addressService := services.NewAddressService(db)
	addressHandler := handlers.NewAddressHandler(addressService)

	addressGroup := router.Group("/api")
	addressGroup.Use(middleware.AuthMiddleware())
	{
		// Create address
		addressGroup.POST("/addresses", addressHandler.CreateAddress)

		// Get all addresses with pagination and filters
		addressGroup.GET("/addresses", addressHandler.GetAllAddresses)

		// Get addresses by user
		addressGroup.GET("/addresses/user/:user_id", addressHandler.GetAddressesByUser)

		// Get primary address by user
		addressGroup.GET("/addresses/user/:user_id/primary", addressHandler.GetPrimaryAddressByUser)

		// Get address by ID
		addressGroup.GET("/addresses/:id", addressHandler.GetAddressByID)

		// Update address
		addressGroup.PUT("/addresses/:id", addressHandler.UpdateAddress)

		// Delete address
		addressGroup.DELETE("/addresses/:id", addressHandler.DeleteAddress)
	}
}