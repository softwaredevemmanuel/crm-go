// services/wallet_ledger.go
package services

import (
	"context"
	"fmt"

	"crm-go/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// debitWalletForDisbursement creates a ledger entry and updates wallet balance
func (s *WalletService) debitWalletForDisbursement(ctx context.Context, walletID uuid.UUID, amount int64, reference string, createdBy uuid.UUID) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&wallet, "id = ?", walletID).Error; err != nil {
			return fmt.Errorf("failed to lock wallet: %w", err)
		}

		if wallet.AvailableBalance < amount {
			return fmt.Errorf("insufficient balance")
		}

		balanceBefore := wallet.Balance
		balanceAfter := wallet.Balance - amount

		// Create transaction record
		transaction := models.WalletTransaction{
			WalletID:      walletID,
			Type:          models.WalletTransactionDebit,
			Amount:        amount,
			BalanceBefore: balanceBefore,
			BalanceAfter:  balanceAfter,
			Reference:     reference,
			Description:   "Payout disbursement",
			Status:        models.WalletTransactionPending,
			CreatedBy:     createdBy,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		// Update wallet balances
		if err := tx.Model(&wallet).Updates(map[string]interface{}{
			"balance":           balanceAfter,
			"available_balance": wallet.AvailableBalance - amount,
		}).Error; err != nil {
			return fmt.Errorf("failed to update wallet: %w", err)
		}

		return nil
	})
}

// creditWallet handles incoming funds
func (s *WalletService) CreditWallet(ctx context.Context, walletID uuid.UUID, amount int64, reference, description string, createdBy *uuid.UUID) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&wallet, "id = ?", walletID).Error; err != nil {
			return fmt.Errorf("failed to lock wallet: %w", err)
		}

		balanceBefore := wallet.Balance
		balanceAfter := wallet.Balance + amount

		transaction := models.WalletTransaction{
			WalletID:      walletID,
			Type:          models.WalletTransactionCredit,
			Amount:        amount,
			BalanceBefore: balanceBefore,
			BalanceAfter:  balanceAfter,
			Reference:     reference,
			Description:   description,
			Status:        models.WalletTransactionCompleted,
			CreatedBy:     *createdBy,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		if err := tx.Model(&wallet).Updates(map[string]interface{}{
			"balance":           balanceAfter,
			"available_balance": wallet.AvailableBalance + amount,
		}).Error; err != nil {
			return fmt.Errorf("failed to update wallet: %w", err)
		}

		return nil
	})
}