// services/test_scheme_item_service.go
package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"crm-go/models"
	"crm-go/dto"
)

type TestSchemeItemService struct {
	db *gorm.DB
}

func NewTestSchemeItemService(db *gorm.DB) *TestSchemeItemService {
	return &TestSchemeItemService{db: db}
}

// CreateTestSchemeItem creates a new test scheme item
func (s *TestSchemeItemService) CreateTestSchemeItem(req *dto.CreateTestSchemeItemRequest) (*dto.TestSchemeItemResponse, error) {
	// Validate input
	if err := s.validateTestSchemeItemRequest(req); err != nil {
		return nil, err
	}

	// Parse UUIDs
	testID, err := uuid.Parse(req.TestID)
	if err != nil {
		return nil, errors.New("invalid test ID format")
	}

	schemeOfWorkItemID, err := uuid.Parse(req.SchemeOfWorkItemID)
	if err != nil {
		return nil, errors.New("invalid scheme of work item ID format")
	}

	// Check if test exists
	var test models.Test
	if err := s.db.Where("id = ? AND deleted_at IS NULL", testID).First(&test).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("test not found")
		}
		return nil, errors.New("failed to verify test: " + err.Error())
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
	var existing models.TestSchemeItem
	if err := s.db.Where("test_id = ? AND scheme_of_work_item_id = ?", testID, schemeOfWorkItemID).
		First(&existing).Error; err == nil {
		return nil, errors.New("test scheme item already exists")
	}

	// Create test scheme item
	testSchemeItem := &models.TestSchemeItem{
		TestID:             testID,
		SchemeOfWorkItemID: schemeOfWorkItemID,
	}

	if err := s.db.Create(testSchemeItem).Error; err != nil {
		return nil, errors.New("failed to create test scheme item: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("Test").Preload("SchemeOfWorkItem").First(testSchemeItem, 
		models.TestSchemeItem{TestID: testID, SchemeOfWorkItemID: schemeOfWorkItemID}).Error; err != nil {
		return nil, errors.New("failed to load test scheme item details: " + err.Error())
	}

	return s.toTestSchemeItemResponse(testSchemeItem), nil
}

// BulkCreateTestSchemeItems creates multiple test scheme items
func (s *TestSchemeItemService) BulkCreateTestSchemeItems(req *dto.BulkCreateTestSchemeItemsRequest) (*dto.BulkTestSchemeItemResult, error) {
	// Parse UUIDs
	testID, err := uuid.Parse(req.TestID)
	if err != nil {
		return nil, errors.New("invalid test ID format")
	}

	// Check if test exists
	var test models.Test
	if err := s.db.Where("id = ? AND deleted_at IS NULL", testID).First(&test).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("test not found")
		}
		return nil, errors.New("failed to verify test: " + err.Error())
	}

	result := &dto.BulkTestSchemeItemResult{
		Created: []dto.TestSchemeItemResponse{},
		Errors:  []dto.BulkTestSchemeItemError{},
	}

	// Get existing associations to avoid duplicates
	var existingItems []models.TestSchemeItem
	if err := s.db.Where("test_id = ?", testID).Find(&existingItems).Error; err != nil {
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
			result.Errors = append(result.Errors, dto.BulkTestSchemeItemError{
				SchemeOfWorkItemID: schemeItemIDStr,
				Error:              "test scheme item already exists",
			})
			continue
		}

		// Parse scheme of work item ID
		schemeItemID, err := uuid.Parse(schemeItemIDStr)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkTestSchemeItemError{
				SchemeOfWorkItemID: schemeItemIDStr,
				Error:              "invalid scheme of work item ID format",
			})
			continue
		}

		// Check if scheme of work item exists
		var schemeItem models.SchemeOfWorkItem
		if err := s.db.Where("id = ? AND deleted_at IS NULL", schemeItemID).First(&schemeItem).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkTestSchemeItemError{
				SchemeOfWorkItemID: schemeItemIDStr,
				Error:              "scheme of work item not found",
			})
			continue
		}

		// Create test scheme item
		testSchemeItem := &models.TestSchemeItem{
			TestID:             testID,
			SchemeOfWorkItemID: schemeItemID,
		}

		if err := s.db.Create(testSchemeItem).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkTestSchemeItemError{
				SchemeOfWorkItemID: schemeItemIDStr,
				Error:              "failed to create test scheme item: " + err.Error(),
			})
			continue
		}

		// Preload relationships
		if err := s.db.Preload("Test").Preload("SchemeOfWorkItem").First(testSchemeItem, 
			models.TestSchemeItem{TestID: testID, SchemeOfWorkItemID: schemeItemID}).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkTestSchemeItemError{
				SchemeOfWorkItemID: schemeItemIDStr,
				Error:              "failed to load test scheme item details",
			})
			continue
		}

		existingMap[schemeItemIDStr] = true
		result.SuccessCount++
		result.Created = append(result.Created, *s.toTestSchemeItemResponse(testSchemeItem))
	}

	return result, nil
}

// GetAllTestSchemeItems retrieves all test scheme items with pagination and filters
func (s *TestSchemeItemService) GetAllTestSchemeItems(params *dto.TestSchemeItemQueryParams) (*dto.TestSchemeItemListResponse, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}

	// Build query
	query := s.db.Model(&models.TestSchemeItem{})

	// Apply filters
	if params.TestID != "" {
		testID, err := uuid.Parse(params.TestID)
		if err == nil {
			query = query.Where("test_id = ?", testID)
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
		return nil, fmt.Errorf("failed to count test scheme items: %w", err)
	}

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	// Execute with preloads
	var items []models.TestSchemeItem
	if err := query.Preload("Test").Preload("SchemeOfWorkItem").Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch test scheme items: %w", err)
	}

	// Convert to response
	responses := make([]dto.TestSchemeItemResponse, len(items))
	for i, item := range items {
		responses[i] = *s.toTestSchemeItemResponse(&item)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.TestSchemeItemListResponse{
		Items:      responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetTestSchemeItemsByTest retrieves all scheme items for a specific test
func (s *TestSchemeItemService) GetTestSchemeItemsByTest(testID string) ([]dto.TestSchemeItemResponse, error) {
	tID, err := uuid.Parse(testID)
	if err != nil {
		return nil, errors.New("invalid test ID")
	}

	var items []models.TestSchemeItem
	if err := s.db.Where("test_id = ?", tID).
		Preload("Test").
		Preload("SchemeOfWorkItem").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch test scheme items: %w", err)
	}

	responses := make([]dto.TestSchemeItemResponse, len(items))
	for i, item := range items {
		responses[i] = *s.toTestSchemeItemResponse(&item)
	}

	return responses, nil
}

// GetTestSchemeItemsBySchemeItem retrieves all tests for a specific scheme of work item
func (s *TestSchemeItemService) GetTestSchemeItemsBySchemeItem(schemeItemID string) ([]dto.TestSchemeItemResponse, error) {
	sID, err := uuid.Parse(schemeItemID)
	if err != nil {
		return nil, errors.New("invalid scheme of work item ID")
	}

	var items []models.TestSchemeItem
	if err := s.db.Where("scheme_of_work_item_id = ?", sID).
		Preload("Test").
		Preload("SchemeOfWorkItem").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch test scheme items: %w", err)
	}

	responses := make([]dto.TestSchemeItemResponse, len(items))
	for i, item := range items {
		responses[i] = *s.toTestSchemeItemResponse(&item)
	}

	return responses, nil
}

// DeleteTestSchemeItem deletes a test scheme item
func (s *TestSchemeItemService) DeleteTestSchemeItem(testID, schemeItemID string) error {
	tID, err := uuid.Parse(testID)
	if err != nil {
		return errors.New("invalid test ID")
	}

	sID, err := uuid.Parse(schemeItemID)
	if err != nil {
		return errors.New("invalid scheme of work item ID")
	}

	// Check if the association exists
	var item models.TestSchemeItem
	if err := s.db.Where("test_id = ? AND scheme_of_work_item_id = ?", tID, sID).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("test scheme item not found")
		}
		return errors.New("failed to fetch test scheme item: " + err.Error())
	}

	// Delete the association
	if err := s.db.Delete(&item).Error; err != nil {
		return errors.New("failed to delete test scheme item: " + err.Error())
	}

	return nil
}

// DeleteAllTestSchemeItemsByTest deletes all scheme items for a test
func (s *TestSchemeItemService) DeleteAllTestSchemeItemsByTest(testID string) error {
	tID, err := uuid.Parse(testID)
	if err != nil {
		return errors.New("invalid test ID")
	}

	if err := s.db.Where("test_id = ?", tID).Delete(&models.TestSchemeItem{}).Error; err != nil {
		return errors.New("failed to delete test scheme items: " + err.Error())
	}

	return nil
}

// validateTestSchemeItemRequest validates the test scheme item request
func (s *TestSchemeItemService) validateTestSchemeItemRequest(req *dto.CreateTestSchemeItemRequest) error {
	if req.TestID == "" {
		return errors.New("test ID is required")
	}
	if req.SchemeOfWorkItemID == "" {
		return errors.New("scheme of work item ID is required")
	}
	return nil
}

// toTestSchemeItemResponse converts model to response DTO
func (s *TestSchemeItemService) toTestSchemeItemResponse(item *models.TestSchemeItem) *dto.TestSchemeItemResponse {
	response := &dto.TestSchemeItemResponse{
		TestID:             item.TestID.String(),
		SchemeOfWorkItemID: item.SchemeOfWorkItemID.String(),
	}

	// Add test details if preloaded
	if item.Test.ID != uuid.Nil {
		response.Test = &dto.TestResponse{
			ID:                item.Test.ID.String(),
			Title:             item.Test.Title,
			TestType:          item.Test.TestType,
			TestDate:          item.Test.TestDate,
			Duration:          item.Test.Duration,
			TotalMarks:        item.Test.TotalMarks,
			Status:            item.Test.Status,
			AcademicSessionID: item.Test.AcademicSessionID.String(),
			TermID:            item.Test.TermID.String(),
			SubjectID:         item.Test.SubjectID.String(),
			ClassID:           item.Test.ClassID.String(),
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