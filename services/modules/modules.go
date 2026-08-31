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

	// Parse subject ID
	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		return nil, errors.New("invalid subject ID format")
	}

	// Check if subject exists
	var subject models.Subject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subject not found")
		}
		return nil, errors.New("failed to verify subject: " + err.Error())
	}

	// Check if module name already exists for this subject
	var existing models.Module
	if err := s.db.Where("subject_id = ? AND name = ? AND deleted_at IS NULL", 
		subjectID, req.Name).First(&existing).Error; err == nil {
		return nil, errors.New("module name already exists for this subject")
	}

	// Check if module code already exists for this subject
	if req.Code != "" {
		if err := s.db.Where("subject_id = ? AND code = ? AND deleted_at IS NULL", 
			subjectID, req.Code).First(&existing).Error; err == nil {
			return nil, errors.New("module code already exists for this subject")
		}
	}

	// Set default sequence if not provided
	sequence := req.Sequence
	if sequence == 0 {
		// Get the highest sequence number for this subject
		var maxSequence int
		if err := s.db.Model(&models.Module{}).
			Where("subject_id = ? AND deleted_at IS NULL", subjectID).
			Select("COALESCE(MAX(sequence), 0)").
			Scan(&maxSequence).Error; err != nil {
			return nil, errors.New("failed to determine sequence: " + err.Error())
		}
		sequence = maxSequence + 1
	}

	// Create module
	module := &models.Module{
		ID:          uuid.New(),
		SubjectID:   subjectID,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Sequence:    sequence,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(module).Error; err != nil {
		return nil, errors.New("failed to create module: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("Subject").First(module, module.ID).Error; err != nil {
		return nil, errors.New("failed to load module details: " + err.Error())
	}

	return s.toModuleResponse(module), nil
}

// BulkCreateModules creates multiple modules
func (s *ModuleService) BulkCreateModules(req *dto.BulkCreateModulesRequest) (*dto.BulkModuleResult, error) {
	// Parse subject ID
	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		return nil, errors.New("invalid subject ID format")
	}

	// Check if subject exists
	var subject models.Subject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subject not found")
		}
		return nil, errors.New("failed to verify subject: " + err.Error())
	}

	result := &dto.BulkModuleResult{
		Created: []dto.ModuleResponse{},
		Errors:  []dto.BulkModuleError{},
	}

	// Get existing module names and codes to avoid duplicates
	var existingModules []models.Module
	if err := s.db.Where("subject_id = ? AND deleted_at IS NULL", subjectID).Find(&existingModules).Error; err != nil {
		return nil, fmt.Errorf("failed to check existing modules: %w", err)
	}

	existingNameMap := make(map[string]bool)
	existingCodeMap := make(map[string]bool)
	for _, m := range existingModules {
		existingNameMap[m.Name] = true
		if m.Code != "" {
			existingCodeMap[m.Code] = true
		}
	}

	// Get max sequence
	var maxSequence int
	if err := s.db.Model(&models.Module{}).
		Where("subject_id = ? AND deleted_at IS NULL", subjectID).
		Select("COALESCE(MAX(sequence), 0)").
		Scan(&maxSequence).Error; err != nil {
		return nil, errors.New("failed to determine sequence: " + err.Error())
	}
	currentSequence := maxSequence

	for _, moduleReq := range req.Modules {
		// Validate
		if moduleReq.Name == "" {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkModuleError{
				Name:  moduleReq.Name,
				Code:  moduleReq.Code,
				Error: "name is required",
			})
			continue
		}

		// Check for duplicate name
		if existingNameMap[moduleReq.Name] {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkModuleError{
				Name:  moduleReq.Name,
				Code:  moduleReq.Code,
				Error: "module name already exists for this subject",
			})
			continue
		}

		// Check for duplicate code
		if moduleReq.Code != "" && existingCodeMap[moduleReq.Code] {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkModuleError{
				Name:  moduleReq.Name,
				Code:  moduleReq.Code,
				Error: "module code already exists for this subject",
			})
			continue
		}

		// Set sequence
		sequence := moduleReq.Sequence
		if sequence == 0 {
			currentSequence++
			sequence = currentSequence
		}

		// Create module
		module := &models.Module{
			ID:          uuid.New(),
			SubjectID:   subjectID,
			Name:        moduleReq.Name,
			Code:        moduleReq.Code,
			Description: moduleReq.Description,
			Sequence:    sequence,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := s.db.Create(module).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkModuleError{
				Name:  moduleReq.Name,
				Code:  moduleReq.Code,
				Error: "failed to create module: " + err.Error(),
			})
			continue
		}

		// Preload relationships
		if err := s.db.Preload("Subject").First(module, module.ID).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkModuleError{
				Name:  moduleReq.Name,
				Code:  moduleReq.Code,
				Error: "failed to load module details",
			})
			continue
		}

		existingNameMap[moduleReq.Name] = true
		if moduleReq.Code != "" {
			existingCodeMap[moduleReq.Code] = true
		}
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
		params.SortBy = "sequence"
	}
	if params.SortOrder == "" {
		params.SortOrder = "asc"
	}

	// Build query
	query := s.db.Model(&models.Module{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.SubjectID != "" {
		subjectID, err := uuid.Parse(params.SubjectID)
		if err == nil {
			query = query.Where("subject_id = ?", subjectID)
		}
	}

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ? OR LOWER(description) LIKE ?",
			search, search, search)
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
	if err := query.Preload("Subject").Find(&modules).Error; err != nil {
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
		Preload("Subject").
		First(&module).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("module not found")
		}
		return nil, errors.New("failed to fetch module: " + err.Error())
	}

	return s.toModuleResponse(&module), nil
}

// GetModulesBySubject retrieves all modules for a specific subject
func (s *ModuleService) GetModulesBySubject(subjectID string) ([]dto.ModuleResponse, error) {
	sID, err := uuid.Parse(subjectID)
	if err != nil {
		return nil, errors.New("invalid subject ID")
	}

	var modules []models.Module
	if err := s.db.Where("subject_id = ? AND deleted_at IS NULL", sID).
		Preload("Subject").
		Order("sequence ASC").
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
	if req.Name != "" {
		// Check if name already exists for another module in the same subject
		var existing models.Module
		if err := s.db.Where("subject_id = ? AND name = ? AND id != ? AND deleted_at IS NULL",
			module.SubjectID, req.Name, moduleID).First(&existing).Error; err == nil {
			return nil, errors.New("module name already exists for this subject")
		}
		module.Name = req.Name
	}

	if req.Code != "" {
		// Check if code already exists for another module in the same subject
		var existing models.Module
		if err := s.db.Where("subject_id = ? AND code = ? AND id != ? AND deleted_at IS NULL",
			module.SubjectID, req.Code, moduleID).First(&existing).Error; err == nil {
			return nil, errors.New("module code already exists for this subject")
		}
		module.Code = req.Code
	}

	if req.Description != "" {
		module.Description = req.Description
	}

	if req.Sequence > 0 {
		// If sequence is changing, handle reordering
		if req.Sequence != module.Sequence {
			if err := s.reorderModules(module.SubjectID, module.Sequence, req.Sequence, moduleID); err != nil {
				return nil, err
			}
			module.Sequence = req.Sequence
		}
	}

	module.UpdatedAt = time.Now()

	if err := s.db.Save(&module).Error; err != nil {
		return nil, errors.New("failed to update module: " + err.Error())
	}

	// Preload relationships
	if err := s.db.Preload("Subject").First(&module, module.ID).Error; err != nil {
		return nil, errors.New("failed to load module details: " + err.Error())
	}

	return s.toModuleResponse(&module), nil
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

// reorderModules reorders modules when sequence changes
func (s *ModuleService) reorderModules(subjectID uuid.UUID, oldSeq, newSeq int, excludeID uuid.UUID) error {
	if oldSeq == newSeq {
		return nil
	}

	var modules []models.Module
	if err := s.db.Where("subject_id = ? AND deleted_at IS NULL AND id != ?", subjectID, excludeID).
		Order("sequence ASC").
		Find(&modules).Error; err != nil {
		return errors.New("failed to fetch modules for reordering: " + err.Error())
	}

	// Update sequences
	if oldSeq < newSeq {
		// Moving down: decrease sequences between oldSeq+1 and newSeq
		for _, m := range modules {
			if m.Sequence > oldSeq && m.Sequence <= newSeq {
				m.Sequence--
				if err := s.db.Save(&m).Error; err != nil {
					return errors.New("failed to update module sequence: " + err.Error())
				}
			}
		}
	} else {
		// Moving up: increase sequences between newSeq and oldSeq-1
		for _, m := range modules {
			if m.Sequence >= newSeq && m.Sequence < oldSeq {
				m.Sequence++
				if err := s.db.Save(&m).Error; err != nil {
					return errors.New("failed to update module sequence: " + err.Error())
				}
			}
		}
	}

	return nil
}

// validateModuleRequest validates the module request
func (s *ModuleService) validateModuleRequest(req *dto.CreateModuleRequest) error {
	if req.SubjectID == "" {
		return errors.New("subject ID is required")
	}
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.Sequence < 0 {
		return errors.New("sequence cannot be negative")
	}
	return nil
}

// toModuleResponse converts model to response DTO
func (s *ModuleService) toModuleResponse(module *models.Module) *dto.ModuleResponse {
	response := &dto.ModuleResponse{
		ID:          module.ID.String(),
		SubjectID:   module.SubjectID.String(),
		Name:        module.Name,
		Code:        module.Code,
		Description: module.Description,
		Sequence:    module.Sequence,
		CreatedAt:   module.CreatedAt,
		UpdatedAt:   module.UpdatedAt,
	}

	// Add subject details if preloaded
	if module.Subject.ID != uuid.Nil {
		response.Subject = &dto.SubjectResponse{
			ID:          module.Subject.ID.String(),
			Name:        module.Subject.Name,
			Code:        module.Subject.Code,
			Description: module.Subject.Description,
		}
	}

	return response
}