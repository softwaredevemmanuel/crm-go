// services/wallet_verify.go
package services

import (
	"bytes"
	"context"
	"crm-go/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
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

	// 🔍 CRITICAL: Log the raw response and status code
	log.Printf("Flutterwave VerifyAccount → status: %d, body: %s", resp.StatusCode, string(body))
log.Printf("VerifyAccount → baseURL=%s, key_prefix=%s, key_length=%d",
    s.baseURL,
    s.secretKey[:min(20, len(s.secretKey))], // log first 20 chars only
    len(s.secretKey),
)
	// 🚨 If it's not a 2xx, return the raw error immediately
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Flutterwave returned %d: %s", resp.StatusCode, string(body))
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
		AccountBank:   bankCode,
		BankName:      bankName,
		IsVerified:    true,
		CreatedBy:     createdBy,
	}

	if err := s.db.Create(beneficiary).Error; err != nil {
		return nil, fmt.Errorf("failed to save beneficiary: %w", err)
	}

	return beneficiary, nil
}
