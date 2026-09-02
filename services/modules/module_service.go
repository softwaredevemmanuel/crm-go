// services/module_service.go
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

type ModuleService struct {
	db *gorm.DB
}

func NewModuleService(db *gorm.DB) *ModuleService {
	return &ModuleService{db: db}
}

// CreateModule creates a new module
func (s *ModuleService) CreateModule(req *dto.CreateModuleRequest) (*dto.ModuleResponse, error) {
	// Validate input
	if err := s.validateModuleRequest(req); err != nil {
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

	// Check if module title already exists for this scheme of work
	var existing models.Module
	if err := s.db.Where("scheme_of_work_id = ? AND title = ? AND deleted_at IS NULL",
		schemeOfWorkID, req.Title).First(&existing).Error; err == nil {
		return nil, errors.New("module title already exists for this scheme of work")
	}

	// Set default module order if not provided
	moduleOrder := req.ModuleOrder
	if moduleOrder == 0 {
		// Get the highest module order for this scheme of work
		var maxOrder int
		if err := s.db.Model(&models.Module{}).
			Where("scheme_of_work_id = ? AND deleted_at IS NULL", schemeOfWorkID).
			Select("COALESCE(MAX(module_order), 0)").
			Scan(&maxOrder).Error; err != nil {
			return nil, errors.New("failed to determine module order: " + err.Error())
		}
		moduleOrder = maxOrder + 1
	}

	// Create module
	module := &models.Module{
		ID:             uuid.New(),
		SchemeOfWorkID: schemeOfWorkID,
		Title:          req.Title,
		Description:    req.Description,
		ModuleOrder:    moduleOrder,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.db.Create(module).Error; err != nil {
		return nil, errors.New("failed to create module: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("SchemeOfWork").First(module, module.ID).Error; err != nil {
		return nil, errors.New("failed to load module details: " + err.Error())
	}

	return s.toModuleResponse(module), nil
}

// BulkCreateModules creates multiple modules
func (s *ModuleService) BulkCreateModules(req *dto.BulkCreateModulesRequest) (*dto.BulkModuleResult, error) {
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

	result := &dto.BulkModuleResult{
		Created: []dto.ModuleResponse{},
		Errors:  []dto.BulkModuleError{},
	}

	// Get existing module titles to avoid duplicates
	var existingModules []models.Module
	if err := s.db.Where("scheme_of_work_id = ? AND deleted_at IS NULL", schemeOfWorkID).Find(&existingModules).Error; err != nil {
		return nil, fmt.Errorf("failed to check existing modules: %w", err)
	}

	existingTitleMap := make(map[string]bool)
	for _, m := range existingModules {
		existingTitleMap[m.Title] = true
	}

	// Get max module order
	var maxOrder int
	if err := s.db.Model(&models.Module{}).
		Where("scheme_of_work_id = ? AND deleted_at IS NULL", schemeOfWorkID).
		Select("COALESCE(MAX(module_order), 0)").
		Scan(&maxOrder).Error; err != nil {
		return nil, errors.New("failed to determine module order: " + err.Error())
	}
	currentOrder := maxOrder

	for _, moduleReq := range req.Modules {
		// Validate
		if moduleReq.Title == "" {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkModuleError{
				Title: moduleReq.Title,
				Error: "title is required",
			})
			continue
		}

		// Check for duplicate title
		if existingTitleMap[moduleReq.Title] {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkModuleError{
				Title: moduleReq.Title,
				Error: "module title already exists for this scheme of work",
			})
			continue
		}

		// Set module order
		moduleOrder := moduleReq.ModuleOrder
		if moduleOrder == 0 {
			currentOrder++
			moduleOrder = currentOrder
		}

		// Create module
		module := &models.Module{
			ID:             uuid.New(),
			SchemeOfWorkID: schemeOfWorkID,
			Title:          moduleReq.Title,
			Description:    moduleReq.Description,
			ModuleOrder:    moduleOrder,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := s.db.Create(module).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkModuleError{
				Title: moduleReq.Title,
				Error: "failed to create module: " + err.Error(),
			})
			continue
		}

		// Preload relationships
		if err := s.db.Preload("SchemeOfWork").First(module, module.ID).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkModuleError{
				Title: moduleReq.Title,
				Error: "failed to load module details",
			})
			continue
		}

		existingTitleMap[moduleReq.Title] = true
		result.SuccessCount++
		result.Created = append(result.Created, *s.toModuleResponse(module))
	}

	return result, nil
}

// GetAllModules retrieves all modules with pagination and filters
func (s *ModuleService) GetAllModules(params *dto.ModuleQueryParams) (*dto.ModuleListResponse, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = "module_order"
	}
	if params.SortOrder == "" {
		params.SortOrder = "asc"
	}

	// Build query
	query := s.db.Model(&models.Module{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.SchemeOfWorkID != "" {
		schemeOfWorkID, err := uuid.Parse(params.SchemeOfWorkID)
		if err == nil {
			query = query.Where("scheme_of_work_id = ?", schemeOfWorkID)
		}
	}

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", search, search)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count modules: %w", err)
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
	var modules []models.Module
	if err := query.Preload("SchemeOfWork").Find(&modules).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch modules: %w", err)
	}

	// Convert to response
	responses := make([]dto.ModuleResponse, len(modules))
	for i, module := range modules {
		responses[i] = *s.toModuleResponse(&module)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.ModuleListResponse{
		Modules:    responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetModuleByID retrieves a single module by ID
func (s *ModuleService) GetModuleByID(id string) (*dto.ModuleResponse, error) {
	moduleID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid module ID")
	}

	var module models.Module
	if err := s.db.Where("id = ? AND deleted_at IS NULL", moduleID).
		Preload("SchemeOfWork").
		Preload("Lessons").
		First(&module).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("module not found")
		}
		return nil, errors.New("failed to fetch module: " + err.Error())
	}

	return s.toModuleResponse(&module), nil
}

// GetModulesBySchemeOfWork retrieves all modules for a specific scheme of work
func (s *ModuleService) GetModulesBySchemeOfWork(schemeOfWorkID string) ([]dto.ModuleResponse, error) {
	sID, err := uuid.Parse(schemeOfWorkID)
	if err != nil {
		return nil, errors.New("invalid scheme of work ID")
	}

	var modules []models.Module
	if err := s.db.Where("scheme_of_work_id = ? AND deleted_at IS NULL", sID).
		Preload("SchemeOfWork").
		Order("module_order ASC").
		Find(&modules).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch modules: %w", err)
	}

	responses := make([]dto.ModuleResponse, len(modules))
	for i, module := range modules {
		responses[i] = *s.toModuleResponse(&module)
	}

	return responses, nil
}

// UpdateModule updates an existing module
func (s *ModuleService) UpdateModule(id string, req *dto.UpdateModuleRequest) (*dto.ModuleResponse, error) {
	moduleID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid module ID")
	}

	// Find existing module
	var module models.Module
	if err := s.db.Where("id = ? AND deleted_at IS NULL", moduleID).First(&module).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("module not found")
		}
		return nil, errors.New("failed to fetch module: " + err.Error())
	}

	// Update fields
	if req.Title != "" {
		// Check if title already exists for another module in the same scheme of work
		var existing models.Module
		if err := s.db.Where("scheme_of_work_id = ? AND title = ? AND id != ? AND deleted_at IS NULL",
			module.SchemeOfWorkID, req.Title, moduleID).First(&existing).Error; err == nil {
			return nil, errors.New("module title already exists for this scheme of work")
		}
		module.Title = req.Title
	}

	if req.Description != "" {
		module.Description = req.Description
	}

	if req.ModuleOrder > 0 && req.ModuleOrder != module.ModuleOrder {
		// Handle reordering
		if err := s.reorderModules(module.SchemeOfWorkID, module.ModuleOrder, req.ModuleOrder, moduleID); err != nil {
			return nil, err
		}
		module.ModuleOrder = req.ModuleOrder
	}

	module.UpdatedAt = time.Now()

	if err := s.db.Save(&module).Error; err != nil {
		return nil, errors.New("failed to update module: " + err.Error())
	}

	// Preload relationships
	if err := s.db.Preload("SchemeOfWork").First(&module, module.ID).Error; err != nil {
		return nil, errors.New("failed to load module details: " + err.Error())
	}

	return s.toModuleResponse(&module), nil
}

// ReorderModules reorders multiple modules
func (s *ModuleService) ReorderModules(req *dto.ReorderModulesRequest) error {
	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, orderReq := range req.ModuleOrders {
		moduleID, err := uuid.Parse(orderReq.ID)
		if err != nil {
			tx.Rollback()
			return errors.New("invalid module ID format: " + err.Error())
		}

		if err := tx.Model(&models.Module{}).
			Where("id = ? AND deleted_at IS NULL", moduleID).
			Update("module_order", orderReq.ModuleOrder).Error; err != nil {
			tx.Rollback()
			return errors.New("failed to update module order: " + err.Error())
		}
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("failed to commit reorder: " + err.Error())
	}

	return nil
}

// reorderModules reorders modules when a single module's order changes
func (s *ModuleService) reorderModules(schemeOfWorkID uuid.UUID, oldOrder, newOrder int, excludeID uuid.UUID) error {
	if oldOrder == newOrder {
		return nil
	}

	var modules []models.Module
	if err := s.db.Where("scheme_of_work_id = ? AND deleted_at IS NULL AND id != ?", schemeOfWorkID, excludeID).
		Order("module_order ASC").
		Find(&modules).Error; err != nil {
		return errors.New("failed to fetch modules for reordering: " + err.Error())
	}

	// Update orders
	if oldOrder < newOrder {
		// Moving down: decrease orders between oldOrder+1 and newOrder
		for _, m := range modules {
			if m.ModuleOrder > oldOrder && m.ModuleOrder <= newOrder {
				m.ModuleOrder--
				if err := s.db.Save(&m).Error; err != nil {
					return errors.New("failed to update module order: " + err.Error())
				}
			}
		}
	} else {
		// Moving up: increase orders between newOrder and oldOrder-1
		for _, m := range modules {
			if m.ModuleOrder >= newOrder && m.ModuleOrder < oldOrder {
				m.ModuleOrder++
				if err := s.db.Save(&m).Error; err != nil {
					return errors.New("failed to update module order: " + err.Error())
				}
			}
		}
	}

	return nil
}

// DeleteModule soft deletes a module
func (s *ModuleService) DeleteModule(id string) error {
	moduleID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid module ID")
	}

	var module models.Module
	if err := s.db.Where("id = ? AND deleted_at IS NULL", moduleID).First(&module).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("module not found")
		}
		return errors.New("failed to fetch module: " + err.Error())
	}

	if err := s.db.Delete(&module).Error; err != nil {
		return errors.New("failed to delete module: " + err.Error())
	}

	return nil
}

// validateModuleRequest validates the module request
func (s *ModuleService) validateModuleRequest(req *dto.CreateModuleRequest) error {
	if req.SchemeOfWorkID == "" {
		return errors.New("scheme of work ID is required")
	}
	if req.Title == "" {
		return errors.New("title is required")
	}
	if req.ModuleOrder < 0 {
		return errors.New("module order cannot be negative")
	}
	return nil
}

// toModuleResponse converts model to response DTO
func (s *ModuleService) toModuleResponse(module *models.Module) *dto.ModuleResponse {
	response := &dto.ModuleResponse{
		ID:             module.ID.String(),
		SchemeOfWorkID: module.SchemeOfWorkID.String(),
		Title:          module.Title,
		Description:    module.Description,
		ModuleOrder:    module.ModuleOrder,
		CreatedAt:      module.CreatedAt,
		UpdatedAt:      module.UpdatedAt,
	}

	// Add scheme of work details if preloaded
	if module.SchemeOfWork.ID != uuid.Nil {
		response.SchemeOfWork = &dto.SchemeOfWorkResponse{
			ID:    module.SchemeOfWork.ID.String(),
			Title: module.SchemeOfWork.Title,
			Term:  module.SchemeOfWork.Term,
			Status: module.SchemeOfWork.Status,
		}
	}

	// Add lessons if preloaded
	if len(module.Lessons) > 0 {
		lessons := make([]dto.LessonResponse, len(module.Lessons))
		for i, lesson := range module.Lessons {
			lessons[i] = dto.LessonResponse{
				ID:    lesson.ID.String(),
				Title: lesson.Title,
				Week:  lesson.Week,
				Status: lesson.Status,
			}
		}
		response.Lessons = lessons
	}

	return response
}