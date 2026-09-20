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

	// Public webhook (no auth — verified by signature)
	router.POST("/api/webhooks/flutterwave", walletController.Webhook)

	// Authenticated routes
	walletGroup := router.Group("/api")
	walletGroup.Use(middleware.AuthMiddleware())
	{
		// ---------- WALLETS ----------
		walletGroup.POST("/wallets", walletController.CreateWallet)
		walletGroup.GET("/wallets", walletController.GetAllWallets)
		walletGroup.GET("/wallets/stats", walletController.GetWalletStats)

		// Beneficiary-specific routes MUST come before /:id to avoid conflicts
		walletGroup.POST("/wallets/verify-account", walletController.VerifyAccount)
		walletGroup.POST("/wallets/beneficiaries", walletController.CreateBeneficiary)
		walletGroup.GET("/wallets/beneficiaries", walletController.GetAllBeneficiaries)
		walletGroup.GET("/wallets/beneficiaries/:id", walletController.GetBeneficiaryByID)
		walletGroup.DELETE("/wallets/beneficiaries/:id", walletController.DeleteBeneficiary)

		// Wallet by ID (wildcard — must be last)
		walletGroup.GET("/wallets/:id", walletController.GetWalletByID)
		walletGroup.PUT("/wallets/:id", walletController.UpdateWallet)
		walletGroup.DELETE("/wallets/:id", walletController.DeleteWallet)

		// Wallet-scoped transactions
		walletGroup.GET("/wallets/:id/transactions", walletController.GetWalletTransactions)

		// Disbursement
		walletGroup.POST("/wallets/:id/disburse", walletController.Disburse)
	}

	// Global transactions & disbursements
	adminGroup := router.Group("/api")
	adminGroup.Use(middleware.AuthMiddleware())
	{
		adminGroup.GET("/transactions", walletController.GetAllTransactions)
		adminGroup.GET("/disbursements", walletController.GetAllDisbursements)
		adminGroup.GET("/disbursements/stats", walletController.GetDisbursementStats)
		adminGroup.GET("/disbursements/:id", walletController.GetDisbursementByID)
	}
}