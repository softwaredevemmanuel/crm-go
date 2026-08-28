// services/address_service.go
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

type AddressService struct {
	db *gorm.DB
}

func NewAddressService(db *gorm.DB) *AddressService {
	return &AddressService{db: db}
}

// CreateAddress creates a new address
func (s *AddressService) CreateAddress(req *dto.CreateAddressRequest, userID uuid.UUID) (*dto.AddressResponse, error) {
	// Validate input
	if err := s.validateAddressRequest(req); err != nil {
		return nil, err
	}

	// Parse User ID
	studentID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Check if user exists
	var user models.User
	if err := s.db.Where("id = ? AND deleted_at IS NULL", studentID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to verify user: " + err.Error())
	}

	// Set default values
	status := req.Status
	if status == "" {
		status = "active"
	}

	addressType := req.AddressType
	if addressType == "" {
		addressType = "home"
	}

	country := req.Country
	if country == "" {
		country = "Nigeria"
	}

	// If this is primary, unset any existing primary addresses for this user
	if req.IsPrimary {
		if err := s.db.Model(&models.Address{}).
			Where("user_id = ? AND is_primary = ? AND deleted_at IS NULL", studentID, true).
			Update("is_primary", false).Error; err != nil {
			return nil, errors.New("failed to update primary addresses: " + err.Error())
		}
	}

	// Create new address
	address := &models.Address{
		ID:          uuid.New(),
		UserID:      studentID,
		Address:     strings.TrimSpace(req.Address),
		City:        strings.TrimSpace(req.City),
		State:       strings.TrimSpace(req.State),
		Country:     strings.TrimSpace(country),
		PostalCode:  strings.TrimSpace(req.PostalCode),
		AddressType: addressType,
		IsPrimary:   req.IsPrimary,
		Status:      status,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(address).Error; err != nil {
		return nil, errors.New("failed to create address: " + err.Error())
	}

	// Preload user for response
	if err := s.db.Preload("User").First(address, address.ID).Error; err != nil {
		return nil, errors.New("failed to load address details: " + err.Error())
	}

	return s.toAddressResponse(address), nil
}

// GetAllAddresses retrieves all addresses with pagination and filters
func (s *AddressService) GetAllAddresses(params *dto.AddressQueryParams) (*dto.AddressListResponse, error) {
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
	query := s.db.Model(&models.Address{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(params.Search) + "%"
		query = query.Joins("JOIN users ON users.id = addresses.user_id").
			Where("LOWER(addresses.address) LIKE ? OR LOWER(addresses.city) LIKE ? OR LOWER(addresses.state) LIKE ? OR LOWER(users.first_name) LIKE ? OR LOWER(users.last_name) LIKE ?",
				searchTerm, searchTerm, searchTerm, searchTerm, searchTerm)
	}

	if params.UserID != "" {
		userID, err := uuid.Parse(params.UserID)
		if err == nil {
			query = query.Where("user_id = ?", userID)
		}
	}

	if params.AddressType != "" {
		query = query.Where("address_type = ?", params.AddressType)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.IsPrimary != nil {
		query = query.Where("is_primary = ?", *params.IsPrimary)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count addresses: %w", err)
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
	var addresses []models.Address
	if err := query.Preload("User").Find(&addresses).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch addresses: %w", err)
	}

	// Convert to response
	responses := make([]dto.AddressResponse, len(addresses))
	for i, address := range addresses {
		responses[i] = *s.toAddressResponse(&address)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.AddressListResponse{
		Addresses:  responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetAddressByID retrieves a single address by ID
func (s *AddressService) GetAddressByID(id string) (*dto.AddressResponse, error) {
	addressID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid address ID")
	}

	var address models.Address
	if err := s.db.Where("id = ? AND deleted_at IS NULL", addressID).
		Preload("User").
		First(&address).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("address not found")
		}
		return nil, errors.New("failed to fetch address: " + err.Error())
	}

	return s.toAddressResponse(&address), nil
}

// GetAddressesByUser retrieves all addresses for a specific user
func (s *AddressService) GetAddressesByUser(userID string) ([]dto.AddressResponse, error) {
	uID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	var addresses []models.Address
	if err := s.db.Where("user_id = ? AND deleted_at IS NULL", uID).
		Preload("User").
		Order("is_primary DESC, created_at ASC").
		Find(&addresses).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch addresses for user: %w", err)
	}

	responses := make([]dto.AddressResponse, len(addresses))
	for i, address := range addresses {
		responses[i] = *s.toAddressResponse(&address)
	}

	return responses, nil
}

// GetPrimaryAddressByUser retrieves the primary address for a user
func (s *AddressService) GetPrimaryAddressByUser(userID string) (*dto.AddressResponse, error) {
	uID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	var address models.Address
	if err := s.db.Where("user_id = ? AND is_primary = ? AND deleted_at IS NULL", uID, true).
		Preload("User").
		First(&address).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no primary address found for this user")
		}
		return nil, errors.New("failed to fetch primary address: " + err.Error())
	}

	return s.toAddressResponse(&address), nil
}

// UpdateAddress updates an existing address
func (s *AddressService) UpdateAddress(id string, req *dto.UpdateAddressRequest, userID uuid.UUID) (*dto.AddressResponse, error) {
	addressID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid address ID")
	}

	// Find existing address
	var address models.Address
	if err := s.db.Where("id = ? AND deleted_at IS NULL", addressID).First(&address).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("address not found")
		}
		return nil, errors.New("failed to fetch address: " + err.Error())
	}

	// If this is being set as primary, unset any existing primary addresses for this user
	if req.IsPrimary != nil && *req.IsPrimary && !address.IsPrimary {
		if err := s.db.Model(&models.Address{}).
			Where("user_id = ? AND is_primary = ? AND id != ? AND deleted_at IS NULL", 
				address.UserID, true, addressID).
			Update("is_primary", false).Error; err != nil {
			return nil, errors.New("failed to update primary addresses: " + err.Error())
		}
	}

	// Update fields
	if req.Address != "" {
		address.Address = strings.TrimSpace(req.Address)
	}
	if req.City != "" {
		address.City = strings.TrimSpace(req.City)
	}
	if req.State != "" {
		address.State = strings.TrimSpace(req.State)
	}
	if req.Country != "" {
		address.Country = strings.TrimSpace(req.Country)
	}
	if req.PostalCode != "" {
		address.PostalCode = strings.TrimSpace(req.PostalCode)
	}
	if req.AddressType != "" {
		address.AddressType = req.AddressType
	}
	if req.IsPrimary != nil {
		address.IsPrimary = *req.IsPrimary
	}
	if req.Status != "" {
		address.Status = req.Status
	}

	address.UpdatedAt = time.Now()

	if err := s.db.Save(&address).Error; err != nil {
		return nil, errors.New("failed to update address: " + err.Error())
	}

	// Preload user for response
	if err := s.db.Preload("User").First(&address, address.ID).Error; err != nil {
		return nil, errors.New("failed to load address details: " + err.Error())
	}

	return s.toAddressResponse(&address), nil
}

// DeleteAddress soft deletes an address
func (s *AddressService) DeleteAddress(id string) error {
	addressID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid address ID")
	}

	var address models.Address
	if err := s.db.Where("id = ? AND deleted_at IS NULL", addressID).First(&address).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("address not found")
		}
		return errors.New("failed to fetch address: " + err.Error())
	}

	if err := s.db.Delete(&address).Error; err != nil {
		return errors.New("failed to delete address: " + err.Error())
	}

	return nil
}

// validateAddressRequest validates the address request
func (s *AddressService) validateAddressRequest(req *dto.CreateAddressRequest) error {
	if req.UserID == "" {
		return errors.New("user ID is required")
	}
	if req.Address == "" {
		return errors.New("address is required")
	}
	if req.AddressType != "" && req.AddressType != "home" && req.AddressType != "school" && 
		req.AddressType != "office" && req.AddressType != "other" {
		return errors.New("address type must be 'home', 'school', 'office', or 'other'")
	}
	if req.Status != "" && req.Status != "active" && req.Status != "inactive" {
		return errors.New("status must be 'active' or 'inactive'")
	}
	return nil
}

// toAddressResponse converts model to response DTO
func (s *AddressService) toAddressResponse(address *models.Address) *dto.AddressResponse {
	response := &dto.AddressResponse{
		ID:          address.ID.String(),
		UserID:      address.UserID.String(),
		Address:     address.Address,
		City:        address.City,
		State:       address.State,
		Country:     address.Country,
		PostalCode:  address.PostalCode,
		AddressType: address.AddressType,
		IsPrimary:   address.IsPrimary,
		Status:      address.Status,
		CreatedBy:   address.CreatedBy.String(),
		CreatedAt:   address.CreatedAt,
		UpdatedAt:   address.UpdatedAt,
	}

	// Add user details if preloaded
	if address.User.ID != uuid.Nil {
		response.User = &dto.UserResponse{
			ID:        address.User.ID.String(),
			FirstName: address.User.FirstName,
			LastName:  address.User.LastName,
			Email:     address.User.Email,
			Phone:     address.User.Phone,
			Role:      address.User.Role,
			Position:  address.User.Position,
		}
	}

	return response
}