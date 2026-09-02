// services/scheme_of_work_service.go
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

type SchemeOfWorkService struct {
	db *gorm.DB
}

func NewSchemeOfWorkService(db *gorm.DB) *SchemeOfWorkService {
	return &SchemeOfWorkService{db: db}
}

// CreateSchemeOfWork creates a new scheme of work
func (s *SchemeOfWorkService) CreateSchemeOfWork(req *dto.CreateSchemeOfWorkRequest, userID uuid.UUID) (*dto.SchemeOfWorkResponse, error) {
	// Validate input
	if err := s.validateSchemeRequest(req); err != nil {
		return nil, err
	}

	// Parse UUIDs
	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		return nil, errors.New("invalid subject ID format")
	}

	gradeID, err := uuid.Parse(req.GradeID)
	if err != nil {
		return nil, errors.New("invalid grade ID format")
	}

	// Check if subject exists
	var subject models.Subject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subject not found")
		}
		return nil, errors.New("failed to verify subject: " + err.Error())
	}

	// Check if grade exists
	var grade models.ClassGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", gradeID).First(&grade).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("grade not found")
		}
		return nil, errors.New("failed to verify grade: " + err.Error())
	}

	// Check if scheme already exists for this combination
	var existing models.SchemeOfWork
	if err := s.db.Where("subject_id = ? AND grade_id = ? AND term = ? AND deleted_at IS NULL",
		subjectID, gradeID, req.Term).First(&existing).Error; err == nil {
		return nil, errors.New("scheme of work already exists for this subject, grade, and term")
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "draft"
	}

	// Create scheme of work
	scheme := &models.SchemeOfWork{
		ID:          uuid.New(),
		SubjectID:   subjectID,
		GradeID:     gradeID,
		Term:        req.Term,
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(scheme).Error; err != nil {
		return nil, errors.New("failed to create scheme of work: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("Subject").Preload("Grade").Preload("Creator").First(scheme, scheme.ID).Error; err != nil {
		return nil, errors.New("failed to load scheme details: " + err.Error())
	}

	return s.toSchemeResponse(scheme), nil
}

// BulkCreateSchemes creates multiple schemes of work
func (s *SchemeOfWorkService) BulkCreateSchemes(req *dto.BulkCreateSchemesRequest, userID uuid.UUID) (*dto.BulkSchemeResult, error) {
	// Parse UUIDs
	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		return nil, errors.New("invalid subject ID format")
	}

	gradeID, err := uuid.Parse(req.GradeID)
	if err != nil {
		return nil, errors.New("invalid grade ID format")
	}

	// Verify subject exists
	var subject models.Subject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subject not found")
		}
		return nil, errors.New("failed to verify subject: " + err.Error())
	}

	// Verify grade exists
	var grade models.ClassGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", gradeID).First(&grade).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("grade not found")
		}
		return nil, errors.New("failed to verify grade: " + err.Error())
	}

	// Check if scheme already exists
	var existing models.SchemeOfWork
	if err := s.db.Where("subject_id = ? AND grade_id = ? AND term = ? AND deleted_at IS NULL",
		subjectID, gradeID, req.Term).First(&existing).Error; err == nil {
		return nil, errors.New("scheme of work already exists for this subject, grade, and term")
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "draft"
	}

	result := &dto.BulkSchemeResult{
		Created: []dto.SchemeOfWorkResponse{},
		Errors:  []dto.BulkSchemeError{},
	}

	for _, schemeReq := range req.Schemes {
		if schemeReq.Title == "" {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkSchemeError{
				Title: schemeReq.Title,
				Error: "title is required",
			})
			continue
		}

		// Create scheme
		scheme := &models.SchemeOfWork{
			ID:          uuid.New(),
			SubjectID:   subjectID,
			GradeID:     gradeID,
			Term:        req.Term,
			Title:       schemeReq.Title,
			Description: schemeReq.Description,
			Status:      status,
			CreatedBy:   userID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := s.db.Create(scheme).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkSchemeError{
				Title: schemeReq.Title,
				Error: "failed to create scheme: " + err.Error(),
			})
			continue
		}

		// Preload relationships
		if err := s.db.Preload("Subject").Preload("Grade").Preload("Creator").First(scheme, scheme.ID).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkSchemeError{
				Title: schemeReq.Title,
				Error: "failed to load scheme details",
			})
			continue
		}

		result.SuccessCount++
		result.Created = append(result.Created, *s.toSchemeResponse(scheme))
	}

	return result, nil
}

// GetAllSchemes retrieves all schemes of work with pagination and filters
func (s *SchemeOfWorkService) GetAllSchemes(params *dto.SchemeOfWorkQueryParams) (*dto.SchemeOfWorkListResponse, error) {
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
	query := s.db.Model(&models.SchemeOfWork{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.SubjectID != "" {
		subjectID, err := uuid.Parse(params.SubjectID)
		if err == nil {
			query = query.Where("subject_id = ?", subjectID)
		}
	}

	if params.GradeID != "" {
		gradeID, err := uuid.Parse(params.GradeID)
		if err == nil {
			query = query.Where("grade_id = ?", gradeID)
		}
	}

	if params.Term != "" {
		query = query.Where("term = ?", params.Term)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", search, search)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count schemes: %w", err)
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
	var schemes []models.SchemeOfWork
	if err := query.Preload("Subject").Preload("Grade").Preload("Creator").Find(&schemes).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch schemes: %w", err)
	}

	// Convert to response
	responses := make([]dto.SchemeOfWorkResponse, len(schemes))
	for i, scheme := range schemes {
		responses[i] = *s.toSchemeResponse(&scheme)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.SchemeOfWorkListResponse{
		Schemes:    responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetSchemeByID retrieves a single scheme by ID
func (s *SchemeOfWorkService) GetSchemeByID(id string) (*dto.SchemeOfWorkResponse, error) {
	schemeID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid scheme ID")
	}

	var scheme models.SchemeOfWork
	if err := s.db.Where("id = ? AND deleted_at IS NULL", schemeID).
		Preload("Subject").
		Preload("Grade").
		Preload("Creator").
		Preload("Lessons").
		First(&scheme).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("scheme of work not found")
		}
		return nil, errors.New("failed to fetch scheme: " + err.Error())
	}

	return s.toSchemeResponse(&scheme), nil
}

// GetSchemesBySubject retrieves all schemes for a specific subject
func (s *SchemeOfWorkService) GetSchemesBySubject(subjectID string) ([]dto.SchemeOfWorkResponse, error) {
	sID, err := uuid.Parse(subjectID)
	if err != nil {
		return nil, errors.New("invalid subject ID")
	}

	var schemes []models.SchemeOfWork
	if err := s.db.Where("subject_id = ? AND deleted_at IS NULL", sID).
		Preload("Subject").
		Preload("Grade").
		Preload("Creator").
		Order("created_at DESC").
		Find(&schemes).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch schemes: %w", err)
	}

	responses := make([]dto.SchemeOfWorkResponse, len(schemes))
	for i, scheme := range schemes {
		responses[i] = *s.toSchemeResponse(&scheme)
	}

	return responses, nil
}

// GetSchemesByGrade retrieves all schemes for a specific grade
func (s *SchemeOfWorkService) GetSchemesByGrade(gradeID string) ([]dto.SchemeOfWorkResponse, error) {
	gID, err := uuid.Parse(gradeID)
	if err != nil {
		return nil, errors.New("invalid grade ID")
	}

	var schemes []models.SchemeOfWork
	if err := s.db.Where("grade_id = ? AND deleted_at IS NULL", gID).
		Preload("Subject").
		Preload("Grade").
		Preload("Creator").
		Order("created_at DESC").
		Find(&schemes).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch schemes: %w", err)
	}

	responses := make([]dto.SchemeOfWorkResponse, len(schemes))
	for i, scheme := range schemes {
		responses[i] = *s.toSchemeResponse(&scheme)
	}

	return responses, nil
}

// GetSchemesByGradeAndTerm retrieves all schemes for a specific grade and term
func (s *SchemeOfWorkService) GetSchemesByGradeAndTerm(gradeID, term string) ([]dto.SchemeOfWorkResponse, error) {
	gID, err := uuid.Parse(gradeID)
	if err != nil {
		return nil, errors.New("invalid grade ID")
	}

	var schemes []models.SchemeOfWork
	if err := s.db.Where("grade_id = ? AND term = ? AND deleted_at IS NULL", gID, term).
		Preload("Subject").
		Preload("Grade").
		Preload("Creator").
		Order("created_at DESC").
		Find(&schemes).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch schemes: %w", err)
	}

	responses := make([]dto.SchemeOfWorkResponse, len(schemes))
	for i, scheme := range schemes {
		responses[i] = *s.toSchemeResponse(&scheme)
	}

	return responses, nil
}

// GetSchemesBySubjectAndGrade retrieves all schemes for a specific subject and grade
func (s *SchemeOfWorkService) GetSchemesBySubjectAndGrade(subjectID, gradeID string) ([]dto.SchemeOfWorkResponse, error) {
	sID, err := uuid.Parse(subjectID)
	if err != nil {
		return nil, errors.New("invalid subject ID")
	}

	gID, err := uuid.Parse(gradeID)
	if err != nil {
		return nil, errors.New("invalid grade ID")
	}

	var schemes []models.SchemeOfWork
	if err := s.db.Where("subject_id = ? AND grade_id = ? AND deleted_at IS NULL", sID, gID).
		Preload("Subject").
		Preload("Grade").
		Preload("Creator").
		Order("term, created_at DESC").
		Find(&schemes).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch schemes: %w", err)
	}

	responses := make([]dto.SchemeOfWorkResponse, len(schemes))
	for i, scheme := range schemes {
		responses[i] = *s.toSchemeResponse(&scheme)
	}

	return responses, nil
}

// UpdateSchemeOfWork updates an existing scheme
func (s *SchemeOfWorkService) UpdateSchemeOfWork(id string, req *dto.UpdateSchemeOfWorkRequest) (*dto.SchemeOfWorkResponse, error) {
	schemeID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid scheme ID")
	}

	// Find existing scheme
	var scheme models.SchemeOfWork
	if err := s.db.Where("id = ? AND deleted_at IS NULL", schemeID).First(&scheme).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("scheme of work not found")
		}
		return nil, errors.New("failed to fetch scheme: " + err.Error())
	}

	// Update fields
	if req.SubjectID != "" {
		subjectID, err := uuid.Parse(req.SubjectID)
		if err != nil {
			return nil, errors.New("invalid subject ID format")
		}
		var subject models.Subject
		if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("subject not found")
			}
			return nil, errors.New("failed to verify subject: " + err.Error())
		}
		scheme.SubjectID = subjectID
	}

	if req.GradeID != "" {
		gradeID, err := uuid.Parse(req.GradeID)
		if err != nil {
			return nil, errors.New("invalid grade ID format")
		}
		var grade models.ClassGrade
		if err := s.db.Where("id = ? AND deleted_at IS NULL", gradeID).First(&grade).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("grade not found")
			}
			return nil, errors.New("failed to verify grade: " + err.Error())
		}
		scheme.GradeID = gradeID
	}

	if req.Term != "" {
		if req.Term != "first" && req.Term != "second" && req.Term != "third" {
			return nil, errors.New("term must be 'first', 'second', or 'third'")
		}
		scheme.Term = req.Term
	}

	if req.Title != "" {
		scheme.Title = req.Title
	}

	if req.Description != "" {
		scheme.Description = req.Description
	}

	if req.Status != "" {
		if req.Status != "draft" && req.Status != "published" && req.Status != "archived" {
			return nil, errors.New("status must be 'draft', 'published', or 'archived'")
		}
		scheme.Status = req.Status
	}

	scheme.UpdatedAt = time.Now()

	if err := s.db.Save(&scheme).Error; err != nil {
		return nil, errors.New("failed to update scheme: " + err.Error())
	}

	// Preload relationships
	if err := s.db.Preload("Subject").Preload("Grade").Preload("Creator").First(&scheme, scheme.ID).Error; err != nil {
		return nil, errors.New("failed to load scheme details: " + err.Error())
	}

	return s.toSchemeResponse(&scheme), nil
}

// DeleteSchemeOfWork soft deletes a scheme
func (s *SchemeOfWorkService) DeleteSchemeOfWork(id string) error {
	schemeID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid scheme ID")
	}

	var scheme models.SchemeOfWork
	if err := s.db.Where("id = ? AND deleted_at IS NULL", schemeID).First(&scheme).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("scheme of work not found")
		}
		return errors.New("failed to fetch scheme: " + err.Error())
	}

	if err := s.db.Delete(&scheme).Error; err != nil {
		return errors.New("failed to delete scheme: " + err.Error())
	}

	return nil
}

// validateSchemeRequest validates the scheme request
func (s *SchemeOfWorkService) validateSchemeRequest(req *dto.CreateSchemeOfWorkRequest) error {
	if req.SubjectID == "" {
		return errors.New("subject ID is required")
	}
	if req.GradeID == "" {
		return errors.New("grade ID is required")
	}
	if req.Term == "" {
		return errors.New("term is required")
	}
	if req.Term != "first" && req.Term != "second" && req.Term != "third" {
		return errors.New("term must be 'first', 'second', or 'third'")
	}
	if req.Title == "" {
		return errors.New("title is required")
	}
	if req.Status != "" && req.Status != "draft" && req.Status != "published" && req.Status != "archived" {
		return errors.New("status must be 'draft', 'published', or 'archived'")
	}
	return nil
}

// toSchemeResponse converts model to response DTO
func (s *SchemeOfWorkService) toSchemeResponse(scheme *models.SchemeOfWork) *dto.SchemeOfWorkResponse {
	response := &dto.SchemeOfWorkResponse{
		ID:          scheme.ID.String(),
		SubjectID:   scheme.SubjectID.String(),
		GradeID:     scheme.GradeID.String(),
		Term:        scheme.Term,
		Title:       scheme.Title,
		Description: scheme.Description,
		Status:      scheme.Status,
		CreatedBy:   scheme.CreatedBy.String(),
		CreatedAt:   scheme.CreatedAt,
		UpdatedAt:   scheme.UpdatedAt,
	}

	// Add subject details if preloaded
	if scheme.Subject.ID != uuid.Nil {
		response.Subject = &dto.SubjectResponse{
			ID:          scheme.Subject.ID.String(),
			Name:        scheme.Subject.Name,
			Code:        scheme.Subject.Code,
			Description: scheme.Subject.Description,
		}
	}

	// Add grade details if preloaded
	if scheme.Grade.ID != uuid.Nil {
		response.Grade = &dto.ClassGradeResponse{
			ID:          scheme.Grade.ID.String(),
			Name:        scheme.Grade.Name,
			Code:        scheme.Grade.Code,
			Level:       scheme.Grade.Level,
			Description: scheme.Grade.Description,
		}
	}

	// Add creator details if preloaded
	if scheme.Creator.ID != uuid.Nil {
		response.Creator = &dto.UserResponse{
			ID:        scheme.Creator.ID.String(),
			FirstName: scheme.Creator.FirstName,
			LastName:  scheme.Creator.LastName,
			Email:     scheme.Creator.Email,
			Phone:     scheme.Creator.Phone,
			Role:      scheme.Creator.Role,
		}
	}

	// Add lessons if preloaded
	if len(scheme.Lessons) > 0 {
		lessons := make([]dto.LessonResponse, len(scheme.Lessons))
		for i, lesson := range scheme.Lessons {
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