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

	// Check if subject exists
	var subject models.Subject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subject not found")
		}
		return nil, errors.New("failed to verify subject: " + err.Error())
	}

	// Check if scheme already exists for this week, grade, term, and subject
	var existing models.SchemeOfWork
	if err := s.db.Where("subject_id = ? AND grade = ? AND term = ? AND week = ? AND deleted_at IS NULL",
		subjectID, req.Grade, req.Term, req.Week).First(&existing).Error; err == nil {
		return nil, errors.New("scheme for this week already exists")
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "draft"
	}

	// Create scheme of work
	scheme := &models.SchemeOfWork{
		ID:                 uuid.New(),
		SubjectID:          subjectID,
		Grade:              req.Grade,
		Term:               req.Term,
		Week:               req.Week,
		Topic:              strings.TrimSpace(req.Topic),
		Subtopics:          strings.TrimSpace(req.Subtopics),
		Objectives:         strings.TrimSpace(req.Objectives),
		Activities:         strings.TrimSpace(req.Activities),
		TeachingResources:  strings.TrimSpace(req.TeachingResources),
		Assessment:         strings.TrimSpace(req.Assessment),
		Status:             status,
		CreatedBy:          userID,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := s.db.Create(scheme).Error; err != nil {
		return nil, errors.New("failed to create scheme of work: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("Subject").Preload("Creator").First(scheme, scheme.ID).Error; err != nil {
		return nil, errors.New("failed to load scheme details: " + err.Error())
	}

	return s.toSchemeResponse(scheme), nil
}

// BulkCreateSchemes creates multiple schemes of work
func (s *SchemeOfWorkService) BulkCreateSchemes(req *dto.BulkCreateSchemeOfWorkRequest, userID uuid.UUID) (*dto.BulkSchemeResult, error) {
	// Parse UUIDs
	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		return nil, errors.New("invalid subject ID format")
	}

	// Check if subject exists
	var subject models.Subject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subject not found")
		}
		return nil, errors.New("failed to verify subject: " + err.Error())
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

	// Get existing weeks to avoid duplicates
	var existingWeeks []int
	if err := s.db.Model(&models.SchemeOfWork{}).
		Where("subject_id = ? AND grade = ? AND term = ? AND deleted_at IS NULL",
			subjectID, req.Grade, req.Term).
		Pluck("week", &existingWeeks).Error; err != nil {
		return nil, fmt.Errorf("failed to check existing weeks: %w", err)
	}

	existingWeekMap := make(map[int]bool)
	for _, week := range existingWeeks {
		existingWeekMap[week] = true
	}

	for _, schemeReq := range req.Schemes {
		// Check for duplicate week
		if existingWeekMap[schemeReq.Week] {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkSchemeError{
				Week:  schemeReq.Week,
				Topic: schemeReq.Topic,
				Error: "week already has a scheme",
			})
			continue
		}

		// Create scheme
		scheme := &models.SchemeOfWork{
			ID:                 uuid.New(),
			SubjectID:          subjectID,
			Grade:              req.Grade,
			Term:               req.Term,
			Week:               schemeReq.Week,
			Topic:              strings.TrimSpace(schemeReq.Topic),
			Subtopics:          strings.TrimSpace(schemeReq.Subtopics),
			Objectives:         strings.TrimSpace(schemeReq.Objectives),
			Activities:         strings.TrimSpace(schemeReq.Activities),
			TeachingResources:  strings.TrimSpace(schemeReq.TeachingResources),
			Assessment:         strings.TrimSpace(schemeReq.Assessment),
			Status:             status,
			CreatedBy:          userID,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}

		if err := s.db.Create(scheme).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkSchemeError{
				Week:  schemeReq.Week,
				Topic: schemeReq.Topic,
				Error: "failed to create scheme: " + err.Error(),
			})
			continue
		}

		// Preload relationships
		if err := s.db.Preload("Subject").Preload("Creator").First(scheme, scheme.ID).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkSchemeError{
				Week:  schemeReq.Week,
				Topic: schemeReq.Topic,
				Error: "failed to load scheme details",
			})
			continue
		}

		existingWeekMap[schemeReq.Week] = true
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
		params.SortBy = "week"
	}
	if params.SortOrder == "" {
		params.SortOrder = "asc"
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

	if params.Grade != "" {
		query = query.Where("grade = ?", params.Grade)
	}

	if params.Term != "" {
		query = query.Where("term = ?", params.Term)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.Week > 0 {
		query = query.Where("week = ?", params.Week)
	}

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(topic) LIKE ? OR LOWER(objectives) LIKE ?", search, search)
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
	if err := query.Preload("Subject").Preload("Creator").Find(&schemes).Error; err != nil {
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
		Preload("Subject").
		Preload("Creator").
		Order("term, week ASC").
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
func (s *SchemeOfWorkService) GetSchemesByGrade(grade string) ([]dto.SchemeOfWorkResponse, error) {
	var schemes []models.SchemeOfWork
	if err := s.db.Where("grade = ? AND deleted_at IS NULL", grade).
		Preload("Subject").
		Preload("Creator").
		Order("term, week ASC").
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
func (s *SchemeOfWorkService) GetSchemesByGradeAndTerm(grade, term string) ([]dto.SchemeOfWorkResponse, error) {
	var schemes []models.SchemeOfWork
	if err := s.db.Where("grade = ? AND term = ? AND deleted_at IS NULL", grade, term).
		Preload("Subject").
		Preload("Creator").
		Order("week ASC").
		Find(&schemes).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch schemes: %w", err)
	}

	responses := make([]dto.SchemeOfWorkResponse, len(schemes))
	for i, scheme := range schemes {
		responses[i] = *s.toSchemeResponse(&scheme)
	}

	return responses, nil
}

// GetSchemesByTeacher retrieves all schemes for subjects taught by a teacher
func (s *SchemeOfWorkService) GetSchemesByTeacher(teacherID string) ([]dto.SchemeOfWorkResponse, error) {
	tID, err := uuid.Parse(teacherID)
	if err != nil {
		return nil, errors.New("invalid teacher ID")
	}

	// Get all subject assignments for this teacher
	var assignments []models.TeacherSubjectAssignment
	if err := s.db.Where("teacher_id = ? AND status = ? AND deleted_at IS NULL", tID, "active").
		Find(&assignments).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch teacher assignments: %w", err)
	}

	if len(assignments) == 0 {
		return []dto.SchemeOfWorkResponse{}, nil
	}

	// Get all subject IDs
	subjectIDs := make([]uuid.UUID, len(assignments))
	for i, assignment := range assignments {
		subjectIDs[i] = assignment.SubjectID
	}

	// Get schemes for these subjects
	var schemes []models.SchemeOfWork
	if err := s.db.Where("subject_id IN ? AND deleted_at IS NULL", subjectIDs).
		Preload("Subject").
		Preload("Creator").
		Order("grade, term, week ASC").
		Find(&schemes).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch schemes: %w", err)
	}

	responses := make([]dto.SchemeOfWorkResponse, len(schemes))
	for i, scheme := range schemes {
		responses[i] = *s.toSchemeResponse(&scheme)
	}

	return responses, nil
}

// GetSchemeStats retrieves statistics for schemes
func (s *SchemeOfWorkService) GetSchemeStats(filter map[string]interface{}) (*dto.SchemeOfWorkStats, error) {
	query := s.db.Model(&models.SchemeOfWork{}).Where("deleted_at IS NULL")

	// Apply filters
	if subjectID, ok := filter["subject_id"].(string); ok && subjectID != "" {
		if id, err := uuid.Parse(subjectID); err == nil {
			query = query.Where("subject_id = ?", id)
		}
	}
	if grade, ok := filter["grade"].(string); ok && grade != "" {
		query = query.Where("grade = ?", grade)
	}
	if term, ok := filter["term"].(string); ok && term != "" {
		query = query.Where("term = ?", term)
	}

	var stats dto.SchemeOfWorkStats
	if err := query.Count(&stats.TotalSchemes).Error; err != nil {
		return nil, fmt.Errorf("failed to count total schemes: %w", err)
	}

	// Count by status
	if err := query.Where("status = ?", "draft").Count(&stats.DraftSchemes).Error; err != nil {
		return nil, fmt.Errorf("failed to count draft schemes: %w", err)
	}
	if err := query.Where("status = ?", "published").Count(&stats.PublishedSchemes).Error; err != nil {
		return nil, fmt.Errorf("failed to count published schemes: %w", err)
	}
	if err := query.Where("status = ?", "archived").Count(&stats.ArchivedSchemes).Error; err != nil {
		return nil, fmt.Errorf("failed to count archived schemes: %w", err)
	}

	// Count unique weeks
	if err := query.Distinct("week").Count(&stats.TotalWeeks).Error; err != nil {
		return nil, fmt.Errorf("failed to count weeks: %w", err)
	}

	return &stats, nil
}

// GetSchemeOverview retrieves an overview of schemes
func (s *SchemeOfWorkService) GetSchemeOverview(subjectID, grade, term string) (*dto.SchemeOverview, error) {
	sID, err := uuid.Parse(subjectID)
	if err != nil {
		return nil, errors.New("invalid subject ID")
	}

	var subject models.Subject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", sID).First(&subject).Error; err != nil {
		return nil, errors.New("subject not found")
	}

	var schemes []models.SchemeOfWork
	if err := s.db.Where("subject_id = ? AND grade = ? AND term = ? AND deleted_at IS NULL",
		sID, grade, term).
		Find(&schemes).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch schemes: %w", err)
	}

	// Calculate weeks covered
	weeksCovered := 0
	status := "not_started"
	if len(schemes) > 0 {
		weeksCovered = len(schemes)
		// Check if all weeks are published
		allPublished := true
		for _, scheme := range schemes {
			if scheme.Status != "published" {
				allPublished = false
				break
			}
		}
		if allPublished {
			status = "complete"
		} else {
			status = "in_progress"
		}
	}

	// Calculate progress (assuming 13 weeks per term)
	totalWeeks := 13
	progress := float64(weeksCovered) / float64(totalWeeks) * 100

	overview := &dto.SchemeOverview{
		SubjectID:    subjectID,
		SubjectName:  subject.Name,
		Grade:        grade,
		Term:         term,
		TotalWeeks:   totalWeeks,
		WeeksCovered: weeksCovered,
		Progress:     fmt.Sprintf("%.0f%%", progress),
		Status:       status,
	}

	return overview, nil
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

	if req.Grade != "" {
		scheme.Grade = req.Grade
	}

	if req.Term != "" {
		scheme.Term = req.Term
	}

	if req.Week > 0 {
		scheme.Week = req.Week
	}

	if req.Topic != "" {
		scheme.Topic = strings.TrimSpace(req.Topic)
	}

	if req.Subtopics != "" {
		scheme.Subtopics = strings.TrimSpace(req.Subtopics)
	}

	if req.Objectives != "" {
		scheme.Objectives = strings.TrimSpace(req.Objectives)
	}

	if req.Activities != "" {
		scheme.Activities = strings.TrimSpace(req.Activities)
	}

	if req.TeachingResources != "" {
		scheme.TeachingResources = strings.TrimSpace(req.TeachingResources)
	}

	if req.Assessment != "" {
		scheme.Assessment = strings.TrimSpace(req.Assessment)
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
	if err := s.db.Preload("Subject").Preload("Creator").First(&scheme, scheme.ID).Error; err != nil {
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
	if req.Grade == "" {
		return errors.New("grade is required")
	}
	if req.Term == "" {
		return errors.New("term is required")
	}
	if req.Term != "first" && req.Term != "second" && req.Term != "third" {
		return errors.New("term must be 'first', 'second', or 'third'")
	}
	if req.Week < 1 || req.Week > 52 {
		return errors.New("week must be between 1 and 52")
	}
	if req.Topic == "" {
		return errors.New("topic is required")
	}
	if req.Objectives == "" {
		return errors.New("objectives are required")
	}
	if req.Status != "" && req.Status != "draft" && req.Status != "published" && req.Status != "archived" {
		return errors.New("status must be 'draft', 'published', or 'archived'")
	}
	return nil
}

// toSchemeResponse converts model to response DTO
func (s *SchemeOfWorkService) toSchemeResponse(scheme *models.SchemeOfWork) *dto.SchemeOfWorkResponse {
	response := &dto.SchemeOfWorkResponse{
		ID:                 scheme.ID.String(),
		SubjectID:          scheme.SubjectID.String(),
		Grade:              scheme.Grade,
		Term:               scheme.Term,
		Week:               scheme.Week,
		Topic:              scheme.Topic,
		Subtopics:          scheme.Subtopics,
		Objectives:         scheme.Objectives,
		Activities:         scheme.Activities,
		TeachingResources:  scheme.TeachingResources,
		Assessment:         scheme.Assessment,
		Status:             scheme.Status,
		CreatedBy:          scheme.CreatedBy.String(),
		CreatedAt:          scheme.CreatedAt,
		UpdatedAt:          scheme.UpdatedAt,
	}

	// Add subject details if preloaded
	if scheme.Subject.ID != uuid.Nil {
		response.Subject = &dto.SubjectResponse{
			ID:          scheme.Subject.ID.String(),
			Name:        scheme.Subject.Name,
			Code:        scheme.Subject.Code,
			Description: scheme.Subject.Description,
			CreatedAt:   scheme.Subject.CreatedAt,
			UpdatedAt:   scheme.Subject.UpdatedAt,
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