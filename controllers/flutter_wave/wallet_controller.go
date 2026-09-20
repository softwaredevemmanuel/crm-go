// controllers/wallet_controller.go
package controllers

import (
	"net/http"
	"strconv"

	_ "crm-go/docs" // swagger docs (generated)
	services "crm-go/services/flutter_wave"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WalletController struct {
	walletService *services.WalletService
}

func NewWalletController(walletService *services.WalletService) *WalletController {
	return &WalletController{walletService: walletService}
}

// ---------- REQUEST / RESPONSE TYPES ----------

// CreateWalletRequest represents the payload to create a wallet
type CreateWalletRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=255" example:"Main Operations Wallet"`
	Currency string `json:"currency" example:"NGN"`
}

// UpdateWalletRequest represents the payload to update a wallet
type UpdateWalletRequest struct {
	Name   string `json:"name" example:"Renamed Wallet"`
	Status string `json:"status" example:"active" Enums(active,inactive,frozen)`
}

// DisburseRequest represents the payload to disburse funds
type DisburseRequest struct {
	BeneficiaryID uuid.UUID `json:"beneficiary_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440003"`
	Amount        int64     `json:"amount" binding:"required,gt=0" example:"5000"`
	Reason        string    `json:"reason" example:"Vendor payment"`
}

// CreateBeneficiaryRequest represents the payload to create a beneficiary
type CreateBeneficiaryRequest struct {
	AccountNumber string `json:"account_number" binding:"required" example:"0690000031"`
	AccountBank   string `json:"account_bank" binding:"required" example:"044"`
	BankName      string `json:"bank_name" example:"Access Bank"`
}

// VerifyAccountRequest represents the payload to verify a bank account
type VerifyAccountRequest struct {
	AccountNumber string `json:"account_number" binding:"required" example:"0690000031"`
	AccountBank   string `json:"account_bank" binding:"required" example:"044"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page  int   `json:"page" example:"1"`
	Limit int   `json:"limit" example:"20"`
	Total int64 `json:"total" example:"42"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"wallet not found"`
}

// MessageResponse represents a simple success message
type MessageResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Wallet deleted successfully"`
}

// ---------- WALLETS ----------

// CreateWallet godoc
// @Summary      Create a new wallet
// @Description  Creates a new wallet with a name and currency. Balance starts at 0.
// @Tags         Wallets
// @Accept       json
// @Produce      json
// @Param        request  body      CreateWalletRequest  true  "Wallet creation payload"
// @Success      201      {object}  map[string]interface{}  "Wallet created"
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /api/wallets [post]
func (ctrl *WalletController) CreateWallet(c *gin.Context) {
	var req CreateWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	wallet, err := ctrl.walletService.CreateWallet(c.Request.Context(), req.Name, req.Currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": wallet})
}

// GetAllWallets godoc
// @Summary      List all wallets
// @Description  Returns a paginated list of wallets. Supports filtering by status and searching by name.
// @Tags         Wallets
// @Accept       json
// @Produce      json
// @Param        page    query     int     false  "Page number"      default(1)
// @Param        limit   query     int     false  "Items per page"   default(20)
// @Param        status  query     string  false  "Wallet status"    Enums(active,inactive,frozen)
// @Param        search  query     string  false  "Search by name"
// @Success      200  {object}  map[string]interface{}  "List of wallets"
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /api/wallets [get]
func (ctrl *WalletController) GetAllWallets(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")
	search := c.Query("search")

	wallets, total, err := ctrl.walletService.GetAllWallets(c.Request.Context(), page, limit, status, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"wallets":    wallets,
			"pagination": Pagination{Page: page, Limit: limit, Total: total},
		},
	})
}

// GetWalletByID godoc
// @Summary      Get wallet by ID
// @Description  Returns a single wallet with its details and current balance.
// @Tags         Wallets
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Wallet ID (UUID)"  Format(uuid)
// @Success      200  {object}  map[string]interface{}  "Wallet details"
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /api/wallets/{id} [get]
func (ctrl *WalletController) GetWalletByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid wallet ID"})
		return
	}

	wallet, err := ctrl.walletService.GetWalletByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": wallet})
}

// UpdateWallet godoc
// @Summary      Update a wallet
// @Description  Updates the wallet's name and/or status.
// @Tags         Wallets
// @Accept       json
// @Produce      json
// @Param        id       path      string               true  "Wallet ID (UUID)"  Format(uuid)
// @Param        request  body      UpdateWalletRequest  true  "Wallet update payload"
// @Success      200      {object}  map[string]interface{}  "Updated wallet"
// @Failure      400      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /api/wallets/{id} [put]
func (ctrl *WalletController) UpdateWallet(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid wallet ID"})
		return
	}

	var req UpdateWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	wallet, err := ctrl.walletService.UpdateWallet(c.Request.Context(), id, req.Name, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": wallet})
}

// DeleteWallet godoc
// @Summary      Delete a wallet
// @Description  Soft-deletes a wallet. Only allowed if wallet balance is zero.
// @Tags         Wallets
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Wallet ID (UUID)"  Format(uuid)
// @Success      200  {object}  MessageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /api/wallets/{id} [delete]
func (ctrl *WalletController) DeleteWallet(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid wallet ID"})
		return
	}

	if err := ctrl.walletService.DeleteWallet(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Wallet deleted successfully"})
}

// GetWalletStats godoc
// @Summary      Wallet aggregate stats
// @Description  Returns total wallets, active wallets, and total balance across all wallets.
// @Tags         Wallets
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Wallet stats"
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /api/wallets/stats [get]
func (ctrl *WalletController) GetWalletStats(c *gin.Context) {
	stats, err := ctrl.walletService.GetWalletStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}

// ---------- TRANSACTIONS ----------

// GetWalletTransactions godoc
// @Summary      List transactions for a wallet
// @Description  Returns paginated transactions belonging to a specific wallet.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        id      path      string  true   "Wallet ID (UUID)"  Format(uuid)
// @Param        page    query     int     false  "Page number"       default(1)
// @Param        limit   query     int     false  "Items per page"    default(20)
// @Param        type    query     string  false  "Transaction type"  Enums(credit,debit)
// @Param        status  query     string  false  "Transaction status" Enums(pending,completed,failed,reversed)
// @Success      200  {object}  map[string]interface{}  "Transaction list"
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /api/wallets/{id}/transactions [get]
func (ctrl *WalletController) GetWalletTransactions(c *gin.Context) {
	walletID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid wallet ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	txType := c.Query("type")
	status := c.Query("status")

	txns, total, err := ctrl.walletService.GetWalletTransactions(c.Request.Context(), walletID, page, limit, txType, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"transactions": txns,
			"pagination":   Pagination{Page: page, Limit: limit, Total: total},
		},
	})
}

// GetAllTransactions godoc
// @Summary      List all transactions
// @Description  Returns paginated transactions across all wallets. Filterable by wallet, type, and status.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "Page number"       default(1)
// @Param        limit      query     int     false  "Items per page"    default(20)
// @Param        wallet_id  query     string  false  "Filter by wallet ID"  Format(uuid)
// @Param        type       query     string  false  "Transaction type"  Enums(credit,debit)
// @Param        status     query     string  false  "Transaction status" Enums(pending,completed,failed,reversed)
// @Success      200  {object}  map[string]interface{}  "Transaction list"
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /api/transactions [get]
func (ctrl *WalletController) GetAllTransactions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	walletID := c.Query("wallet_id")
	txType := c.Query("type")
	status := c.Query("status")

	txns, total, err := ctrl.walletService.GetAllTransactions(c.Request.Context(), page, limit, walletID, txType, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"transactions": txns,
			"pagination":   Pagination{Page: page, Limit: limit, Total: total},
		},
	})
}

// ---------- DISBURSEMENTS ----------

// Disburse godoc
// @Summary      Disburse funds to a beneficiary
// @Description  Sends money from a wallet to a verified beneficiary's bank account via Flutterwave.
// @Tags         Disbursements
// @Accept       json
// @Produce      json
// @Param        id       path      string            true  "Wallet ID (UUID)"  Format(uuid)
// @Param        request  body      DisburseRequest   true  "Disbursement payload"
// @Success      200      {object}  map[string]interface{}  "Disbursement created"
// @Failure      400      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /api/wallets/{id}/disburse [post]
func (ctrl *WalletController) Disburse(c *gin.Context) {
	walletID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid wallet ID"})
		return
	}

	var req DisburseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	userID, _ := uuid.Parse(c.GetString("user_id"))

	disbursement, err := ctrl.walletService.Disburse(
		c.Request.Context(),
		walletID,
		req.BeneficiaryID,
		req.Amount,
		userID,
		req.Reason,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": disbursement})
}

// GetAllDisbursements godoc
// @Summary      List all disbursements
// @Description  Returns paginated disbursements across all wallets. Filterable by wallet and status.
// @Tags         Disbursements
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "Page number"       default(1)
// @Param        limit      query     int     false  "Items per page"    default(20)
// @Param        wallet_id  query     string  false  "Filter by wallet ID"  Format(uuid)
// @Param        status     query     string  false  "Disbursement status"  Enums(pending,processing,successful,failed,reversed)
// @Success      200  {object}  map[string]interface{}  "Disbursement list"
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /api/disbursements [get]
func (ctrl *WalletController) GetAllDisbursements(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	walletID := c.Query("wallet_id")
	status := c.Query("status")

	disbursements, total, err := ctrl.walletService.GetAllDisbursements(c.Request.Context(), page, limit, walletID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"disbursements": disbursements,
			"pagination":    Pagination{Page: page, Limit: limit, Total: total},
		},
	})
}

// GetDisbursementByID godoc
// @Summary      Get disbursement by ID
// @Description  Returns a single disbursement with related wallet and beneficiary details.
// @Tags         Disbursements
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Disbursement ID (UUID)"  Format(uuid)
// @Success      200  {object}  map[string]interface{}  "Disbursement details"
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /api/disbursements/{id} [get]
func (ctrl *WalletController) GetDisbursementByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid disbursement ID"})
		return
	}

	d, err := ctrl.walletService.GetDisbursementByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": d})
}

// GetDisbursementStats godoc
// @Summary      Disbursement aggregate stats
// @Description  Returns totals by status plus total amount and fees for successful disbursements.
// @Tags         Disbursements
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Disbursement stats"
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /api/disbursements/stats [get]
func (ctrl *WalletController) GetDisbursementStats(c *gin.Context) {
	stats, err := ctrl.walletService.GetDisbursementStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}

// ---------- BENEFICIARIES ----------

// CreateBeneficiary godoc
// @Summary      Create a beneficiary
// @Description  Verifies the bank account via Flutterwave and stores it as a beneficiary.
// @Tags         Beneficiaries
// @Accept       json
// @Produce      json
// @Param        request  body      CreateBeneficiaryRequest  true  "Beneficiary payload"
// @Success      201      {object}  map[string]interface{}  "Beneficiary created"
// @Failure      400      {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /api/wallets/beneficiaries [post]
func (ctrl *WalletController) CreateBeneficiary(c *gin.Context) {
	var req CreateBeneficiaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	userID, _ := uuid.Parse(c.GetString("user_id"))

	beneficiary, err := ctrl.walletService.CreateBeneficiary(
		c.Request.Context(),
		"",
		req.AccountNumber,
		req.AccountBank,
		req.BankName,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": beneficiary})
}

// GetAllBeneficiaries godoc
// @Summary      List all beneficiaries
// @Description  Returns paginated beneficiaries. Supports search by name, account number, or bank.
// @Tags         Beneficiaries
// @Accept       json
// @Produce      json
// @Param        page    query     int     false  "Page number"     default(1)
// @Param        limit   query     int     false  "Items per page"  default(20)
// @Param        search  query     string  false  "Search term"
// @Success      200  {object}  map[string]interface{}  "Beneficiary list"
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /api/wallets/beneficiaries [get]
func (ctrl *WalletController) GetAllBeneficiaries(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")

	beneficiaries, total, err := ctrl.walletService.GetAllBeneficiaries(c.Request.Context(), page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"beneficiaries": beneficiaries,
			"pagination":    Pagination{Page: page, Limit: limit, Total: total},
		},
	})
}

// GetBeneficiaryByID godoc
// @Summary      Get beneficiary by ID
// @Description  Returns a single beneficiary.
// @Tags         Beneficiaries
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Beneficiary ID (UUID)"  Format(uuid)
// @Success      200  {object}  map[string]interface{}  "Beneficiary details"
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /api/wallets/beneficiaries/{id} [get]
func (ctrl *WalletController) GetBeneficiaryByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid beneficiary ID"})
		return
	}

	b, err := ctrl.walletService.GetBeneficiaryByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": b})
}

// DeleteBeneficiary godoc
// @Summary      Delete a beneficiary
// @Description  Deletes a beneficiary. Fails if there are pending disbursements.
// @Tags         Beneficiaries
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Beneficiary ID (UUID)"  Format(uuid)
// @Success      200  {object}  MessageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /api/wallets/beneficiaries/{id} [delete]
func (ctrl *WalletController) DeleteBeneficiary(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid beneficiary ID"})
		return
	}

	if err := ctrl.walletService.DeleteBeneficiary(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Beneficiary deleted successfully"})
}

// VerifyAccount godoc
// @Summary      Verify a bank account
// @Description  Resolves a bank account via Flutterwave and returns the account name.
// @Tags         Beneficiaries
// @Accept       json
// @Produce      json
// @Param        request  body      VerifyAccountRequest  true  "Account verification payload"
// @Success      200      {object}  map[string]interface{}  "Account details"
// @Failure      400      {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /api/wallets/verify-account [post]
func (ctrl *WalletController) VerifyAccount(c *gin.Context) {
	var req VerifyAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	result, err := ctrl.walletService.VerifyAccount(c.Request.Context(), req.AccountNumber, req.AccountBank)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result.Data})
}

// Webhook godoc
// @Summary      Flutterwave webhook receiver
// @Description  Receives transfer status updates from Flutterwave. Verified via verif-hash header. Not authenticated by JWT.
// @Tags         Webhooks
// @Accept       json
// @Produce      json
// @Success      200  {string}  string  "OK"
// @Failure      401  {object}  ErrorResponse
// @Router       /api/webhooks/flutterwave [post]
func (ctrl *WalletController) Webhook(c *gin.Context) {
	secretHash := c.GetString("FLUTTERWAVE_WEBHOOK_HASH")

	if err := ctrl.walletService.HandleWebhook(c, secretHash); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
