// services/subject_service.go
package services

import (
	"errors"
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

	// Check if subject with same name or code already exists
	existingSubject := &models.Subject{}
	if err := s.db.Where("name = ? OR code = ?", req.Name, req.Code).First(existingSubject).Error; err == nil {
		if existingSubject.Name == req.Name {
			return nil, errors.New("subject with this name already exists")
		}
		if existingSubject.Code == req.Code {
			return nil, errors.New("subject with this code already exists")
		}
	}

	// Create new subject
	subject := &models.Subject{
		ID:          uuid.New(),
		Name:        strings.TrimSpace(req.Name),
		Code:        strings.ToUpper(strings.TrimSpace(req.Code)),
		Description: strings.TrimSpace(req.Description),
		Department:  strings.TrimSpace(req.Department),
		Credits:     req.Credits,
		Status:      req.Status,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Set default status if not provided
	if subject.Status == "" {
		subject.Status = "active"
	}

	// Save to database
	if err := s.db.Create(subject).Error; err != nil {
		return nil, errors.New("failed to create subject: " + err.Error())
	}

	// Convert to response DTO
	return s.toSubjectResponse(subject), nil
}

// validateSubjectRequest validates the subject request
func (s *SubjectService) validateSubjectRequest(req *dto.CreateSubjectRequest) error {
	if req.Name == "" {
		return errors.New("subject name is required")
	}
	if len(req.Name) < 3 {
		return errors.New("subject name must be at least 3 characters")
	}
	if req.Code == "" {
		return errors.New("subject code is required")
	}
	if len(req.Code) < 2 {
		return errors.New("subject code must be at least 2 characters")
	}
	if req.Credits < 0 {
		return errors.New("credits cannot be negative")
	}
	if req.Credits > 10 {
		return errors.New("credits cannot exceed 10")
	}
	if req.Status != "" && req.Status != "active" && req.Status != "inactive" {
		return errors.New("status must be either 'active' or 'inactive'")
	}
	return nil
}

// toSubjectResponse converts model to response DTO
func (s *SubjectService) toSubjectResponse(subject *models.Subject) *dto.SubjectResponse {
	return &dto.SubjectResponse{
		ID:          subject.ID.String(),
		Name:        subject.Name,
		Code:        subject.Code,
		Description: subject.Description,
		Department:  subject.Department,
		Credits:     subject.Credits,
		Status:      subject.Status,
		CreatedBy:   subject.CreatedBy.String(),
		CreatedAt:   subject.CreatedAt,
		UpdatedAt:   subject.UpdatedAt,
	}
}


// GetAllSubjects retrieves all subjects with pagination and filters
func (s *SubjectService) GetAllSubjects(params *dto.SubjectQueryParams) (*dto.SubjectListResponse, error) {
	// Set default values
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}

	// Build query
	query := s.db.Model(&models.Subject{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where(
			"LOWER(name) LIKE ? OR LOWER(code) LIKE ? OR LOWER(description) LIKE ?",
			searchTerm, searchTerm, searchTerm,
		)
	}

	if params.Department != "" {
		query = query.Where("LOWER(department) = LOWER(?)", params.Department)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("failed to count subjects: " + err.Error())
	}

	// Apply sorting
	sortColumn := params.SortBy
	if sortColumn == "" {
		sortColumn = "created_at"
	}
	sortDirection := "DESC"
	if strings.ToLower(params.SortOrder) == "asc" {
		sortDirection = "ASC"
	}
	query = query.Order(sortColumn + " " + sortDirection)

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	// Execute query
	var subjects []models.Subject
	if err := query.Find(&subjects).Error; err != nil {
		return nil, errors.New("failed to fetch subjects: " + err.Error())
	}

	// Convert to response DTO
	var subjectResponses []dto.SubjectResponse
	for _, subject := range subjects {
		subjectResponses = append(subjectResponses, *s.toSubjectResponse(&subject))
	}

	// Calculate total pages
	totalPages := int(total) / params.Limit
	if int(total)%params.Limit != 0 {
		totalPages++
	}

	return &dto.SubjectListResponse{
		Subjects:   subjectResponses,
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
	if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subject not found")
		}
		return nil, errors.New("failed to fetch subject: " + err.Error())
	}

	return s.toSubjectResponse(&subject), nil
}

// GetSubjectDepartments retrieves all unique departments
func (s *SubjectService) GetSubjectDepartments() ([]string, error) {
	var departments []string
	if err := s.db.Model(&models.Subject{}).
		Where("deleted_at IS NULL").
		Distinct("department").
		Where("department IS NOT NULL AND department != ''").
		Pluck("department", &departments).Error; err != nil {
		return nil, errors.New("failed to fetch departments: " + err.Error())
	}
	return departments, nil
}


// UpdateSubject updates an existing subject
func (s *SubjectService) UpdateSubject(id string, req *dto.UpdateSubjectRequest, userID uuid.UUID) (*dto.SubjectResponse, error) {
	// Parse subject ID
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

	// Check for duplicate name/code if being updated
	if req.Name != "" && req.Name != subject.Name {
		var existing models.Subject
		if err := s.db.Where("name = ? AND id != ? AND deleted_at IS NULL", req.Name, subjectID).First(&existing).Error; err == nil {
			return nil, errors.New("subject with this name already exists")
		}
	}

	if req.Code != "" && req.Code != subject.Code {
		var existing models.Subject
		if err := s.db.Where("code = ? AND id != ? AND deleted_at IS NULL", req.Code, subjectID).First(&existing).Error; err == nil {
			return nil, errors.New("subject with this code already exists")
		}
	}

	// Update fields
	if req.Name != "" {
		subject.Name = strings.TrimSpace(req.Name)
	}
	if req.Code != "" {
		subject.Code = strings.ToUpper(strings.TrimSpace(req.Code))
	}
	if req.Description != "" {
		subject.Description = strings.TrimSpace(req.Description)
	}
	if req.Department != "" {
		subject.Department = strings.TrimSpace(req.Department)
	}
	if req.Credits > 0 {
		subject.Credits = req.Credits
	}
	if req.Status != "" {
		subject.Status = req.Status
	}

	// Update timestamp
	subject.UpdatedAt = time.Now()

	// Save to database
	if err := s.db.Save(&subject).Error; err != nil {
		return nil, errors.New("failed to update subject: " + err.Error())
	}

	// Convert to response DTO
	return s.toSubjectResponse(&subject), nil
}

// DeleteSubject soft deletes a subject
func (s *SubjectService) DeleteSubject(id string, userID uuid.UUID) error {
	// Parse subject ID
	subjectID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid subject ID")
	}

	// Find existing subject
	var subject models.Subject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("subject not found")
		}
		return errors.New("failed to fetch subject: " + err.Error())
	}

	// Check if subject is in use (has courses or assignments)
	// This is optional - implement if you have relationships
	// Example: check if subject has courses
	// var courseCount int64
	// if err := s.db.Model(&models.Course{}).Where("subject_id = ?", subjectID).Count(&courseCount).Error; err != nil {
	// 	return errors.New("failed to check subject usage: " + err.Error())
	// }
	// if courseCount > 0 {
	// 	return errors.New("cannot delete subject: it is being used by courses")
	// }

	// Perform soft delete
	if err := s.db.Delete(&subject).Error; err != nil {
		return errors.New("failed to delete subject: " + err.Error())
	}

	return nil
}


