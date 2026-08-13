// services/department_service.go
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

type DepartmentService struct {
	db *gorm.DB
}

func NewDepartmentService(db *gorm.DB) *DepartmentService {
	return &DepartmentService{db: db}
}

// ============ 1. CREATE DEPARTMENT ============
func (s *DepartmentService) CreateDepartment(req *dto.CreateDepartmentRequest, userID uuid.UUID) (*dto.DepartmentWithSubjectsResponse, error) {
	// Validate input
	if err := s.validateDepartmentRequest(req); err != nil {
		return nil, err
	}

	// Check duplicate name
	var existing models.Department
	if err := s.db.Where("name = ?", req.Name).First(&existing).Error; err == nil {
		return nil, errors.New("department with this name already exists")
	}

	// Check duplicate code
	if err := s.db.Where("code = ?", req.Code).First(&existing).Error; err == nil {
		return nil, errors.New("department with this code already exists")
	}

	// Parse HeadID if provided
	var headID *uuid.UUID
	if req.HeadID != "" {
		parsed, err := uuid.Parse(req.HeadID)
		if err != nil {
			return nil, errors.New("invalid head ID format")
		}

		// Verify user exists
		var user models.User
		if err := s.db.Where("id = ?", parsed).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("head user not found")
			}
			return nil, fmt.Errorf("failed to verify head user: %w", err)
		}
		headID = &parsed
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	// Create department
	department := &models.Department{
		ID:          uuid.New(),
		Name:        strings.TrimSpace(req.Name),
		Code:        strings.ToUpper(strings.TrimSpace(req.Code)),
		Description: strings.TrimSpace(req.Description),
		HeadOfDept:  headID,
		Status:      status,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save
	if err := s.db.Create(department).Error; err != nil {
		return nil, fmt.Errorf("failed to create department: %w", err)
	}

	// Reload with relationships
	return s.getDepartmentWithDetails(department.ID.String())
}

// ============ 2. GET ALL DEPARTMENTS ============
func (s *DepartmentService) GetAllDepartments(params *dto.DepartmentQueryParams) (*dto.DepartmentListResponse, error) {
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
	query := s.db.Model(&models.Department{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where(
			"LOWER(name) LIKE ? OR LOWER(code) LIKE ? OR LOWER(description) LIKE ?",
			searchTerm, searchTerm, searchTerm,
		)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count departments: %w", err)
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
	var departments []models.Department
	if err := query.
		Preload("Head").
		Preload("Subjects", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL")
		}).
		Find(&departments).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch departments: %w", err)
	}

	// Convert to response
	responses := make([]dto.DepartmentWithSubjectsResponse, len(departments))
	for i, dept := range departments {
		responses[i] = *s.toDepartmentWithSubjectsResponse(&dept)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.DepartmentListResponse{
		Departments: responses,
		Total:       total,
		Page:        params.Page,
		Limit:       params.Limit,
		TotalPages:  totalPages,
	}, nil
}

// ============ 3. GET DEPARTMENT WITH SUBJECTS ============
func (s *DepartmentService) GetDepartmentWithSubjects(id string) (*dto.DepartmentWithSubjectsResponse, error) {
	return s.getDepartmentWithDetails(id)
}

// ============ 4. GET DEPARTMENT WITH HEAD AND SUBJECTS ============
func (s *DepartmentService) GetDepartmentWithHeadAndSubjects(id string) (*dto.DepartmentWithSubjectsResponse, error) {
	return s.getDepartmentWithDetails(id)
}

// ============ 5. GET DEPARTMENT BY ID WITH DETAILS ============
func (s *DepartmentService) GetDepartmentByID(id string) (*dto.DepartmentWithSubjectsResponse, error) {
	return s.getDepartmentWithDetails(id)
}

// ============ 6. UPDATE DEPARTMENT ============
func (s *DepartmentService) UpdateDepartment(id string, req *dto.UpdateDepartmentRequest, userID uuid.UUID) (*dto.DepartmentWithSubjectsResponse, error) {
	// Parse ID
	deptID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid department ID format")
	}

	// Find existing
	var department models.Department
	if err := s.db.Where("id = ? AND deleted_at IS NULL", deptID).First(&department).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("department not found")
		}
		return nil, fmt.Errorf("failed to fetch department: %w", err)
	}

	// Check duplicate name
	if req.Name != "" && req.Name != department.Name {
		var existing models.Department
		if err := s.db.Where("name = ? AND id != ? AND deleted_at IS NULL", req.Name, deptID).First(&existing).Error; err == nil {
			return nil, errors.New("department with this name already exists")
		}
		department.Name = strings.TrimSpace(req.Name)
	}

	// Check duplicate code
	if req.Code != "" && req.Code != department.Code {
		var existing models.Department
		if err := s.db.Where("code = ? AND id != ? AND deleted_at IS NULL", req.Code, deptID).First(&existing).Error; err == nil {
			return nil, errors.New("department with this code already exists")
		}
		department.Code = strings.ToUpper(strings.TrimSpace(req.Code))
	}

	// Update description
	if req.Description != "" {
		department.Description = strings.TrimSpace(req.Description)
	}

	// Update HeadID
	if req.HeadID != "" {
		parsed, err := uuid.Parse(req.HeadID)
		if err != nil {
			return nil, errors.New("invalid head ID format")
		}

		var user models.User
		if err := s.db.Where("id = ? AND deleted_at IS NULL", parsed).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("head user not found")
			}
			return nil, fmt.Errorf("failed to verify head user: %w", err)
		}
		department.HeadOfDept = &parsed
	}

	// Update status
	if req.Status != "" {
		department.Status = req.Status
	}

	// Update timestamp
	department.UpdatedAt = time.Now()

	// Save
	if err := s.db.Save(&department).Error; err != nil {
		return nil, fmt.Errorf("failed to update department: %w", err)
	}

	// Reload with details
	return s.getDepartmentWithDetails(id)
}

// ============ 7. DELETE DEPARTMENT ============
func (s *DepartmentService) DeleteDepartment(id string, userID uuid.UUID) error {
	// Parse ID
	deptID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid department ID format")
	}

	// Find existing
	var department models.Department
	if err := s.db.Where("id = ? AND deleted_at IS NULL", deptID).First(&department).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("department not found")
		}
		return fmt.Errorf("failed to fetch department: %w", err)
	}

	// Check if department has subjects
	var subjectCount int64
	if err := s.db.Model(&models.Subject{}).Where("department_id = ? AND deleted_at IS NULL", deptID).Count(&subjectCount).Error; err != nil {
		return fmt.Errorf("failed to check department usage: %w", err)
	}
	if subjectCount > 0 {
		return errors.New("cannot delete department: it has associated subjects")
	}

	// Soft delete
	if err := s.db.Delete(&department).Error; err != nil {
		return fmt.Errorf("failed to delete department: %w", err)
	}

	return nil
}

// ============ HELPER METHODS ============

func (s *DepartmentService) getDepartmentWithDetails(id string) (*dto.DepartmentWithSubjectsResponse, error) {
	deptID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid department ID format")
	}

	var department models.Department
	if err := s.db.
		Where("id = ? AND deleted_at IS NULL", deptID).
		Preload("Head").
		Preload("Subjects", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL")
		}).
		First(&department).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("department not found")
		}
		return nil, fmt.Errorf("failed to fetch department: %w", err)
	}

	return s.toDepartmentWithSubjectsResponse(&department), nil
}

func (s *DepartmentService) validateDepartmentRequest(req *dto.CreateDepartmentRequest) error {
	if req.Name == "" {
		return errors.New("department name is required")
	}
	if len(req.Name) < 2 {
		return errors.New("department name must be at least 2 characters")
	}
	if len(req.Name) > 100 {
		return errors.New("department name must be less than 100 characters")
	}
	if req.Code == "" {
		return errors.New("department code is required")
	}
	if len(req.Code) < 2 {
		return errors.New("department code must be at least 2 characters")
	}
	if len(req.Code) > 20 {
		return errors.New("department code must be less than 20 characters")
	}
	if req.Status != "" && req.Status != "active" && req.Status != "inactive" && req.Status != "archived" {
		return errors.New("status must be 'active', 'inactive', or 'archived'")
	}
	return nil
}

func (s *DepartmentService) toDepartmentWithSubjectsResponse(dept *models.Department) *dto.DepartmentWithSubjectsResponse {
	response := &dto.DepartmentWithSubjectsResponse{
		ID:          dept.ID.String(),
		Name:        dept.Name,
		Code:        dept.Code,
		Description: dept.Description,
		Status:      dept.Status,
		CreatedBy:   dept.CreatedBy.String(),
		CreatedAt:   dept.CreatedAt,
		UpdatedAt:   dept.UpdatedAt,
		SubjectCount: len(dept.Subjects),
	}

	if dept.HeadOfDept != nil {
		headID := dept.HeadOfDept.String()
		response.HeadID = &headID
	}

	if dept.Head.ID != uuid.Nil {
		response.Head = &dto.UserBrief{
			ID:        dept.Head.ID.String(),
			FirstName: dept.Head.FirstName,
			LastName:  dept.Head.LastName,
			Email:     dept.Head.Email,
		}
	}

	// Add subjects
	if len(dept.Subjects) > 0 {
		subjects := make([]dto.SubjectBriefResponse, len(dept.Subjects))
		for i, subject := range dept.Subjects {
			subjects[i] = dto.SubjectBriefResponse{
				ID:          subject.ID.String(),
				Name:        subject.Name,
				Code:        subject.Code,
				Description: subject.Description,
				Credits:     subject.Credits,
				Status:      subject.Status,
			}
		}
		response.Subjects = subjects
	}

	return response
}