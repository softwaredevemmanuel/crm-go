// services/lesson_service.go
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

type LessonService struct {
	db *gorm.DB
}

func NewLessonService(db *gorm.DB) *LessonService {
	return &LessonService{db: db}
}

// CreateLesson creates a new lesson
func (s *LessonService) CreateLesson(req *dto.CreateLessonRequest, userID uuid.UUID) (*dto.LessonResponse, error) {
	// Validate input
	if err := s.validateLessonRequest(req); err != nil {
		return nil, err
	}

	// Parse UUIDs
	schemeOfWorkID, err := uuid.Parse(req.SchemeOfWorkID)
	if err != nil {
		return nil, errors.New("invalid scheme of work ID format")
	}

	moduleID, err := uuid.Parse(req.ModuleID)
	if err != nil {
		return nil, errors.New("invalid module ID format")
	}

	topicID, err := uuid.Parse(req.TopicID)
	if err != nil {
		return nil, errors.New("invalid topic ID format")
	}

	// Verify entities exist
	if err := s.verifyEntities(schemeOfWorkID, moduleID, topicID); err != nil {
		return nil, err
	}

	// Parse lesson date if provided
	var lessonDate *time.Time
	if req.LessonDate != "" {
		date, err := time.Parse("2006-01-02", req.LessonDate)
		if err != nil {
			return nil, errors.New("invalid lesson date format. Use YYYY-MM-DD")
		}
		lessonDate = &date
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "draft"
	}

	// Set default lesson order if not provided
	lessonOrder := req.LessonOrder
	if lessonOrder == 0 {
		var maxOrder int
		if err := s.db.Model(&models.Lesson{}).
			Where("topic_id = ? AND deleted_at IS NULL", topicID).
			Select("COALESCE(MAX(lesson_order), 0)").
			Scan(&maxOrder).Error; err != nil {
			return nil, errors.New("failed to determine lesson order: " + err.Error())
		}
		lessonOrder = maxOrder + 1
	}



	// Create lesson
	lesson := &models.Lesson{
		ID:             uuid.New(),
		SchemeOfWorkID: schemeOfWorkID,
		ModuleID:       moduleID,
		TopicID:        topicID,
		LessonOrder:    lessonOrder,
		Title:          req.Title,
		Description:    req.Description,
		LessonDate:     lessonDate,
		Week:           req.Week,
		Duration:       req.Duration,
		Content:        req.Content,
		Objectives:     req.Objectives,
		Activities:     req.Activities,
		Resources:      req.Resources,
		Assessment:     req.Assessment,
		Status:         status,
		CreatedBy:      userID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.db.Create(lesson).Error; err != nil {
		return nil, errors.New("failed to create lesson: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("SchemeOfWork").Preload("Module").Preload("Topic").Preload("Creator").First(lesson, lesson.ID).Error; err != nil {
		return nil, errors.New("failed to load lesson details: " + err.Error())
	}

	return s.toLessonResponse(lesson), nil
}

// BulkCreateLessons creates multiple lessons
func (s *LessonService) BulkCreateLessons(req *dto.BulkCreateLessonsRequest, userID uuid.UUID) (*dto.BulkLessonResult, error) {
	// Parse UUIDs
	schemeOfWorkID, err := uuid.Parse(req.SchemeOfWorkID)
	if err != nil {
		return nil, errors.New("invalid scheme of work ID format")
	}

	moduleID, err := uuid.Parse(req.ModuleID)
	if err != nil {
		return nil, errors.New("invalid module ID format")
	}

	topicID, err := uuid.Parse(req.TopicID)
	if err != nil {
		return nil, errors.New("invalid topic ID format")
	}

	// Verify entities exist
	if err := s.verifyEntities(schemeOfWorkID, moduleID, topicID); err != nil {
		return nil, err
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "draft"
	}

	result := &dto.BulkLessonResult{
		Created: []dto.LessonResponse{},
		Errors:  []dto.BulkLessonError{},
	}

	// Get max lesson order for the topic
	var maxOrder int
	if err := s.db.Model(&models.Lesson{}).
		Where("topic_id = ? AND deleted_at IS NULL", topicID).
		Select("COALESCE(MAX(lesson_order), 0)").
		Scan(&maxOrder).Error; err != nil {
		return nil, errors.New("failed to determine lesson order: " + err.Error())
	}
	currentOrder := maxOrder

	for _, lessonReq := range req.Lessons {
		// Validate
		if lessonReq.Title == "" {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkLessonError{
				Title: lessonReq.Title,
				Error: "title is required",
			})
			continue
		}

		// Parse lesson date if provided
		var lessonDate *time.Time
		if lessonReq.LessonDate != "" {
			date, err := time.Parse("2006-01-02", lessonReq.LessonDate)
			if err != nil {
				result.FailedCount++
				result.Errors = append(result.Errors, dto.BulkLessonError{
					Title: lessonReq.Title,
					Error: "invalid lesson date format. Use YYYY-MM-DD",
				})
				continue
			}
			lessonDate = &date
		}

		// Set lesson order
		lessonOrder := lessonReq.LessonOrder
		if lessonOrder == 0 {
			currentOrder++
			lessonOrder = currentOrder
		}


		// Create lesson
		lesson := &models.Lesson{
			ID:             uuid.New(),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:       moduleID,
			TopicID:        topicID,
			LessonOrder:    lessonOrder,
			Title:          lessonReq.Title,
			Description:    lessonReq.Description,
			LessonDate:     lessonDate,
			Week:           lessonReq.Week,
			Duration:       lessonReq.Duration,
			Content:        lessonReq.Content,
			Objectives:     lessonReq.Objectives,
			Activities:     lessonReq.Activities,
			Resources:      lessonReq.Resources,
			Assessment:     lessonReq.Assessment,
			Status:         status,
			CreatedBy:      userID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := s.db.Create(lesson).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkLessonError{
				Title: lessonReq.Title,
				Error: "failed to create lesson: " + err.Error(),
			})
			continue
		}

		// Preload relationships
		if err := s.db.Preload("SchemeOfWork").Preload("Module").Preload("Topic").Preload("Creator").First(lesson, lesson.ID).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkLessonError{
				Title: lessonReq.Title,
				Error: "failed to load lesson details",
			})
			continue
		}

		result.SuccessCount++
		result.Created = append(result.Created, *s.toLessonResponse(lesson))
	}

	return result, nil
}

// GetAllLessons retrieves all lessons with pagination and filters
func (s *LessonService) GetAllLessons(params *dto.LessonQueryParams) (*dto.LessonListResponse, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = "lesson_order"
	}
	if params.SortOrder == "" {
		params.SortOrder = "asc"
	}

	// Build query
	query := s.db.Model(&models.Lesson{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.SchemeOfWorkID != "" {
		schemeID, err := uuid.Parse(params.SchemeOfWorkID)
		if err == nil {
			query = query.Where("scheme_of_work_id = ?", schemeID)
		}
	}

	if params.ModuleID != "" {
		moduleID, err := uuid.Parse(params.ModuleID)
		if err == nil {
			query = query.Where("module_id = ?", moduleID)
		}
	}

	if params.TopicID != "" {
		topicID, err := uuid.Parse(params.TopicID)
		if err == nil {
			query = query.Where("topic_id = ?", topicID)
		}
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.Week > 0 {
		query = query.Where("week = ?", params.Week)
	}

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ? OR LOWER(objectives) LIKE ?",
			search, search, search)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count lessons: %w", err)
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
	var lessons []models.Lesson
	if err := query.Preload("SchemeOfWork").Preload("Module").Preload("Topic").Preload("Creator").Find(&lessons).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch lessons: %w", err)
	}

	// Convert to response
	responses := make([]dto.LessonResponse, len(lessons))
	for i, lesson := range lessons {
		responses[i] = *s.toLessonResponse(&lesson)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.LessonListResponse{
		Lessons:    responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetLessonByID retrieves a single lesson by ID
func (s *LessonService) GetLessonByID(id string) (*dto.LessonResponse, error) {
	lessonID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid lesson ID")
	}

	var lesson models.Lesson
	if err := s.db.Where("id = ? AND deleted_at IS NULL", lessonID).
		Preload("SchemeOfWork").
		Preload("Module").
		Preload("Topic").
		Preload("Creator").
		First(&lesson).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lesson not found")
		}
		return nil, errors.New("failed to fetch lesson: " + err.Error())
	}

	return s.toLessonResponse(&lesson), nil
}

// GetLessonsBySchemeOfWork retrieves all lessons for a scheme of work
func (s *LessonService) GetLessonsBySchemeOfWork(schemeOfWorkID string) ([]dto.LessonResponse, error) {
	sID, err := uuid.Parse(schemeOfWorkID)
	if err != nil {
		return nil, errors.New("invalid scheme of work ID")
	}

	var lessons []models.Lesson
	if err := s.db.Where("scheme_of_work_id = ? AND deleted_at IS NULL", sID).
		Preload("SchemeOfWork").
		Preload("Module").
		Preload("Topic").
		Preload("Creator").
		Order("lesson_order ASC").
		Find(&lessons).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch lessons: %w", err)
	}

	responses := make([]dto.LessonResponse, len(lessons))
	for i, lesson := range lessons {
		responses[i] = *s.toLessonResponse(&lesson)
	}

	return responses, nil
}

// GetLessonsByModule retrieves all lessons for a module
func (s *LessonService) GetLessonsByModule(moduleID string) ([]dto.LessonResponse, error) {
	mID, err := uuid.Parse(moduleID)
	if err != nil {
		return nil, errors.New("invalid module ID")
	}

	var lessons []models.Lesson
	if err := s.db.Where("module_id = ? AND deleted_at IS NULL", mID).
		Preload("SchemeOfWork").
		Preload("Module").
		Preload("Topic").
		Preload("Creator").
		Order("lesson_order ASC").
		Find(&lessons).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch lessons: %w", err)
	}

	responses := make([]dto.LessonResponse, len(lessons))
	for i, lesson := range lessons {
		responses[i] = *s.toLessonResponse(&lesson)
	}

	return responses, nil
}

// GetLessonsByTopic retrieves all lessons for a topic
func (s *LessonService) GetLessonsByTopic(topicID string) ([]dto.LessonResponse, error) {
	tID, err := uuid.Parse(topicID)
	if err != nil {
		return nil, errors.New("invalid topic ID")
	}

	var lessons []models.Lesson
	if err := s.db.Where("topic_id = ? AND deleted_at IS NULL", tID).
		Preload("SchemeOfWork").
		Preload("Module").
		Preload("Topic").
		Preload("Creator").
		Order("lesson_order ASC").
		Find(&lessons).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch lessons: %w", err)
	}

	responses := make([]dto.LessonResponse, len(lessons))
	for i, lesson := range lessons {
		responses[i] = *s.toLessonResponse(&lesson)
	}

	return responses, nil
}

// UpdateLesson updates an existing lesson
func (s *LessonService) UpdateLesson(id string, req *dto.UpdateLessonRequest) (*dto.LessonResponse, error) {
	lessonID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid lesson ID")
	}

	// Find existing lesson
	var lesson models.Lesson
	if err := s.db.Where("id = ? AND deleted_at IS NULL", lessonID).First(&lesson).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lesson not found")
		}
		return nil, errors.New("failed to fetch lesson: " + err.Error())
	}

	// Update fields
	if req.SchemeOfWorkID != "" {
		schemeID, err := uuid.Parse(req.SchemeOfWorkID)
		if err != nil {
			return nil, errors.New("invalid scheme of work ID format")
		}
		if err := s.verifySchemeOfWork(schemeID); err != nil {
			return nil, err
		}
		lesson.SchemeOfWorkID = schemeID
	}

	if req.ModuleID != "" {
		moduleID, err := uuid.Parse(req.ModuleID)
		if err != nil {
			return nil, errors.New("invalid module ID format")
		}
		if err := s.verifyModule(moduleID); err != nil {
			return nil, err
		}
		lesson.ModuleID = moduleID
	}

	if req.TopicID != "" {
		topicID, err := uuid.Parse(req.TopicID)
		if err != nil {
			return nil, errors.New("invalid topic ID format")
		}
		if err := s.verifyTopic(topicID); err != nil {
			return nil, err
		}
		lesson.TopicID = topicID
	}

	if req.Title != "" {
		lesson.Title = req.Title
	}

	if req.Description != "" {
		lesson.Description = req.Description
	}

	if req.LessonDate != "" {
		date, err := time.Parse("2006-01-02", req.LessonDate)
		if err != nil {
			return nil, errors.New("invalid lesson date format. Use YYYY-MM-DD")
		}
		lesson.LessonDate = &date
	}

	if req.Week > 0 {
		lesson.Week = req.Week
	}

	if req.Duration > 0 {
		lesson.Duration = req.Duration
	}


	if req.Objectives != "" {
		lesson.Objectives = req.Objectives
	}

	if req.Activities != "" {
		lesson.Activities = req.Activities
	}

	if req.Resources != "" {
		lesson.Resources = req.Resources
	}

	if req.Assessment != "" {
		lesson.Assessment = req.Assessment
	}

	if req.Status != "" {
		validStatuses := map[string]bool{"draft": true, "published": true, "in review": true, "archived": true}
		if !validStatuses[req.Status] {
			return nil, errors.New("status must be 'draft', 'published', 'in review', or 'archived'")
		}
		lesson.Status = req.Status
	}

	if req.LessonOrder > 0 && req.LessonOrder != lesson.LessonOrder {
		// Handle reordering
		if err := s.reorderLessons(lesson.TopicID, lesson.LessonOrder, req.LessonOrder, lessonID); err != nil {
			return nil, err
		}
		lesson.LessonOrder = req.LessonOrder
	}

	lesson.UpdatedAt = time.Now()

	if err := s.db.Save(&lesson).Error; err != nil {
		return nil, errors.New("failed to update lesson: " + err.Error())
	}

	// Preload relationships
	if err := s.db.Preload("SchemeOfWork").Preload("Module").Preload("Topic").Preload("Creator").First(&lesson, lesson.ID).Error; err != nil {
		return nil, errors.New("failed to load lesson details: " + err.Error())
	}

	return s.toLessonResponse(&lesson), nil
}

// ReorderLessons reorders multiple lessons
func (s *LessonService) ReorderLessons(req *dto.ReorderLessonsRequest) error {
	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, orderReq := range req.LessonOrders {
		lessonID, err := uuid.Parse(orderReq.ID)
		if err != nil {
			tx.Rollback()
			return errors.New("invalid lesson ID format: " + err.Error())
		}

		if err := tx.Model(&models.Lesson{}).
			Where("id = ? AND deleted_at IS NULL", lessonID).
			Update("lesson_order", orderReq.LessonOrder).Error; err != nil {
			tx.Rollback()
			return errors.New("failed to update lesson order: " + err.Error())
		}
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("failed to commit reorder: " + err.Error())
	}

	return nil
}

// reorderLessons reorders lessons when a single lesson's order changes
func (s *LessonService) reorderLessons(topicID uuid.UUID, oldOrder, newOrder int, excludeID uuid.UUID) error {
	if oldOrder == newOrder {
		return nil
	}

	var lessons []models.Lesson
	if err := s.db.Where("topic_id = ? AND deleted_at IS NULL AND id != ?", topicID, excludeID).
		Order("lesson_order ASC").
		Find(&lessons).Error; err != nil {
		return errors.New("failed to fetch lessons for reordering: " + err.Error())
	}

	if oldOrder < newOrder {
		// Moving down: decrease orders between oldOrder+1 and newOrder
		for _, l := range lessons {
			if l.LessonOrder > oldOrder && l.LessonOrder <= newOrder {
				l.LessonOrder--
				if err := s.db.Save(&l).Error; err != nil {
					return errors.New("failed to update lesson order: " + err.Error())
				}
			}
		}
	} else {
		// Moving up: increase orders between newOrder and oldOrder-1
		for _, l := range lessons {
			if l.LessonOrder >= newOrder && l.LessonOrder < oldOrder {
				l.LessonOrder++
				if err := s.db.Save(&l).Error; err != nil {
					return errors.New("failed to update lesson order: " + err.Error())
				}
			}
		}
	}

	return nil
}

// DeleteLesson soft deletes a lesson
func (s *LessonService) DeleteLesson(id string) error {
	lessonID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid lesson ID")
	}

	var lesson models.Lesson
	if err := s.db.Where("id = ? AND deleted_at IS NULL", lessonID).First(&lesson).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("lesson not found")
		}
		return errors.New("failed to fetch lesson: " + err.Error())
	}

	if err := s.db.Delete(&lesson).Error; err != nil {
		return errors.New("failed to delete lesson: " + err.Error())
	}

	return nil
}

// verifyEntities verifies that all referenced entities exist
func (s *LessonService) verifyEntities(schemeOfWorkID, moduleID, topicID uuid.UUID) error {
	if err := s.verifySchemeOfWork(schemeOfWorkID); err != nil {
		return err
	}
	if err := s.verifyModule(moduleID); err != nil {
		return err
	}
	if err := s.verifyTopic(topicID); err != nil {
		return err
	}
	return nil
}

func (s *LessonService) verifySchemeOfWork(id uuid.UUID) error {
	var scheme models.SchemeOfWork
	if err := s.db.Where("id = ? AND deleted_at IS NULL", id).First(&scheme).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("scheme of work not found")
		}
		return errors.New("failed to verify scheme of work: " + err.Error())
	}
	return nil
}

func (s *LessonService) verifyModule(id uuid.UUID) error {
	var module models.Module
	if err := s.db.Where("id = ? AND deleted_at IS NULL", id).First(&module).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("module not found")
		}
		return errors.New("failed to verify module: " + err.Error())
	}
	return nil
}

func (s *LessonService) verifyTopic(id uuid.UUID) error {
	var topic models.Topic
	if err := s.db.Where("id = ? AND deleted_at IS NULL", id).First(&topic).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("topic not found")
		}
		return errors.New("failed to verify topic: " + err.Error())
	}
	return nil
}

// validateLessonRequest validates the lesson request
func (s *LessonService) validateLessonRequest(req *dto.CreateLessonRequest) error {
	if req.SchemeOfWorkID == "" {
		return errors.New("scheme of work ID is required")
	}
	if req.ModuleID == "" {
		return errors.New("module ID is required")
	}
	if req.TopicID == "" {
		return errors.New("topic ID is required")
	}
	if req.Title == "" {
		return errors.New("title is required")
	}
	if req.Status != "" {
		validStatuses := map[string]bool{"draft": true, "published": true, "in review": true, "archived": true}
		if !validStatuses[req.Status] {
			return errors.New("status must be 'draft', 'published', 'in review', or 'archived'")
		}
	}
	if req.LessonOrder < 0 {
		return errors.New("lesson order cannot be negative")
	}
	if req.Duration < 0 {
		return errors.New("duration cannot be negative")
	}
	if req.Week < 0 {
		return errors.New("week cannot be negative")
	}
	return nil
}

// toLessonResponse converts model to response DTO
func (s *LessonService) toLessonResponse(lesson *models.Lesson) *dto.LessonResponse {
	response := &dto.LessonResponse{
		ID:             lesson.ID.String(),
		SchemeOfWorkID: lesson.SchemeOfWorkID.String(),
		ModuleID:       lesson.ModuleID.String(),
		TopicID:        lesson.TopicID.String(),
		LessonOrder:    lesson.LessonOrder,
		Title:          lesson.Title,
		Description:    lesson.Description,
		LessonDate:     lesson.LessonDate,
		Week:           lesson.Week,
		Duration:       lesson.Duration,
		Objectives:     lesson.Objectives,
		Activities:     lesson.Activities,
		Resources:      lesson.Resources,
		Assessment:     lesson.Assessment,
		Status:         lesson.Status,
		CreatedBy:      lesson.CreatedBy.String(),
		CreatedAt:      lesson.CreatedAt,
		UpdatedAt:      lesson.UpdatedAt,
	}

	// Convert Content JSON to string
	if len(lesson.Content) > 0 {
		response.Content = string(lesson.Content)
	}

	// Add scheme of work details if preloaded
	if lesson.SchemeOfWork.ID != uuid.Nil {
		response.SchemeOfWork = &dto.SchemeOfWorkResponse{
			ID:    lesson.SchemeOfWork.ID.String(),
			Title: lesson.SchemeOfWork.Title,
			Term:  lesson.SchemeOfWork.Term,
		}
	}

	// Add module details if preloaded
	if lesson.Module.ID != uuid.Nil {
		response.Module = &dto.ModuleResponse{
			ID:    lesson.Module.ID.String(),
			Title: lesson.Module.Title,
		}
	}

	// Add topic details if preloaded
	if lesson.Topic.ID != uuid.Nil {
		response.Topic = &dto.TopicResponse{
			ID:    lesson.Topic.ID.String(),
			Title: lesson.Topic.Title,
		}
	}

	// Add creator details if preloaded
	if lesson.Creator.ID != uuid.Nil {
		response.Creator = &dto.UserResponse{
			ID:        lesson.Creator.ID.String(),
			FirstName: lesson.Creator.FirstName,
			LastName:  lesson.Creator.LastName,
			Email:     lesson.Creator.Email,
		}
	}

	return response
}