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
	academicSessionID, err := uuid.Parse(req.AcademicSessionID)
	if err != nil {
		return nil, errors.New("invalid academic session ID format")
	}

	termID, err := uuid.Parse(req.TermID)
	if err != nil {
		return nil, errors.New("invalid term ID format")
	}

	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		return nil, errors.New("invalid subject ID format")
	}

	classID, err := uuid.Parse(req.ClassID)
	if err != nil {
		return nil, errors.New("invalid class ID format")
	}

	// Check if academic session exists
	var academicSession models.AcademicSession
	if err := s.db.Where("id = ? AND deleted_at IS NULL", academicSessionID).First(&academicSession).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("academic session not found")
		}
		return nil, errors.New("failed to verify academic session: " + err.Error())
	}

	// Check if term exists
	var term models.Term
	if err := s.db.Where("id = ? AND deleted_at IS NULL", termID).First(&term).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("term not found")
		}
		return nil, errors.New("failed to verify term: " + err.Error())
	}

	// Check if subject exists
	var subject models.Subject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subject not found")
		}
		return nil, errors.New("failed to verify subject: " + err.Error())
	}

	// Check if class exists
	var class models.ClassGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", classID).First(&class).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("class not found")
		}
		return nil, errors.New("failed to verify class: " + err.Error())
	}

	// Check if scheme already exists for this combination
	var existing models.SchemeOfWork
	if err := s.db.Where("academic_session_id = ? AND term_id = ? AND subject_id = ? AND class_id = ? AND deleted_at IS NULL",
		academicSessionID, termID, subjectID, classID).First(&existing).Error; err == nil {
		return nil, errors.New("scheme of work already exists for this combination")
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "draft"
	}

	// Create scheme of work
	scheme := &models.SchemeOfWork{
		ID:                uuid.New(),
		AcademicSessionID: academicSessionID,
		TermID:            termID,
		SubjectID:         subjectID,
		ClassID:           classID,
		Title:             req.Title,
		Description:       req.Description,
		Status:            status,
		CreatedBy:         userID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.db.Create(scheme).Error; err != nil {
		return nil, errors.New("failed to create scheme of work: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("AcademicSession").Preload("Term").Preload("Subject").Preload("Class").Preload("Creator").First(scheme, scheme.ID).Error; err != nil {
		return nil, errors.New("failed to load scheme details: " + err.Error())
	}

	return s.toSchemeResponse(scheme), nil
}

// BulkCreateSchemes creates multiple schemes of work
func (s *SchemeOfWorkService) BulkCreateSchemes(req *dto.BulkCreateSchemesRequest, userID uuid.UUID) (*dto.BulkSchemeResult, error) {
	// Parse UUIDs
	academicSessionID, err := uuid.Parse(req.AcademicSessionID)
	if err != nil {
		return nil, errors.New("invalid academic session ID format")
	}

	termID, err := uuid.Parse(req.TermID)
	if err != nil {
		return nil, errors.New("invalid term ID format")
	}

	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		return nil, errors.New("invalid subject ID format")
	}

	classID, err := uuid.Parse(req.ClassID)
	if err != nil {
		return nil, errors.New("invalid class ID format")
	}

	// Verify all entities exist
	if err := s.verifyEntities(academicSessionID, termID, subjectID, classID); err != nil {
		return nil, err
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

	// Check if scheme already exists
	var existing models.SchemeOfWork
	if err := s.db.Where("academic_session_id = ? AND term_id = ? AND subject_id = ? AND class_id = ? AND deleted_at IS NULL",
		academicSessionID, termID, subjectID, classID).First(&existing).Error; err == nil {
		return nil, errors.New("scheme of work already exists for this combination")
	}

	// Create schemes
	for _, schemeReq := range req.Schemes {
		if schemeReq.Title == "" {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkSchemeError{
				Title: schemeReq.Title,
				Error: "title is required",
			})
			continue
		}

		scheme := &models.SchemeOfWork{
			ID:                uuid.New(),
			AcademicSessionID: academicSessionID,
			TermID:            termID,
			SubjectID:         subjectID,
			ClassID:           classID,
			Title:             schemeReq.Title,
			Description:       schemeReq.Description,
			Status:            status,
			CreatedBy:         userID,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
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
		if err := s.db.Preload("AcademicSession").Preload("Term").Preload("Subject").Preload("Class").Preload("Creator").First(scheme, scheme.ID).Error; err != nil {
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
	if params.AcademicSessionID != "" {
		sessionID, err := uuid.Parse(params.AcademicSessionID)
		if err == nil {
			query = query.Where("academic_session_id = ?", sessionID)
		}
	}

	if params.TermID != "" {
		termID, err := uuid.Parse(params.TermID)
		if err == nil {
			query = query.Where("term_id = ?", termID)
		}
	}

	if params.SubjectID != "" {
		subjectID, err := uuid.Parse(params.SubjectID)
		if err == nil {
			query = query.Where("subject_id = ?", subjectID)
		}
	}

	if params.ClassID != "" {
		classID, err := uuid.Parse(params.ClassID)
		if err == nil {
			query = query.Where("class_id = ?", classID)
		}
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
	if err := query.Preload("AcademicSession").Preload("Term").Preload("Subject").Preload("Class").Preload("Creator").Find(&schemes).Error; err != nil {
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
		Preload("AcademicSession").
		Preload("Term").
		Preload("Subject").
		Preload("Class").
		Preload("Creator").
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
		Preload("AcademicSession").
		Preload("Term").
		Preload("Subject").
		Preload("Class").
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

// GetSchemesByClass retrieves all schemes for a specific class
func (s *SchemeOfWorkService) GetSchemesByClass(classID string) ([]dto.SchemeOfWorkResponse, error) {
	cID, err := uuid.Parse(classID)
	if err != nil {
		return nil, errors.New("invalid class ID")
	}

	var schemes []models.SchemeOfWork
	if err := s.db.Where("class_id = ? AND deleted_at IS NULL", cID).
		Preload("AcademicSession").
		Preload("Term").
		Preload("Subject").
		Preload("Class").
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
	if req.AcademicSessionID != "" {
		sessionID, err := uuid.Parse(req.AcademicSessionID)
		if err != nil {
			return nil, errors.New("invalid academic session ID format")
		}
		var session models.AcademicSession
		if err := s.db.Where("id = ? AND deleted_at IS NULL", sessionID).First(&session).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("academic session not found")
			}
			return nil, errors.New("failed to verify academic session: " + err.Error())
		}
		scheme.AcademicSessionID = sessionID
	}

	if req.TermID != "" {
		termID, err := uuid.Parse(req.TermID)
		if err != nil {
			return nil, errors.New("invalid term ID format")
		}
		var term models.Term
		if err := s.db.Where("id = ? AND deleted_at IS NULL", termID).First(&term).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("term not found")
			}
			return nil, errors.New("failed to verify term: " + err.Error())
		}
		scheme.TermID = termID
	}

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

	if req.ClassID != "" {
		classID, err := uuid.Parse(req.ClassID)
		if err != nil {
			return nil, errors.New("invalid class ID format")
		}
		var class models.ClassGrade
		if err := s.db.Where("id = ? AND deleted_at IS NULL", classID).First(&class).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("class not found")
			}
			return nil, errors.New("failed to verify class: " + err.Error())
		}
		scheme.ClassID = classID
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
	if err := s.db.Preload("AcademicSession").Preload("Term").Preload("Subject").Preload("Class").Preload("Creator").First(&scheme, scheme.ID).Error; err != nil {
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

// verifyEntities verifies that all referenced entities exist
func (s *SchemeOfWorkService) verifyEntities(academicSessionID, termID, subjectID, classID uuid.UUID) error {
	// Check academic session
	var academicSession models.AcademicSession
	if err := s.db.Where("id = ? AND deleted_at IS NULL", academicSessionID).First(&academicSession).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("academic session not found")
		}
		return errors.New("failed to verify academic session: " + err.Error())
	}

	// Check term
	var term models.Term
	if err := s.db.Where("id = ? AND deleted_at IS NULL", termID).First(&term).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("term not found")
		}
		return errors.New("failed to verify term: " + err.Error())
	}

	// Check subject
	var subject models.Subject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("subject not found")
		}
		return errors.New("failed to verify subject: " + err.Error())
	}

	// Check class
	var class models.ClassGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", classID).First(&class).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("class not found")
		}
		return errors.New("failed to verify class: " + err.Error())
	}

	return nil
}

// validateSchemeRequest validates the scheme request
func (s *SchemeOfWorkService) validateSchemeRequest(req *dto.CreateSchemeOfWorkRequest) error {
	if req.AcademicSessionID == "" {
		return errors.New("academic session ID is required")
	}
	if req.TermID == "" {
		return errors.New("term ID is required")
	}
	if req.SubjectID == "" {
		return errors.New("subject ID is required")
	}
	if req.ClassID == "" {
		return errors.New("class ID is required")
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
		ID:                scheme.ID.String(),
		AcademicSessionID: scheme.AcademicSessionID.String(),
		TermID:            scheme.TermID.String(),
		SubjectID:         scheme.SubjectID.String(),
		ClassID:           scheme.ClassID.String(),
		Title:             scheme.Title,
		Description:       scheme.Description,
		Status:            scheme.Status,
		CreatedBy:         scheme.CreatedBy.String(),
		CreatedAt:         scheme.CreatedAt,
		UpdatedAt:         scheme.UpdatedAt,
	}

	// Add academic session details if preloaded
	if scheme.AcademicSession.ID != uuid.Nil {
		response.AcademicSession = &dto.AcademicSessionResponse{
			ID:           scheme.AcademicSession.ID.String(),
			AcademicYear: scheme.AcademicSession.AcademicYear,
			Code:         scheme.AcademicSession.Code,
			StartDate:    scheme.AcademicSession.StartDate,
			EndDate:      scheme.AcademicSession.EndDate,
			Status:       scheme.AcademicSession.Status,
			IsCurrent:    scheme.AcademicSession.IsCurrent,
			Description:  scheme.AcademicSession.Description,
		}
	}

	// Add term details if preloaded
	if scheme.Term.ID != uuid.Nil {
		response.Term = &dto.TermResponse{
			ID:         scheme.Term.ID.String(),
			Name:       scheme.Term.Name,
			Code:       scheme.Term.Code,
			TermNumber: scheme.Term.TermNumber,
			StartDate:  scheme.Term.StartDate,
			EndDate:    scheme.Term.EndDate,
			IsCurrent:  scheme.Term.IsCurrent,
			Status:     scheme.Term.Status,
		}
	}

	// Add subject details if preloaded
	if scheme.Subject.ID != uuid.Nil {
		response.Subject = &dto.SubjectResponse{
			ID:           scheme.Subject.ID.String(),
			Name:         scheme.Subject.Name,
			Code:         scheme.Subject.Code,
			Description:  scheme.Subject.Description,
			DepartmentID: scheme.Subject.DepartmentID.String(),
		}
	}

	// Add class details if preloaded
	if scheme.Class.ID != uuid.Nil {
		response.Class = &dto.ClassGradeResponse{
			ID:          scheme.Class.ID.String(),
			Name:        scheme.Class.Name,
			Code:        scheme.Class.Code,
			Level:       scheme.Class.Level,
			Description: scheme.Class.Description,
			Status:      scheme.Class.Status,
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

	return response
}