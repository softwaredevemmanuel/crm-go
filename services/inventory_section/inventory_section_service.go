// services/inventory_section_service.go
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

type InventorySectionService struct {
	db *gorm.DB
}

func NewInventorySectionService(db *gorm.DB) *InventorySectionService {
	return &InventorySectionService{db: db}
}

// CreateInventorySection creates a new inventory section
func (s *InventorySectionService) CreateInventorySection(req *dto.CreateInventorySectionRequest) (*dto.InventorySectionResponse, error) {
	// Validate input
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Check if code already exists
	var existingSection models.InventorySection
	if err := s.db.Where("code = ? AND deleted_at IS NULL", req.Code).First(&existingSection).Error; err == nil {
		return nil, errors.New("inventory section with this code already exists")
	}

	// Parse HeadOfDeptID if provided
	var headOfDeptID *uuid.UUID
	if req.HeadOfDeptID != "" {
		id, err := uuid.Parse(req.HeadOfDeptID)
		if err != nil {
			return nil, errors.New("invalid head of department ID format")
		}

		// Verify user exists
		var user models.User
		if err := s.db.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("head of department user not found")
			}
			return nil, fmt.Errorf("failed to verify user: %w", err)
		}
		headOfDeptID = &id
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	// Create the section
	section := &models.InventorySection{
		ID:           uuid.New(),
		Name:         strings.TrimSpace(req.Name),
		Code:         strings.ToUpper(strings.TrimSpace(req.Code)),
		Description:  req.Description,
		Location:     req.Location,
		HeadOfDeptID: headOfDeptID,
		Status:       models.InventoryStatus(status),
		Notes:        req.Notes,
	}

	if err := s.db.Create(section).Error; err != nil {
		return nil, fmt.Errorf("failed to create inventory section: %w", err)
	}

	// Preload relationships
	if err := s.db.Preload("HeadOfDept").First(section, section.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load section details: %w", err)
	}

	return s.toSectionResponse(section), nil
}

// GetInventorySectionByID retrieves a section by ID
func (s *InventorySectionService) GetInventorySectionByID(id string) (*dto.InventorySectionResponse, error) {
	sectionID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid section ID")
	}

	var section models.InventorySection
	if err := s.db.Where("id = ? AND deleted_at IS NULL", sectionID).
		Preload("HeadOfDept").
		First(&section).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory section not found")
		}
		return nil, fmt.Errorf("failed to fetch section: %w", err)
	}

	// Get item count
	var itemCount int64
	s.db.Model(&models.InventoryItem{}).
		Where("inventory_id = ? AND deleted_at IS NULL", sectionID).
		Count(&itemCount)

	response := s.toSectionResponse(&section)
	response.ItemCount = itemCount

	return response, nil
}

// GetInventorySections retrieves sections with filters and pagination
func (s *InventorySectionService) GetInventorySections(params *dto.InventorySectionQueryParams) (*dto.InventorySectionListResponse, error) {
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
	query := s.db.Model(&models.InventorySection{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ? OR description ILIKE ?", search, search, search)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.Location != "" {
		location := "%" + params.Location + "%"
		query = query.Where("location ILIKE ?", location)
	}

	// Count using explicit model without preload
	var total int64
	if err := s.db.Model(&models.InventorySection{}).
		Where("deleted_at IS NULL").
		Scopes(func(db *gorm.DB) *gorm.DB {
			if params.Search != "" {
				search := "%" + params.Search + "%"
				db = db.Where("name ILIKE ? OR code ILIKE ? OR description ILIKE ?", search, search, search)
			}
			if params.Status != "" {
				db = db.Where("status = ?", params.Status)
			}
			if params.Location != "" {
				location := "%" + params.Location + "%"
				db = db.Where("location ILIKE ?", location)
			}
			return db
		}).
		Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count sections: %w", err)
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
	// Execute query with preloads
	var sections []models.InventorySection
	if err := query.
		Preload("HeadOfDept").
		Find(&sections).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch sections: %w", err)
	}

	// Convert to response
	responses := make([]dto.InventorySectionResponse, len(sections))
	for i, section := range sections {
		response := s.toSectionResponse(&section)

		// Get item count for each section
		var itemCount int64
		s.db.Model(&models.InventoryItem{}).
			Where("inventory_id = ? AND deleted_at IS NULL", section.ID).
			Count(&itemCount)
		response.ItemCount = itemCount

		responses[i] = *response
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.InventorySectionListResponse{
		Sections:   responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// UpdateInventorySection updates an existing section
func (s *InventorySectionService) UpdateInventorySection(id string, req *dto.UpdateInventorySectionRequest) (*dto.InventorySectionResponse, error) {
	sectionID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid section ID")
	}

	// Find existing section
	var section models.InventorySection
	if err := s.db.Where("id = ? AND deleted_at IS NULL", sectionID).First(&section).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory section not found")
		}
		return nil, fmt.Errorf("failed to fetch section: %w", err)
	}

	// Update name
	if req.Name != "" {
		section.Name = strings.TrimSpace(req.Name)
	}

	// Update code (check for duplicates)
	if req.Code != "" {
		newCode := strings.ToUpper(strings.TrimSpace(req.Code))
		if newCode != section.Code {
			var existingSection models.InventorySection
			if err := s.db.Where("code = ? AND id != ? AND deleted_at IS NULL", newCode, sectionID).First(&existingSection).Error; err == nil {
				return nil, errors.New("inventory section with this code already exists")
			}
			section.Code = newCode
		}
	}

	// Update description
	if req.Description != "" {
		section.Description = req.Description
	}

	// Update location
	if req.Location != "" {
		section.Location = req.Location
	}

	// Update head of department
	if req.HeadOfDeptID != "" {
		id, err := uuid.Parse(req.HeadOfDeptID)
		if err != nil {
			return nil, errors.New("invalid head of department ID format")
		}

		// Verify user exists
		var user models.User
		if err := s.db.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("head of department user not found")
			}
			return nil, fmt.Errorf("failed to verify user: %w", err)
		}
		section.HeadOfDeptID = &id
	}

	// Update status
	if req.Status != "" {
		if req.Status != "active" && req.Status != "inactive" && req.Status != "closed" && req.Status != "renovation" {
			return nil, errors.New("status must be 'active', 'inactive', 'closed', or 'renovation'")
		}
		section.Status = models.InventoryStatus(req.Status)
	}

	// Update notes
	if req.Notes != "" {
		section.Notes = req.Notes
	}

	if err := s.db.Save(&section).Error; err != nil {
		return nil, fmt.Errorf("failed to update section: %w", err)
	}

	// Preload relationships
	if err := s.db.Preload("HeadOfDept").First(&section, section.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load section details: %w", err)
	}

	return s.toSectionResponse(&section), nil
}

// DeleteInventorySection soft deletes a section
func (s *InventorySectionService) DeleteInventorySection(id string) error {

	sectionID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid section ID")
	}
	// Find existing section
	var section models.InventorySection
	if err := s.db.Where("id = ? AND deleted_at IS NULL", sectionID).First(&section).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("inventory section not found")
		}
		return fmt.Errorf("failed to fetch section: %w", err)
	}

	// Check if section has items
	var itemCount int64
	if err := s.db.Model(&models.InventoryItem{}).
		Where("inventory_section_id = ? AND deleted_at IS NULL", sectionID).
		Count(&itemCount).Error; err != nil {
		return fmt.Errorf("failed to check items: %w", err)
	}

	if itemCount > 0 {
		return errors.New("cannot delete section with existing items. Please remove or transfer items first")
	}

	// Soft delete
	if err := s.db.Delete(&section).Error; err != nil {
		return fmt.Errorf("failed to delete section: %w", err)
	}

	return nil
}

// validateCreateRequest validates the create request
func (s *InventorySectionService) validateCreateRequest(req *dto.CreateInventorySectionRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(req.Code) == "" {
		return errors.New("code is required")
	}
	if req.Status != "" && req.Status != "active" && req.Status != "inactive" && req.Status != "closed" && req.Status != "renovation" {
		return errors.New("status must be 'active', 'inactive', 'closed', or 'renovation'")
	}
	return nil
}

// toSectionResponse converts model to response DTO
func (s *InventorySectionService) toSectionResponse(section *models.InventorySection) *dto.InventorySectionResponse {
	response := &dto.InventorySectionResponse{
		ID:          section.ID.String(),
		Name:        section.Name,
		Code:        section.Code,
		Description: section.Description,
		Location:    section.Location,
		Status:      string(section.Status),
		Notes:       section.Notes,
		CreatedAt:   section.CreatedAt,
		UpdatedAt:   section.UpdatedAt,
	}

	if section.HeadOfDeptID != nil {
		response.HeadOfDeptID = section.HeadOfDeptID.String()
	}

	// Add head of department details if preloaded
	if section.HeadOfDept != nil && section.HeadOfDept.ID != uuid.Nil {
		response.HeadOfDept = &dto.UserResponse{
			ID:        section.HeadOfDept.ID.String(),
			FirstName: section.HeadOfDept.FirstName,
			LastName:  section.HeadOfDept.LastName,
			Email:     section.HeadOfDept.Email,
			Phone:     section.HeadOfDept.Phone,
			Role:      section.HeadOfDept.Role,
			Position:  section.HeadOfDept.Position,
		}
	}

	return response
}