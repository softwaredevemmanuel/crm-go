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

// CreateSubject creates a new subject
func (s *SubjectService) CreateSubject(req *dto.CreateSubjectRequest, userID uuid.UUID) (*dto.SubjectResponse, error) {
	// Validate input
	if err := s.validateSubjectRequest(req); err != nil {
		return nil, err
	}

	// Check if name already exists
	var existing models.Subject
	if err := s.db.Where("name = ? AND deleted_at IS NULL", req.Name).First(&existing).Error; err == nil {
		return nil, errors.New("subject name already exists")
	}

	// Check if code already exists
	if err := s.db.Where("code = ? AND deleted_at IS NULL", req.Code).First(&existing).Error; err == nil {
		return nil, errors.New("subject code already exists")
	}

	// Parse department ID
	departmentID, err := uuid.Parse(req.DepartmentID)
	if err != nil {
		return nil, errors.New("invalid department ID format")
	}

	// Check if department exists
	var department models.Department
	if err := s.db.Where("id = ? AND deleted_at IS NULL", departmentID).First(&department).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("department not found")
		}
		return nil, errors.New("failed to verify department: " + err.Error())
	}

	// Set default values
	credits := req.Credits
	if credits == 0 {
		credits = 3
	}

	status := req.Status
	if status == "" {
		status = "active"
	}

	// Create subject
	subject := &models.Subject{
		ID:           uuid.New(),
		Name:         req.Name,
		Code:         req.Code,
		Description:  req.Description,
		DepartmentID: departmentID,
		Credits:      credits,
		Status:       status,
		CreatedBy:    userID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.db.Create(subject).Error; err != nil {
		return nil, errors.New("failed to create subject: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("Department").Preload("Creator").First(subject, subject.ID).Error; err != nil {
		return nil, errors.New("failed to load subject details: " + err.Error())
	}

	return s.toSubjectResponse(subject), nil
}

// BulkCreateSubjects creates multiple subjects
func (s *SubjectService) BulkCreateSubjects(req *dto.BulkCreateSubjectsRequest, userID uuid.UUID) (*dto.BulkSubjectResult, error) {
	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	result := &dto.BulkSubjectResult{
		Created: []dto.SubjectResponse{},
		Errors:  []dto.BulkSubjectError{},
	}

	for _, subjectReq := range req.Subjects {
		// Validate
		if subjectReq.Name == "" {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkSubjectError{
				Name:  subjectReq.Name,
				Code:  subjectReq.Code,
				Error: "name is required",
			})
			continue
		}
		if subjectReq.Code == "" {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkSubjectError{
				Name:  subjectReq.Name,
				Code:  subjectReq.Code,
				Error: "code is required",
			})
			continue
		}

		// Check if name already exists
		var existing models.Subject
		if err := s.db.Where("name = ? AND deleted_at IS NULL", subjectReq.Name).First(&existing).Error; err == nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkSubjectError{
				Name:  subjectReq.Name,
				Code:  subjectReq.Code,
				Error: "subject name already exists",
			})
			continue
		}

		// Check if code already exists
		if err := s.db.Where("code = ? AND deleted_at IS NULL", subjectReq.Code).First(&existing).Error; err == nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkSubjectError{
				Name:  subjectReq.Name,
				Code:  subjectReq.Code,
				Error: "subject code already exists",
			})
			continue
		}

		// Parse department ID
		departmentID, err := uuid.Parse(subjectReq.DepartmentID)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkSubjectError{
				Name:  subjectReq.Name,
				Code:  subjectReq.Code,
				Error: "invalid department ID format",
			})
			continue
		}

		// Check if department exists
		var department models.Department
		if err := s.db.Where("id = ? AND deleted_at IS NULL", departmentID).First(&department).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkSubjectError{
				Name:  subjectReq.Name,
				Code:  subjectReq.Code,
				Error: "department not found",
			})
			continue
		}

		// Set default credits
		credits := subjectReq.Credits
		if credits == 0 {
			credits = 3
		}

		// Create subject
		subject := &models.Subject{
			ID:           uuid.New(),
			Name:         subjectReq.Name,
			Code:         subjectReq.Code,
			Description:  subjectReq.Description,
			DepartmentID: departmentID,
			Credits:      credits,
			Status:       status,
			CreatedBy:    userID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := s.db.Create(subject).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkSubjectError{
				Name:  subjectReq.Name,
				Code:  subjectReq.Code,
				Error: "failed to create subject: " + err.Error(),
			})
			continue
		}

		// Preload relationships
		if err := s.db.Preload("Department").Preload("Creator").First(subject, subject.ID).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkSubjectError{
				Name:  subjectReq.Name,
				Code:  subjectReq.Code,
				Error: "failed to load subject details",
			})
			continue
		}

		result.SuccessCount++
		result.Created = append(result.Created, *s.toSubjectResponse(subject))
	}

	return result, nil
}

// GetAllSubjects retrieves all subjects with pagination and filters
func (s *SubjectService) GetAllSubjects(params *dto.SubjectQueryParams) (*dto.SubjectListResponse, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = "name"
	}
	if params.SortOrder == "" {
		params.SortOrder = "asc"
	}

	// Build query
	query := s.db.Model(&models.Subject{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.DepartmentID != "" {
		departmentID, err := uuid.Parse(params.DepartmentID)
		if err == nil {
			query = query.Where("department_id = ?", departmentID)
		}
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ? OR LOWER(description) LIKE ?",
			search, search, search)
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
	query = query.Order(params.SortBy + " " + sortDirection)

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	// Execute with preloads
	var subjects []models.Subject
	if err := query.Preload("Department").Preload("Creator").Find(&subjects).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch subjects: %w", err)
	}

	// Convert to response
	responses := make([]dto.SubjectResponse, len(subjects))
	for i, subject := range subjects {
		responses[i] = *s.toSubjectResponse(&subject)
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

// GetSubjectByID retrieves a single subject by ID
func (s *SubjectService) GetSubjectByID(id string) (*dto.SubjectResponse, error) {
	subjectID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid subject ID")
	}

	var subject models.Subject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).
		Preload("Department").
		Preload("Creator").
		First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subject not found")
		}
		return nil, errors.New("failed to fetch subject: " + err.Error())
	}

	return s.toSubjectResponse(&subject), nil
}

// GetSubjectsByDepartment retrieves all subjects for a specific department
func (s *SubjectService) GetSubjectsByDepartment(departmentID string) ([]dto.SubjectResponse, error) {
	dID, err := uuid.Parse(departmentID)
	if err != nil {
		return nil, errors.New("invalid department ID")
	}

	var subjects []models.Subject
	if err := s.db.Where("department_id = ? AND deleted_at IS NULL", dID).
		Preload("Department").
		Preload("Creator").
		Order("name ASC").
		Find(&subjects).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch subjects: %w", err)
	}

	responses := make([]dto.SubjectResponse, len(subjects))
	for i, subject := range subjects {
		responses[i] = *s.toSubjectResponse(&subject)
	}

	return responses, nil
}

// GetActiveSubjects retrieves all active subjects
func (s *SubjectService) GetActiveSubjects() ([]dto.SubjectResponse, error) {
	var subjects []models.Subject
	if err := s.db.Where("status = ? AND deleted_at IS NULL", "active").
		Preload("Department").
		Preload("Creator").
		Order("name ASC").
		Find(&subjects).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch active subjects: %w", err)
	}

	responses := make([]dto.SubjectResponse, len(subjects))
	for i, subject := range subjects {
		responses[i] = *s.toSubjectResponse(&subject)
	}

	return responses, nil
}

// GetSubjectStats retrieves statistics for subjects
func (s *SubjectService) GetSubjectStats() (*dto.SubjectStats, error) {
	var stats dto.SubjectStats

	// Count total subjects
	if err := s.db.Model(&models.Subject{}).Where("deleted_at IS NULL").Count(&stats.TotalSubjects).Error; err != nil {
		return nil, fmt.Errorf("failed to count total subjects: %w", err)
	}

	// Count active subjects
	if err := s.db.Model(&models.Subject{}).Where("status = ? AND deleted_at IS NULL", "active").Count(&stats.ActiveSubjects).Error; err != nil {
		return nil, fmt.Errorf("failed to count active subjects: %w", err)
	}

	// Count inactive subjects
	if err := s.db.Model(&models.Subject{}).Where("status = ? AND deleted_at IS NULL", "inactive").Count(&stats.InactiveSubjects).Error; err != nil {
		return nil, fmt.Errorf("failed to count inactive subjects: %w", err)
	}

	// Count total departments with subjects
	if err := s.db.Model(&models.Subject{}).
		Where("deleted_at IS NULL").
		Distinct("department_id").
		Count(&stats.TotalDepartments).Error; err != nil {
		return nil, fmt.Errorf("failed to count departments: %w", err)
	}

	return &stats, nil
}

// UpdateSubject updates an existing subject
func (s *SubjectService) UpdateSubject(id string, req *dto.UpdateSubjectRequest) (*dto.SubjectResponse, error) {
	subjectID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid subject ID")
	}

	// Find existing subject
	var subject models.Subject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subject not found")
		}
		return nil, errors.New("failed to fetch subject: " + err.Error())
	}

	// Update fields
	if req.Name != "" {
		// Check if name already exists for another subject
		var existing models.Subject
		if err := s.db.Where("name = ? AND id != ? AND deleted_at IS NULL", req.Name, subjectID).First(&existing).Error; err == nil {
			return nil, errors.New("subject name already exists")
		}
		subject.Name = req.Name
	}

	if req.Code != "" {
		// Check if code already exists for another subject
		var existing models.Subject
		if err := s.db.Where("code = ? AND id != ? AND deleted_at IS NULL", req.Code, subjectID).First(&existing).Error; err == nil {
			return nil, errors.New("subject code already exists")
		}
		subject.Code = req.Code
	}

	if req.Description != "" {
		subject.Description = req.Description
	}

	if req.DepartmentID != "" {
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
			return nil, errors.New("failed to verify department: " + err.Error())
		}
		subject.DepartmentID = departmentID
	}

	if req.Credits > 0 {
		subject.Credits = req.Credits
	}

	if req.Status != "" {
		if req.Status != "active" && req.Status != "inactive" {
			return nil, errors.New("status must be 'active' or 'inactive'")
		}
		subject.Status = req.Status
	}

	subject.UpdatedAt = time.Now()

	if err := s.db.Save(&subject).Error; err != nil {
		return nil, errors.New("failed to update subject: " + err.Error())
	}

	// Preload relationships
	if err := s.db.Preload("Department").Preload("Creator").First(&subject, subject.ID).Error; err != nil {
		return nil, errors.New("failed to load subject details: " + err.Error())
	}

	return s.toSubjectResponse(&subject), nil
}

// DeleteSubject soft deletes a subject
func (s *SubjectService) DeleteSubject(id string) error {
	subjectID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid subject ID")
	}

	var subject models.Subject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("subject not found")
		}
		return errors.New("failed to fetch subject: " + err.Error())
	}

	if err := s.db.Delete(&subject).Error; err != nil {
		return errors.New("failed to delete subject: " + err.Error())
	}

	return nil
}

// validateSubjectRequest validates the subject request
func (s *SubjectService) validateSubjectRequest(req *dto.CreateSubjectRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.Code == "" {
		return errors.New("code is required")
	}
	if req.DepartmentID == "" {
		return errors.New("department ID is required")
	}
	if req.Status != "" && req.Status != "active" && req.Status != "inactive" {
		return errors.New("status must be 'active' or 'inactive'")
	}
	if req.Credits < 0 {
		return errors.New("credits cannot be negative")
	}
	return nil
}

// toSubjectResponse converts model to response DTO
func (s *SubjectService) toSubjectResponse(subject *models.Subject) *dto.SubjectResponse {
	response := &dto.SubjectResponse{
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

	// Add department details if preloaded
	if subject.Department.ID != uuid.Nil {
		response.Department = &dto.DepartmentResponse{
			ID:          subject.Department.ID.String(),
			Name:        subject.Department.Name,
			Code:        subject.Department.Code,
			Description: subject.Department.Description,
		}
	}

	// Add creator details if preloaded
	if subject.Creator.ID != uuid.Nil {
		response.Creator = &dto.UserResponse{
			ID:        subject.Creator.ID.String(),
			FirstName: subject.Creator.FirstName,
			LastName:  subject.Creator.LastName,
			Email:     subject.Creator.Email,
			Phone:     subject.Creator.Phone,
			Role:      subject.Creator.Role,
		}
	}

	return response
}