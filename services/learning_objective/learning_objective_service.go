// services/learning_objective_service.go
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

type LearningObjectiveService struct {
	db *gorm.DB
}

func NewLearningObjectiveService(db *gorm.DB) *LearningObjectiveService {
	return &LearningObjectiveService{db: db}
}

// CreateLearningObjective creates a new learning objective
func (s *LearningObjectiveService) CreateLearningObjective(req *dto.CreateLearningObjectiveRequest) (*dto.LearningObjectiveResponse, error) {
	// Validate input
	if err := s.validateObjectiveRequest(req); err != nil {
		return nil, err
	}

	// Parse UUIDs
	schemeOfWorkItemID, err := uuid.Parse(req.SchemeOfWorkItemID)
	if err != nil {
		return nil, errors.New("invalid scheme of work item ID format")
	}

	// Check if scheme of work item exists
	var schemeItem models.SchemeOfWorkItem
	if err := s.db.Where("id = ? AND deleted_at IS NULL", schemeOfWorkItemID).First(&schemeItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("scheme of work item not found")
		}
		return nil, errors.New("failed to verify scheme of work item: " + err.Error())
	}

	// Set default sequence if not provided
	sequence := req.Sequence
	if sequence == 0 {
		var maxSequence int
		if err := s.db.Model(&models.LearningObjective{}).
			Where("scheme_of_work_item_id = ? AND deleted_at IS NULL", schemeOfWorkItemID).
			Select("COALESCE(MAX(sequence), 0)").
			Scan(&maxSequence).Error; err != nil {
			return nil, errors.New("failed to determine sequence: " + err.Error())
		}
		sequence = maxSequence + 1
	}

	// Create learning objective
	objective := &models.LearningObjective{
		ID:                 uuid.New(),
		SchemeOfWorkItemID: schemeOfWorkItemID,
		Objective:          req.Objective,
		Sequence:           sequence,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := s.db.Create(objective).Error; err != nil {
		return nil, errors.New("failed to create learning objective: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("SchemeOfWorkItem").First(objective, objective.ID).Error; err != nil {
		return nil, errors.New("failed to load objective details: " + err.Error())
	}

	return s.toObjectiveResponse(objective), nil
}

// BulkCreateLearningObjectives creates multiple learning objectives
func (s *LearningObjectiveService) BulkCreateLearningObjectives(req *dto.BulkCreateLearningObjectivesRequest) (*dto.BulkLearningObjectiveResult, error) {
	// Parse UUIDs
	schemeOfWorkItemID, err := uuid.Parse(req.SchemeOfWorkItemID)
	if err != nil {
		return nil, errors.New("invalid scheme of work item ID format")
	}

	// Check if scheme of work item exists
	var schemeItem models.SchemeOfWorkItem
	if err := s.db.Where("id = ? AND deleted_at IS NULL", schemeOfWorkItemID).First(&schemeItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("scheme of work item not found")
		}
		return nil, errors.New("failed to verify scheme of work item: " + err.Error())
	}

	result := &dto.BulkLearningObjectiveResult{
		Created: []dto.LearningObjectiveResponse{},
		Errors:  []dto.BulkLearningObjectiveError{},
	}

	// Get max sequence
	var maxSequence int
	if err := s.db.Model(&models.LearningObjective{}).
		Where("scheme_of_work_item_id = ? AND deleted_at IS NULL", schemeOfWorkItemID).
		Select("COALESCE(MAX(sequence), 0)").
		Scan(&maxSequence).Error; err != nil {
		return nil, errors.New("failed to determine sequence: " + err.Error())
	}
	currentSequence := maxSequence

	for _, objReq := range req.Objectives {
		// Validate
		if objReq.Objective == "" {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkLearningObjectiveError{
				Objective: objReq.Objective,
				Error:     "objective is required",
			})
			continue
		}

		// Set sequence
		sequence := objReq.Sequence
		if sequence == 0 {
			currentSequence++
			sequence = currentSequence
		}

		// Create objective
		objective := &models.LearningObjective{
			ID:                 uuid.New(),
			SchemeOfWorkItemID: schemeOfWorkItemID,
			Objective:          objReq.Objective,
			Sequence:           sequence,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}

		if err := s.db.Create(objective).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkLearningObjectiveError{
				Objective: objReq.Objective,
				Error:     "failed to create objective: " + err.Error(),
			})
			continue
		}

		// Preload relationships
		if err := s.db.Preload("SchemeOfWorkItem").First(objective, objective.ID).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkLearningObjectiveError{
				Objective: objReq.Objective,
				Error:     "failed to load objective details",
			})
			continue
		}

		result.SuccessCount++
		result.Created = append(result.Created, *s.toObjectiveResponse(objective))
	}

	return result, nil
}

// GetAllObjectives retrieves all learning objectives with pagination and filters
func (s *LearningObjectiveService) GetAllObjectives(params *dto.LearningObjectiveQueryParams) (*dto.LearningObjectiveListResponse, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = "sequence"
	}
	if params.SortOrder == "" {
		params.SortOrder = "asc"
	}

	// Build query
	query := s.db.Model(&models.LearningObjective{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.SchemeOfWorkItemID != "" {
		itemID, err := uuid.Parse(params.SchemeOfWorkItemID)
		if err == nil {
			query = query.Where("scheme_of_work_item_id = ?", itemID)
		}
	}

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(objective) LIKE ?", search)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count objectives: %w", err)
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
	var objectives []models.LearningObjective
	if err := query.Preload("SchemeOfWorkItem").Find(&objectives).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch objectives: %w", err)
	}

	// Convert to response
	responses := make([]dto.LearningObjectiveResponse, len(objectives))
	for i, objective := range objectives {
		responses[i] = *s.toObjectiveResponse(&objective)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.LearningObjectiveListResponse{
		Objectives: responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetObjectiveByID retrieves a single learning objective by ID
func (s *LearningObjectiveService) GetObjectiveByID(id string) (*dto.LearningObjectiveResponse, error) {
	objectiveID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid objective ID")
	}

	var objective models.LearningObjective
	if err := s.db.Where("id = ? AND deleted_at IS NULL", objectiveID).
		Preload("SchemeOfWorkItem").
		First(&objective).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("learning objective not found")
		}
		return nil, errors.New("failed to fetch objective: " + err.Error())
	}

	return s.toObjectiveResponse(&objective), nil
}

// GetObjectivesBySchemeItem retrieves all objectives for a specific scheme of work item
func (s *LearningObjectiveService) GetObjectivesBySchemeItem(schemeItemID string) ([]dto.LearningObjectiveResponse, error) {
	sID, err := uuid.Parse(schemeItemID)
	if err != nil {
		return nil, errors.New("invalid scheme of work item ID")
	}

	var objectives []models.LearningObjective
	if err := s.db.Where("scheme_of_work_item_id = ? AND deleted_at IS NULL", sID).
		Preload("SchemeOfWorkItem").
		Order("sequence ASC").
		Find(&objectives).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch objectives: %w", err)
	}

	responses := make([]dto.LearningObjectiveResponse, len(objectives))
	for i, objective := range objectives {
		responses[i] = *s.toObjectiveResponse(&objective)
	}

	return responses, nil
}

// UpdateLearningObjective updates an existing learning objective
func (s *LearningObjectiveService) UpdateLearningObjective(id string, req *dto.UpdateLearningObjectiveRequest) (*dto.LearningObjectiveResponse, error) {
	objectiveID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid objective ID")
	}

	// Find existing objective
	var objective models.LearningObjective
	if err := s.db.Where("id = ? AND deleted_at IS NULL", objectiveID).First(&objective).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("learning objective not found")
		}
		return nil, errors.New("failed to fetch objective: " + err.Error())
	}

	// Update fields
	if req.Objective != "" {
		objective.Objective = req.Objective
	}

	if req.Sequence > 0 && req.Sequence != objective.Sequence {
		// Handle reordering
		if err := s.reorderObjectives(objective.SchemeOfWorkItemID, objective.Sequence, req.Sequence, objectiveID); err != nil {
			return nil, err
		}
		objective.Sequence = req.Sequence
	}

	objective.UpdatedAt = time.Now()

	if err := s.db.Save(&objective).Error; err != nil {
		return nil, errors.New("failed to update objective: " + err.Error())
	}

	// Preload relationships
	if err := s.db.Preload("SchemeOfWorkItem").First(&objective, objective.ID).Error; err != nil {
		return nil, errors.New("failed to load objective details: " + err.Error())
	}

	return s.toObjectiveResponse(&objective), nil
}

// DeleteLearningObjective soft deletes a learning objective
func (s *LearningObjectiveService) DeleteLearningObjective(id string) error {
	objectiveID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid objective ID")
	}

	var objective models.LearningObjective
	if err := s.db.Where("id = ? AND deleted_at IS NULL", objectiveID).First(&objective).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("learning objective not found")
		}
		return errors.New("failed to fetch objective: " + err.Error())
	}

	if err := s.db.Delete(&objective).Error; err != nil {
		return errors.New("failed to delete objective: " + err.Error())
	}

	return nil
}

// reorderObjectives reorders objectives when sequence changes
func (s *LearningObjectiveService) reorderObjectives(schemeItemID uuid.UUID, oldSeq, newSeq int, excludeID uuid.UUID) error {
	if oldSeq == newSeq {
		return nil
	}

	var objectives []models.LearningObjective
	if err := s.db.Where("scheme_of_work_item_id = ? AND deleted_at IS NULL AND id != ?", schemeItemID, excludeID).
		Order("sequence ASC").
		Find(&objectives).Error; err != nil {
		return errors.New("failed to fetch objectives for reordering: " + err.Error())
	}

	// Update sequences
	if oldSeq < newSeq {
		// Moving down: decrease sequences between oldSeq+1 and newSeq
		for _, obj := range objectives {
			if obj.Sequence > oldSeq && obj.Sequence <= newSeq {
				obj.Sequence--
				if err := s.db.Save(&obj).Error; err != nil {
					return errors.New("failed to update objective sequence: " + err.Error())
				}
			}
		}
	} else {
		// Moving up: increase sequences between newSeq and oldSeq-1
		for _, obj := range objectives {
			if obj.Sequence >= newSeq && obj.Sequence < oldSeq {
				obj.Sequence++
				if err := s.db.Save(&obj).Error; err != nil {
					return errors.New("failed to update objective sequence: " + err.Error())
				}
			}
		}
	}

	return nil
}

// validateObjectiveRequest validates the objective request
func (s *LearningObjectiveService) validateObjectiveRequest(req *dto.CreateLearningObjectiveRequest) error {
	if req.SchemeOfWorkItemID == "" {
		return errors.New("scheme of work item ID is required")
	}
	if req.Objective == "" {
		return errors.New("objective is required")
	}
	return nil
}

// toObjectiveResponse converts model to response DTO
func (s *LearningObjectiveService) toObjectiveResponse(objective *models.LearningObjective) *dto.LearningObjectiveResponse {
	response := &dto.LearningObjectiveResponse{
		ID:                 objective.ID.String(),
		SchemeOfWorkItemID: objective.SchemeOfWorkItemID.String(),
		Objective:          objective.Objective,
		Sequence:           objective.Sequence,
		CreatedAt:          objective.CreatedAt,
		UpdatedAt:          objective.UpdatedAt,
	}

	// Add scheme of work item details if preloaded
	if objective.SchemeOfWorkItem.ID != uuid.Nil {
		response.SchemeOfWorkItem = &dto.SchemeOfWorkItemResponse{
			ID:        objective.SchemeOfWorkItem.ID.String(),
			Topic:     objective.SchemeOfWorkItem.Topic,
			Subtopic:  objective.SchemeOfWorkItem.Subtopic,
			WeekStart: objective.SchemeOfWorkItem.WeekStart,
			WeekEnd:   objective.SchemeOfWorkItem.WeekEnd,
			Sequence:  objective.SchemeOfWorkItem.Sequence,
		}
	}

	return response
}