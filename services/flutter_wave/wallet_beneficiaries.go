// services/wallet_beneficiaries.go
package services

import (
	"context"
	"errors"
	"fmt"

	"crm-go/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetBeneficiaryByID fetches a beneficiary
func (s *WalletService) GetBeneficiaryByID(ctx context.Context, id uuid.UUID) (*models.Beneficiary, error) {
	var b models.Beneficiary
	if err := s.db.WithContext(ctx).First(&b, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("beneficiary not found")
		}
		return nil, fmt.Errorf("failed to fetch beneficiary: %w", err)
	}
	return &b, nil
}

// GetAllBeneficiaries lists beneficiaries with pagination
func (s *WalletService) GetAllBeneficiaries(ctx context.Context, page, limit int, search string) ([]models.Beneficiary, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := s.db.WithContext(ctx).Model(&models.Beneficiary{})

	if search != "" {
		query = query.Where(
			"account_name ILIKE ? OR account_number ILIKE ? OR bank_name ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%",
		)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count beneficiaries: %w", err)
	}

	var beneficiaries []models.Beneficiary
	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&beneficiaries).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch beneficiaries: %w", err)
	}

	return beneficiaries, total, nil
}

// DeleteBeneficiary removes a beneficiary
func (s *WalletService) DeleteBeneficiary(ctx context.Context, id uuid.UUID) error {
	beneficiary, err := s.GetBeneficiaryByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if there are pending disbursements using this beneficiary
	var pendingCount int64
	if err := s.db.WithContext(ctx).
		Model(&models.WalletDisbursement{}).
		Where("beneficiary_id = ? AND status IN ?", id, []string{"pending", "processing"}).
		Count(&pendingCount).Error; err != nil {
		return fmt.Errorf("failed to check pending disbursements: %w", err)
	}

	if pendingCount > 0 {
		return errors.New("cannot delete beneficiary with pending disbursements")
	}

	if err := s.db.WithContext(ctx).Delete(beneficiary).Error; err != nil {
		return fmt.Errorf("failed to delete beneficiary: %w", err)
	}

	return nil
}