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
func (s *LessonService) CreateLesson(req *dto.CreateLessonRequest) (*dto.LessonResponse, error) {
	// Validate input
	if err := s.validateLessonRequest(req); err != nil {
		return nil, err
	}

	// Parse UUIDs
	schemeOfWorkItemID, err := uuid.Parse(req.SchemeOfWorkItemID)
	if err != nil {
		return nil, errors.New("invalid scheme of work item ID format")
	}

	classID, err := uuid.Parse(req.ClassID)
	if err != nil {
		return nil, errors.New("invalid class ID format")
	}

	var armID uuid.UUID
	if req.ArmID != "" {
		armID, err = uuid.Parse(req.ArmID)
		if err != nil {
			return nil, errors.New("invalid arm ID format")
		}
	}

	// Check if scheme of work item exists
	var schemeItem models.SchemeOfWorkItem
	if err := s.db.Where("id = ? AND deleted_at IS NULL", schemeOfWorkItemID).First(&schemeItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("scheme of work item not found")
		}
		return nil, errors.New("failed to verify scheme of work item: " + err.Error())
	}

	// Check if class exists
	var class models.ClassGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", classID).First(&class).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("class not found")
		}
		return nil, errors.New("failed to verify class: " + err.Error())
	}

	// Check if arm exists (if provided)
	if req.ArmID != "" {
		var arm models.Arm
		if err := s.db.Where("id = ? AND deleted_at IS NULL", armID).First(&arm).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("arm not found")
			}
			return nil, errors.New("failed to verify arm: " + err.Error())
		}
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
		status = "planned"
	}

	// Create lesson
	lesson := &models.Lesson{
		ID:                 uuid.New(),
		SchemeOfWorkItemID: schemeOfWorkItemID,
		ClassID:            classID,
		ArmID:              armID,
		Title:              req.Title,
		LessonDate:         lessonDate,
		Week:               req.Week,
		Period:             req.Period,
		Duration:           req.Duration,
		Status:             status,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := s.db.Create(lesson).Error; err != nil {
		return nil, errors.New("failed to create lesson: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("SchemeOfWorkItem").Preload("Class").Preload("Arm").First(lesson, lesson.ID).Error; err != nil {
		return nil, errors.New("failed to load lesson details: " + err.Error())
	}

	return s.toLessonResponse(lesson), nil
}

// BulkCreateLessons creates multiple lessons
func (s *LessonService) BulkCreateLessons(req *dto.BulkCreateLessonsRequest) (*dto.BulkLessonResult, error) {
	// Parse UUIDs

	classID, err := uuid.Parse(req.ClassID)
	if err != nil {
		return nil, errors.New("invalid class ID format")
	}

	var armID uuid.UUID
	if req.ArmID != "" {
		armID, err = uuid.Parse(req.ArmID)
		if err != nil {
			return nil, errors.New("invalid arm ID format")
		}
	}


	// Check if class exists
	var class models.ClassGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", classID).First(&class).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("class not found")
		}
		return nil, errors.New("failed to verify class: " + err.Error())
	}

	// Check if arm exists (if provided)
	if req.ArmID != "" {
		var arm models.Arm
		if err := s.db.Where("id = ? AND deleted_at IS NULL", armID).First(&arm).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("arm not found")
			}
			return nil, errors.New("failed to verify arm: " + err.Error())
		}
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "planned"
	}

	result := &dto.BulkLessonResult{
		Created: []dto.LessonResponse{},
		Errors:  []dto.BulkLessonError{},
	}

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

		if lessonReq.SchemeOfWorkItemID == "" {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkLessonError{
				Title: lessonReq.Title,
				Error: "scheme of work item ID is required",
			})
			continue
		}

		// Parse scheme of work item ID
		schemeItemID, err := uuid.Parse(lessonReq.SchemeOfWorkItemID)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkLessonError{
				Title: lessonReq.Title,
				Error: "invalid scheme of work item ID format",
			})
			continue
		}

		// Check if scheme of work item exists
		var schemeItem models.SchemeOfWorkItem
		if err := s.db.Where("id = ? AND deleted_at IS NULL", schemeItemID).First(&schemeItem).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkLessonError{
				Title: lessonReq.Title,
				Error: "scheme of work item not found",
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

		// Create lesson
		lesson := &models.Lesson{
			ID:                 uuid.New(),
			SchemeOfWorkItemID: schemeItemID,
			ClassID:            classID,
			ArmID:              armID,
			Title:              lessonReq.Title,
			LessonDate:         lessonDate,
			Week:               lessonReq.Week,
			Period:             lessonReq.Period,
			Duration:           lessonReq.Duration,
			Status:             status,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
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
		if err := s.db.Preload("SchemeOfWorkItem").Preload("Class").Preload("Arm").First(lesson, lesson.ID).Error; err != nil {
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
		params.SortBy = "lesson_date"
	}
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}

	// Build query
	query := s.db.Model(&models.Lesson{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.SchemeOfWorkItemID != "" {
		itemID, err := uuid.Parse(params.SchemeOfWorkItemID)
		if err == nil {
			query = query.Where("scheme_of_work_item_id = ?", itemID)
		}
	}


	if params.ClassID != "" {
		classID, err := uuid.Parse(params.ClassID)
		if err == nil {
			query = query.Where("class_id = ?", classID)
		}
	}

	if params.ArmID != "" {
		armID, err := uuid.Parse(params.ArmID)
		if err == nil {
			query = query.Where("arm_id = ?", armID)
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
		query = query.Where("LOWER(title) LIKE ?", search)
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
	if err := query.Preload("SchemeOfWorkItem").Preload("Class").Preload("Arm").Find(&lessons).Error; err != nil {
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
		Preload("SchemeOfWorkItem").
		Preload("Class").
		Preload("Arm").
		First(&lesson).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lesson not found")
		}
		return nil, errors.New("failed to fetch lesson: " + err.Error())
	}

	return s.toLessonResponse(&lesson), nil
}


// GetLessonsByClass retrieves all lessons for a specific class
func (s *LessonService) GetLessonsByClass(classID string) ([]dto.LessonResponse, error) {
	cID, err := uuid.Parse(classID)
	if err != nil {
		return nil, errors.New("invalid class ID")
	}

	var lessons []models.Lesson
	if err := s.db.Where("class_id = ? AND deleted_at IS NULL", cID).
		Preload("SchemeOfWorkItem").
		Preload("Class").
		Preload("Arm").
		Order("lesson_date DESC").
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
	if req.SchemeOfWorkItemID != "" {
		itemID, err := uuid.Parse(req.SchemeOfWorkItemID)
		if err != nil {
			return nil, errors.New("invalid scheme of work item ID format")
		}
		var schemeItem models.SchemeOfWorkItem
		if err := s.db.Where("id = ? AND deleted_at IS NULL", itemID).First(&schemeItem).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("scheme of work item not found")
			}
			return nil, errors.New("failed to verify scheme of work item: " + err.Error())
		}
		lesson.SchemeOfWorkItemID = itemID
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
		lesson.ClassID = classID
	}

	if req.ArmID != "" {
		armID, err := uuid.Parse(req.ArmID)
		if err != nil {
			return nil, errors.New("invalid arm ID format")
		}
		var arm models.Arm
		if err := s.db.Where("id = ? AND deleted_at IS NULL", armID).First(&arm).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("arm not found")
			}
			return nil, errors.New("failed to verify arm: " + err.Error())
		}
		lesson.ArmID = armID
	}

	if req.Title != "" {
		lesson.Title = req.Title
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

	if req.Period > 0 {
		lesson.Period = req.Period
	}

	if req.Duration > 0 {
		lesson.Duration = req.Duration
	}

	if req.Status != "" {
		validStatuses := map[string]bool{"planned": true, "ongoing": true, "completed": true, "cancelled": true}
		if !validStatuses[req.Status] {
			return nil, errors.New("status must be 'planned', 'ongoing', 'completed', or 'cancelled'")
		}
		lesson.Status = req.Status
	}

	lesson.UpdatedAt = time.Now()

	if err := s.db.Save(&lesson).Error; err != nil {
		return nil, errors.New("failed to update lesson: " + err.Error())
	}

	// Preload relationships
	if err := s.db.Preload("SchemeOfWorkItem").Preload("Class").Preload("Arm").First(&lesson, lesson.ID).Error; err != nil {
		return nil, errors.New("failed to load lesson details: " + err.Error())
	}

	return s.toLessonResponse(&lesson), nil
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

// validateLessonRequest validates the lesson request
func (s *LessonService) validateLessonRequest(req *dto.CreateLessonRequest) error {
	if req.SchemeOfWorkItemID == "" {
		return errors.New("scheme of work item ID is required")
	}
	if req.ClassID == "" {
		return errors.New("class ID is required")
	}
	if req.Title == "" {
		return errors.New("title is required")
	}
	if req.Status != "" {
		validStatuses := map[string]bool{"planned": true, "ongoing": true, "completed": true, "cancelled": true}
		if !validStatuses[req.Status] {
			return errors.New("status must be 'planned', 'ongoing', 'completed', or 'cancelled'")
		}
	}
	return nil
}

// toLessonResponse converts model to response DTO
func (s *LessonService) toLessonResponse(lesson *models.Lesson) *dto.LessonResponse {
	response := &dto.LessonResponse{
		ID:                 lesson.ID.String(),
		SchemeOfWorkItemID: lesson.SchemeOfWorkItemID.String(),
		ClassID:            lesson.ClassID.String(),
		ArmID:              lesson.ArmID.String(),
		Title:              lesson.Title,
		LessonDate:         lesson.LessonDate,
		Week:               lesson.Week,
		Period:             lesson.Period,
		Duration:           lesson.Duration,
		Status:             lesson.Status,
		CreatedAt:          lesson.CreatedAt,
		UpdatedAt:          lesson.UpdatedAt,
	}

	// Add scheme of work item details if preloaded
	if lesson.SchemeOfWorkItem.ID != uuid.Nil {
		response.SchemeOfWorkItem = &dto.SchemeOfWorkItemResponse{
			ID:        lesson.SchemeOfWorkItem.ID.String(),
			Topic:     lesson.SchemeOfWorkItem.Topic,
			Subtopic:  lesson.SchemeOfWorkItem.Subtopic,
			WeekStart: lesson.SchemeOfWorkItem.WeekStart,
			WeekEnd:   lesson.SchemeOfWorkItem.WeekEnd,
			Sequence:  lesson.SchemeOfWorkItem.Sequence,
		}
	}


	// Add class details if preloaded
	if lesson.Class.ID != uuid.Nil {
		response.Class = &dto.ClassGradeResponse{
			ID:          lesson.Class.ID.String(),
			Name:        lesson.Class.Name,
			Code:        lesson.Class.Code,
			Level:       lesson.Class.Level,
			Description: lesson.Class.Description,
			Status:      lesson.Class.Status,
		}
	}

	// Add arm details if preloaded
	if lesson.Arm.ID != uuid.Nil {
		response.Arm = &dto.ArmResponse{
			ID:          lesson.Arm.ID.String(),
			Name:        lesson.Arm.Name,
			Code:        lesson.Arm.Code,
			Description: lesson.Arm.Description,
			Capacity:    lesson.Arm.Capacity,
			Status:      lesson.Arm.Status,
		}
	}

	return response
}