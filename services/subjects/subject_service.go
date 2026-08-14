// services/subject_service.go
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

type SubjectService struct {
	db *gorm.DB
}

func NewSubjectService(db *gorm.DB) *SubjectService {
	return &SubjectService{db: db}
}

// ============ 1. CREATE SUBJECT ============
func (s *SubjectService) CreateSubject(req *dto.CreateSubjectRequest, userID uuid.UUID) (*dto.SubjectWithDepartmentResponse, error) {
	// Validate input
	if err := s.validateSubjectRequest(req); err != nil {
		return nil, err
	}

	// Parse Department ID
	departmentID, err := uuid.Parse(req.DepartmentID)
	if err != nil {
		return nil, errors.New("invalid department ID format")
	}

	// Verify department exists
	var department models.Department
	if err := s.db.Where("id = ? AND deleted_at IS NULL", departmentID).First(&department).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("department not found")
		}
		return nil, fmt.Errorf("failed to verify department: %w", err)
	}

	// Check duplicate name
	var existingSubject models.Subject
	if err := s.db.Where("name = ?", req.Name).First(&existingSubject).Error; err == nil {
		return nil, errors.New("subject with this name already exists")
	}

	// Check duplicate code
	if err := s.db.Where("code = ?", req.Code).First(&existingSubject).Error; err == nil {
		return nil, errors.New("subject with this code already exists")
	}

	// Set default values
	status := req.Status
	if status == "" {
		status = "active"
	}
	credits := req.Credits
	if credits == 0 {
		credits = 3
	}

	// Create subject
	subject := &models.Subject{
		ID:           uuid.New(),
		Name:         strings.TrimSpace(req.Name),
		Code:         strings.ToUpper(strings.TrimSpace(req.Code)),
		Description:  strings.TrimSpace(req.Description),
		DepartmentID: departmentID,
		Credits:      credits,
		Status:       status,
		CreatedBy:    userID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save
	if err := s.db.Create(subject).Error; err != nil {
		return nil, fmt.Errorf("failed to create subject: %w", err)
	}

	// Reload with relationships
	return s.getSubjectWithDetails(subject.ID.String())
}

// ============ 2. GET ALL SUBJECTS ============
func (s *SubjectService) GetAllSubjects(params *dto.SubjectQueryParams) (*dto.SubjectListResponse, error) {
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
	query := s.db.Model(&models.Subject{}).Where("subjects.deleted_at IS NULL")

	// Apply filters
	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where(
			"LOWER(subjects.name) LIKE ? OR LOWER(subjects.code) LIKE ? OR LOWER(subjects.description) LIKE ?",
			searchTerm, searchTerm, searchTerm,
		)
	}

	if params.Status != "" {
		query = query.Where("subjects.status = ?", params.Status)
	}

	// ✅ Handle department filtering - support both ID and Name
	if params.Department != "" {
		// Filter by department name
		query = query.Joins("LEFT JOIN departments ON departments.id = subjects.department_id").
			Where("LOWER(departments.name) = LOWER(?)", params.Department)
	} else if params.DepartmentID != "" {
		// Filter by department ID (UUID)
		query = query.Where("subjects.department_id = ?", params.DepartmentID)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count subjects: %w", err)
	}

	// Apply sorting
	sortDirection := "DESC"
	if strings.ToLower(params.SortOrder) == "asc" {
		sortDirection = "ASC"
	}
	query = query.Order("subjects." + params.SortBy + " " + sortDirection)

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	// Execute with preloads
	var subjects []models.Subject
	if err := query.
		Preload("Department").
		Preload("Department.Head").
		Find(&subjects).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch subjects: %w", err)
	}

	// Convert to response
	responses := make([]dto.SubjectWithDepartmentResponse, len(subjects))
	for i, subject := range subjects {
		responses[i] = *s.toSubjectWithDepartmentResponse(&subject)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.SubjectListResponse{
		Subjects:   responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// ============ 5. GET SUBJECT BY ID WITH DEPARTMENT AND HEAD ============
func (s *SubjectService) GetSubjectWithDepartmentAndHead(id string) (*dto.SubjectWithDepartmentResponse, error) {
	return s.getSubjectWithDetails(id)
}

// ============ 6. UPDATE SUBJECT ============
func (s *SubjectService) UpdateSubject(id string, req *dto.UpdateSubjectRequest, userID uuid.UUID) (*dto.SubjectWithDepartmentResponse, error) {
	// Parse ID
	subjectID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid subject ID format")
	}

	// Find existing
	var subject models.Subject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subject not found")
		}
		return nil, fmt.Errorf("failed to fetch subject: %w", err)
	}

	// Check duplicate name
	if req.Name != "" && req.Name != subject.Name {
		var existing models.Subject
		if err := s.db.Where("name = ? AND id != ? AND deleted_at IS NULL", req.Name, subjectID).First(&existing).Error; err == nil {
			return nil, errors.New("subject with this name already exists")
		}
		subject.Name = strings.TrimSpace(req.Name)
	}

	// Check duplicate code
	if req.Code != "" && req.Code != subject.Code {
		var existing models.Subject
		if err := s.db.Where("code = ? AND id != ? AND deleted_at IS NULL", req.Code, subjectID).First(&existing).Error; err == nil {
			return nil, errors.New("subject with this code already exists")
		}
		subject.Code = strings.ToUpper(strings.TrimSpace(req.Code))
	}

	// Update description
	if req.Description != "" {
		subject.Description = strings.TrimSpace(req.Description)
	}

	// Update DepartmentID
	if req.DepartmentID != "" {
		departmentID, err := uuid.Parse(req.DepartmentID)
		if err != nil {
			return nil, errors.New("invalid department ID format")
		}

		var department models.Department
		if err := s.db.Where("id = ? AND deleted_at IS NULL", departmentID).First(&department).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("department not found")
			}
			return nil, fmt.Errorf("failed to verify department: %w", err)
		}
		subject.DepartmentID = departmentID
	}

	// Update credits
	if req.Credits > 0 {
		subject.Credits = req.Credits
	}

	// Update status
	if req.Status != "" {
		subject.Status = req.Status
	}

	// Update timestamp
	subject.UpdatedAt = time.Now()

	// Save
	if err := s.db.Save(&subject).Error; err != nil {
		return nil, fmt.Errorf("failed to update subject: %w", err)
	}

	// Reload with details
	return s.getSubjectWithDetails(id)
}

// ============ 7. DELETE SUBJECT ============
func (s *SubjectService) DeleteSubject(id string, userID uuid.UUID) error {
	// Parse ID
	subjectID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid subject ID format")
	}

	// Find existing
	var subject models.Subject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("subject not found")
		}
		return fmt.Errorf("failed to fetch subject: %w", err)
	}

	// Soft delete
	if err := s.db.Delete(&subject).Error; err != nil {
		return fmt.Errorf("failed to delete subject: %w", err)
	}

	return nil
}

// ============ HELPER METHODS ============

func (s *SubjectService) getSubjectWithDetails(id string) (*dto.SubjectWithDepartmentResponse, error) {
	subjectID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid subject ID format")
	}

	var subject models.Subject
	if err := s.db.
		Where("id = ? AND deleted_at IS NULL", subjectID).
		Preload("Department").
		Preload("Department.Head").
		First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subject not found")
		}
		return nil, fmt.Errorf("failed to fetch subject: %w", err)
	}

	return s.toSubjectWithDepartmentResponse(&subject), nil
}

func (s *SubjectService) validateSubjectRequest(req *dto.CreateSubjectRequest) error {
	if req.Name == "" {
		return errors.New("subject name is required")
	}
	if len(req.Name) < 3 {
		return errors.New("subject name must be at least 3 characters")
	}
	if len(req.Name) > 255 {
		return errors.New("subject name must be less than 255 characters")
	}
	if req.Code == "" {
		return errors.New("subject code is required")
	}
	if len(req.Code) < 2 {
		return errors.New("subject code must be at least 2 characters")
	}
	if len(req.Code) > 50 {
		return errors.New("subject code must be less than 50 characters")
	}
	if req.DepartmentID == "" {
		return errors.New("department ID is required")
	}
	if req.Credits < 0 {
		return errors.New("credits cannot be negative")
	}
	if req.Credits > 10 {
		return errors.New("credits cannot exceed 10")
	}
	if req.Status != "" && req.Status != "active" && req.Status != "inactive" && req.Status != "archived" {
		return errors.New("status must be 'active', 'inactive', or 'archived'")
	}
	return nil
}

func (s *SubjectService) toSubjectWithDepartmentResponse(subject *models.Subject) *dto.SubjectWithDepartmentResponse {
	response := &dto.SubjectWithDepartmentResponse{
		ID:           subject.ID.String(),
		Name:         subject.Name,
		Code:         subject.Code,
		Description:  subject.Description,
		DepartmentID: subject.DepartmentID.String(),
		Credits:      subject.Credits,
		Status:       subject.Status,
		CreatedBy:    subject.CreatedBy.String(),
		CreatedAt:    subject.CreatedAt,
		UpdatedAt:    subject.UpdatedAt,
	}

	// Add department details if loaded
	if subject.Department.ID != uuid.Nil {
		response.Department = &dto.DepartmentWithHeadBriefResponse{
			ID:          subject.Department.ID.String(),
			Name:        subject.Department.Name,
			Code:        subject.Department.Code,
			Description: subject.Department.Description,
		}

		if subject.Department.HeadOfDept != nil {
			headID := subject.Department.HeadOfDept.String()
			response.Department.HeadID = &headID
		}

		if subject.Department.Head.ID != uuid.Nil {
			response.Department.Head = &dto.UserBrief{
				ID:        subject.Department.Head.ID.String(),
				FirstName: subject.Department.Head.FirstName,
				LastName:  subject.Department.Head.LastName,
				Email:     subject.Department.Head.Email,
			}
		}
	}

	return response
}