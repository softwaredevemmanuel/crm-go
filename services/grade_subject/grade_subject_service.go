// services/grade_subject_service.go
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
	"log"
)

type GradeSubjectService struct {
	db *gorm.DB
}

func NewGradeSubjectService(db *gorm.DB) *GradeSubjectService {
	return &GradeSubjectService{db: db}
}

// CreateGradeSubject creates a new grade-subject mapping
func (s *GradeSubjectService) CreateGradeSubject(req *dto.CreateGradeSubjectRequest, userID uuid.UUID) (*dto.GradeSubjectResponse, error) {
	// Validate input
	if err := s.validateGradeSubjectRequest(req); err != nil {
		return nil, err
	}

	// Parse UUIDs
	gradeID, err := uuid.Parse(req.GradeID)
	if err != nil {
		return nil, errors.New("invalid grade ID format")
	}

	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		return nil, errors.New("invalid subject ID format")
	}

	// Check if grade exists
	var grade models.ClassGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", gradeID).First(&grade).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("grade not found")
		}
		return nil, errors.New("failed to verify grade: " + err.Error())
	}

	// Check if subject exists
	var subject models.Subject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subject not found")
		}
		return nil, errors.New("failed to verify subject: " + err.Error())
	}

	// Check if mapping already exists
	var existing models.GradeSubject
	if err := s.db.Where("grade_id = ? AND subject_id = ? AND deleted_at IS NULL", gradeID, subjectID).First(&existing).Error; err == nil {
		return nil, errors.New("this grade-subject mapping already exists")
	}

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = "active"
	}
	compulsory := req.IsCompulsory
	if compulsory == true{
			log.Printf("true")

	}else{
					log.Printf("false")

	}
	


	// Create new grade-subject mapping
	gradeSubject := &models.GradeSubject{
		ID:           uuid.New(),
		GradeID:      gradeID,
		SubjectID:    subjectID,
		Status:       status,
		IsCompulsory: compulsory,
		CreatedBy:    userID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.db.Create(gradeSubject).Error; err != nil {
		return nil, errors.New("failed to create grade-subject mapping: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("Grade").Preload("Subject").First(gradeSubject, gradeSubject.ID).Error; err != nil {
		return nil, errors.New("failed to load mapping details: " + err.Error())
	}

	return s.toGradeSubjectResponse(gradeSubject), nil
}

// BulkCreateGradeSubjects creates multiple grade-subject mappings
func (s *GradeSubjectService) BulkCreateGradeSubjects(req *dto.BulkCreateGradeSubjectRequest, userID uuid.UUID) (*dto.BulkCreateResult, error) {
	// Parse Grade ID
	gradeID, err := uuid.Parse(req.GradeID)
	if err != nil {
		return nil, errors.New("invalid grade ID format")
	}

	// Check if grade exists
	var grade models.ClassGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", gradeID).First(&grade).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("grade not found")
		}
		return nil, errors.New("failed to verify grade: " + err.Error())
	}

	result := &dto.BulkCreateResult{
		Created: []dto.GradeSubjectResponse{},
		Errors:  []dto.BulkCreateError{},
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	for _, subjectIDStr := range req.SubjectIDs {
		subjectID, err := uuid.Parse(subjectIDStr)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkCreateError{
				SubjectID: subjectIDStr,
				Error:     "invalid subject ID format",
			})
			continue
		}

		// Check if subject exists
		var subject models.Subject
		if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkCreateError{
				SubjectID: subjectIDStr,
				Error:     "subject not found",
			})
			continue
		}

		// Check if mapping already exists
		var existing models.GradeSubject
		if err := s.db.Where("grade_id = ? AND subject_id = ? AND deleted_at IS NULL", gradeID, subjectID).First(&existing).Error; err == nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkCreateError{
				SubjectID: subjectIDStr,
				Error:     "mapping already exists",
			})
			continue
		}

		// Create mapping
		gradeSubject := &models.GradeSubject{
			ID:           uuid.New(),
			GradeID:      gradeID,
			SubjectID:    subjectID,
			Status:       status,
			IsCompulsory: req.IsCompulsory,
			CreatedBy:    userID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := s.db.Create(gradeSubject).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkCreateError{
				SubjectID: subjectIDStr,
				Error:     "failed to create mapping: " + err.Error(),
			})
			continue
		}

		// Preload relationships
		if err := s.db.Preload("Grade").Preload("Subject").First(gradeSubject, gradeSubject.ID).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkCreateError{
				SubjectID: subjectIDStr,
				Error:     "failed to load mapping details",
			})
			continue
		}

		result.SuccessCount++
		result.Created = append(result.Created, *s.toGradeSubjectResponse(gradeSubject))
	}

	return result, nil
}

// GetAllGradeSubjects retrieves all grade-subject mappings with pagination and filters
func (s *GradeSubjectService) GetAllGradeSubjects(params *dto.GradeSubjectQueryParams) (*dto.GradeSubjectListResponse, error) {
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
	query := s.db.Model(&models.GradeSubject{}).Where("grade_subjects.deleted_at IS NULL")

	// Apply filters
	if params.GradeID != "" {
		gradeID, err := uuid.Parse(params.GradeID)
		if err == nil {
			query = query.Where("grade_subjects.grade_id = ?", gradeID)
		}
	}

	if params.SubjectID != "" {
		subjectID, err := uuid.Parse(params.SubjectID)
		if err == nil {
			query = query.Where("grade_subjects.subject_id = ?", subjectID)
		}
	}

	if params.Status != "" {
		query = query.Where("grade_subjects.status = ?", params.Status)
	}

	if params.IsCompulsory != nil {
		query = query.Where("grade_subjects.is_compulsory = ?", *params.IsCompulsory)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count grade-subject mappings: %w", err)
	}

	// Apply sorting
	sortDirection := "DESC"
	if strings.ToLower(params.SortOrder) == "asc" {
		sortDirection = "ASC"
	}
	query = query.Order("grade_subjects." + params.SortBy + " " + sortDirection)

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	// Execute with preloads
	var gradeSubjects []models.GradeSubject
	if err := query.Preload("Grade").Preload("Subject").Find(&gradeSubjects).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch grade-subject mappings: %w", err)
	}

	// Convert to response
	responses := make([]dto.GradeSubjectResponse, len(gradeSubjects))
	for i, gs := range gradeSubjects {
		responses[i] = *s.toGradeSubjectResponse(&gs)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.GradeSubjectListResponse{
		GradeSubjects: responses,
		Total:         total,
		Page:          params.Page,
		Limit:         params.Limit,
		TotalPages:    totalPages,
	}, nil
}

// GetGradeSubjectByID retrieves a single grade-subject mapping by ID
func (s *GradeSubjectService) GetGradeSubjectByID(id string) (*dto.GradeSubjectResponse, error) {
	gsID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid grade-subject ID")
	}

	var gradeSubject models.GradeSubject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", gsID).
		Preload("Grade").
		Preload("Subject").
		First(&gradeSubject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("grade-subject mapping not found")
		}
		return nil, errors.New("failed to fetch grade-subject mapping: " + err.Error())
	}

	return s.toGradeSubjectResponse(&gradeSubject), nil
}

// GetSubjectsByGrade retrieves all subjects for a specific grade
func (s *GradeSubjectService) GetSubjectsByGrade(gradeID string) ([]dto.GradeSubjectResponse, error) {
	gID, err := uuid.Parse(gradeID)
	if err != nil {
		return nil, errors.New("invalid grade ID")
	}

	var gradeSubjects []models.GradeSubject
	if err := s.db.Where("grade_id = ? AND deleted_at IS NULL", gID).
		Preload("Grade").
		Preload("Subject").
		Order("is_compulsory DESC, created_at ASC").
		Find(&gradeSubjects).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch subjects for grade: %w", err)
	}

	responses := make([]dto.GradeSubjectResponse, len(gradeSubjects))
	for i, gs := range gradeSubjects {
		responses[i] = *s.toGradeSubjectResponse(&gs)
	}

	return responses, nil
}

// GetGradesBySubject retrieves all grades for a specific subject
func (s *GradeSubjectService) GetGradesBySubject(subjectID string) ([]dto.GradeSubjectResponse, error) {
	sID, err := uuid.Parse(subjectID)
	if err != nil {
		return nil, errors.New("invalid subject ID")
	}

	var gradeSubjects []models.GradeSubject
	if err := s.db.Where("subject_id = ? AND deleted_at IS NULL", sID).
		Preload("Grade").
		Preload("Subject").
		Find(&gradeSubjects).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch grades for subject: %w", err)
	}

	responses := make([]dto.GradeSubjectResponse, len(gradeSubjects))
	for i, gs := range gradeSubjects {
		responses[i] = *s.toGradeSubjectResponse(&gs)
	}

	return responses, nil
}

// UpdateGradeSubject updates an existing grade-subject mapping
func (s *GradeSubjectService) UpdateGradeSubject(id string, req *dto.UpdateGradeSubjectRequest, userID uuid.UUID) (*dto.GradeSubjectResponse, error) {
	gsID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid grade-subject ID")
	}

	// Find existing mapping
	var gradeSubject models.GradeSubject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", gsID).First(&gradeSubject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("grade-subject mapping not found")
		}
		return nil, errors.New("failed to fetch grade-subject mapping: " + err.Error())
	}

	// Update fields
	if req.Status != "" {
		gradeSubject.Status = req.Status
	}
	if req.IsCompulsory != nil {
		gradeSubject.IsCompulsory = *req.IsCompulsory
	}

	gradeSubject.UpdatedAt = time.Now()

	if err := s.db.Save(&gradeSubject).Error; err != nil {
		return nil, errors.New("failed to update grade-subject mapping: " + err.Error())
	}

	// Preload relationships
	if err := s.db.Preload("Grade").Preload("Subject").First(&gradeSubject, gradeSubject.ID).Error; err != nil {
		return nil, errors.New("failed to load mapping details: " + err.Error())
	}

	return s.toGradeSubjectResponse(&gradeSubject), nil
}

// DeleteGradeSubject soft deletes a grade-subject mapping
func (s *GradeSubjectService) DeleteGradeSubject(id string) error {
	gsID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid grade-subject ID")
	}

	var gradeSubject models.GradeSubject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", gsID).First(&gradeSubject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("grade-subject mapping not found")
		}
		return errors.New("failed to fetch grade-subject mapping: " + err.Error())
	}

	if err := s.db.Delete(&gradeSubject).Error; err != nil {
		return errors.New("failed to delete grade-subject mapping: " + err.Error())
	}

	return nil
}

// validateGradeSubjectRequest validates the grade-subject request
func (s *GradeSubjectService) validateGradeSubjectRequest(req *dto.CreateGradeSubjectRequest) error {
	if req.GradeID == "" {
		return errors.New("grade ID is required")
	}
	if req.SubjectID == "" {
		return errors.New("subject ID is required")
	}
	if req.Status != "" && req.Status != "active" && req.Status != "inactive" {
		return errors.New("status must be 'active' or 'inactive'")
	}
	return nil
}

// toGradeSubjectResponse converts model to response DTO
func (s *GradeSubjectService) toGradeSubjectResponse(gs *models.GradeSubject) *dto.GradeSubjectResponse {
	response := &dto.GradeSubjectResponse{
		ID:           gs.ID.String(),
		GradeID:      gs.GradeID.String(),
		SubjectID:    gs.SubjectID.String(),
		Status:       gs.Status,
		IsCompulsory: gs.IsCompulsory,
		CreatedBy:    gs.CreatedBy.String(),
		CreatedAt:    gs.CreatedAt,
		UpdatedAt:    gs.UpdatedAt,
	}

	// Add grade details if preloaded
	if gs.Grade.ID != uuid.Nil {
		response.Grade = &dto.ClassGradeResponse{
			ID:          gs.Grade.ID.String(),
			Name:        gs.Grade.Name,
			Code:        gs.Grade.Code,
			Level:       gs.Grade.Level,
			Description: gs.Grade.Description,
			Status:      gs.Grade.Status,
			CreatedAt:   gs.Grade.CreatedAt,
			UpdatedAt:   gs.Grade.UpdatedAt,
		}
	}

	// Add subject details if preloaded
	if gs.Subject.ID != uuid.Nil {
		response.Subject = &dto.SubjectResponse{
			ID:          gs.Subject.ID.String(),
			Name:        gs.Subject.Name,
			Code:        gs.Subject.Code,
			Description: gs.Subject.Description,
			Credits:     gs.Subject.Credits,
			Status:      gs.Subject.Status,
			CreatedAt:   gs.Subject.CreatedAt,
			UpdatedAt:   gs.Subject.UpdatedAt,
		}
	}

	return response
}