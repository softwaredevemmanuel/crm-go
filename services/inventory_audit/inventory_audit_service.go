// services/inventory_audit_service.go
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

type InventoryAuditService struct {
	db *gorm.DB
}

func NewInventoryAuditService(db *gorm.DB) *InventoryAuditService {
	return &InventoryAuditService{db: db}
}

// generateAuditCode generates a unique audit code
func (s *InventoryAuditService) generateAuditCode() (string, error) {
	prefix := "AUD-" + time.Now().Format("20060102") + "-"

	var count int64
	today := time.Now().Format("2006-01-02")
	s.db.Model(&models.InventoryAudit{}).
		Where("DATE(created_at) = ?", today).
		Count(&count)

	return fmt.Sprintf("%s%04d", prefix, count+1), nil
}

// CreateInventoryAudit creates a new audit with optional items
func (s *InventoryAuditService) CreateInventoryAudit(userID uuid.UUID, req *dto.CreateInventoryAuditRequest) (*dto.InventoryAuditResponse, error) {
	// Parse section ID
	sectionID, err := uuid.Parse(req.InventorySectionID)
	if err != nil {
		return nil, errors.New("invalid inventory section ID")
	}

	// Verify section exists
	var section models.InventorySection
	if err := s.db.Where("id = ? AND deleted_at IS NULL", sectionID).First(&section).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory section not found")
		}
		return nil, fmt.Errorf("failed to verify section: %w", err)
	}

	// Parse audit date
	auditDate, err := time.Parse("2006-01-02", req.AuditDate)
	if err != nil {
		return nil, errors.New("invalid audit date format. Use YYYY-MM-DD")
	}

	// Generate audit code
	auditCode, err := s.generateAuditCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate audit code: %w", err)
	}

	// Begin database transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Create audit
	audit := &models.InventoryAudit{
		ID:                 uuid.New(),
		AuditCode:          auditCode,
		InventorySectionID: sectionID,
		AuditDate:          auditDate,
		ConductedByID:      userID,
		Status:             models.InventoryAuditStatusInProgress,
		Notes:              req.Notes,
		CreatedBy:          userID,
	}

	if err := tx.Create(audit).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create audit: %w", err)
	}

	// Add audit items if provided
	if len(req.Items) > 0 {
		for _, itemInput := range req.Items {
			itemID, err := uuid.Parse(itemInput.ItemID)
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("invalid item ID: %s", itemInput.ItemID)
			}

			// Verify item exists and belongs to section
			var item models.InventoryItem
			if err := tx.Where("id = ? AND inventory_section_id = ? AND deleted_at IS NULL", itemID, sectionID).First(&item).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("item %s not found in section", itemInput.ItemID)
			}

			// Calculate difference
			difference := itemInput.PhysicalQuantity - item.QuantityAvailable

			// Determine flags
			isMissing := itemInput.IsMissing || difference < 0
			isExcess := itemInput.IsExcess || difference > 0
			isDamaged := itemInput.IsDamaged || itemInput.Condition == "damaged" || itemInput.Condition == "poor"

			condition := models.InventoryItemCondition(itemInput.Condition)
			if itemInput.Condition == "" {
				condition = models.InventoryItemCondition("good")
			}

			auditItem := &models.InventoryAuditItem{
				ID:               uuid.New(),
				AuditID:          audit.ID,
				ItemID:           itemID,
				SystemQuantity:   item.QuantityAvailable,
				PhysicalQuantity: itemInput.PhysicalQuantity,
				Difference:       difference,
				Condition:        condition,
				Location:         itemInput.Location,
				IsMissing:        isMissing,
				IsExcess:         isExcess,
				IsDamaged:        isDamaged,
				Notes:            itemInput.Notes,
				ActionTaken:      itemInput.ActionTaken,
				CreatedBy:        userID,
			}

			if err := tx.Create(auditItem).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to create audit item: %w", err)
			}
		}

		// Recalculate audit totals
		if err := s.recalculateAuditTotals(tx, audit.ID); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()

	// Reload with relationships
	if err := s.db.Preload("InventorySection").
		Preload("ConductedBy").
		Preload("ApprovedBy").
		Preload("Creator").
		Preload("AuditItems").
		Preload("AuditItems.Item").
		Preload("AuditItems.ActionTakenBy").
		First(audit, audit.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load audit details: %w", err)
	}

	return s.toAuditResponse(audit), nil
}

// recalculateAuditTotals recalculates totals for an audit
func (s *InventoryAuditService) recalculateAuditTotals(tx *gorm.DB, auditID uuid.UUID) error {
	var auditItems []models.InventoryAuditItem
	if err := tx.Where("audit_id = ? AND deleted_at IS NULL", auditID).Find(&auditItems).Error; err != nil {
		return fmt.Errorf("failed to fetch audit items: %w", err)
	}

	totalChecked := len(auditItems)
	totalDiscrepancies := 0
	totalMissing := 0
	totalExcess := 0
	totalDamaged := 0

	for _, item := range auditItems {
		if item.Difference != 0 {
			totalDiscrepancies++
		}
		if item.IsMissing {
			totalMissing++
		}
		if item.IsExcess {
			totalExcess++
		}
		if item.IsDamaged {
			totalDamaged++
		}
	}

	return tx.Model(&models.InventoryAudit{}).
		Where("id = ?", auditID).
		Updates(map[string]interface{}{
			"total_items_checked":  totalChecked,
			"total_discrepancies":  totalDiscrepancies,
			"total_items_missing":  totalMissing,
			"total_items_excess":   totalExcess,
			"total_items_damaged":  totalDamaged,
		}).Error
}

// GetInventoryAuditByID retrieves an audit by ID
func (s *InventoryAuditService) GetInventoryAuditByID(id string) (*dto.InventoryAuditResponse, error) {
	auditID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid audit ID")
	}

	var audit models.InventoryAudit
	if err := s.db.Where("id = ? AND deleted_at IS NULL", auditID).
		Preload("InventorySection").
		Preload("ConductedBy").
		Preload("ApprovedBy").
		Preload("Creator").
		Preload("AuditItems").
		Preload("AuditItems.Item").
		Preload("AuditItems.ActionTakenBy").
		Preload("AuditItems.Creator").
		First(&audit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory audit not found")
		}
		return nil, fmt.Errorf("failed to fetch audit: %w", err)
	}

	return s.toAuditResponse(&audit), nil
}

// GetInventoryAudits retrieves audits with filters and pagination
func (s *InventoryAuditService) GetInventoryAudits(params *dto.InventoryAuditQueryParams) (*dto.InventoryAuditListResponse, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}

	// Build query
	query := s.db.Model(&models.InventoryAudit{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.InventorySectionID != "" {
		sectionID, err := uuid.Parse(params.InventorySectionID)
		if err == nil {
			query = query.Where("inventory_section_id = ?", sectionID)
		}
	}

	if params.ConductedByID != "" {
		userID, err := uuid.Parse(params.ConductedByID)
		if err == nil {
			query = query.Where("conducted_by_id = ?", userID)
		}
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("audit_code ILIKE ? OR notes ILIKE ?", search, search)
	}

	if params.DateFrom != "" {
		dateFrom, err := time.Parse("2006-01-02", params.DateFrom)
		if err == nil {
			query = query.Where("audit_date >= ?", dateFrom)
		}
	}

	if params.DateTo != "" {
		dateTo, err := time.Parse("2006-01-02", params.DateTo)
		if err == nil {
			query = query.Where("audit_date <= ?", dateTo)
		}
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count audits: %w", err)
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
	var audits []models.InventoryAudit
	if err := query.
		Preload("InventorySection").
		Preload("ConductedBy").
		Preload("ApprovedBy").
		Preload("Creator").
		Preload("AuditItems").
		Find(&audits).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch audits: %w", err)
	}

	// Convert to response
	responses := make([]dto.InventoryAuditResponse, len(audits))
	for i, audit := range audits {
		responses[i] = *s.toAuditResponse(&audit)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.InventoryAuditListResponse{
		Audits:     responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// UpdateInventoryAudit updates an in-progress audit
func (s *InventoryAuditService) UpdateInventoryAudit(id string, userID uuid.UUID, req *dto.UpdateInventoryAuditRequest) (*dto.InventoryAuditResponse, error) {
	auditID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid audit ID")
	}

	var audit models.InventoryAudit
	if err := s.db.Where("id = ? AND deleted_at IS NULL", auditID).First(&audit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory audit not found")
		}
		return nil, fmt.Errorf("failed to fetch audit: %w", err)
	}

	// Only in-progress audits can be updated
	if audit.Status != models.InventoryAuditStatusInProgress {
		return nil, errors.New("only in-progress audits can be updated")
	}

	// Only conductor or creator can update
	if audit.ConductedByID != userID && audit.CreatedBy != userID {
		return nil, errors.New("you can only update your own audits")
	}

	if req.AuditDate != "" {
		parsed, err := time.Parse("2006-01-02", req.AuditDate)
		if err != nil {
			return nil, errors.New("invalid audit date format")
		}
		audit.AuditDate = parsed
	}

	if req.Notes != "" {
		audit.Notes = req.Notes
	}

	if err := s.db.Save(&audit).Error; err != nil {
		return nil, fmt.Errorf("failed to update audit: %w", err)
	}

	// Reload with relationships
	if err := s.db.Preload("InventorySection").
		Preload("ConductedBy").
		Preload("ApprovedBy").
		Preload("Creator").
		Preload("AuditItems").
		First(&audit, audit.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load audit details: %w", err)
	}

	return s.toAuditResponse(&audit), nil
}

// AddAuditItem adds an item to an in-progress audit
func (s *InventoryAuditService) AddAuditItem(auditID string, userID uuid.UUID, req *dto.AuditItemInput) (*dto.InventoryAuditResponse, error) {
	aID, err := uuid.Parse(auditID)
	if err != nil {
		return nil, errors.New("invalid audit ID")
	}

	var audit models.InventoryAudit
	if err := s.db.Where("id = ? AND deleted_at IS NULL", aID).First(&audit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory audit not found")
		}
		return nil, fmt.Errorf("failed to fetch audit: %w", err)
	}

	if audit.Status != models.InventoryAuditStatusInProgress {
		return nil, errors.New("can only add items to in-progress audits")
	}

	itemID, err := uuid.Parse(req.ItemID)
	if err != nil {
		return nil, errors.New("invalid item ID")
	}

	// Verify item exists and belongs to section
	var item models.InventoryItem
	if err := s.db.Where("id = ? AND inventory_section_id = ? AND deleted_at IS NULL", itemID, audit.InventorySectionID).First(&item).Error; err != nil {
		return nil, errors.New("item not found in this section")
	}

	// Check for duplicate entry
	var existingItem models.InventoryAuditItem
	if err := s.db.Where("audit_id = ? AND item_id = ? AND deleted_at IS NULL", aID, itemID).First(&existingItem).Error; err == nil {
		return nil, errors.New("this item has already been added to the audit")
	}

	difference := req.PhysicalQuantity - item.QuantityAvailable
	isMissing := req.IsMissing || difference < 0
	isExcess := req.IsExcess || difference > 0
	isDamaged := req.IsDamaged || req.Condition == "damaged" || req.Condition == "poor"

	condition := models.InventoryItemCondition(req.Condition)
	if req.Condition == "" {
		condition = models.InventoryItemCondition("good")
	}

	auditItem := &models.InventoryAuditItem{
		ID:               uuid.New(),
		AuditID:          aID,
		ItemID:           itemID,
		SystemQuantity:   item.QuantityAvailable,
		PhysicalQuantity: req.PhysicalQuantity,
		Difference:       difference,
		Condition:        condition,
		Location:         req.Location,
		IsMissing:        isMissing,
		IsExcess:         isExcess,
		IsDamaged:        isDamaged,
		Notes:            req.Notes,
		ActionTaken:      req.ActionTaken,
		CreatedBy:        userID,
	}

	if err := s.db.Create(auditItem).Error; err != nil {
		return nil, fmt.Errorf("failed to add audit item: %w", err)
	}

	// Recalculate totals
	if err := s.recalculateAuditTotals(s.db, aID); err != nil {
		return nil, err
	}

	// Reload audit
	return s.GetInventoryAuditByID(auditID)
}

// UpdateAuditItem updates an existing audit item
func (s *InventoryAuditService) UpdateAuditItem(auditID, itemID string, userID uuid.UUID, req *dto.UpdateAuditItemRequest) (*dto.InventoryAuditResponse, error) {
	aID, err := uuid.Parse(auditID)
	if err != nil {
		return nil, errors.New("invalid audit ID")
	}

	iID, err := uuid.Parse(itemID)
	if err != nil {
		return nil, errors.New("invalid item ID")
	}

	var audit models.InventoryAudit
	if err := s.db.Where("id = ? AND deleted_at IS NULL", aID).First(&audit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory audit not found")
		}
		return nil, fmt.Errorf("failed to fetch audit: %w", err)
	}

	if audit.Status != models.InventoryAuditStatusInProgress {
		return nil, errors.New("can only update items in in-progress audits")
	}

	var auditItem models.InventoryAuditItem
	if err := s.db.Where("id = ? AND audit_id = ? AND deleted_at IS NULL", iID, aID).First(&auditItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("audit item not found")
		}
		return nil, fmt.Errorf("failed to fetch audit item: %w", err)
	}

	// Update fields
	if req.PhysicalQuantity >= 0 {
		auditItem.PhysicalQuantity = req.PhysicalQuantity
		auditItem.Difference = req.PhysicalQuantity - auditItem.SystemQuantity
	}

	if req.Condition != "" {
		auditItem.Condition = models.InventoryItemCondition(req.Condition)
		auditItem.IsDamaged = req.Condition == "damaged" || req.Condition == "poor"
	}

	if req.Location != "" {
		auditItem.Location = req.Location
	}

	if req.IsMissing != nil {
		auditItem.IsMissing = *req.IsMissing
	}

	if req.IsExcess != nil {
		auditItem.IsExcess = *req.IsExcess
	}

	if req.IsDamaged != nil {
		auditItem.IsDamaged = *req.IsDamaged
	}

	if req.Notes != "" {
		auditItem.Notes = req.Notes
	}

	if req.ActionTaken != "" {
		auditItem.ActionTaken = req.ActionTaken
		auditItem.ActionTakenByID = &userID
	}

	if err := s.db.Save(&auditItem).Error; err != nil {
		return nil, fmt.Errorf("failed to update audit item: %w", err)
	}

	// Recalculate totals
	if err := s.recalculateAuditTotals(s.db, aID); err != nil {
		return nil, err
	}

	return s.GetInventoryAuditByID(auditID)
}

// RemoveAuditItem removes an item from an in-progress audit
func (s *InventoryAuditService) RemoveAuditItem(auditID, itemID string) error {
	aID, err := uuid.Parse(auditID)
	if err != nil {
		return errors.New("invalid audit ID")
	}

	iID, err := uuid.Parse(itemID)
	if err != nil {
		return errors.New("invalid item ID")
	}

	var audit models.InventoryAudit
	if err := s.db.Where("id = ? AND deleted_at IS NULL", aID).First(&audit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("inventory audit not found")
		}
		return fmt.Errorf("failed to fetch audit: %w", err)
	}

	if audit.Status != models.InventoryAuditStatusInProgress {
		return errors.New("can only remove items from in-progress audits")
	}

	var auditItem models.InventoryAuditItem
	if err := s.db.Where("id = ? AND audit_id = ? AND deleted_at IS NULL", iID, aID).First(&auditItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("audit item not found")
		}
		return fmt.Errorf("failed to fetch audit item: %w", err)
	}

	if err := s.db.Delete(&auditItem).Error; err != nil {
		return fmt.Errorf("failed to remove audit item: %w", err)
	}

	// Recalculate totals
	return s.recalculateAuditTotals(s.db, aID)
}

// CompleteInventoryAudit marks an audit as completed
func (s *InventoryAuditService) CompleteInventoryAudit(id string, userID uuid.UUID) (*dto.InventoryAuditResponse, error) {
	auditID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid audit ID")
	}

	var audit models.InventoryAudit
	if err := s.db.Where("id = ? AND deleted_at IS NULL", auditID).First(&audit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory audit not found")
		}
		return nil, fmt.Errorf("failed to fetch audit: %w", err)
	}

	if audit.Status != models.InventoryAuditStatusInProgress {
		return nil, errors.New("only in-progress audits can be completed")
	}

	// Only conductor or creator can complete
	if audit.ConductedByID != userID && audit.CreatedBy != userID {
		return nil, errors.New("you can only complete your own audits")
	}

	// Check if audit has items
	var itemCount int64
	if err := s.db.Model(&models.InventoryAuditItem{}).
		Where("audit_id = ? AND deleted_at IS NULL", auditID).
		Count(&itemCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count items: %w", err)
	}

	if itemCount == 0 {
		return nil, errors.New("cannot complete an audit with no items")
	}

	now := time.Now()
	audit.Status = models.InventoryAuditStatusCompleted
	audit.CompletedAt = &now

	if err := s.db.Save(&audit).Error; err != nil {
		return nil, fmt.Errorf("failed to complete audit: %w", err)
	}

	return s.GetInventoryAuditByID(id)
}

// ApproveInventoryAudit approves a completed audit and applies discrepancies
func (s *InventoryAuditService) ApproveInventoryAudit(id string, approverID uuid.UUID, req *dto.ApproveInventoryAuditRequest) (*dto.InventoryAuditResponse, error) {
	auditID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid audit ID")
	}

	var audit models.InventoryAudit
	if err := s.db.Where("id = ? AND deleted_at IS NULL", auditID).First(&audit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory audit not found")
		}
		return nil, fmt.Errorf("failed to fetch audit: %w", err)
	}

	if audit.Status != models.InventoryAuditStatusCompleted {
		return nil, errors.New("only completed audits can be approved")
	}

	// Begin database transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Apply stock adjustments based on audit items
	var auditItems []models.InventoryAuditItem
	if err := tx.Where("audit_id = ? AND deleted_at IS NULL", auditID).Find(&auditItems).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to fetch audit items: %w", err)
	}

	for _, auditItem := range auditItems {
		var item models.InventoryItem
		if err := tx.Where("id = ? AND deleted_at IS NULL", auditItem.ItemID).First(&item).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to fetch item: %w", err)
		}

		// Apply the physical quantity as new available quantity
		item.QuantityAvailable = auditItem.PhysicalQuantity
		item.QuantityTotal = auditItem.PhysicalQuantity

		// Auto-update status
		if item.QuantityAvailable <= 0 {
			item.Status = models.InventoryItemStatusOutOfStock
		} else if item.Status == models.InventoryItemStatusOutOfStock {
			item.Status = models.InventoryItemStatusActive
		}

		// Update condition if damaged
		if auditItem.IsDamaged {
			item.Condition = auditItem.Condition
		}

		if err := tx.Save(&item).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update item: %w", err)
		}
	}

	// Update audit status
	now := time.Now()
	audit.Status = models.InventoryAuditStatusApproved
	audit.ApprovedByID = &approverID
	audit.ApprovedAt = &now

	if req.Notes != "" {
		audit.Notes = req.Notes
	}

	if err := tx.Save(&audit).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to approve audit: %w", err)
	}

	tx.Commit()

	return s.GetInventoryAuditByID(id)
}

// RejectInventoryAudit rejects a completed audit
func (s *InventoryAuditService) RejectInventoryAudit(id string, rejectorID uuid.UUID, req *dto.RejectInventoryAuditRequest) (*dto.InventoryAuditResponse, error) {
	auditID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid audit ID")
	}

	var audit models.InventoryAudit
	if err := s.db.Where("id = ? AND deleted_at IS NULL", auditID).First(&audit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory audit not found")
		}
		return nil, fmt.Errorf("failed to fetch audit: %w", err)
	}

	if audit.Status != models.InventoryAuditStatusCompleted {
		return nil, errors.New("only completed audits can be rejected")
	}

	now := time.Now()
	audit.Status = models.InventoryAuditStatusRejected
	audit.ApprovedByID = &rejectorID
	audit.ApprovedAt = &now
	audit.RejectionReason = req.RejectionReason

	if err := s.db.Save(&audit).Error; err != nil {
		return nil, fmt.Errorf("failed to reject audit: %w", err)
	}

	return s.GetInventoryAuditByID(id)
}

// CancelInventoryAudit cancels an in-progress audit
func (s *InventoryAuditService) CancelInventoryAudit(id string, userID uuid.UUID) error {
	auditID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid audit ID")
	}

	var audit models.InventoryAudit
	if err := s.db.Where("id = ? AND deleted_at IS NULL", auditID).First(&audit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("inventory audit not found")
		}
		return fmt.Errorf("failed to fetch audit: %w", err)
	}

	if audit.Status != models.InventoryAuditStatusInProgress {
		return errors.New("only in-progress audits can be cancelled")
	}

	if audit.ConductedByID != userID && audit.CreatedBy != userID {
		return errors.New("you can only cancel your own audits")
	}

	audit.Status = models.InventoryAuditStatusCancelled

	if err := s.db.Save(&audit).Error; err != nil {
		return fmt.Errorf("failed to cancel audit: %w", err)
	}

	return nil
}

// DeleteInventoryAudit soft deletes an audit
func (s *InventoryAuditService) DeleteInventoryAudit(id string) error {
	auditID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid audit ID")
	}

	var audit models.InventoryAudit
	if err := s.db.Where("id = ? AND deleted_at IS NULL", auditID).First(&audit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("inventory audit not found")
		}
		return fmt.Errorf("failed to fetch audit: %w", err)
	}

	if err := s.db.Delete(&audit).Error; err != nil {
		return fmt.Errorf("failed to delete audit: %w", err)
	}

	return nil
}

// GetInventoryAuditStats retrieves statistics for audits
func (s *InventoryAuditService) GetInventoryAuditStats(sectionID string) (*dto.InventoryAuditStats, error) {
	query := s.db.Model(&models.InventoryAudit{}).Where("deleted_at IS NULL")

	if sectionID != "" {
		sID, err := uuid.Parse(sectionID)
		if err != nil {
			return nil, errors.New("invalid section ID")
		}
		query = query.Where("inventory_section_id = ?", sID)
	}

	var stats dto.InventoryAuditStats

	if err := query.Count(&stats.TotalAudits).Error; err != nil {
		return nil, fmt.Errorf("failed to count audits: %w", err)
	}

	if err := query.Where("status = ?", "in_progress").Count(&stats.InProgressAudits).Error; err != nil {
		return nil, fmt.Errorf("failed to count in-progress audits: %w", err)
	}

	if err := query.Where("status = ?", "completed").Count(&stats.CompletedAudits).Error; err != nil {
		return nil, fmt.Errorf("failed to count completed audits: %w", err)
	}

	if err := query.Where("status = ?", "approved").Count(&stats.ApprovedAudits).Error; err != nil {
		return nil, fmt.Errorf("failed to count approved audits: %w", err)
	}

	if err := query.Where("status = ?", "rejected").Count(&stats.RejectedAudits).Error; err != nil {
		return nil, fmt.Errorf("failed to count rejected audits: %w", err)
	}

	// Sum discrepancies
	var totalDiscrepancies int64
	if err := query.Select("COALESCE(SUM(total_discrepancies), 0)").Scan(&totalDiscrepancies).Error; err != nil {
		return nil, fmt.Errorf("failed to sum discrepancies: %w", err)
	}
	stats.TotalDiscrepancies = totalDiscrepancies

	if stats.TotalAudits > 0 {
		stats.AverageDiscrepancy = float64(totalDiscrepancies) / float64(stats.TotalAudits)
	}

	return &stats, nil
}

// toAuditResponse converts model to response DTO
func (s *InventoryAuditService) toAuditResponse(audit *models.InventoryAudit) *dto.InventoryAuditResponse {
	response := &dto.InventoryAuditResponse{
		ID:                 audit.ID.String(),
		AuditCode:          audit.AuditCode,
		InventorySectionID: audit.InventorySectionID.String(),
		AuditDate:          audit.AuditDate,
		ConductedByID:      audit.ConductedByID.String(),
		Status:             string(audit.Status),
		TotalItemsChecked:  audit.TotalItemsChecked,
		TotalDiscrepancies: audit.TotalDiscrepancies,
		TotalItemsMissing:  audit.TotalItemsMissing,
		TotalItemsExcess:   audit.TotalItemsExcess,
		TotalItemsDamaged:  audit.TotalItemsDamaged,
		Notes:              audit.Notes,
		RejectionReason:    audit.RejectionReason,
		ApprovedAt:         audit.ApprovedAt,
		CompletedAt:        audit.CompletedAt,
		CreatedBy:          audit.CreatedBy.String(),
		CreatedAt:          audit.CreatedAt,
		UpdatedAt:          audit.UpdatedAt,
	}

	if audit.ApprovedByID != nil {
		response.ApprovedByID = audit.ApprovedByID.String()
	}

	// Add section details
	if audit.InventorySection.ID != uuid.Nil {
		response.InventorySection = &dto.InventorySectionResponse{
			ID:          audit.InventorySection.ID.String(),
			Name:        audit.InventorySection.Name,
			Code:        audit.InventorySection.Code,
			Description: audit.InventorySection.Description,
			Location:    audit.InventorySection.Location,
			Status:      string(audit.InventorySection.Status),
		}
	}

	// Add conducted by details
	if audit.ConductedBy.ID != uuid.Nil {
		response.ConductedBy = &dto.UserResponse{
			ID:        audit.ConductedBy.ID.String(),
			FirstName: audit.ConductedBy.FirstName,
			LastName:  audit.ConductedBy.LastName,
			Email:     audit.ConductedBy.Email,
			Phone:     audit.ConductedBy.Phone,
		}
	}

	// Add approved by details
	if audit.ApprovedBy != nil && audit.ApprovedBy.ID != uuid.Nil {
		response.ApprovedBy = &dto.UserResponse{
			ID:        audit.ApprovedBy.ID.String(),
			FirstName: audit.ApprovedBy.FirstName,
			LastName:  audit.ApprovedBy.LastName,
			Email:     audit.ApprovedBy.Email,
		}
	}

	// Add creator details
	if audit.Creator.ID != uuid.Nil {
		response.Creator = &dto.UserResponse{
			ID:        audit.Creator.ID.String(),
			FirstName: audit.Creator.FirstName,
			LastName:  audit.Creator.LastName,
			Email:     audit.Creator.Email,
		}
	}

	// Add audit items
	if len(audit.AuditItems) > 0 {
		items := make([]dto.AuditItemResponse, len(audit.AuditItems))
		for i, item := range audit.AuditItems {
			items[i] = dto.AuditItemResponse{
				ID:               item.ID.String(),
				AuditID:          item.AuditID.String(),
				ItemID:           item.ItemID.String(),
				SystemQuantity:   item.SystemQuantity,
				PhysicalQuantity: item.PhysicalQuantity,
				Difference:       item.Difference,
				Condition:        string(item.Condition),
				Location:         item.Location,
				IsMissing:        item.IsMissing,
				IsExcess:         item.IsExcess,
				IsDamaged:        item.IsDamaged,
				Notes:            item.Notes,
				ActionTaken:      item.ActionTaken,
				CreatedBy:        item.CreatedBy.String(),
				CreatedAt:        item.CreatedAt,
				UpdatedAt:        item.UpdatedAt,
			}

			if item.ActionTakenByID != nil {
				items[i].ActionTakenByID = item.ActionTakenByID.String()
			}

			// Add item details
			if item.Item.ID != uuid.Nil {
				items[i].Item = &dto.InventoryItemResponse{
					ID:                item.Item.ID.String(),
					ItemCode:          item.Item.ItemCode,
					ItemName:          item.Item.ItemName,
					Description:       item.Item.Description,
					QuantityTotal:     item.Item.QuantityTotal,
					QuantityAvailable: item.Item.QuantityAvailable,
					Condition:         string(item.Item.Condition),
					Status:            string(item.Item.Status),
				}
			}

			// Add action taken by details
			if item.ActionTakenBy != nil && item.ActionTakenBy.ID != uuid.Nil {
				items[i].ActionTakenBy = &dto.UserResponse{
					ID:        item.ActionTakenBy.ID.String(),
					FirstName: item.ActionTakenBy.FirstName,
					LastName:  item.ActionTakenBy.LastName,
					Email:     item.ActionTakenBy.Email,
				}
			}
		}
		response.AuditItems = items
	}

	return response
}