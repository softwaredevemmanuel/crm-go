// services/topic_service.go
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

type TopicService struct {
	db *gorm.DB
}

func NewTopicService(db *gorm.DB) *TopicService {
	return &TopicService{db: db}
}

// CreateTopic creates a new topic
func (s *TopicService) CreateTopic(req *dto.CreateTopicRequest) (*dto.TopicResponse, error) {
	// Validate input
	if err := s.validateTopicRequest(req); err != nil {
		return nil, err
	}

	// Parse UUIDs
	moduleID, err := uuid.Parse(req.ModuleID)
	if err != nil {
		return nil, errors.New("invalid module ID format")
	}

	// Check if module exists
	var module models.Module
	if err := s.db.Where("id = ? AND deleted_at IS NULL", moduleID).First(&module).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("module not found")
		}
		return nil, errors.New("failed to verify module: " + err.Error())
	}

	// Check if topic title already exists for this module
	var existing models.Topic
	if err := s.db.Where("module_id = ? AND title = ? AND deleted_at IS NULL",
		moduleID, req.Title).First(&existing).Error; err == nil {
		return nil, errors.New("topic title already exists for this module")
	}

	// Set default topic order if not provided
	topicOrder := req.TopicOrder
	if topicOrder == 0 {
		// Get the highest topic order for this module
		var maxOrder int
		if err := s.db.Model(&models.Topic{}).
			Where("module_id = ? AND deleted_at IS NULL", moduleID).
			Select("COALESCE(MAX(topic_order), 0)").
			Scan(&maxOrder).Error; err != nil {
			return nil, errors.New("failed to determine topic order: " + err.Error())
		}
		topicOrder = maxOrder + 1
	}

	// Create topic
	topic := &models.Topic{
		ID:          uuid.New(),
		ModuleID:    moduleID,
		Title:       req.Title,
		Description: req.Description,
		TopicOrder:  topicOrder,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(topic).Error; err != nil {
		return nil, errors.New("failed to create topic: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("Module").First(topic, topic.ID).Error; err != nil {
		return nil, errors.New("failed to load topic details: " + err.Error())
	}

	return s.toTopicResponse(topic), nil
}

// BulkCreateTopics creates multiple topics
func (s *TopicService) BulkCreateTopics(req *dto.BulkCreateTopicsRequest) (*dto.BulkTopicResult, error) {
	// Parse UUIDs
	moduleID, err := uuid.Parse(req.ModuleID)
	if err != nil {
		return nil, errors.New("invalid module ID format")
	}

	// Check if module exists
	var module models.Module
	if err := s.db.Where("id = ? AND deleted_at IS NULL", moduleID).First(&module).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("module not found")
		}
		return nil, errors.New("failed to verify module: " + err.Error())
	}

	result := &dto.BulkTopicResult{
		Created: []dto.TopicResponse{},
		Errors:  []dto.BulkTopicError{},
	}

	// Get existing topic titles to avoid duplicates
	var existingTopics []models.Topic
	if err := s.db.Where("module_id = ? AND deleted_at IS NULL", moduleID).Find(&existingTopics).Error; err != nil {
		return nil, fmt.Errorf("failed to check existing topics: %w", err)
	}

	existingTitleMap := make(map[string]bool)
	for _, t := range existingTopics {
		existingTitleMap[t.Title] = true
	}

	// Get max topic order
	var maxOrder int
	if err := s.db.Model(&models.Topic{}).
		Where("module_id = ? AND deleted_at IS NULL", moduleID).
		Select("COALESCE(MAX(topic_order), 0)").
		Scan(&maxOrder).Error; err != nil {
		return nil, errors.New("failed to determine topic order: " + err.Error())
	}
	currentOrder := maxOrder

	for _, topicReq := range req.Topics {
		// Validate
		if topicReq.Title == "" {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkTopicError{
				Title: topicReq.Title,
				Error: "title is required",
			})
			continue
		}

		// Check for duplicate title
		if existingTitleMap[topicReq.Title] {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkTopicError{
				Title: topicReq.Title,
				Error: "topic title already exists for this module",
			})
			continue
		}

		// Set topic order
		topicOrder := topicReq.TopicOrder
		if topicOrder == 0 {
			currentOrder++
			topicOrder = currentOrder
		}

		// Create topic
		topic := &models.Topic{
			ID:          uuid.New(),
			ModuleID:    moduleID,
			Title:       topicReq.Title,
			Description: topicReq.Description,
			TopicOrder:  topicOrder,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := s.db.Create(topic).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkTopicError{
				Title: topicReq.Title,
				Error: "failed to create topic: " + err.Error(),
			})
			continue
		}

		// Preload relationships
		if err := s.db.Preload("Module").First(topic, topic.ID).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkTopicError{
				Title: topicReq.Title,
				Error: "failed to load topic details",
			})
			continue
		}

		existingTitleMap[topicReq.Title] = true
		result.SuccessCount++
		result.Created = append(result.Created, *s.toTopicResponse(topic))
	}

	return result, nil
}

// GetAllTopics retrieves all topics with pagination and filters
func (s *TopicService) GetAllTopics(params *dto.TopicQueryParams) (*dto.TopicListResponse, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = "topic_order"
	}
	if params.SortOrder == "" {
		params.SortOrder = "asc"
	}

	// Build query
	query := s.db.Model(&models.Topic{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.ModuleID != "" {
		moduleID, err := uuid.Parse(params.ModuleID)
		if err == nil {
			query = query.Where("module_id = ?", moduleID)
		}
	}

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", search, search)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count topics: %w", err)
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
	var topics []models.Topic
	if err := query.Preload("Module").Find(&topics).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch topics: %w", err)
	}

	// Convert to response
	responses := make([]dto.TopicResponse, len(topics))
	for i, topic := range topics {
		responses[i] = *s.toTopicResponse(&topic)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.TopicListResponse{
		Topics:     responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetTopicByID retrieves a single topic by ID
func (s *TopicService) GetTopicByID(id string) (*dto.TopicResponse, error) {
	topicID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid topic ID")
	}

	var topic models.Topic
	if err := s.db.Where("id = ? AND deleted_at IS NULL", topicID).
		Preload("Module").
		Preload("Lessons").
		First(&topic).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("topic not found")
		}
		return nil, errors.New("failed to fetch topic: " + err.Error())
	}

	return s.toTopicResponse(&topic), nil
}

// GetTopicsByModule retrieves all topics for a specific module
func (s *TopicService) GetTopicsByModule(moduleID string) ([]dto.TopicResponse, error) {
	mID, err := uuid.Parse(moduleID)
	if err != nil {
		return nil, errors.New("invalid module ID")
	}

	var topics []models.Topic
	if err := s.db.Where("module_id = ? AND deleted_at IS NULL", mID).
		Preload("Module").
		Order("topic_order ASC").
		Find(&topics).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch topics: %w", err)
	}

	responses := make([]dto.TopicResponse, len(topics))
	for i, topic := range topics {
		responses[i] = *s.toTopicResponse(&topic)
	}

	return responses, nil
}

// UpdateTopic updates an existing topic
func (s *TopicService) UpdateTopic(id string, req *dto.UpdateTopicRequest) (*dto.TopicResponse, error) {
	topicID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid topic ID")
	}

	// Find existing topic
	var topic models.Topic
	if err := s.db.Where("id = ? AND deleted_at IS NULL", topicID).First(&topic).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("topic not found")
		}
		return nil, errors.New("failed to fetch topic: " + err.Error())
	}

	// Update fields
	if req.Title != "" {
		// Check if title already exists for another topic in the same module
		var existing models.Topic
		if err := s.db.Where("module_id = ? AND title = ? AND id != ? AND deleted_at IS NULL",
			topic.ModuleID, req.Title, topicID).First(&existing).Error; err == nil {
			return nil, errors.New("topic title already exists for this module")
		}
		topic.Title = req.Title
	}

	if req.Description != "" {
		topic.Description = req.Description
	}

	if req.TopicOrder > 0 && req.TopicOrder != topic.TopicOrder {
		// Handle reordering
		if err := s.reorderTopics(topic.ModuleID, topic.TopicOrder, req.TopicOrder, topicID); err != nil {
			return nil, err
		}
		topic.TopicOrder = req.TopicOrder
	}

	topic.UpdatedAt = time.Now()

	if err := s.db.Save(&topic).Error; err != nil {
		return nil, errors.New("failed to update topic: " + err.Error())
	}

	// Preload relationships
	if err := s.db.Preload("Module").First(&topic, topic.ID).Error; err != nil {
		return nil, errors.New("failed to load topic details: " + err.Error())
	}

	return s.toTopicResponse(&topic), nil
}

// ReorderTopics reorders multiple topics
func (s *TopicService) ReorderTopics(req *dto.ReorderTopicsRequest) error {
	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, orderReq := range req.TopicOrders {
		topicID, err := uuid.Parse(orderReq.ID)
		if err != nil {
			tx.Rollback()
			return errors.New("invalid topic ID format: " + err.Error())
		}

		if err := tx.Model(&models.Topic{}).
			Where("id = ? AND deleted_at IS NULL", topicID).
			Update("topic_order", orderReq.TopicOrder).Error; err != nil {
			tx.Rollback()
			return errors.New("failed to update topic order: " + err.Error())
		}
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("failed to commit reorder: " + err.Error())
	}

	return nil
}

// reorderTopics reorders topics when a single topic's order changes
func (s *TopicService) reorderTopics(moduleID uuid.UUID, oldOrder, newOrder int, excludeID uuid.UUID) error {
	if oldOrder == newOrder {
		return nil
	}

	var topics []models.Topic
	if err := s.db.Where("module_id = ? AND deleted_at IS NULL AND id != ?", moduleID, excludeID).
		Order("topic_order ASC").
		Find(&topics).Error; err != nil {
		return errors.New("failed to fetch topics for reordering: " + err.Error())
	}

	// Update orders
	if oldOrder < newOrder {
		// Moving down: decrease orders between oldOrder+1 and newOrder
		for _, t := range topics {
			if t.TopicOrder > oldOrder && t.TopicOrder <= newOrder {
				t.TopicOrder--
				if err := s.db.Save(&t).Error; err != nil {
					return errors.New("failed to update topic order: " + err.Error())
				}
			}
		}
	} else {
		// Moving up: increase orders between newOrder and oldOrder-1
		for _, t := range topics {
			if t.TopicOrder >= newOrder && t.TopicOrder < oldOrder {
				t.TopicOrder++
				if err := s.db.Save(&t).Error; err != nil {
					return errors.New("failed to update topic order: " + err.Error())
				}
			}
		}
	}

	return nil
}

// DeleteTopic soft deletes a topic
func (s *TopicService) DeleteTopic(id string) error {
	topicID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid topic ID")
	}

	var topic models.Topic
	if err := s.db.Where("id = ? AND deleted_at IS NULL", topicID).First(&topic).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("topic not found")
		}
		return errors.New("failed to fetch topic: " + err.Error())
	}

	if err := s.db.Delete(&topic).Error; err != nil {
		return errors.New("failed to delete topic: " + err.Error())
	}

	return nil
}

// validateTopicRequest validates the topic request
func (s *TopicService) validateTopicRequest(req *dto.CreateTopicRequest) error {
	if req.ModuleID == "" {
		return errors.New("module ID is required")
	}
	if req.Title == "" {
		return errors.New("title is required")
	}
	if req.TopicOrder < 0 {
		return errors.New("topic order cannot be negative")
	}
	return nil
}

// toTopicResponse converts model to response DTO
func (s *TopicService) toTopicResponse(topic *models.Topic) *dto.TopicResponse {
	response := &dto.TopicResponse{
		ID:          topic.ID.String(),
		ModuleID:    topic.ModuleID.String(),
		Title:       topic.Title,
		Description: topic.Description,
		TopicOrder:  topic.TopicOrder,
		CreatedAt:   topic.CreatedAt,
		UpdatedAt:   topic.UpdatedAt,
	}

	// Add module details if preloaded
	if topic.Module.ID != uuid.Nil {
		response.Module = &dto.ModuleResponse{
			ID:          topic.Module.ID.String(),
			Title:       topic.Module.Title,
			Description: topic.Module.Description,
			ModuleOrder: topic.Module.ModuleOrder,
		}
	}

	// Add lessons if preloaded
	if len(topic.Lessons) > 0 {
		lessons := make([]dto.LessonResponse, len(topic.Lessons))
		for i, lesson := range topic.Lessons {
			lessons[i] = dto.LessonResponse{
				ID:     lesson.ID.String(),
				Title:  lesson.Title,
				Week:   lesson.Week,
				Status: lesson.Status,
			}
		}
		response.Lessons = lessons
	}

	return response
}