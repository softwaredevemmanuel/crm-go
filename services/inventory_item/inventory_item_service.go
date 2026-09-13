// services/inventory_item_service.go
package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"crm-go/dto"
	"crm-go/models"
)

type InventoryItemService struct {
	db *gorm.DB
}

func NewInventoryItemService(db *gorm.DB) *InventoryItemService {
	return &InventoryItemService{db: db}
}

// CreateInventoryItem creates a new inventory item
func (s *InventoryItemService) CreateInventoryItem(userID uuid.UUID, req *dto.CreateInventoryItemRequest) (*dto.InventoryItemResponse, error) {
	// Validate input
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Check if item code already exists
	var existingItem models.InventoryItem
	if err := s.db.Where("item_code = ? AND deleted_at IS NULL", strings.ToUpper(req.ItemCode)).First(&existingItem).Error; err == nil {
		return nil, errors.New("item with this code already exists")
	}

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

	// Set defaults
	condition := req.Condition
	if condition == "" {
		condition = "good"
	}
	status := req.Status
	if status == "" {
		status = "active"
	}
	hazardClass := req.HazardClass
	if hazardClass == "" {
		hazardClass = "none"
	}

	// Create item
	item := &models.InventoryItem{
		ID:                   uuid.New(),
		ItemCode:             strings.ToUpper(strings.TrimSpace(req.ItemCode)),
		ItemName:             strings.TrimSpace(req.ItemName),
		Description:          req.Description,
		InventorySectionID:   sectionID,
		Brand:                req.Brand,
		Model:                req.Model,
		SerialNumber:         req.SerialNumber,
		Barcode:              req.Barcode,
		QuantityTotal:        req.QuantityTotal,
		QuantityAvailable:    req.QuantityAvailable,
		QuantityInUse:        req.QuantityInUse,
		QuantityDamaged:      req.QuantityDamaged,
		Condition:            models.InventoryItemCondition(condition),
		Status:               models.InventoryItemStatus(status),
		Location:             req.Location,
		StorageConditions:    req.StorageConditions,
		IsChemical:           req.IsChemical,
		ChemicalFormula:      req.ChemicalFormula,
		CASNumber:            req.CASNumber,
		HazardClass:          models.InventoryItemHazardClass(hazardClass),
		Concentration:        req.Concentration,
		Volume:               req.Volume,
		SafetyDataSheetURL:   req.SafetyDataSheetURL,
		RequiresSupervision:  req.RequiresSupervision,
		SafetyInstructions:   req.SafetyInstructions,
		DisposalInstructions: req.DisposalInstructions,
		RestrictedAccess:     req.RestrictedAccess,
		AllowedRoles:         req.AllowedRoles,
		CreatedBy:            userID,
	}

	if err := s.db.Create(item).Error; err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	// Create images
	for _, img := range req.Images {
		image := &models.InventoryItemImage{
			ID:              uuid.New(),
			InventoryItemID: item.ID,
			ImageURL:        img.ImageURL,
			Caption:         img.Caption,
			IsPrimary:       img.IsPrimary,
			DisplayOrder:    img.DisplayOrder,
		}
		if err := s.db.Create(image).Error; err != nil {
			return nil, fmt.Errorf("failed to create image: %w", err)
		}
	}

	// Preload relationships
	if err := s.db.Preload("InventorySection").
		Preload("Creator").
		Preload("Images").
		First(item, item.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load item details: %w", err)
	}

	return s.toItemResponse(item), nil
}

// GetInventoryItemByID retrieves an item by ID
func (s *InventoryItemService) GetInventoryItemByID(id string) (*dto.InventoryItemResponse, error) {
	itemID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid item ID")
	}

	var item models.InventoryItem
	if err := s.db.Where("id = ? AND deleted_at IS NULL", itemID).
		Preload("InventorySection").
		Preload("Creator").
		Preload("Images").
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory item not found")
		}
		return nil, fmt.Errorf("failed to fetch item: %w", err)
	}

	return s.toItemResponse(&item), nil
}

// GetInventoryItems retrieves items with filters and pagination
func (s *InventoryItemService) GetInventoryItems(params *dto.InventoryItemQueryParams) (*dto.InventoryItemListResponse, error) {
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
	query := s.db.Model(&models.InventoryItem{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.InventorySectionID != "" {
		sectionID, err := uuid.Parse(params.InventorySectionID)
		if err == nil {
			query = query.Where("inventory_section_id = ?", sectionID)
		}
	}

	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("item_name ILIKE ? OR item_code ILIKE ? OR description ILIKE ? OR brand ILIKE ? OR model ILIKE ?", search, search, search, search, search)
	}

	if params.Condition != "" {
		query = query.Where("condition = ?", params.Condition)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.IsChemical != nil {
		query = query.Where("is_chemical = ?", *params.IsChemical)
	}

	if params.HazardClass != "" {
		query = query.Where("hazard_class = ?", params.HazardClass)
	}

	if params.Location != "" {
		location := "%" + params.Location + "%"
		query = query.Where("location ILIKE ?", location)
	}

	if params.LowStock != nil && *params.LowStock {
		query = query.Where("quantity_available <= quantity_total * 0.2") // Low stock threshold
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count items: %w", err)
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
	var items []models.InventoryItem
	if err := query.
		Preload("InventorySection").
		Preload("Creator").
		Preload("Images").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch items: %w", err)
	}

	// Convert to response
	responses := make([]dto.InventoryItemResponse, len(items))
	for i, item := range items {
		responses[i] = *s.toItemResponse(&item)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.InventoryItemListResponse{
		Items:      responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetInventoryItemsBySection retrieves all items in a specific section
func (s *InventoryItemService) GetInventoryItemsBySection(sectionID string) ([]dto.InventoryItemResponse, error) {
	sID, err := uuid.Parse(sectionID)
	if err != nil {
		return nil, errors.New("invalid section ID")
	}

	var items []models.InventoryItem
	if err := s.db.Where("inventory_section_id = ? AND deleted_at IS NULL", sID).
		Preload("InventorySection").
		Preload("Creator").
		Preload("Images").
		Order("item_name ASC").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch items: %w", err)
	}

	responses := make([]dto.InventoryItemResponse, len(items))
	for i, item := range items {
		responses[i] = *s.toItemResponse(&item)
	}

	return responses, nil
}

// UpdateInventoryItem updates an existing item
func (s *InventoryItemService) UpdateInventoryItem(id string, req *dto.UpdateInventoryItemRequest) (*dto.InventoryItemResponse, error) {
	itemID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid item ID")
	}

	// Find existing item
	var item models.InventoryItem
	if err := s.db.Where("id = ? AND deleted_at IS NULL", itemID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory item not found")
		}
		return nil, fmt.Errorf("failed to fetch item: %w", err)
	}

	// Update item code (check for duplicates)
	if req.ItemCode != "" {
		newCode := strings.ToUpper(strings.TrimSpace(req.ItemCode))
		if newCode != item.ItemCode {
			var existingItem models.InventoryItem
			if err := s.db.Where("item_code = ? AND id != ? AND deleted_at IS NULL", newCode, itemID).First(&existingItem).Error; err == nil {
				return nil, errors.New("item with this code already exists")
			}
			item.ItemCode = newCode
		}
	}

	if req.ItemName != "" {
		item.ItemName = strings.TrimSpace(req.ItemName)
	}

	if req.Description != "" {
		item.Description = req.Description
	}

	if req.InventorySectionID != "" {
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
		item.InventorySectionID = sectionID
	}

	if req.Brand != "" {
		item.Brand = req.Brand
	}

	if req.Model != "" {
		item.Model = req.Model
	}

	if req.SerialNumber != "" {
		item.SerialNumber = req.SerialNumber
	}

	if req.Barcode != "" {
		item.Barcode = req.Barcode
	}

	if req.QuantityTotal != nil {
		item.QuantityTotal = *req.QuantityTotal
	}

	if req.QuantityAvailable != nil {
		item.QuantityAvailable = *req.QuantityAvailable
	}

	if req.QuantityInUse != nil {
		item.QuantityInUse = *req.QuantityInUse
	}

	if req.QuantityDamaged != nil {
		item.QuantityDamaged = *req.QuantityDamaged
	}

	if req.Condition != "" {
		item.Condition = models.InventoryItemCondition(req.Condition)
	}

	if req.Status != "" {
		item.Status = models.InventoryItemStatus(req.Status)
	}

	if req.Location != "" {
		item.Location = req.Location
	}

	if req.StorageConditions != "" {
		item.StorageConditions = req.StorageConditions
	}

	if req.IsChemical != nil {
		item.IsChemical = *req.IsChemical
	}

	if req.ChemicalFormula != "" {
		item.ChemicalFormula = req.ChemicalFormula
	}

	if req.CASNumber != "" {
		item.CASNumber = req.CASNumber
	}

	if req.HazardClass != "" {
		item.HazardClass = models.InventoryItemHazardClass(req.HazardClass)
	}

	if req.Concentration != "" {
		item.Concentration = req.Concentration
	}

	if req.Volume != "" {
		item.Volume = req.Volume
	}

	if req.SafetyDataSheetURL != "" {
		item.SafetyDataSheetURL = req.SafetyDataSheetURL
	}

	if req.RequiresSupervision != nil {
		item.RequiresSupervision = *req.RequiresSupervision
	}

	if req.SafetyInstructions != "" {
		item.SafetyInstructions = req.SafetyInstructions
	}

	if req.DisposalInstructions != "" {
		item.DisposalInstructions = req.DisposalInstructions
	}

	if req.RestrictedAccess != nil {
		item.RestrictedAccess = *req.RestrictedAccess
	}

	if req.AllowedRoles != "" {
		item.AllowedRoles = req.AllowedRoles
	}

	if err := s.db.Save(&item).Error; err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	// Update images if provided
	if req.Images != nil {
		// Delete existing images
		if err := s.db.Where("inventory_item_id = ?", item.ID).Delete(&models.InventoryItemImage{}).Error; err != nil {
			return nil, fmt.Errorf("failed to delete old images: %w", err)
		}

		// Create new images
		for _, img := range req.Images {
			image := &models.InventoryItemImage{
				ID:              uuid.New(),
				InventoryItemID: item.ID,
				ImageURL:        img.ImageURL,
				Caption:         img.Caption,
				IsPrimary:       img.IsPrimary,
				DisplayOrder:    img.DisplayOrder,
			}
			if err := s.db.Create(image).Error; err != nil {
				return nil, fmt.Errorf("failed to create image: %w", err)
			}
		}
	}

	// Reload with relationships
	if err := s.db.Preload("InventorySection").
		Preload("Creator").
		Preload("Images").
		First(&item, item.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load item details: %w", err)
	}

	return s.toItemResponse(&item), nil
}

// DeleteInventoryItem soft deletes an item
func (s *InventoryItemService) DeleteInventoryItem(id string) error {
	itemID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid item ID")
	}

	var item models.InventoryItem
	if err := s.db.Where("id = ? AND deleted_at IS NULL", itemID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("inventory item not found")
		}
		return fmt.Errorf("failed to fetch item: %w", err)
	}

	if err := s.db.Delete(&item).Error; err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}

// GetInventoryItemStats retrieves statistics for inventory items
func (s *InventoryItemService) GetInventoryItemStats(sectionID string) (*dto.InventoryItemStats, error) {
	query := s.db.Model(&models.InventoryItem{}).Where("deleted_at IS NULL")

	if sectionID != "" {
		sID, err := uuid.Parse(sectionID)
		if err != nil {
			return nil, errors.New("invalid section ID")
		}
		query = query.Where("inventory_section_id = ?", sID)
	}

	var stats dto.InventoryItemStats

	// Total items
	if err := query.Count(&stats.TotalItems).Error; err != nil {
		return nil, fmt.Errorf("failed to count items: %w", err)
	}

	// Active items
	if err := query.Where("status = ?", "active").Count(&stats.ActiveItems).Error; err != nil {
		return nil, fmt.Errorf("failed to count active items: %w", err)
	}

	// Inactive items
	if err := query.Where("status = ?", "inactive").Count(&stats.InactiveItems).Error; err != nil {
		return nil, fmt.Errorf("failed to count inactive items: %w", err)
	}

	// Out of stock items
	if err := query.Where("status = ?", "out_of_stock").Count(&stats.OutOfStockItems).Error; err != nil {
		return nil, fmt.Errorf("failed to count out of stock items: %w", err)
	}

	// Chemical items
	if err := query.Where("is_chemical = ?", true).Count(&stats.ChemicalItems).Error; err != nil {
		return nil, fmt.Errorf("failed to count chemical items: %w", err)
	}

	// Restricted items
	if err := query.Where("restricted_access = ?", true).Count(&stats.RestrictedItems).Error; err != nil {
		return nil, fmt.Errorf("failed to count restricted items: %w", err)
	}

	// Sum quantities
	type QuantitySum struct {
		TotalQuantity     int64
		AvailableQuantity int64
		DamagedQuantity   int64
	}

	var sums QuantitySum
	if err := query.Select(
		"COALESCE(SUM(quantity_total), 0) as total_quantity, " +
			"COALESCE(SUM(quantity_available), 0) as available_quantity, " +
			"COALESCE(SUM(quantity_damaged), 0) as damaged_quantity",
	).Scan(&sums).Error; err != nil {
		return nil, fmt.Errorf("failed to sum quantities: %w", err)
	}

	stats.TotalQuantity = sums.TotalQuantity
	stats.AvailableQuantity = sums.AvailableQuantity
	stats.DamagedQuantity = sums.DamagedQuantity

	return &stats, nil
}

// validateCreateRequest validates the create request
func (s *InventoryItemService) validateCreateRequest(req *dto.CreateInventoryItemRequest) error {
	if strings.TrimSpace(req.ItemCode) == "" {
		return errors.New("item code is required")
	}
	if strings.TrimSpace(req.ItemName) == "" {
		return errors.New("item name is required")
	}
	if req.InventorySectionID == "" {
		return errors.New("inventory section ID is required")
	}
	if req.Condition != "" && req.Condition != "new" && req.Condition != "good" && req.Condition != "fair" &&
		req.Condition != "poor" && req.Condition != "damaged" && req.Condition != "obsolete" && req.Condition != "expired" {
		return errors.New("invalid condition value")
	}
	if req.Status != "" && req.Status != "active" && req.Status != "inactive" && req.Status != "archived" &&
		req.Status != "disposed" && req.Status != "out_of_stock" {
		return errors.New("invalid status value")
	}
	return nil
}

// toItemResponse converts model to response DTO
func (s *InventoryItemService) toItemResponse(item *models.InventoryItem) *dto.InventoryItemResponse {
	response := &dto.InventoryItemResponse{
		ID:                   item.ID.String(),
		ItemCode:             item.ItemCode,
		ItemName:             item.ItemName,
		Description:          item.Description,
		InventorySectionID:   item.InventorySectionID.String(),
		Brand:                item.Brand,
		Model:                item.Model,
		SerialNumber:         item.SerialNumber,
		Barcode:              item.Barcode,
		QuantityTotal:        item.QuantityTotal,
		QuantityAvailable:    item.QuantityAvailable,
		QuantityInUse:        item.QuantityInUse,
		QuantityDamaged:      item.QuantityDamaged,
		Condition:            string(item.Condition),
		Status:               string(item.Status),
		Location:             item.Location,
		StorageConditions:    item.StorageConditions,
		IsChemical:           item.IsChemical,
		ChemicalFormula:      item.ChemicalFormula,
		CASNumber:            item.CASNumber,
		HazardClass:          string(item.HazardClass),
		Concentration:        item.Concentration,
		Volume:               item.Volume,
		SafetyDataSheetURL:   item.SafetyDataSheetURL,
		RequiresSupervision:  item.RequiresSupervision,
		SafetyInstructions:   item.SafetyInstructions,
		DisposalInstructions: item.DisposalInstructions,
		RestrictedAccess:     item.RestrictedAccess,
		AllowedRoles:         item.AllowedRoles,
		CreatedBy:            item.CreatedBy.String(),
		CreatedAt:            item.CreatedAt,
		UpdatedAt:            item.UpdatedAt,
	}

	// Add images
	if len(item.Images) > 0 {
		images := make([]dto.ImageResponse, len(item.Images))
		for i, img := range item.Images {
			images[i] = dto.ImageResponse{
				ID:           img.ID.String(),
				ImageURL:     img.ImageURL,
				Caption:      img.Caption,
				IsPrimary:    img.IsPrimary,
				DisplayOrder: img.DisplayOrder,
			}
		}
		response.Images = images
	}

	// Add section details if preloaded
	if item.InventorySection.ID != uuid.Nil {
		response.InventorySection = &dto.InventorySectionResponse{
			ID:          item.InventorySection.ID.String(),
			Name:        item.InventorySection.Name,
			Code:        item.InventorySection.Code,
			Description: item.InventorySection.Description,
			Location:    item.InventorySection.Location,
			Status:      string(item.InventorySection.Status),
			Notes:       item.InventorySection.Notes,
			CreatedAt:   item.InventorySection.CreatedAt,
			UpdatedAt:   item.InventorySection.UpdatedAt,
		}
	}

	// Add creator details if preloaded
	if item.Creator.ID != uuid.Nil {
		response.Creator = &dto.UserResponse{
			ID:        item.Creator.ID.String(),
			FirstName: item.Creator.FirstName,
			LastName:  item.Creator.LastName,
			Email:     item.Creator.Email,
			Phone:     item.Creator.Phone,
			Role:      item.Creator.Role,
			Position:  item.Creator.Position,
		}
	}

	return response
}