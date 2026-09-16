// routes/wallet_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/flutter_wave"
	"crm-go/middleware"
	"crm-go/services/flutter_wave"
)

func WalletRoutes(router *gin.RouterGroup, db *gorm.DB, flutterwaveSecretKey string) {
	walletService := services.NewWalletService(db, flutterwaveSecretKey)
	walletController := controllers.NewWalletController(walletService)

	walletGroup := router.Group("/api")
	walletGroup.Use(middleware.AuthMiddleware())
	{
		// Disbursement
		walletGroup.POST("/wallets/:wallet_id/disburse", walletController.Disburse)

		// Account verification
		walletGroup.POST("/wallets/verify-account", walletController.VerifyAccount)

		// Beneficiaries
		walletGroup.POST("/wallets/beneficiaries", walletController.CreateBeneficiary)
	}

	// Webhook (no auth middleware - verified via signature)
	router.POST("/api/webhooks/flutterwave", walletController.Webhook)
}