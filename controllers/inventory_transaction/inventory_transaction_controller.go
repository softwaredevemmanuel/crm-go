// controllers/inventory_transaction_controller.go
package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/inventory_transaction"
)

type InventoryTransactionController struct {
	transactionService *services.InventoryTransactionService
}

func NewInventoryTransactionController(transactionService *services.InventoryTransactionService) *InventoryTransactionController {
	return &InventoryTransactionController{
		transactionService: transactionService,
	}
}

// CreateInventoryTransaction creates a new transaction
// @Summary Create inventory transaction
// @Description Create a new inventory transaction (stock movement)
// @Tags Inventory Transactions
// @Accept json
// @Produce json
// @Param request body dto.CreateInventoryTransactionRequest true "Transaction data"
// @Success 201 {object} map[string]interface{} "Transaction created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or insufficient stock"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Item or section not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-transactions [post]
func (c *InventoryTransactionController) CreateInventoryTransaction(ctx *gin.Context) {
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dto.CreateInventoryTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	transaction, err := c.transactionService.CreateInventoryTransaction(userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "insufficient") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Inventory transaction created successfully",
		"data":    transaction,
	})
}

// GetAllInventoryTransactions retrieves all transactions
// @Summary Get all inventory transactions
// @Description Get a paginated list of inventory transactions with filters
// @Tags Inventory Transactions
// @Accept json
// @Produce json
// @Param item_id query string false "Filter by Item ID"
// @Param inventory_section_id query string false "Filter by Section ID"
// @Param transaction_type query string false "Filter by transaction type"
// @Param status query string false "Filter by status" Enums(pending, approved, rejected, completed, cancelled)
// @Param issued_to_id query string false "Filter by recipient"
// @Param issued_by_id query string false "Filter by issuer"
// @Param approved_by_id query string false "Filter by approver"
// @Param class_id query string false "Filter by class"
// @Param search query string false "Search by code, reference, purpose, or notes"
// @Param date_from query string false "Filter by date from (YYYY-MM-DD)"
// @Param date_to query string false "Filter by date to (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param sort_by query string false "Sort field" default(transaction_date)
// @Param sort_order query string false "Sort direction" default(desc) Enums(asc, desc)
// @Success 200 {object} map[string]interface{} "Transactions retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-transactions [get]
func (c *InventoryTransactionController) GetAllInventoryTransactions(ctx *gin.Context) {
	var params dto.InventoryTransactionQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := c.transactionService.GetInventoryTransactions(&params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Transactions retrieved successfully",
		"data":    response,
	})
}

// GetInventoryTransactionByID retrieves a transaction by ID
// @Summary Get inventory transaction by ID
// @Description Get a specific inventory transaction by its ID
// @Tags Inventory Transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} map[string]interface{} "Transaction retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Transaction ID is required"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-transactions/{id} [get]
func (c *InventoryTransactionController) GetInventoryTransactionByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Transaction ID is required"})
		return
	}

	transaction, err := c.transactionService.GetInventoryTransactionByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "invalid") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Transaction retrieved successfully",
		"data":    transaction,
	})
}

// GetTransactionsByItem retrieves all transactions for an item
// @Summary Get transactions by item
// @Description Get all transactions for a specific inventory item
// @Tags Inventory Transactions
// @Accept json
// @Produce json
// @Param item_id path string true "Item ID"
// @Success 200 {object} map[string]interface{} "Transactions retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Item ID is required"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-transactions/item/{item_id} [get]
func (c *InventoryTransactionController) GetTransactionsByItem(ctx *gin.Context) {
	itemID := ctx.Param("item_id")
	if itemID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Item ID is required"})
		return
	}

	transactions, err := c.transactionService.GetTransactionsByItem(itemID)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Transactions retrieved successfully",
		"data":    transactions,
	})
}

// UpdateInventoryTransaction updates an existing transaction
// @Summary Update inventory transaction
// @Description Update a pending transaction (only by creator)
// @Tags Inventory Transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Param request body dto.UpdateInventoryTransactionRequest true "Update data"
// @Success 200 {object} map[string]interface{} "Transaction updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or not pending"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "You can only update your own transactions"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-transactions/{id} [put]
func (c *InventoryTransactionController) UpdateInventoryTransaction(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Transaction ID is required"})
		return
	}

	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dto.UpdateInventoryTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	transaction, err := c.transactionService.UpdateInventoryTransaction(id, userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "only update your own") {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "only pending") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Transaction updated successfully",
		"data":    transaction,
	})
}

// ApproveInventoryTransaction approves a transaction
// @Summary Approve inventory transaction
// @Description Approve a pending inventory transaction (admin only)
// @Tags Inventory Transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} map[string]interface{} "Transaction approved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or not pending"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-transactions/{id}/approve [post]
func (c *InventoryTransactionController) ApproveInventoryTransaction(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Transaction ID is required"})
		return
	}

	approverIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	approverID, err := uuid.Parse(approverIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	transaction, err := c.transactionService.ApproveInventoryTransaction(id, approverID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "only pending") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Transaction approved successfully",
		"data":    transaction,
	})
}

// CompleteInventoryTransaction completes a transaction
// @Summary Complete inventory transaction
// @Description Complete an approved transaction (e.g., mark item as returned)
// @Tags Inventory Transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Param request body dto.CompleteTransactionRequest true "Completion data"
// @Success 200 {object} map[string]interface{} "Transaction completed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or status"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-transactions/{id}/complete [post]
func (c *InventoryTransactionController) CompleteInventoryTransaction(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Transaction ID is required"})
		return
	}

	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dto.CompleteTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	transaction, err := c.transactionService.CompleteInventoryTransaction(id, userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "only approved") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Transaction completed successfully",
		"data":    transaction,
	})
}

// RejectInventoryTransaction rejects a transaction
// @Summary Reject inventory transaction
// @Description Reject a pending inventory transaction (admin only)
// @Tags Inventory Transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} map[string]interface{} "Transaction rejected successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or not pending"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-transactions/{id}/reject [post]
func (c *InventoryTransactionController) RejectInventoryTransaction(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Transaction ID is required"})
		return
	}

	rejectorIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	rejectorID, err := uuid.Parse(rejectorIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	transaction, err := c.transactionService.RejectInventoryTransaction(id, rejectorID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "only pending") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Transaction rejected successfully",
		"data":    transaction,
	})
}

// CancelInventoryTransaction cancels a transaction
// @Summary Cancel inventory transaction
// @Description Cancel a pending transaction (by creator)
// @Tags Inventory Transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} map[string]interface{} "Transaction cancelled successfully"
// @Failure 400 {object} map[string]interface{} "Cannot cancel this transaction"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "You can only cancel your own transactions"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-transactions/{id}/cancel [post]
func (c *InventoryTransactionController) CancelInventoryTransaction(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Transaction ID is required"})
		return
	}

	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := c.transactionService.CancelInventoryTransaction(id, userID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "only cancel your own") {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Transaction cancelled successfully",
	})
}

// DeleteInventoryTransaction deletes a transaction
// @Summary Delete inventory transaction
// @Description Soft delete an inventory transaction
// @Tags Inventory Transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} map[string]interface{} "Transaction deleted successfully"
// @Failure 400 {object} map[string]interface{} "Transaction ID is required"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-transactions/{id} [delete]
func (c *InventoryTransactionController) DeleteInventoryTransaction(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Transaction ID is required"})
		return
	}

	if err := c.transactionService.DeleteInventoryTransaction(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Transaction deleted successfully",
	})
}

// GetInventoryTransactionStats retrieves statistics for transactions
// @Summary Get inventory transaction statistics
// @Description Get statistics for inventory transactions
// @Tags Inventory Transactions
// @Accept json
// @Produce json
// @Param section_id query string false "Filter by section ID"
// @Success 200 {object} map[string]interface{} "Statistics retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid section ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-transactions/stats [get]
func (c *InventoryTransactionController) GetInventoryTransactionStats(ctx *gin.Context) {
	sectionID := ctx.Query("section_id")

	stats, err := c.transactionService.GetInventoryTransactionStats(sectionID)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Statistics retrieved successfully",
		"data":    stats,
	})
}