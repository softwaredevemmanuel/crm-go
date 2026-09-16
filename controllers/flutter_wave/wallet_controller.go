// controllers/wallet_controller.go
package controllers

import (
	"net/http"
	"crm-go/services/flutter_wave"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WalletController struct {
	walletService *services.WalletService
}

func NewWalletController(walletService *services.WalletService) *WalletController {
	return &WalletController{walletService: walletService}
}

// Disburse handles POST /api/wallets/:wallet_id/disburse
func (ctrl *WalletController) Disburse(c *gin.Context) {
	walletID, err := uuid.Parse(c.Param("wallet_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet ID"})
		return
	}

	var req struct {
		BeneficiaryID uuid.UUID `json:"beneficiary_id" binding:"required"`
		Amount        int64     `json:"amount" binding:"required,gt=0"`
		Reason        string    `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := uuid.Parse(c.GetString("user_id"))

	disbursement, err := ctrl.walletService.Disburse(c.Request.Context(), walletID, req.BeneficiaryID, req.Amount, userID, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    disbursement,
	})
}

// VerifyAccount handles POST /api/wallets/verify-account
func (ctrl *WalletController) VerifyAccount(c *gin.Context) {
	var req struct {
		AccountNumber string `json:"account_number" binding:"required"`
		BankCode      string `json:"bank_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := ctrl.walletService.VerifyAccount(c.Request.Context(), req.AccountNumber, req.BankCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result.Data,
	})
}

// CreateBeneficiary handles POST /api/wallets/beneficiaries
func (ctrl *WalletController) CreateBeneficiary(c *gin.Context) {
	var req struct {
		AccountNumber string `json:"account_number" binding:"required"`
		BankCode      string `json:"bank_code" binding:"required"`
		BankName      string `json:"bank_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := uuid.Parse(c.GetString("user_id"))

	beneficiary, err := ctrl.walletService.CreateBeneficiary(c.Request.Context(), "", req.AccountNumber, req.BankCode, req.BankName, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    beneficiary,
	})
}

// Webhook handles POST /api/webhooks/flutterwave
func (ctrl *WalletController) Webhook(c *gin.Context) {
	secretHash := c.GetString("flutterwave_webhook_hash") // Set via env or config

	if err := ctrl.walletService.HandleWebhook(c, secretHash); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Always return 200 to prevent Flutterwave retries
	c.Status(http.StatusOK)
}