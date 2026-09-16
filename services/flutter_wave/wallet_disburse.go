// services/wallet_disburse.go
package services

import (
	"bytes"
	"context"
	"crm-go/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// Disburse initiates a payout from a wallet to a beneficiary's bank account
func (s *WalletService) Disburse(ctx context.Context, walletID uuid.UUID, beneficiaryID uuid.UUID, amount int64, initiatedBy uuid.UUID, reason string) (*models.WalletDisbursement, error) {
	// 1. Fetch wallet
	var wallet models.Wallet
	if err := s.db.First(&wallet, "id = ?", walletID).Error; err != nil {
		return nil, fmt.Errorf("wallet not found: %w", err)
	}

	if wallet.Status != models.WalletStatusActive {
		return nil, fmt.Errorf("wallet is not active")
	}

	// 2. Fetch beneficiary
	var beneficiary models.Beneficiary
	if err := s.db.First(&beneficiary, "id = ?", beneficiaryID).Error; err != nil {
		return nil, fmt.Errorf("beneficiary not found: %w", err)
	}

	if !beneficiary.IsVerified {
		return nil, fmt.Errorf("beneficiary account is not verified")
	}

	// 3. Check balance
	if wallet.AvailableBalance < amount {
		return nil, fmt.Errorf("insufficient wallet balance")
	}

	// 4. Generate unique reference
	reference := fmt.Sprintf("WD-%s", strings.ToUpper(uuid.New().String()[:8]))

	// 5. Create disbursement record (pending)
	disbursement := &models.WalletDisbursement{
		WalletID:      walletID,
		BeneficiaryID: beneficiaryID,
		Amount:        amount,
		Fee:           0, // Will update after Flutterwave response
		TotalAmount:   amount,
		Reference:     reference,
		Status:        models.DisbursementPending,
		Reason:        reason,
		InitiatedBy:   initiatedBy,
	}

	if err := s.db.Create(disbursement).Error; err != nil {
		return nil, fmt.Errorf("failed to create disbursement record: %w", err)
	}

	// 6. Call Flutterwave API
	transferResp, err := s.initiateFlutterwaveTransfer(ctx, &beneficiary, amount, reference)
	if err != nil {
		// Mark as failed
		s.db.Model(disbursement).Updates(map[string]interface{}{
			"status": models.DisbursementFailed,
			"reason": fmt.Sprintf("Flutterwave error: %v", err),
		})
		return disbursement, fmt.Errorf("transfer failed: %w", err)
	}

	// 7. Update disbursement with provider reference
	feeAmount := int64(transferResp.Data.Fee * 100) // Convert to kobo
	totalAmount := amount + feeAmount

	s.db.Model(disbursement).Updates(map[string]interface{}{
		"status":             models.DisbursementProcessing,
		"provider_reference": fmt.Sprintf("%d", transferResp.Data.ID),
		"fee":                feeAmount,
		"total_amount":       totalAmount,
	})

	// 8. Debit wallet balance (hold funds)
	if err := s.debitWalletForDisbursement(ctx, walletID, totalAmount, reference, initiatedBy); err != nil {
		return disbursement, fmt.Errorf("failed to debit wallet: %w", err)
	}

	// Refresh disbursement
	s.db.First(disbursement, "id = ?", disbursement.ID)

	return disbursement, nil
}

// initiateFlutterwaveTransfer makes the API call to Flutterwave
func (s *WalletService) initiateFlutterwaveTransfer(ctx context.Context, beneficiary *models.Beneficiary, amount int64, reference string) (*TransferResponse, error) {
	payload := TransferRequest{
		AccountBank:   beneficiary.BankCode,
		AccountNumber: beneficiary.AccountNumber,
		Amount:        amount,
		Narration:     fmt.Sprintf("Payout to %s", beneficiary.AccountName),
		Currency:      "NGN",
		Reference:     reference,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/transfers", s.baseURL), bytes.NewBuffer(jsonBody))
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

	var transferResp TransferResponse
	if err := json.Unmarshal(body, &transferResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if transferResp.Status != "success" {
		return nil, fmt.Errorf("Flutterwave error: %s - %s", transferResp.Message, transferResp.Data.CompleteMessage)
	}

	return &transferResp, nil
}
