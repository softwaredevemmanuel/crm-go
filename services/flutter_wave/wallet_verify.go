// services/wallet_verify.go
package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/google/uuid"
	"crm-go/models"
)

// VerifyAccount resolves and verifies a bank account before creating a beneficiary
func (s *WalletService) VerifyAccount(ctx context.Context, accountNumber, bankCode string) (*AccountResolutionResponse, error) {
	payload := AccountResolutionRequest{
		AccountNumber: accountNumber,
		AccountBank:   bankCode,
		Country:       "NG",
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/accounts/resolve", s.baseURL), bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.secretKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result AccountResolutionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("account resolution failed: %s", result.Message)
	}

	return &result, nil
}

// CreateBeneficiary verifies and stores a beneficiary
func (s *WalletService) CreateBeneficiary(ctx context.Context, accountName, accountNumber, bankCode, bankName string, createdBy uuid.UUID) (*models.Beneficiary, error) {
	// Verify account first
	resolved, err := s.VerifyAccount(ctx, accountNumber, bankCode)
	if err != nil {
		return nil, fmt.Errorf("account verification failed: %w", err)
	}

	beneficiary := &models.Beneficiary{
		AccountName:   resolved.Data.AccountName,
		AccountNumber: accountNumber,
		BankCode:      bankCode,
		BankName:      bankName,
		IsVerified:    true,
		CreatedBy:     createdBy,
	}

	if err := s.db.Create(beneficiary).Error; err != nil {
		return nil, fmt.Errorf("failed to save beneficiary: %w", err)
	}

	return beneficiary, nil
}