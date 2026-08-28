// services/user_service.go
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

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// GetAllUsers retrieves all users with pagination and filters
func (s *UserService) GetAllUsers(params *dto.UserQueryParams) (*dto.UserListResponse, error) {
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
	query := s.db.Model(&models.User{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where(
			"LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(email) LIKE ? OR LOWER(phone) LIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	if params.Role != "" {
		query = query.Where("role = ?", params.Role)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	if params.IsVerified != nil {
		query = query.Where("is_verified = ?", *params.IsVerified)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
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

	// Execute query
	var users []models.User
	if err := query.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	// Convert to response
	responses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = s.toUserResponse(&user)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.UserListResponse{
		Users:      responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetUserByID retrieves a single user by ID
func (s *UserService) GetUserByID(id string) (*dto.UserResponse, error) {
	// Parse UUID
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Find user
	var user models.User
	if err := s.db.Where("id = ? AND deleted_at IS NULL", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Convert to response
	response := s.toUserResponse(&user)
	return &response, nil
}

// GetUsersByRole retrieves all users with a specific role
func (s *UserService) GetUsersByRole(role string) ([]dto.UserResponse, error) {
	var users []models.User
	if err := s.db.Where("role = ? AND deleted_at IS NULL", role).
		Order("first_name ASC").
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch users by role: %w", err)
	}

	responses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = s.toUserResponse(&user)
	}

	return responses, nil
}


// GetUsersByPosition retrieves all users with a specific position
func (s *UserService) GetUsersByPosition(position string) ([]dto.UserResponse, error) {
	var users []models.User
	if err := s.db.Where("position = ? AND deleted_at IS NULL", position).
		Order("first_name ASC").
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch users by position: %w", err)
	}

	responses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = s.toUserResponse(&user)
	}

	return responses, nil
}

// ✅ UpdateUser updates an existing user
func (s *UserService) UpdateUser(id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// Parse UUID
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Find existing user
	var user models.User
	if err := s.db.Where("id = ? AND deleted_at IS NULL", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to fetch user: " + err.Error())
	}

	// Update fields
	if req.FirstName != "" {
		user.FirstName = strings.TrimSpace(req.FirstName)
	}
	if req.LastName != "" {
		user.LastName = strings.TrimSpace(req.LastName)
	}
	if req.MiddleName != "" {
		user.MiddleName = strings.TrimSpace(req.MiddleName)
	}
	if req.Phone != "" {
		user.Phone = strings.TrimSpace(req.Phone)
	}
	if req.Role != "" {
		user.Role = req.Role
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if req.IsVerified != nil {
		user.IsVerified = *req.IsVerified
	}
	if req.Location != "" {
		user.Location = strings.TrimSpace(req.Location)
	}

	// Update timestamp
	user.UpdatedAt = time.Now()

	// Save to database
	if err := s.db.Save(&user).Error; err != nil {
		return nil, errors.New("failed to update user: " + err.Error())
	}

	// Convert to response
	response := s.toUserResponse(&user)
	return &response, nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid user ID")
	}

	var user models.User
	if err := s.db.Where("id = ? AND deleted_at IS NULL", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return errors.New("failed to fetch user: " + err.Error())
	}

	// Check if user has associated guardians
	var guardianCount int64
	if err := s.db.Model(&models.Guardian{}).Where("user_id = ? AND deleted_at IS NULL", userID).Count(&guardianCount).Error; err != nil {
		return errors.New("failed to check user associations: " + err.Error())
	}
	if guardianCount > 0 {
		return errors.New("cannot delete user: user is associated with guardians")
	}

	// Check if user is a student with guardians
	var studentGuardianCount int64
	if err := s.db.Model(&models.Guardian{}).Where("student_id = ? AND deleted_at IS NULL", userID).Count(&studentGuardianCount).Error; err != nil {
		return errors.New("failed to check student associations: " + err.Error())
	}
	if studentGuardianCount > 0 {
		return errors.New("cannot delete user: student has guardians assigned")
	}

	if err := s.db.Delete(&user).Error; err != nil {
		return errors.New("failed to delete user: " + err.Error())
	}

	return nil
}

// toUserResponse converts model to response DTO
func (s *UserService) toUserResponse(user *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:          user.ID.String(),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		MiddleName:  user.MiddleName,
		FullName:    strings.TrimSpace(user.FirstName + " " + user.MiddleName + " " + user.LastName),
		Email:       user.Email,
		Phone:       user.Phone,
		Role:        user.Role,
		Position:        user.Position,
		Picture:     user.Picture,
		IsVerified:  user.IsVerified,
		IsActive:    user.IsActive,
		Location:    user.Location,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}