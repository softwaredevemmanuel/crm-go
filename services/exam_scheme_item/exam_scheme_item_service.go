// services/exam_scheme_item_service.go
package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"crm-go/models"
	"crm-go/dto"
)

type ExamSchemeItemService struct {
	db *gorm.DB
}

func NewExamSchemeItemService(db *gorm.DB) *ExamSchemeItemService {
	return &ExamSchemeItemService{db: db}
}

// CreateExamSchemeItem creates a new exam scheme item
func (s *ExamSchemeItemService) CreateExamSchemeItem(req *dto.CreateExamSchemeItemRequest) (*dto.ExamSchemeItemResponse, error) {
	// Validate input
	if err := s.validateExamSchemeItemRequest(req); err != nil {
		return nil, err
	}

	// Parse UUIDs
	examID, err := uuid.Parse(req.ExamID)
	if err != nil {
		return nil, errors.New("invalid exam ID format")
	}

	schemeOfWorkItemID, err := uuid.Parse(req.SchemeOfWorkItemID)
	if err != nil {
		return nil, errors.New("invalid scheme of work item ID format")
	}

	// Check if exam exists
	var exam models.Exam
	if err := s.db.Where("id = ? AND deleted_at IS NULL", examID).First(&exam).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("exam not found")
		}
		return nil, errors.New("failed to verify exam: " + err.Error())
	}

	// Check if scheme of work item exists
	var schemeItem models.SchemeOfWorkItem
	if err := s.db.Where("id = ? AND deleted_at IS NULL", schemeOfWorkItemID).First(&schemeItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("scheme of work item not found")
		}
		return nil, errors.New("failed to verify scheme of work item: " + err.Error())
	}

	// Check if the association already exists
	var existing models.ExamSchemeItem
	if err := s.db.Where("exam_id = ? AND scheme_of_work_item_id = ?", examID, schemeOfWorkItemID).
		First(&existing).Error; err == nil {
		return nil, errors.New("exam scheme item already exists")
	}

	// Create exam scheme item
	examSchemeItem := &models.ExamSchemeItem{
		ExamID:             examID,
		SchemeOfWorkItemID: schemeOfWorkItemID,
	}

	if err := s.db.Create(examSchemeItem).Error; err != nil {
		return nil, errors.New("failed to create exam scheme item: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("Exam").Preload("SchemeOfWorkItem").First(examSchemeItem,
		models.ExamSchemeItem{ExamID: examID, SchemeOfWorkItemID: schemeOfWorkItemID}).Error; err != nil {
		return nil, errors.New("failed to load exam scheme item details: " + err.Error())
	}

	return s.toExamSchemeItemResponse(examSchemeItem), nil
}

// BulkCreateExamSchemeItems creates multiple exam scheme items
func (s *ExamSchemeItemService) BulkCreateExamSchemeItems(req *dto.BulkCreateExamSchemeItemsRequest) (*dto.BulkExamSchemeItemResult, error) {
	// Parse UUIDs
	examID, err := uuid.Parse(req.ExamID)
	if err != nil {
		return nil, errors.New("invalid exam ID format")
	}

	// Check if exam exists
	var exam models.Exam
	if err := s.db.Where("id = ? AND deleted_at IS NULL", examID).First(&exam).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("exam not found")
		}
		return nil, errors.New("failed to verify exam: " + err.Error())
	}

	result := &dto.BulkExamSchemeItemResult{
		Created: []dto.ExamSchemeItemResponse{},
		Errors:  []dto.BulkExamSchemeItemError{},
	}

	// Get existing associations to avoid duplicates
	var existingItems []models.ExamSchemeItem
	if err := s.db.Where("exam_id = ?", examID).Find(&existingItems).Error; err != nil {
		return nil, fmt.Errorf("failed to check existing items: %w", err)
	}

	existingMap := make(map[string]bool)
	for _, item := range existingItems {
		existingMap[item.SchemeOfWorkItemID.String()] = true
	}

	for _, schemeItemIDStr := range req.SchemeOfWorkItemIDs {
		// Check if already exists
		if existingMap[schemeItemIDStr] {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkExamSchemeItemError{
				SchemeOfWorkItemID: schemeItemIDStr,
				Error:              "exam scheme item already exists",
			})
			continue
		}

		// Parse scheme of work item ID
		schemeItemID, err := uuid.Parse(schemeItemIDStr)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkExamSchemeItemError{
				SchemeOfWorkItemID: schemeItemIDStr,
				Error:              "invalid scheme of work item ID format",
			})
			continue
		}

		// Check if scheme of work item exists
		var schemeItem models.SchemeOfWorkItem
		if err := s.db.Where("id = ? AND deleted_at IS NULL", schemeItemID).First(&schemeItem).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkExamSchemeItemError{
				SchemeOfWorkItemID: schemeItemIDStr,
				Error:              "scheme of work item not found",
			})
			continue
		}

		// Create exam scheme item
		examSchemeItem := &models.ExamSchemeItem{
			ExamID:             examID,
			SchemeOfWorkItemID: schemeItemID,
		}

		if err := s.db.Create(examSchemeItem).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkExamSchemeItemError{
				SchemeOfWorkItemID: schemeItemIDStr,
				Error:              "failed to create exam scheme item: " + err.Error(),
			})
			continue
		}

		// Preload relationships
		if err := s.db.Preload("Exam").Preload("SchemeOfWorkItem").First(examSchemeItem,
			models.ExamSchemeItem{ExamID: examID, SchemeOfWorkItemID: schemeItemID}).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkExamSchemeItemError{
				SchemeOfWorkItemID: schemeItemIDStr,
				Error:              "failed to load exam scheme item details",
			})
			continue
		}

		existingMap[schemeItemIDStr] = true
		result.SuccessCount++
		result.Created = append(result.Created, *s.toExamSchemeItemResponse(examSchemeItem))
	}

	return result, nil
}

// GetAllExamSchemeItems retrieves all exam scheme items with pagination and filters
func (s *ExamSchemeItemService) GetAllExamSchemeItems(params *dto.ExamSchemeItemQueryParams) (*dto.ExamSchemeItemListResponse, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}

	// Build query
	query := s.db.Model(&models.ExamSchemeItem{})

	// Apply filters
	if params.ExamID != "" {
		examID, err := uuid.Parse(params.ExamID)
		if err == nil {
			query = query.Where("exam_id = ?", examID)
		}
	}

	if params.SchemeOfWorkItemID != "" {
		schemeItemID, err := uuid.Parse(params.SchemeOfWorkItemID)
		if err == nil {
			query = query.Where("scheme_of_work_item_id = ?", schemeItemID)
		}
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count exam scheme items: %w", err)
	}

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	// Execute with preloads
	var items []models.ExamSchemeItem
	if err := query.Preload("Exam").Preload("SchemeOfWorkItem").Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch exam scheme items: %w", err)
	}

	// Convert to response
	responses := make([]dto.ExamSchemeItemResponse, len(items))
	for i, item := range items {
		responses[i] = *s.toExamSchemeItemResponse(&item)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.ExamSchemeItemListResponse{
		Items:      responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetExamSchemeItemsByExam retrieves all scheme items for a specific exam
func (s *ExamSchemeItemService) GetExamSchemeItemsByExam(examID string) ([]dto.ExamSchemeItemResponse, error) {
	eID, err := uuid.Parse(examID)
	if err != nil {
		return nil, errors.New("invalid exam ID")
	}

	var items []models.ExamSchemeItem
	if err := s.db.Where("exam_id = ?", eID).
		Preload("Exam").
		Preload("SchemeOfWorkItem").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch exam scheme items: %w", err)
	}

	responses := make([]dto.ExamSchemeItemResponse, len(items))
	for i, item := range items {
		responses[i] = *s.toExamSchemeItemResponse(&item)
	}

	return responses, nil
}

// GetExamSchemeItemsBySchemeItem retrieves all exams for a specific scheme of work item
func (s *ExamSchemeItemService) GetExamSchemeItemsBySchemeItem(schemeItemID string) ([]dto.ExamSchemeItemResponse, error) {
	sID, err := uuid.Parse(schemeItemID)
	if err != nil {
		return nil, errors.New("invalid scheme of work item ID")
	}

	var items []models.ExamSchemeItem
	if err := s.db.Where("scheme_of_work_item_id = ?", sID).
		Preload("Exam").
		Preload("SchemeOfWorkItem").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch exam scheme items: %w", err)
	}

	responses := make([]dto.ExamSchemeItemResponse, len(items))
	for i, item := range items {
		responses[i] = *s.toExamSchemeItemResponse(&item)
	}

	return responses, nil
}

// DeleteExamSchemeItem deletes an exam scheme item
func (s *ExamSchemeItemService) DeleteExamSchemeItem(examID, schemeItemID string) error {
	eID, err := uuid.Parse(examID)
	if err != nil {
		return errors.New("invalid exam ID")
	}

	sID, err := uuid.Parse(schemeItemID)
	if err != nil {
		return errors.New("invalid scheme of work item ID")
	}

	// Check if the association exists
	var item models.ExamSchemeItem
	if err := s.db.Where("exam_id = ? AND scheme_of_work_item_id = ?", eID, sID).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("exam scheme item not found")
		}
		return errors.New("failed to fetch exam scheme item: " + err.Error())
	}

	// Delete the association
	if err := s.db.Delete(&item).Error; err != nil {
		return errors.New("failed to delete exam scheme item: " + err.Error())
	}

	return nil
}

// DeleteAllExamSchemeItemsByExam deletes all scheme items for an exam
func (s *ExamSchemeItemService) DeleteAllExamSchemeItemsByExam(examID string) error {
	eID, err := uuid.Parse(examID)
	if err != nil {
		return errors.New("invalid exam ID")
	}

	if err := s.db.Where("exam_id = ?", eID).Delete(&models.ExamSchemeItem{}).Error; err != nil {
		return errors.New("failed to delete exam scheme items: " + err.Error())
	}

	return nil
}

// validateExamSchemeItemRequest validates the exam scheme item request
func (s *ExamSchemeItemService) validateExamSchemeItemRequest(req *dto.CreateExamSchemeItemRequest) error {
	if req.ExamID == "" {
		return errors.New("exam ID is required")
	}
	if req.SchemeOfWorkItemID == "" {
		return errors.New("scheme of work item ID is required")
	}
	return nil
}

// toExamSchemeItemResponse converts model to response DTO
func (s *ExamSchemeItemService) toExamSchemeItemResponse(item *models.ExamSchemeItem) *dto.ExamSchemeItemResponse {
	response := &dto.ExamSchemeItemResponse{
		ExamID:             item.ExamID.String(),
		SchemeOfWorkItemID: item.SchemeOfWorkItemID.String(),
	}

	// Add exam details if preloaded
	if item.Exam.ID != uuid.Nil {
		response.Exam = &dto.ExamResponse{
			ID:                item.Exam.ID.String(),
			Title:             item.Exam.Title,
			ExamType:          item.Exam.ExamType,
			ExamDate:          item.Exam.ExamDate,
			Duration:          item.Exam.Duration,
			TotalMarks:        item.Exam.TotalMarks,
			Status:            item.Exam.Status,
			AcademicSessionID: item.Exam.AcademicSessionID.String(),
			TermID:            item.Exam.TermID.String(),
			SubjectID:         item.Exam.SubjectID.String(),
			ClassID:           item.Exam.ClassID.String(),
		}
	}

	// Add scheme of work item details if preloaded
	if item.SchemeOfWorkItem.ID != uuid.Nil {
		response.SchemeOfWorkItem = &dto.SchemeOfWorkItemResponse{
			ID:        item.SchemeOfWorkItem.ID.String(),
			Topic:     item.SchemeOfWorkItem.Topic,
			Subtopic:  item.SchemeOfWorkItem.Subtopic,
			WeekStart: item.SchemeOfWorkItem.WeekStart,
			WeekEnd:   item.SchemeOfWorkItem.WeekEnd,
			Sequence:  item.SchemeOfWorkItem.Sequence,
		}
	}

	return response
}