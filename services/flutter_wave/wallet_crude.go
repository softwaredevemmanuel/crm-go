// services/wallet_crud.go
package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"crm-go/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateWallet creates a new wallet
func (s *WalletService) CreateWallet(ctx context.Context, name, currency string) (*models.Wallet, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("wallet name is required")
	}
	if currency == "" {
		currency = "NGN"
	}

	wallet := &models.Wallet{
		Name:     strings.TrimSpace(name),
		Currency: strings.ToUpper(currency),
		Status:   models.WalletStatusActive,
		Balance:  0,
	}

	if err := s.db.WithContext(ctx).Create(wallet).Error; err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	return wallet, nil
}

// GetWalletByID fetches a single wallet by ID
func (s *WalletService) GetWalletByID(ctx context.Context, id uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := s.db.WithContext(ctx).First(&wallet, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet not found")
		}
		return nil, fmt.Errorf("failed to fetch wallet: %w", err)
	}
	return &wallet, nil
}

// GetAllWallets lists wallets with pagination and optional filters
func (s *WalletService) GetAllWallets(ctx context.Context, page, limit int, status, search string) ([]models.Wallet, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := s.db.WithContext(ctx).Model(&models.Wallet{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count wallets: %w", err)
	}

	var wallets []models.Wallet
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&wallets).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch wallets: %w", err)
	}

	return wallets, total, nil
}

// UpdateWallet updates allowed wallet fields (name, status)
func (s *WalletService) UpdateWallet(ctx context.Context, id uuid.UUID, name, status string) (*models.Wallet, error) {
	wallet, err := s.GetWalletByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if strings.TrimSpace(name) != "" {
		updates["name"] = strings.TrimSpace(name)
	}
	if status != "" {
		// Validate status
		valid := map[string]bool{
			string(models.WalletStatusActive):   true,
			string(models.WalletStatusInactive): true,
			string(models.WalletStatusFrozen):   true,
		}
		if !valid[status] {
			return nil, errors.New("invalid wallet status")
		}
		updates["status"] = status
	}

	if len(updates) == 0 {
		return wallet, nil
	}

	if err := s.db.WithContext(ctx).Model(wallet).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update wallet: %w", err)
	}

	return s.GetWalletByID(ctx, id)
}

// DeleteWallet soft-deletes a wallet (only if balance is zero)
func (s *WalletService) DeleteWallet(ctx context.Context, id uuid.UUID) error {
	wallet, err := s.GetWalletByID(ctx, id)
	if err != nil {
		return err
	}

	if wallet.Balance != 0 {
		return errors.New("cannot delete wallet with non-zero balance")
	}

	if err := s.db.WithContext(ctx).Delete(wallet).Error; err != nil {
		return fmt.Errorf("failed to delete wallet: %w", err)
	}

	return nil
}

// GetWalletStats returns aggregate stats for all wallets
func (s *WalletService) GetWalletStats(ctx context.Context) (map[string]interface{}, error) {
	type Result struct {
		TotalWallets  int64
		ActiveWallets int64
		TotalBalance  int64
	}

	var r Result
	err := s.db.WithContext(ctx).
		Model(&models.Wallet{}).
		Select(`
			COUNT(*) as total_wallets,
			COUNT(*) FILTER (WHERE status = 'active') as active_wallets,
			COALESCE(SUM(balance), 0) as total_balance
		`).
		Scan(&r).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet stats: %w", err)
	}

	return map[string]interface{}{
		"total_wallets":  r.TotalWallets,
		"active_wallets": r.ActiveWallets,
		"total_balance":  r.TotalBalance,
	}, nil
}