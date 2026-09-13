// services/inventory_transaction_service.go
package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"crm-go/dto"
	"crm-go/models"
)

type InventoryTransactionService struct {
	db *gorm.DB
}

func NewInventoryTransactionService(db *gorm.DB) *InventoryTransactionService {
	return &InventoryTransactionService{db: db}
}

// generateTransactionCode generates a unique transaction code
func (s *InventoryTransactionService) generateTransactionCode() (string, error) {
	prefix := "TXN-" + time.Now().Format("20060102") + "-"

	var count int64
	today := time.Now().Format("2006-01-02")
	s.db.Model(&models.InventoryTransaction{}).
		Where("DATE(created_at) = ?", today).
		Count(&count)

	return fmt.Sprintf("%s%04d", prefix, count+1), nil
}

// CreateInventoryTransaction creates a new transaction
func (s *InventoryTransactionService) CreateInventoryTransaction(userID uuid.UUID, req *dto.CreateInventoryTransactionRequest) (*dto.InventoryTransactionResponse, error) {
	// Parse IDs
	itemID, err := uuid.Parse(req.ItemID)
	if err != nil {
		return nil, errors.New("invalid item ID")
	}

	sectionID, err := uuid.Parse(req.InventorySectionID)
	if err != nil {
		return nil, errors.New("invalid inventory section ID")
	}

	// Verify item exists
	var item models.InventoryItem
	if err := s.db.Where("id = ? AND deleted_at IS NULL", itemID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory item not found")
		}
		return nil, fmt.Errorf("failed to verify item: %w", err)
	}

	// Verify section exists
	var section models.InventorySection
	if err := s.db.Where("id = ? AND deleted_at IS NULL", sectionID).First(&section).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory section not found")
		}
		return nil, fmt.Errorf("failed to verify section: %w", err)
	}

	// Verify item belongs to section
	if item.InventorySectionID != sectionID {
		return nil, errors.New("item does not belong to the specified section")
	}

	// Parse optional IDs
	var fromSectionID, toSectionID, issuedToID, issuedByID, approvedByID, receivedByID, supervisorID, classID *uuid.UUID

	if req.FromSectionID != "" {
		id, err := uuid.Parse(req.FromSectionID)
		if err != nil {
			return nil, errors.New("invalid from section ID")
		}
		fromSectionID = &id
	}

	if req.ToSectionID != "" {
		id, err := uuid.Parse(req.ToSectionID)
		if err != nil {
			return nil, errors.New("invalid to section ID")
		}
		toSectionID = &id
	}

	if req.IssuedToID != "" {
		id, err := uuid.Parse(req.IssuedToID)
		if err != nil {
			return nil, errors.New("invalid issued to ID")
		}
		issuedToID = &id
	}

	if req.IssuedByID != "" {
		id, err := uuid.Parse(req.IssuedByID)
		if err != nil {
			return nil, errors.New("invalid issued by ID")
		}
		issuedByID = &id
	}

	if req.ApprovedByID != "" {
		id, err := uuid.Parse(req.ApprovedByID)
		if err != nil {
			return nil, errors.New("invalid approved by ID")
		}
		approvedByID = &id
	}

	if req.ReceivedByID != "" {
		id, err := uuid.Parse(req.ReceivedByID)
		if err != nil {
			return nil, errors.New("invalid received by ID")
		}
		receivedByID = &id
	}

	if req.SupervisorID != "" {
		id, err := uuid.Parse(req.SupervisorID)
		if err != nil {
			return nil, errors.New("invalid supervisor ID")
		}
		supervisorID = &id
	}

	if req.ClassID != "" {
		id, err := uuid.Parse(req.ClassID)
		if err != nil {
			return nil, errors.New("invalid class ID")
		}
		classID = &id
	}

	// Parse reference type
	var referenceType *models.InventoryReferenceType
	if req.ReferenceType != "" {
		rt := models.InventoryReferenceType(req.ReferenceType)
		referenceType = &rt
	}

	// Parse conditions
	var conditionBefore, conditionAfter *models.InventoryItemCondition
	if req.ConditionBefore != "" {
		cb := models.InventoryItemCondition(req.ConditionBefore)
		conditionBefore = &cb
	}
	if req.ConditionAfter != "" {
		ca := models.InventoryItemCondition(req.ConditionAfter)
		conditionAfter = &ca
	}

	// Parse dates
	transactionDate := time.Now()
	if req.TransactionDate != "" {
		parsed, err := time.Parse(time.RFC3339, req.TransactionDate)
		if err != nil {
			// Try simple date format
			parsed, err = time.Parse("2006-01-02", req.TransactionDate)
			if err != nil {
				return nil, errors.New("invalid transaction date format")
			}
		}
		transactionDate = parsed
	}

	var expectedReturnDate *time.Time
	if req.ExpectedReturnDate != "" {
		parsed, err := time.Parse("2006-01-02", req.ExpectedReturnDate)
		if err != nil {
			return nil, errors.New("invalid expected return date format")
		}
		expectedReturnDate = &parsed
	}

	// Set defaults
	status := req.Status
	if status == "" {
		status = "pending"
	}

	// Generate transaction code
	transactionCode, err := s.generateTransactionCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate transaction code: %w", err)
	}

	// Begin database transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Record current quantity
	quantityBefore := item.QuantityAvailable
	quantityAfter := quantityBefore

	// Apply stock changes based on transaction type
	quantityAfter = s.calculateQuantityAfter(item, req.TransactionType, req.Quantity)

	// Update item stock
	if err := s.updateItemStock(tx, &item, req.TransactionType, req.Quantity); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create transaction record
	transaction := &models.InventoryTransaction{
		ID:                 uuid.New(),
		TransactionCode:    transactionCode,
		ItemID:             itemID,
		InventorySectionID: sectionID,
		TransactionType:    models.InventoryTransactionType(req.TransactionType),
		Quantity:           req.Quantity,
		QuantityBefore:     quantityBefore,
		QuantityAfter:      quantityAfter,
		UnitCost:           req.UnitCost,
		TotalCost:          float64(req.Quantity) * req.UnitCost,
		ReferenceNumber:    req.ReferenceNumber,
		ReferenceType:      referenceType,
		FromLocation:       req.FromLocation,
		ToLocation:         req.ToLocation,
		FromSectionID:      fromSectionID,
		ToSectionID:        toSectionID,
		IssuedToID:         issuedToID,
		IssuedByID:         issuedByID,
		ApprovedByID:       approvedByID,
		ReceivedByID:       receivedByID,
		SupervisorID:       supervisorID,
		ClassID:            classID,
		Purpose:            req.Purpose,
		ExperimentName:     req.ExperimentName,
		StudentCount:       req.StudentCount,
		Notes:              req.Notes,
		Status:             models.InventoryTransactionStatus(status),
		TransactionDate:    transactionDate,
		ExpectedReturnDate: expectedReturnDate,
		ConditionBefore:    conditionBefore,
		ConditionAfter:     conditionAfter,
		AttachmentURL:      req.AttachmentURL,
		CreatedBy:          userID,
	}

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	tx.Commit()

	// Reload with relationships
	if err := s.db.Preload("Item").
		Preload("InventorySection").
		Preload("FromSection").
		Preload("ToSection").
		Preload("IssuedTo").
		Preload("IssuedBy").
		Preload("ApprovedBy").
		Preload("ReceivedBy").
		Preload("Supervisor").
		Preload("Creator").
		First(transaction, transaction.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load transaction details: %w", err)
	}

	return s.toTransactionResponse(transaction), nil
}

// calculateQuantityAfter calculates the quantity after a transaction
func (s *InventoryTransactionService) calculateQuantityAfter(item models.InventoryItem, txnType string, quantity int) int {
	switch txnType {
	// Transactions that INCREASE stock
	case "purchase", "donation", "return", "repair_return":
		return item.QuantityAvailable + quantity

	// Transactions that DECREASE stock
	case "issue", "usage", "damage", "loss", "disposal", "expiry", "transfer", "borrow", "experiment", "repair":
		return item.QuantityAvailable - quantity

	// Adjustments - direct set
	case "adjustment", "audit":
		return quantity

	default:
		return item.QuantityAvailable
	}
}

// updateItemStock updates the item's stock based on transaction type
func (s *InventoryTransactionService) updateItemStock(tx *gorm.DB, item *models.InventoryItem, txnType string, quantity int) error {
	switch txnType {
	// INCREASE stock
	case "purchase", "donation", "return", "repair_return":
		item.QuantityTotal += quantity
		item.QuantityAvailable += quantity

	// DECREASE available
	case "issue", "usage", "experiment", "borrow", "repair", "transfer":
		if item.QuantityAvailable < quantity {
			return fmt.Errorf("insufficient stock: available %d, requested %d", item.QuantityAvailable, quantity)
		}
		item.QuantityAvailable -= quantity
		item.QuantityInUse += quantity

	// DECREASE permanently
	case "damage", "loss", "disposal", "expiry":
		if item.QuantityAvailable < quantity {
			return fmt.Errorf("insufficient stock: available %d, requested %d", item.QuantityAvailable, quantity)
		}
		item.QuantityAvailable -= quantity
		if txnType == "damage" {
			item.QuantityDamaged += quantity
		}
		item.QuantityTotal -= quantity

	// ADJUSTMENT - direct set
	case "adjustment", "audit":
		item.QuantityAvailable = quantity
		item.QuantityTotal = quantity

	default:
		return fmt.Errorf("unsupported transaction type: %s", txnType)
	}

	// Auto-update status
	if item.QuantityAvailable <= 0 {
		item.Status = models.InventoryItemStatusOutOfStock
	} else if item.Status == models.InventoryItemStatusOutOfStock {
		item.Status = models.InventoryItemStatusActive
	}

	return tx.Save(item).Error
}

// GetInventoryTransactionByID retrieves a transaction by ID
func (s *InventoryTransactionService) GetInventoryTransactionByID(id string) (*dto.InventoryTransactionResponse, error) {
	transactionID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid transaction ID")
	}

	var transaction models.InventoryTransaction
	if err := s.db.Where("id = ? AND deleted_at IS NULL", transactionID).
		Preload("Item").
		Preload("InventorySection").
		Preload("FromSection").
		Preload("ToSection").
		Preload("IssuedTo").
		Preload("IssuedBy").
		Preload("ApprovedBy").
		Preload("ReceivedBy").
		Preload("Supervisor").
		Preload("Creator").
		First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory transaction not found")
		}
		return nil, fmt.Errorf("failed to fetch transaction: %w", err)
	}

	return s.toTransactionResponse(&transaction), nil
}

// GetInventoryTransactions retrieves transactions with filters and pagination
func (s *InventoryTransactionService) GetInventoryTransactions(params *dto.InventoryTransactionQueryParams) (*dto.InventoryTransactionListResponse, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = "transaction_date"
	}
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}

	// Build query
	query := s.db.Model(&models.InventoryTransaction{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.ItemID != "" {
		itemID, err := uuid.Parse(params.ItemID)
		if err == nil {
			query = query.Where("item_id = ?", itemID)
		}
	}

	if params.InventorySectionID != "" {
		sectionID, err := uuid.Parse(params.InventorySectionID)
		if err == nil {
			query = query.Where("inventory_section_id = ?", sectionID)
		}
	}

	if params.TransactionType != "" {
		query = query.Where("transaction_type = ?", params.TransactionType)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.IssuedToID != "" {
		id, err := uuid.Parse(params.IssuedToID)
		if err == nil {
			query = query.Where("issued_to_id = ?", id)
		}
	}

	if params.IssuedByID != "" {
		id, err := uuid.Parse(params.IssuedByID)
		if err == nil {
			query = query.Where("issued_by_id = ?", id)
		}
	}

	if params.ApprovedByID != "" {
		id, err := uuid.Parse(params.ApprovedByID)
		if err == nil {
			query = query.Where("approved_by_id = ?", id)
		}
	}

	if params.ClassID != "" {
		id, err := uuid.Parse(params.ClassID)
		if err == nil {
			query = query.Where("class_id = ?", id)
		}
	}

	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("transaction_code ILIKE ? OR reference_number ILIKE ? OR purpose ILIKE ? OR notes ILIKE ? OR experiment_name ILIKE ?", search, search, search, search, search)
	}

	if params.DateFrom != "" {
		dateFrom, err := time.Parse("2006-01-02", params.DateFrom)
		if err == nil {
			query = query.Where("transaction_date >= ?", dateFrom)
		}
	}

	if params.DateTo != "" {
		dateTo, err := time.Parse("2006-01-02", params.DateTo)
		if err == nil {
			query = query.Where("transaction_date <= ?", dateTo.Add(24*time.Hour))
		}
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count transactions: %w", err)
	}

	// Apply sorting
	sortDirection := "DESC"
	if strings.ToLower(params.SortOrder) == "asc" {
		sortDirection = "ASC"
	}
	query = query.Order(fmt.Sprintf("%s %s", params.SortBy, sortDirection))

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	// Execute query
	var transactions []models.InventoryTransaction
	if err := query.
		Preload("Item").
		Preload("InventorySection").
		Preload("FromSection").
		Preload("ToSection").
		Preload("IssuedTo").
		Preload("IssuedBy").
		Preload("ApprovedBy").
		Preload("ReceivedBy").
		Preload("Supervisor").
		Preload("Creator").
		Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	// Convert to response
	responses := make([]dto.InventoryTransactionResponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = *s.toTransactionResponse(&transaction)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.InventoryTransactionListResponse{
		Transactions: responses,
		Total:        total,
		Page:         params.Page,
		Limit:        params.Limit,
		TotalPages:   totalPages,
	}, nil
}

// GetTransactionsByItem retrieves all transactions for a specific item
func (s *InventoryTransactionService) GetTransactionsByItem(itemID string) ([]dto.InventoryTransactionResponse, error) {
	iID, err := uuid.Parse(itemID)
	if err != nil {
		return nil, errors.New("invalid item ID")
	}

	var transactions []models.InventoryTransaction
	if err := s.db.Where("item_id = ? AND deleted_at IS NULL", iID).
		Preload("Item").
		Preload("InventorySection").
		Preload("IssuedTo").
		Preload("IssuedBy").
		Preload("ApprovedBy").
		Preload("Creator").
		Order("transaction_date DESC").
		Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	responses := make([]dto.InventoryTransactionResponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = *s.toTransactionResponse(&transaction)
	}

	return responses, nil
}

// UpdateInventoryTransaction updates an existing transaction (only pending)
func (s *InventoryTransactionService) UpdateInventoryTransaction(id string, userID uuid.UUID, req *dto.UpdateInventoryTransactionRequest) (*dto.InventoryTransactionResponse, error) {
	transactionID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid transaction ID")
	}

	var transaction models.InventoryTransaction
	if err := s.db.Where("id = ? AND deleted_at IS NULL", transactionID).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory transaction not found")
		}
		return nil, fmt.Errorf("failed to fetch transaction: %w", err)
	}

	// Only pending transactions can be updated
	if transaction.Status != models.InventoryTxnStatusPending {
		return nil, errors.New("only pending transactions can be updated")
	}

	// Only creator can update
	if transaction.CreatedBy != userID {
		return nil, errors.New("you can only update your own transactions")
	}

	// Update fields
	if req.Quantity != nil {
		transaction.Quantity = *req.Quantity
	}

	if req.UnitCost != nil {
		transaction.UnitCost = *req.UnitCost
		transaction.TotalCost = float64(transaction.Quantity) * *req.UnitCost
	}

	if req.ReferenceNumber != "" {
		transaction.ReferenceNumber = req.ReferenceNumber
	}

	if req.ReferenceType != "" {
		rt := models.InventoryReferenceType(req.ReferenceType)
		transaction.ReferenceType = &rt
	}

	if req.FromLocation != "" {
		transaction.FromLocation = req.FromLocation
	}

	if req.ToLocation != "" {
		transaction.ToLocation = req.ToLocation
	}

	if req.Purpose != "" {
		transaction.Purpose = req.Purpose
	}

	if req.ExperimentName != "" {
		transaction.ExperimentName = req.ExperimentName
	}

	if req.StudentCount != nil {
		transaction.StudentCount = *req.StudentCount
	}

	if req.Notes != "" {
		transaction.Notes = req.Notes
	}

	if req.Status != "" {
		transaction.Status = models.InventoryTransactionStatus(req.Status)
	}

	if req.TransactionDate != "" {
		parsed, err := time.Parse(time.RFC3339, req.TransactionDate)
		if err != nil {
			parsed, err = time.Parse("2006-01-02", req.TransactionDate)
			if err != nil {
				return nil, errors.New("invalid transaction date format")
			}
		}
		transaction.TransactionDate = parsed
	}

	if req.ExpectedReturnDate != "" {
		parsed, err := time.Parse("2006-01-02", req.ExpectedReturnDate)
		if err != nil {
			return nil, errors.New("invalid expected return date format")
		}
		transaction.ExpectedReturnDate = &parsed
	}

	if req.ConditionAfter != "" {
		ca := models.InventoryItemCondition(req.ConditionAfter)
		transaction.ConditionAfter = &ca
	}

	if req.AttachmentURL != "" {
		transaction.AttachmentURL = req.AttachmentURL
	}

	if err := s.db.Save(&transaction).Error; err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	// Reload with relationships
	if err := s.db.Preload("Item").
		Preload("InventorySection").
		Preload("FromSection").
		Preload("ToSection").
		Preload("IssuedTo").
		Preload("IssuedBy").
		Preload("ApprovedBy").
		Preload("ReceivedBy").
		Preload("Supervisor").
		Preload("Creator").
		First(&transaction, transaction.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load transaction details: %w", err)
	}

	return s.toTransactionResponse(&transaction), nil
}

// ApproveInventoryTransaction approves a pending transaction
func (s *InventoryTransactionService) ApproveInventoryTransaction(id string, approverID uuid.UUID) (*dto.InventoryTransactionResponse, error) {
	transactionID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid transaction ID")
	}

	var transaction models.InventoryTransaction
	if err := s.db.Where("id = ? AND deleted_at IS NULL", transactionID).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory transaction not found")
		}
		return nil, fmt.Errorf("failed to fetch transaction: %w", err)
	}

	// Only pending transactions can be approved
	if transaction.Status != models.InventoryTxnStatusPending {
		return nil, errors.New("only pending transactions can be approved")
	}

	transaction.Status = models.InventoryTxnStatusApproved
	transaction.ApprovedByID = &approverID

	if err := s.db.Save(&transaction).Error; err != nil {
		return nil, fmt.Errorf("failed to approve transaction: %w", err)
	}

	// Reload with relationships
	if err := s.db.Preload("Item").
		Preload("InventorySection").
		Preload("ApprovedBy").
		Preload("Creator").
		First(&transaction, transaction.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load transaction details: %w", err)
	}

	return s.toTransactionResponse(&transaction), nil
}

// CompleteInventoryTransaction completes an approved transaction
func (s *InventoryTransactionService) CompleteInventoryTransaction(id string, userID uuid.UUID, req *dto.CompleteTransactionRequest) (*dto.InventoryTransactionResponse, error) {
	transactionID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid transaction ID")
	}

	var transaction models.InventoryTransaction
	if err := s.db.Where("id = ? AND deleted_at IS NULL", transactionID).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory transaction not found")
		}
		return nil, fmt.Errorf("failed to fetch transaction: %w", err)
	}

	// Only approved or pending transactions can be completed
	if transaction.Status != models.InventoryTxnStatusApproved && transaction.Status != models.InventoryTxnStatusPending {
		return nil, errors.New("only approved or pending transactions can be completed")
	}

	// Begin database transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// For return transactions, restore the item stock
	if transaction.TransactionType == models.InventoryTxnReturn ||
		transaction.TransactionType == models.InventoryTxnRepairReturn {
		var item models.InventoryItem
		if err := tx.Where("id = ? AND deleted_at IS NULL", transaction.ItemID).First(&item).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("inventory item not found")
		}

		item.QuantityAvailable += transaction.Quantity
		if item.QuantityInUse >= transaction.Quantity {
			item.QuantityInUse -= transaction.Quantity
		} else {
			item.QuantityInUse = 0
		}

		if err := tx.Save(&item).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update item stock: %w", err)
		}
	}

	// Update transaction
	now := time.Now()
	transaction.Status = models.InventoryTxnStatusCompleted

	if req.ActualReturnDate != "" {
		parsed, err := time.Parse("2006-01-02", req.ActualReturnDate)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("invalid actual return date format")
		}
		transaction.ActualReturnDate = &parsed
	} else if transaction.ExpectedReturnDate != nil {
		// Set actual return date to now if not provided but there's an expected return
		transaction.ActualReturnDate = &now
	}

	if req.ConditionAfter != "" {
		ca := models.InventoryItemCondition(req.ConditionAfter)
		transaction.ConditionAfter = &ca
	}

	if req.Notes != "" {
		transaction.Notes = req.Notes
	}

	if req.AttachmentURL != "" {
		transaction.AttachmentURL = req.AttachmentURL
	}

	if err := tx.Save(&transaction).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to complete transaction: %w", err)
	}

	tx.Commit()

	// Reload with relationships
	if err := s.db.Preload("Item").
		Preload("InventorySection").
		Preload("FromSection").
		Preload("ToSection").
		Preload("IssuedTo").
		Preload("IssuedBy").
		Preload("ApprovedBy").
		Preload("ReceivedBy").
		Preload("Supervisor").
		Preload("Creator").
		First(&transaction, transaction.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load transaction details: %w", err)
	}

	return s.toTransactionResponse(&transaction), nil
}

// RejectInventoryTransaction rejects a pending transaction
func (s *InventoryTransactionService) RejectInventoryTransaction(id string, rejectorID uuid.UUID) (*dto.InventoryTransactionResponse, error) {
	transactionID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid transaction ID")
	}

	var transaction models.InventoryTransaction
	if err := s.db.Where("id = ? AND deleted_at IS NULL", transactionID).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory transaction not found")
		}
		return nil, fmt.Errorf("failed to fetch transaction: %w", err)
	}

	if transaction.Status != models.InventoryTxnStatusPending {
		return nil, errors.New("only pending transactions can be rejected")
	}

	transaction.Status = models.InventoryTxnStatusRejected
	transaction.ApprovedByID = &rejectorID

	if err := s.db.Save(&transaction).Error; err != nil {
		return nil, fmt.Errorf("failed to reject transaction: %w", err)
	}

	// Reload with relationships
	if err := s.db.Preload("Item").
		Preload("InventorySection").
		Preload("ApprovedBy").
		Preload("Creator").
		First(&transaction, transaction.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load transaction details: %w", err)
	}

	return s.toTransactionResponse(&transaction), nil
}

// CancelInventoryTransaction cancels a pending transaction
func (s *InventoryTransactionService) CancelInventoryTransaction(id string, userID uuid.UUID) error {
	transactionID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid transaction ID")
	}

	var transaction models.InventoryTransaction
	if err := s.db.Where("id = ? AND deleted_at IS NULL", transactionID).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("inventory transaction not found")
		}
		return fmt.Errorf("failed to fetch transaction: %w", err)
	}

	// Only pending transactions can be cancelled
	if transaction.Status != models.InventoryTxnStatusPending {
		return errors.New("only pending transactions can be cancelled")
	}

	// Only creator can cancel
	if transaction.CreatedBy != userID {
		return errors.New("you can only cancel your own transactions")
	}

	transaction.Status = models.InventoryTxnStatusCancelled

	if err := s.db.Save(&transaction).Error; err != nil {
		return fmt.Errorf("failed to cancel transaction: %w", err)
	}

	return nil
}

// DeleteInventoryTransaction soft deletes a transaction
func (s *InventoryTransactionService) DeleteInventoryTransaction(id string) error {
	transactionID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid transaction ID")
	}

	var transaction models.InventoryTransaction
	if err := s.db.Where("id = ? AND deleted_at IS NULL", transactionID).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("inventory transaction not found")
		}
		return fmt.Errorf("failed to fetch transaction: %w", err)
	}

	if err := s.db.Delete(&transaction).Error; err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	return nil
}

// GetInventoryTransactionStats retrieves statistics for inventory transactions
func (s *InventoryTransactionService) GetInventoryTransactionStats(sectionID string) (*dto.InventoryTransactionStats, error) {
	query := s.db.Model(&models.InventoryTransaction{}).Where("deleted_at IS NULL")

	if sectionID != "" {
		sID, err := uuid.Parse(sectionID)
		if err != nil {
			return nil, errors.New("invalid section ID")
		}
		query = query.Where("inventory_section_id = ?", sID)
	}

	var stats dto.InventoryTransactionStats

	// Total
	if err := query.Count(&stats.TotalTransactions).Error; err != nil {
		return nil, fmt.Errorf("failed to count transactions: %w", err)
	}

	// By status
	if err := query.Where("status = ?", "pending").Count(&stats.PendingTransactions).Error; err != nil {
		return nil, fmt.Errorf("failed to count pending: %w", err)
	}

	if err := query.Where("status = ?", "completed").Count(&stats.CompletedTransactions).Error; err != nil {
		return nil, fmt.Errorf("failed to count completed: %w", err)
	}

	if err := query.Where("status = ?", "cancelled").Count(&stats.CancelledTransactions).Error; err != nil {
		return nil, fmt.Errorf("failed to count cancelled: %w", err)
	}

	// Sum quantities in (purchase, donation, return, repair_return)
	inTypes := []string{"purchase", "donation", "return", "repair_return"}
	var quantityIn int64
	if err := query.Where("transaction_type IN ?", inTypes).
		Select("COALESCE(SUM(quantity), 0)").Scan(&quantityIn).Error; err != nil {
		return nil, fmt.Errorf("failed to sum quantities in: %w", err)
	}
	stats.TotalQuantityIn = quantityIn

	// Sum quantities out (issue, usage, damage, loss, disposal, expiry, transfer)
	outTypes := []string{"issue", "usage", "damage", "loss", "disposal", "expiry", "transfer", "experiment"}
	var quantityOut int64
	if err := query.Where("transaction_type IN ?", outTypes).
		Select("COALESCE(SUM(quantity), 0)").Scan(&quantityOut).Error; err != nil {
		return nil, fmt.Errorf("failed to sum quantities out: %w", err)
	}
	stats.TotalQuantityOut = quantityOut

	// Sum total cost in
	var valueIn float64
	if err := query.Where("transaction_type IN ?", inTypes).
		Select("COALESCE(SUM(total_cost), 0)").Scan(&valueIn).Error; err != nil {
		return nil, fmt.Errorf("failed to sum value in: %w", err)
	}
	stats.TotalValueIn = valueIn

	// Sum total cost out
	var valueOut float64
	if err := query.Where("transaction_type IN ?", outTypes).
		Select("COALESCE(SUM(total_cost), 0)").Scan(&valueOut).Error; err != nil {
		return nil, fmt.Errorf("failed to sum value out: %w", err)
	}
	stats.TotalValueOut = valueOut

	return &stats, nil
}

// toTransactionResponse converts model to response DTO
func (s *InventoryTransactionService) toTransactionResponse(transaction *models.InventoryTransaction) *dto.InventoryTransactionResponse {
	response := &dto.InventoryTransactionResponse{
		ID:                 transaction.ID.String(),
		TransactionCode:    transaction.TransactionCode,
		ItemID:             transaction.ItemID.String(),
		InventorySectionID: transaction.InventorySectionID.String(),
		TransactionType:    string(transaction.TransactionType),
		Quantity:           transaction.Quantity,
		QuantityBefore:     transaction.QuantityBefore,
		QuantityAfter:      transaction.QuantityAfter,
		UnitCost:           transaction.UnitCost,
		TotalCost:          transaction.TotalCost,
		ReferenceNumber:    transaction.ReferenceNumber,
		FromLocation:       transaction.FromLocation,
		ToLocation:         transaction.ToLocation,
		Purpose:            transaction.Purpose,
		ExperimentName:     transaction.ExperimentName,
		StudentCount:       transaction.StudentCount,
		Notes:              transaction.Notes,
		Status:             string(transaction.Status),
		TransactionDate:    transaction.TransactionDate,
		ExpectedReturnDate: transaction.ExpectedReturnDate,
		ActualReturnDate:   transaction.ActualReturnDate,
		AttachmentURL:      transaction.AttachmentURL,
		CreatedBy:          transaction.CreatedBy.String(),
		CreatedAt:          transaction.CreatedAt,
		UpdatedAt:          transaction.UpdatedAt,
	}

	if transaction.ReferenceType != nil {
		response.ReferenceType = string(*transaction.ReferenceType)
	}

	if transaction.FromSectionID != nil {
		response.FromSectionID = transaction.FromSectionID.String()
	}

	if transaction.ToSectionID != nil {
		response.ToSectionID = transaction.ToSectionID.String()
	}

	if transaction.IssuedToID != nil {
		response.IssuedToID = transaction.IssuedToID.String()
	}

	if transaction.IssuedByID != nil {
		response.IssuedByID = transaction.IssuedByID.String()
	}

	if transaction.ApprovedByID != nil {
		response.ApprovedByID = transaction.ApprovedByID.String()
	}

	if transaction.ReceivedByID != nil {
		response.ReceivedByID = transaction.ReceivedByID.String()
	}

	if transaction.SupervisorID != nil {
		response.SupervisorID = transaction.SupervisorID.String()
	}

	if transaction.ClassID != nil {
		response.ClassID = transaction.ClassID.String()
	}

	if transaction.ConditionBefore != nil {
		response.ConditionBefore = string(*transaction.ConditionBefore)
	}

	if transaction.ConditionAfter != nil {
		response.ConditionAfter = string(*transaction.ConditionAfter)
	}

	// Add item details
	if transaction.Item.ID != uuid.Nil {
		response.Item = &dto.InventoryItemResponse{
			ID:                transaction.Item.ID.String(),
			ItemCode:          transaction.Item.ItemCode,
			ItemName:          transaction.Item.ItemName,
			Description:       transaction.Item.Description,
			QuantityTotal:     transaction.Item.QuantityTotal,
			QuantityAvailable: transaction.Item.QuantityAvailable,
			Condition:         string(transaction.Item.Condition),
			Status:            string(transaction.Item.Status),
		}
	}

	// Add section details
	if transaction.InventorySection.ID != uuid.Nil {
		response.InventorySection = &dto.InventorySectionResponse{
			ID:          transaction.InventorySection.ID.String(),
			Name:        transaction.InventorySection.Name,
			Code:        transaction.InventorySection.Code,
			Description: transaction.InventorySection.Description,
			Location:    transaction.InventorySection.Location,
			Status:      string(transaction.InventorySection.Status),
		}
	}

	// Add user relationships
	if transaction.IssuedTo != nil && transaction.IssuedTo.ID != uuid.Nil {
		response.IssuedTo = &dto.UserResponse{
			ID:        transaction.IssuedTo.ID.String(),
			FirstName: transaction.IssuedTo.FirstName,
			LastName:  transaction.IssuedTo.LastName,
			Email:     transaction.IssuedTo.Email,
			Phone:     transaction.IssuedTo.Phone,
		}
	}

	if transaction.IssuedBy != nil && transaction.IssuedBy.ID != uuid.Nil {
		response.IssuedBy = &dto.UserResponse{
			ID:        transaction.IssuedBy.ID.String(),
			FirstName: transaction.IssuedBy.FirstName,
			LastName:  transaction.IssuedBy.LastName,
			Email:     transaction.IssuedBy.Email,
		}
	}

	if transaction.ApprovedBy != nil && transaction.ApprovedBy.ID != uuid.Nil {
		response.ApprovedBy = &dto.UserResponse{
			ID:        transaction.ApprovedBy.ID.String(),
			FirstName: transaction.ApprovedBy.FirstName,
			LastName:  transaction.ApprovedBy.LastName,
			Email:     transaction.ApprovedBy.Email,
		}
	}

	if transaction.ReceivedBy != nil && transaction.ReceivedBy.ID != uuid.Nil {
		response.ReceivedBy = &dto.UserResponse{
			ID:        transaction.ReceivedBy.ID.String(),
			FirstName: transaction.ReceivedBy.FirstName,
			LastName:  transaction.ReceivedBy.LastName,
			Email:     transaction.ReceivedBy.Email,
		}
	}

	if transaction.Supervisor != nil && transaction.Supervisor.ID != uuid.Nil {
		response.Supervisor = &dto.UserResponse{
			ID:        transaction.Supervisor.ID.String(),
			FirstName: transaction.Supervisor.FirstName,
			LastName:  transaction.Supervisor.LastName,
			Email:     transaction.Supervisor.Email,
		}
	}

	if transaction.Creator.ID != uuid.Nil {
		response.Creator = &dto.UserResponse{
			ID:        transaction.Creator.ID.String(),
			FirstName: transaction.Creator.FirstName,
			LastName:  transaction.Creator.LastName,
			Email:     transaction.Creator.Email,
		}
	}

	return response
}