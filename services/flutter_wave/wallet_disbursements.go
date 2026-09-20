// services/wallet_disbursements.go
package services

import (
	"context"
	"errors"
	"fmt"

	"crm-go/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetDisbursementByID fetches a single disbursement with relations
func (s *WalletService) GetDisbursementByID(ctx context.Context, id uuid.UUID) (*models.WalletDisbursement, error) {
	var d models.WalletDisbursement
	if err := s.db.WithContext(ctx).
		Preload("Wallet").
		Preload("Beneficiary").
		First(&d, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("disbursement not found")
		}
		return nil, fmt.Errorf("failed to fetch disbursement: %w", err)
	}
	return &d, nil
}

// GetAllDisbursements lists disbursements with filters
func (s *WalletService) GetAllDisbursements(ctx context.Context, page, limit int, walletID, status string) ([]models.WalletDisbursement, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := s.db.WithContext(ctx).Model(&models.WalletDisbursement{})

	if walletID != "" {
		if id, err := uuid.Parse(walletID); err == nil {
			query = query.Where("wallet_id = ?", id)
		}
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count disbursements: %w", err)
	}

	var disbursements []models.WalletDisbursement
	if err := query.
		Preload("Wallet").
		Preload("Beneficiary").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&disbursements).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch disbursements: %w", err)
	}

	return disbursements, total, nil
}

// GetDisbursementStats returns aggregate stats
func (s *WalletService) GetDisbursementStats(ctx context.Context) (map[string]interface{}, error) {
	type Result struct {
		Total        int64
		Pending      int64
		Processing   int64
		Successful   int64
		Failed       int64
		TotalAmount  int64
		TotalFees    int64
	}

	var r Result
	err := s.db.WithContext(ctx).
		Model(&models.WalletDisbursement{}).
		Select(`
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'pending') as pending,
			COUNT(*) FILTER (WHERE status = 'processing') as processing,
			COUNT(*) FILTER (WHERE status = 'successful') as successful,
			COUNT(*) FILTER (WHERE status = 'failed') as failed,
			COALESCE(SUM(amount) FILTER (WHERE status = 'successful'), 0) as total_amount,
			COALESCE(SUM(fee) FILTER (WHERE status = 'successful'), 0) as total_fees
		`).
		Scan(&r).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get disbursement stats: %w", err)
	}

	return map[string]interface{}{
		"total":        r.Total,
		"pending":      r.Pending,
		"processing":   r.Processing,
		"successful":   r.Successful,
		"failed":       r.Failed,
		"total_amount": r.TotalAmount,
		"total_fees":   r.TotalFees,
	}, nil
}