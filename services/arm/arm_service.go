// services/arm_service.go
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

type ArmService struct {
	db *gorm.DB
}

func NewArmService(db *gorm.DB) *ArmService {
	return &ArmService{db: db}
}

// CreateArm creates a new arm
func (s *ArmService) CreateArm(req *dto.CreateArmRequest, userID uuid.UUID) (*dto.ArmResponse, error) {
	// Validate input
	if err := s.validateArmRequest(req); err != nil {
		return nil, err
	}

	// Parse Grade ID
	gradeID, err := uuid.Parse(req.GradeID)
	if err != nil {
		return nil, errors.New("invalid grade ID")
	}

	// Check if grade exists
	var grade models.ClassGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", gradeID).First(&grade).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("grade not found")
		}
		return nil, errors.New("failed to verify grade: " + err.Error())
	}

	// Generate code if not provided
	code := req.Code
	if code == "" {
		code = strings.ToUpper(req.Name[:3])
	}

	// Check if arm with same code already exists
	var existing models.Arm

	// Check if arm with same name and grade already exists
	if err := s.db.Where("name = ? AND grade_id = ?", req.Name, gradeID).First(&existing).Error; err == nil {
		return nil, errors.New("Grade exists with same name")
	}
	// Check if arm with same code and grade already exists
	if err := s.db.Where("code = ? AND grade_id = ?", req.Code, gradeID).First(&existing).Error; err == nil {
		return nil, errors.New("Grade exits with the same code")
	}

	// Check if arm with same name, code, and grade already exists
	if err := s.db.
		Where(
			"name = ? AND code = ? AND grade_id = ?",
			req.Name,
			req.Code,
			gradeID,
		).
		First(&existing).Error; err == nil {

		return nil, errors.New("arm with this name and code already exists for this grade")
	}

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = "active"
	}

	// Set default capacity if not provided
	capacity := req.Capacity
	if capacity == 0 {
		capacity = 30
	}

	// Create new arm
	arm := &models.Arm{
		ID:          uuid.New(),
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Code:        strings.ToUpper(strings.TrimSpace(code)),
		GradeID:     gradeID,
		Status:      status,
		Capacity:    capacity,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to database
	if err := s.db.Create(arm).Error; err != nil {
		return nil, errors.New("failed to create arm: " + err.Error())
	}

	// Preload grade for response
	if err := s.db.Preload("Grade").First(arm, arm.ID).Error; err != nil {
		return nil, errors.New("failed to load arm details: " + err.Error())
	}

	// Convert to response DTO
	return s.toArmResponse(arm), nil
}

// GetAllArms retrieves all arms with pagination and filters
func (s *ArmService) GetAllArms(params *dto.ArmQueryParams) (*dto.ArmListResponse, error) {
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
	query := s.db.Model(&models.Arm{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where(
			"LOWER(name) LIKE ? OR LOWER(code) LIKE ? OR LOWER(description) LIKE ?",
			searchTerm, searchTerm, searchTerm,
		)
	}

	if params.GradeID != "" {
		gradeID, err := uuid.Parse(params.GradeID)
		if err == nil {
			query = query.Where("grade_id = ?", gradeID)
		}
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count arms: %w", err)
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
	var arms []models.Arm
	if err := query.Preload("Grade").Find(&arms).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch arms: %w", err)
	}

	// Convert to response
	responses := make([]dto.ArmResponse, len(arms))
	for i, arm := range arms {
		responses[i] = *s.toArmResponse(&arm)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.ArmListResponse{
		Arms:       responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetArmByID retrieves a single arm by ID
func (s *ArmService) GetArmByID(id string) (*dto.ArmResponse, error) {
	armID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid arm ID")
	}

	var arm models.Arm
	if err := s.db.Where("id = ? AND deleted_at IS NULL", armID).Preload("Grade").First(&arm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("arm not found")
		}
		return nil, errors.New("failed to fetch arm: " + err.Error())
	}

	return s.toArmResponse(&arm), nil
}

// GetArmsByGrade retrieves all arms for a specific grade
func (s *ArmService) GetArmsByGrade(gradeID string) ([]dto.ArmResponse, error) {
	gID, err := uuid.Parse(gradeID)
	if err != nil {
		return nil, errors.New("invalid grade ID")
	}

	var arms []models.Arm
	if err := s.db.Where("grade_id = ? AND deleted_at IS NULL", gID).
		Preload("Grade").
		Order("name ASC").
		Find(&arms).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch arms for grade: %w", err)
	}

	responses := make([]dto.ArmResponse, len(arms))
	for i, arm := range arms {
		responses[i] = *s.toArmResponse(&arm)
	}

	return responses, nil
}

// UpdateArm updates an existing arm
func (s *ArmService) UpdateArm(id string, req *dto.UpdateArmRequest, userID uuid.UUID) (*dto.ArmResponse, error) {
	armID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid arm ID")
	}

	// Find existing arm
	var arm models.Arm
	if err := s.db.Where("id = ? AND deleted_at IS NULL", armID).First(&arm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("arm not found")
		}
		return nil, errors.New("failed to fetch arm: " + err.Error())
	}

	// Check for duplicate name within same grade if being updated
	if req.Name != "" && req.Name != arm.Name {
		gradeID := arm.GradeID
		if req.GradeID != "" {
			gID, err := uuid.Parse(req.GradeID)
			if err != nil {
				return nil, errors.New("invalid grade ID")
			}
			gradeID = gID
		}

		var existing models.Arm
		if err := s.db.Where("name = ? AND grade_id = ? AND id != ? AND deleted_at IS NULL", req.Name, gradeID, armID).First(&existing).Error; err == nil {
			return nil, errors.New("arm with this name already exists for this grade")
		}
	}

	// Check for duplicate name within same grade if being updated
	if req.Code != "" && req.Code != arm.Code {
		gradeID := arm.GradeID
		if req.GradeID != "" {
			gID, err := uuid.Parse(req.GradeID)
			if err != nil {
				return nil, errors.New("invalid grade ID")
			}
			gradeID = gID
		}

		var existing models.Arm
		if err := s.db.Where("code = ? AND grade_id = ? AND id != ? AND deleted_at IS NULL", req.Code, gradeID, armID).First(&existing).Error; err == nil {
			return nil, errors.New("arm with this code already exists for this grade")
		}
	}

	// Update fields
	if req.Name != "" {
		arm.Name = strings.TrimSpace(req.Name)
	}
	if req.Description != "" {
		arm.Description = strings.TrimSpace(req.Description)
	}
	if req.Code != "" {
		arm.Code = strings.ToUpper(strings.TrimSpace(req.Code))
	}
	if req.GradeID != "" {
		gID, err := uuid.Parse(req.GradeID)
		if err != nil {
			return nil, errors.New("invalid grade ID")
		}
		// Verify grade exists
		var grade models.ClassGrade
		if err := s.db.Where("id = ? AND deleted_at IS NULL", gID).First(&grade).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("grade not found")
			}
			return nil, errors.New("failed to verify grade: " + err.Error())
		}
		arm.GradeID = gID
	}
	if req.Capacity > 0 {
		arm.Capacity = req.Capacity
	}
	if req.Status != "" {
		arm.Status = req.Status
	}

	// Update timestamp
	arm.UpdatedAt = time.Now()

	// Save to database
	if err := s.db.Save(&arm).Error; err != nil {
		return nil, errors.New("failed to update arm: " + err.Error())
	}

	// Preload grade for response
	if err := s.db.Preload("Grade").First(&arm, arm.ID).Error; err != nil {
		return nil, errors.New("failed to load arm details: " + err.Error())
	}

	return s.toArmResponse(&arm), nil
}

// DeleteArm soft deletes an arm
func (s *ArmService) DeleteArm(id string) error {
	armID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid arm ID")
	}

	var arm models.Arm
	if err := s.db.Where("id = ? AND deleted_at IS NULL", armID).First(&arm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("arm not found")
		}
		return errors.New("failed to fetch arm: " + err.Error())
	}

	// Check if arm has students (if you have a Student model)
	// Uncomment if you have a Student model with arm_id
	// var studentCount int64
	// if err := s.db.Model(&models.Student{}).Where("arm_id = ?", armID).Count(&studentCount).Error; err != nil {
	// 	return errors.New("failed to check arm usage: " + err.Error())
	// }
	// if studentCount > 0 {
	// 	return errors.New("cannot delete arm: it has students assigned")
	// }

	if err := s.db.Delete(&arm).Error; err != nil {
		return errors.New("failed to delete arm: " + err.Error())
	}

	return nil
}

// DeleteArmPermanently permanently deletes an arm from the database
func (s *ArmService) DeleteArmPermanently(id string) error {
	armID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid arm ID")
	}

	var arm models.Arm

	// Find the arm, including soft-deleted records
	if err := s.db.
		Unscoped().
		Where("id = ?", armID).
		First(&arm).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("arm not found")
		}

		return errors.New("failed to fetch arm: " + err.Error())
	}

	// Permanently delete the arm
	if err := s.db.
		Unscoped().
		Delete(&arm).Error; err != nil {

		return errors.New("failed to permanently delete arm: " + err.Error())
	}

	return nil
}

// validateArmRequest validates the arm request
func (s *ArmService) validateArmRequest(req *dto.CreateArmRequest) error {
	if req.Name == "" {
		return errors.New("arm name is required")
	}
	if len(req.Name) < 2 {
		return errors.New("arm name must be at least 2 characters")
	}
	if req.GradeID == "" {
		return errors.New("grade ID is required")
	}
	if req.Capacity < 1 {
		return errors.New("capacity must be at least 1")
	}
	if req.Capacity > 100 {
		return errors.New("capacity cannot exceed 100")
	}
	if req.Status != "" && req.Status != "active" && req.Status != "inactive" && req.Status != "archived" {
		return errors.New("status must be 'active', 'inactive', or 'archived'")
	}
	return nil
}

// toArmResponse converts model to response DTO
func (s *ArmService) toArmResponse(arm *models.Arm) *dto.ArmResponse {
	response := &dto.ArmResponse{
		ID:          arm.ID.String(),
		Name:        arm.Name,
		Description: arm.Description,
		Code:        arm.Code,
		GradeID:     arm.GradeID.String(),
		Status:      arm.Status,
		Capacity:    arm.Capacity,
		CreatedBy:   arm.CreatedBy.String(),
		CreatedAt:   arm.CreatedAt,
		UpdatedAt:   arm.UpdatedAt,
	}

	// Add grade details if preloaded
	if arm.Grade.ID != uuid.Nil {
		response.Grade = &dto.ClassGradeResponse{
			ID:           arm.Grade.ID.String(),
			Name:         arm.Grade.Name,
			Code:         arm.Grade.Code,
			Level:        arm.Grade.Level,
			Description:  arm.Grade.Description,
			AcademicSessionID: string(arm.Grade.AcademicSessionID.String()),
			Capacity:     arm.Grade.Capacity,
			Status:       arm.Grade.Status,
			CreatedAt:    arm.Grade.CreatedAt,
			UpdatedAt:    arm.Grade.UpdatedAt,
		}
	}

	return response
}
