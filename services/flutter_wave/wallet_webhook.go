// services/wallet_webhook.go
package services

import (
	"crypto/subtle"
	"fmt"
	"log"
	"crm-go/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

)

type FlutterwaveWebhookPayload struct {
	Event string `json:"event"`
	Data  struct {
		ID              int64   `json:"id"`
		TxRef           string  `json:"tx_ref"`
		Reference       string  `json:"reference"`
		Amount          int64   `json:"amount"`
		Currency        string  `json:"currency"`
		Status          string  `json:"status"`
		CompleteMessage string  `json:"complete_message"`
		Fee             float64 `json:"fee"`
	} `json:"data"`
}

// HandleWebhook processes Flutterwave transfer webhooks
func (s *WalletService) HandleWebhook(c *gin.Context, webhookSecretHash string) error {
	// 1. Verify signature
	verifHash := c.GetHeader("verif-hash")
	if verifHash == "" {
		return fmt.Errorf("missing verif-hash header")
	}

	// Use constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare([]byte(verifHash), []byte(webhookSecretHash)) != 1 {
		return fmt.Errorf("invalid webhook signature")
	}

	// 2. Parse payload
	var payload FlutterwaveWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	// 3. Handle transfer events
	switch payload.Event {
	case "transfer.completed":
		return s.handleTransferCompleted(&payload)
	case "transfer.failed":
		return s.handleTransferFailed(&payload)
	case "transfer.reversal":
		return s.handleTransferReversal(&payload)
	default:
		log.Printf("Unhandled webhook event: %s", payload.Event)
	}

	return nil
}

func (s *WalletService) handleTransferCompleted(payload *FlutterwaveWebhookPayload) error {
	// Find disbursement by provider reference
	var disbursement models.WalletDisbursement
	if err := s.db.First(&disbursement, "provider_reference = ?", fmt.Sprintf("%d", payload.Data.ID)).Error; err != nil {
		return fmt.Errorf("disbursement not found: %w", err)
	}

	// Update status
	return s.db.Model(&disbursement).Updates(map[string]interface{}{
		"status": models.DisbursementSuccessful,
	}).Error
}

func (s *WalletService) handleTransferFailed(payload *FlutterwaveWebhookPayload) error {
	var disbursement models.WalletDisbursement
	if err := s.db.First(&disbursement, "provider_reference = ?", fmt.Sprintf("%d", payload.Data.ID)).Error; err != nil {
		return fmt.Errorf("disbursement not found: %w", err)
	}

	// Update status and reason
	if err := s.db.Model(&disbursement).Updates(map[string]interface{}{
		"status": models.DisbursementFailed,
		"reason": payload.Data.CompleteMessage,
	}).Error; err != nil {
		return err
	}

	// Reverse the wallet debit
	return s.reverseWalletDebit(&disbursement)
}

func (s *WalletService) reverseWalletDebit(disbursement *models.WalletDisbursement) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		if err := tx.First(&wallet, "id = ?", disbursement.WalletID).Error; err != nil {
			return err
		}

		// Create reversal transaction
		transaction := models.WalletTransaction{
			WalletID:      disbursement.WalletID,
			Type:          models.WalletTransactionCredit,
			Amount:        disbursement.TotalAmount,
			BalanceBefore: wallet.Balance,
			BalanceAfter:  wallet.Balance + disbursement.TotalAmount,
			Reference:     fmt.Sprintf("REV-%s", disbursement.Reference),
			Description:   fmt.Sprintf("Reversal for failed payout: %s", disbursement.Reference),
			Status:        models.WalletTransactionCompleted,
			CreatedBy:     disbursement.InitiatedBy,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		// Credit back the wallet
		return tx.Model(&wallet).Updates(map[string]interface{}{
			"balance":           wallet.Balance + disbursement.TotalAmount,
			"available_balance": wallet.AvailableBalance + disbursement.TotalAmount,
		}).Error
	})
}

func (s *WalletService) handleTransferReversal(payload *FlutterwaveWebhookPayload) error {
	// Handle reversal webhooks
	log.Printf("Transfer reversal received: %+v", payload)
	return nil
}