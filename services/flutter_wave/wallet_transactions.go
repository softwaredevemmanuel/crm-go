// services/wallet_transactions.go
package services

import (
	"context"
	"fmt"

	"crm-go/models"

	"github.com/google/uuid"
)

// GetWalletTransactions lists transactions for a wallet with pagination
func (s *WalletService) GetWalletTransactions(ctx context.Context, walletID uuid.UUID, page, limit int, txType, status string) ([]models.WalletTransaction, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := s.db.WithContext(ctx).Model(&models.WalletTransaction{}).
		Where("wallet_id = ?", walletID)

	if txType != "" {
		query = query.Where("type = ?", txType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	var transactions []models.WalletTransaction
	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	return transactions, total, nil
}

// GetAllTransactions lists all transactions across all wallets
func (s *WalletService) GetAllTransactions(ctx context.Context, page, limit int, walletID, txType, status string) ([]models.WalletTransaction, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := s.db.WithContext(ctx).Model(&models.WalletTransaction{})

	if walletID != "" {
		if id, err := uuid.Parse(walletID); err == nil {
			query = query.Where("wallet_id = ?", id)
		}
	}
	if txType != "" {
		query = query.Where("type = ?", txType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	var transactions []models.WalletTransaction
	if err := query.
		Preload("Wallet").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	return transactions, total, nil
}