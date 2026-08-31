// services/scheme_of_work_item_service.go
package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"crm-go/models"
	"crm-go/dto"
)

type SchemeOfWorkItemService struct {
	db *gorm.DB
}

func NewSchemeOfWorkItemService(db *gorm.DB) *SchemeOfWorkItemService {
	return &SchemeOfWorkItemService{db: db}
}

// CreateSchemeOfWorkItem creates a new scheme of work item
func (s *SchemeOfWorkItemService) CreateSchemeOfWorkItem(req *dto.CreateSchemeOfWorkItemRequest) (*dto.SchemeOfWorkItemResponse, error) {
	// Validate input
	if err := s.validateItemRequest(req); err != nil {
		return nil, err
	}

	// Parse UUIDs
	schemeOfWorkID, err := uuid.Parse(req.SchemeOfWorkID)
	if err != nil {
		return nil, errors.New("invalid scheme of work ID format")
	}

	// Check if scheme of work exists
	var schemeOfWork models.SchemeOfWork
	if err := s.db.Where("id = ? AND deleted_at IS NULL", schemeOfWorkID).First(&schemeOfWork).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("scheme of work not found")
		}
		return nil, errors.New("failed to verify scheme of work: " + err.Error())
	}

	// Validate week range
	if req.WeekEnd > 0 && req.WeekEnd < req.WeekStart {
		return nil, errors.New("week end must be greater than or equal to week start")
	}

	// Check if module exists (if provided)
	var module models.Module
	if req.ModuleID != "" {
		moduleID, err := uuid.Parse(req.ModuleID)
		if err != nil {
			return nil, errors.New("invalid module ID format")
		}
		if err := s.db.Where("id = ? AND deleted_at IS NULL", moduleID).First(&module).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("module not found")
			}
			return nil, errors.New("failed to verify module: " + err.Error())
		}
	}

	// Set default sequence if not provided
	sequence := req.Sequence
	if sequence == 0 {
		var maxSequence int
		if err := s.db.Model(&models.SchemeOfWorkItem{}).
			Where("scheme_of_work_id = ? AND deleted_at IS NULL", schemeOfWorkID).
			Select("COALESCE(MAX(sequence), 0)").
			Scan(&maxSequence).Error; err != nil {
			return nil, errors.New("failed to determine sequence: " + err.Error())
		}
		sequence = maxSequence + 1
	}

	// Create scheme of work item
	item := &models.SchemeOfWorkItem{
		ID:             uuid.New(),
		SchemeOfWorkID: schemeOfWorkID,
		ModuleID:       uuid.Nil,
		WeekStart:      req.WeekStart,
		WeekEnd:        req.WeekEnd,
		Topic:          req.Topic,
		Subtopic:       req.Subtopic,
		Content:        req.Content,
		Sequence:       sequence,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if req.ModuleID != "" {
		moduleID, _ := uuid.Parse(req.ModuleID)
		item.ModuleID = moduleID
	}

	if err := s.db.Create(item).Error; err != nil {
		return nil, errors.New("failed to create scheme of work item: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("SchemeOfWork").Preload("Module").First(item, item.ID).Error; err != nil {
		return nil, errors.New("failed to load item details: " + err.Error())
	}

	return s.toItemResponse(item), nil
}

// BulkCreateSchemeItems creates multiple scheme of work items
func (s *SchemeOfWorkItemService) BulkCreateSchemeItems(req *dto.BulkCreateSchemeItemsRequest) (*dto.BulkSchemeItemResult, error) {
	// Parse UUIDs
	schemeOfWorkID, err := uuid.Parse(req.SchemeOfWorkID)
	if err != nil {
		return nil, errors.New("invalid scheme of work ID format")
	}

	// Check if scheme of work exists
	var schemeOfWork models.SchemeOfWork
	if err := s.db.Where("id = ? AND deleted_at IS NULL", schemeOfWorkID).First(&schemeOfWork).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("scheme of work not found")
		}
		return nil, errors.New("failed to verify scheme of work: " + err.Error())
	}

	result := &dto.BulkSchemeItemResult{
		Created: []dto.SchemeOfWorkItemResponse{},
		Errors:  []dto.BulkSchemeItemError{},
	}

	// Get max sequence
	var maxSequence int
	if err := s.db.Model(&models.SchemeOfWorkItem{}).
		Where("scheme_of_work_id = ? AND deleted_at IS NULL", schemeOfWorkID).
		Select("COALESCE(MAX(sequence), 0)").
		Scan(&maxSequence).Error; err != nil {
		return nil, errors.New("failed to determine sequence: " + err.Error())
	}
	currentSequence := maxSequence

	for _, itemReq := range req.Items {
		// Validate
		if itemReq.Topic == "" {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkSchemeItemError{
				Topic: itemReq.Topic,
				Error: "topic is required",
			})
			continue
		}

		if itemReq.WeekStart < 1 {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkSchemeItemError{
				Topic: itemReq.Topic,
				Error: "week start must be at least 1",
			})
			continue
		}

		if itemReq.WeekEnd > 0 && itemReq.WeekEnd < itemReq.WeekStart {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkSchemeItemError{
				Topic: itemReq.Topic,
				Error: "week end must be greater than or equal to week start",
			})
			continue
		}

		// Check if module exists (if provided)
		var moduleID uuid.UUID
		if itemReq.ModuleID != "" {
			moduleID, err = uuid.Parse(itemReq.ModuleID)
			if err != nil {
				result.FailedCount++
				result.Errors = append(result.Errors, dto.BulkSchemeItemError{
					Topic: itemReq.Topic,
					Error: "invalid module ID format",
				})
				continue
			}
			var module models.Module
			if err := s.db.Where("id = ? AND deleted_at IS NULL", moduleID).First(&module).Error; err != nil {
				result.FailedCount++
				result.Errors = append(result.Errors, dto.BulkSchemeItemError{
					Topic: itemReq.Topic,
					Error: "module not found",
				})
				continue
			}
		}

		// Set sequence
		sequence := itemReq.Sequence
		if sequence == 0 {
			currentSequence++
			sequence = currentSequence
		}

		// Create item
		item := &models.SchemeOfWorkItem{
			ID:             uuid.New(),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:       moduleID,
			WeekStart:      itemReq.WeekStart,
			WeekEnd:        itemReq.WeekEnd,
			Topic:          itemReq.Topic,
			Subtopic:       itemReq.Subtopic,
			Content:        itemReq.Content,
			Sequence:       sequence,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := s.db.Create(item).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkSchemeItemError{
				Topic: itemReq.Topic,
				Error: "failed to create item: " + err.Error(),
			})
			continue
		}

		// Preload relationships
		if err := s.db.Preload("SchemeOfWork").Preload("Module").First(item, item.ID).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkSchemeItemError{
				Topic: itemReq.Topic,
				Error: "failed to load item details",
			})
			continue
		}

		result.SuccessCount++
		result.Created = append(result.Created, *s.toItemResponse(item))
	}

	return result, nil
}

// GetAllItems retrieves all scheme of work items with pagination and filters
func (s *SchemeOfWorkItemService) GetAllItems(params *dto.SchemeOfWorkItemQueryParams) (*dto.SchemeOfWorkItemListResponse, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = "sequence"
	}
	if params.SortOrder == "" {
		params.SortOrder = "asc"
	}

	// Build query
	query := s.db.Model(&models.SchemeOfWorkItem{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.SchemeOfWorkID != "" {
		schemeID, err := uuid.Parse(params.SchemeOfWorkID)
		if err == nil {
			query = query.Where("scheme_of_work_id = ?", schemeID)
		}
	}

	if params.ModuleID != "" {
		moduleID, err := uuid.Parse(params.ModuleID)
		if err == nil {
			query = query.Where("module_id = ?", moduleID)
		}
	}

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(topic) LIKE ? OR LOWER(subtopic) LIKE ? OR LOWER(content) LIKE ?",
			search, search, search)
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
	query = query.Order(params.SortBy + " " + sortDirection)

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	// Execute with preloads
	var items []models.SchemeOfWorkItem
	if err := query.Preload("SchemeOfWork").Preload("Module").Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch items: %w", err)
	}

	// Convert to response
	responses := make([]dto.SchemeOfWorkItemResponse, len(items))
	for i, item := range items {
		responses[i] = *s.toItemResponse(&item)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.SchemeOfWorkItemListResponse{
		Items:      responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetItemByID retrieves a single scheme of work item by ID
func (s *SchemeOfWorkItemService) GetItemByID(id string) (*dto.SchemeOfWorkItemResponse, error) {
	itemID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid item ID")
	}

	var item models.SchemeOfWorkItem
	if err := s.db.Where("id = ? AND deleted_at IS NULL", itemID).
		Preload("SchemeOfWork").
		Preload("Module").
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("scheme of work item not found")
		}
		return nil, errors.New("failed to fetch item: " + err.Error())
	}

	return s.toItemResponse(&item), nil
}

// GetItemsBySchemeOfWork retrieves all items for a specific scheme of work
func (s *SchemeOfWorkItemService) GetItemsBySchemeOfWork(schemeOfWorkID string) ([]dto.SchemeOfWorkItemResponse, error) {
	sID, err := uuid.Parse(schemeOfWorkID)
	if err != nil {
		return nil, errors.New("invalid scheme of work ID")
	}

	var items []models.SchemeOfWorkItem
	if err := s.db.Where("scheme_of_work_id = ? AND deleted_at IS NULL", sID).
		Preload("SchemeOfWork").
		Preload("Module").
		Order("sequence ASC").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch items: %w", err)
	}

	responses := make([]dto.SchemeOfWorkItemResponse, len(items))
	for i, item := range items {
		responses[i] = *s.toItemResponse(&item)
	}

	return responses, nil
}

// GetItemsByModule retrieves all items for a specific module
func (s *SchemeOfWorkItemService) GetItemsByModule(moduleID string) ([]dto.SchemeOfWorkItemResponse, error) {
	mID, err := uuid.Parse(moduleID)
	if err != nil {
		return nil, errors.New("invalid module ID")
	}

	var items []models.SchemeOfWorkItem
	if err := s.db.Where("module_id = ? AND deleted_at IS NULL", mID).
		Preload("SchemeOfWork").
		Preload("Module").
		Order("sequence ASC").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch items: %w", err)
	}

	responses := make([]dto.SchemeOfWorkItemResponse, len(items))
	for i, item := range items {
		responses[i] = *s.toItemResponse(&item)
	}

	return responses, nil
}

// UpdateSchemeOfWorkItem updates an existing scheme of work item
func (s *SchemeOfWorkItemService) UpdateSchemeOfWorkItem(id string, req *dto.UpdateSchemeOfWorkItemRequest) (*dto.SchemeOfWorkItemResponse, error) {
	itemID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid item ID")
	}

	// Find existing item
	var item models.SchemeOfWorkItem
	if err := s.db.Where("id = ? AND deleted_at IS NULL", itemID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("scheme of work item not found")
		}
		return nil, errors.New("failed to fetch item: " + err.Error())
	}

	// Update fields
	if req.ModuleID != "" {
		moduleID, err := uuid.Parse(req.ModuleID)
		if err != nil {
			return nil, errors.New("invalid module ID format")
		}
		if moduleID != uuid.Nil {
			var module models.Module
			if err := s.db.Where("id = ? AND deleted_at IS NULL", moduleID).First(&module).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.New("module not found")
				}
				return nil, errors.New("failed to verify module: " + err.Error())
			}
		}
		item.ModuleID = moduleID
	}

	if req.WeekStart > 0 {
		item.WeekStart = req.WeekStart
	}

	if req.WeekEnd > 0 {
		if req.WeekEnd < item.WeekStart {
			return nil, errors.New("week end must be greater than or equal to week start")
		}
		item.WeekEnd = req.WeekEnd
	}

	if req.Topic != "" {
		item.Topic = req.Topic
	}

	if req.Subtopic != "" {
		item.Subtopic = req.Subtopic
	}

	if req.Content != "" {
		item.Content = req.Content
	}

	if req.Sequence > 0 && req.Sequence != item.Sequence {
		// Handle reordering
		if err := s.reorderItems(item.SchemeOfWorkID, item.Sequence, req.Sequence, itemID); err != nil {
			return nil, err
		}
		item.Sequence = req.Sequence
	}

	item.UpdatedAt = time.Now()

	if err := s.db.Save(&item).Error; err != nil {
		return nil, errors.New("failed to update item: " + err.Error())
	}

	// Preload relationships
	if err := s.db.Preload("SchemeOfWork").Preload("Module").First(&item, item.ID).Error; err != nil {
		return nil, errors.New("failed to load item details: " + err.Error())
	}

	return s.toItemResponse(&item), nil
}

// DeleteSchemeOfWorkItem soft deletes a scheme of work item
func (s *SchemeOfWorkItemService) DeleteSchemeOfWorkItem(id string) error {
	itemID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid item ID")
	}

	var item models.SchemeOfWorkItem
	if err := s.db.Where("id = ? AND deleted_at IS NULL", itemID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("scheme of work item not found")
		}
		return errors.New("failed to fetch item: " + err.Error())
	}

	if err := s.db.Delete(&item).Error; err != nil {
		return errors.New("failed to delete item: " + err.Error())
	}

	return nil
}

// reorderItems reorders items when sequence changes
func (s *SchemeOfWorkItemService) reorderItems(schemeOfWorkID uuid.UUID, oldSeq, newSeq int, excludeID uuid.UUID) error {
	if oldSeq == newSeq {
		return nil
	}

	var items []models.SchemeOfWorkItem
	if err := s.db.Where("scheme_of_work_id = ? AND deleted_at IS NULL AND id != ?", schemeOfWorkID, excludeID).
		Order("sequence ASC").
		Find(&items).Error; err != nil {
		return errors.New("failed to fetch items for reordering: " + err.Error())
	}

	// Update sequences
	if oldSeq < newSeq {
		// Moving down: decrease sequences between oldSeq+1 and newSeq
		for _, it := range items {
			if it.Sequence > oldSeq && it.Sequence <= newSeq {
				it.Sequence--
				if err := s.db.Save(&it).Error; err != nil {
					return errors.New("failed to update item sequence: " + err.Error())
				}
			}
		}
	} else {
		// Moving up: increase sequences between newSeq and oldSeq-1
		for _, it := range items {
			if it.Sequence >= newSeq && it.Sequence < oldSeq {
				it.Sequence++
				if err := s.db.Save(&it).Error; err != nil {
					return errors.New("failed to update item sequence: " + err.Error())
				}
			}
		}
	}

	return nil
}

// validateItemRequest validates the item request
func (s *SchemeOfWorkItemService) validateItemRequest(req *dto.CreateSchemeOfWorkItemRequest) error {
	if req.SchemeOfWorkID == "" {
		return errors.New("scheme of work ID is required")
	}
	if req.WeekStart < 1 {
		return errors.New("week start must be at least 1")
	}
	if req.Topic == "" {
		return errors.New("topic is required")
	}
	if req.WeekEnd > 0 && req.WeekEnd < req.WeekStart {
		return errors.New("week end must be greater than or equal to week start")
	}
	return nil
}

// toItemResponse converts model to response DTO
func (s *SchemeOfWorkItemService) toItemResponse(item *models.SchemeOfWorkItem) *dto.SchemeOfWorkItemResponse {
	response := &dto.SchemeOfWorkItemResponse{
		ID:             item.ID.String(),
		SchemeOfWorkID: item.SchemeOfWorkID.String(),
		ModuleID:       item.ModuleID.String(),
		WeekStart:      item.WeekStart,
		WeekEnd:        item.WeekEnd,
		Topic:          item.Topic,
		Subtopic:       item.Subtopic,
		Content:        item.Content,
		Sequence:       item.Sequence,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}

	// Add scheme of work details if preloaded
	if item.SchemeOfWork.ID != uuid.Nil {
		response.SchemeOfWork = &dto.SchemeOfWorkResponse{
			ID:          item.SchemeOfWork.ID.String(),
			Title:       item.SchemeOfWork.Title,
			Description: item.SchemeOfWork.Description,
			Status:      item.SchemeOfWork.Status,
		}
	}

	// Add module details if preloaded
	if item.Module.ID != uuid.Nil {
		response.Module = &dto.ModuleResponse{
			ID:          item.Module.ID.String(),
			Name:        item.Module.Name,
			Code:        item.Module.Code,
			Description: item.Module.Description,
			Sequence:    item.Module.Sequence,
		}
	}

	return response
}