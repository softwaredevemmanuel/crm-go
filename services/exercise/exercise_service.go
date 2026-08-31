// services/exercise_service.go
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

type ExerciseService struct {
	db *gorm.DB
}

func NewExerciseService(db *gorm.DB) *ExerciseService {
	return &ExerciseService{db: db}
}

// CreateExercise creates a new exercise
func (s *ExerciseService) CreateExercise(req *dto.CreateExerciseRequest, userID uuid.UUID) (*dto.ExerciseResponse, error) {
	// Validate input
	if err := s.validateExerciseRequest(req); err != nil {
		return nil, err
	}

	// Parse UUIDs
	lessonID, err := uuid.Parse(req.LessonID)
	if err != nil {
		return nil, errors.New("invalid lesson ID format")
	}

	// Check if lesson exists
	var lesson models.Lesson
	if err := s.db.Where("id = ? AND deleted_at IS NULL", lessonID).First(&lesson).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lesson not found")
		}
		return nil, errors.New("failed to verify lesson: " + err.Error())
	}

	// Create exercise
	exercise := &models.Exercise{
		ID:           uuid.New(),
		LessonID:     lessonID,
		Title:        req.Title,
		Instructions: req.Instructions,
		Content:      req.Content,
		TotalMarks:   req.TotalMarks,
		CreatedBy:    userID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.db.Create(exercise).Error; err != nil {
		return nil, errors.New("failed to create exercise: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("Lesson").Preload("Creator").First(exercise, exercise.ID).Error; err != nil {
		return nil, errors.New("failed to load exercise details: " + err.Error())
	}

	return s.toExerciseResponse(exercise), nil
}

// BulkCreateExercises creates multiple exercises
func (s *ExerciseService) BulkCreateExercises(req *dto.BulkCreateExercisesRequest, userID uuid.UUID) (*dto.BulkExerciseResult, error) {
	// Parse UUIDs
	lessonID, err := uuid.Parse(req.LessonID)
	if err != nil {
		return nil, errors.New("invalid lesson ID format")
	}

	// Check if lesson exists
	var lesson models.Lesson
	if err := s.db.Where("id = ? AND deleted_at IS NULL", lessonID).First(&lesson).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lesson not found")
		}
		return nil, errors.New("failed to verify lesson: " + err.Error())
	}

	result := &dto.BulkExerciseResult{
		Created: []dto.ExerciseResponse{},
		Errors:  []dto.BulkExerciseError{},
	}

	for _, exerciseReq := range req.Exercises {
		// Validate
		if exerciseReq.Title == "" {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkExerciseError{
				Title: exerciseReq.Title,
				Error: "title is required",
			})
			continue
		}

		// Create exercise
		exercise := &models.Exercise{
			ID:           uuid.New(),
			LessonID:     lessonID,
			Title:        exerciseReq.Title,
			Instructions: exerciseReq.Instructions,
			Content:      exerciseReq.Content,
			TotalMarks:   exerciseReq.TotalMarks,
			CreatedBy:    userID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := s.db.Create(exercise).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkExerciseError{
				Title: exerciseReq.Title,
				Error: "failed to create exercise: " + err.Error(),
			})
			continue
		}

		// Preload relationships
		if err := s.db.Preload("Lesson").Preload("Creator").First(exercise, exercise.ID).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkExerciseError{
				Title: exerciseReq.Title,
				Error: "failed to load exercise details",
			})
			continue
		}

		result.SuccessCount++
		result.Created = append(result.Created, *s.toExerciseResponse(exercise))
	}

	return result, nil
}

// GetAllExercises retrieves all exercises with pagination and filters
func (s *ExerciseService) GetAllExercises(params *dto.ExerciseQueryParams) (*dto.ExerciseListResponse, error) {
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
	query := s.db.Model(&models.Exercise{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.LessonID != "" {
		lessonID, err := uuid.Parse(params.LessonID)
		if err == nil {
			query = query.Where("lesson_id = ?", lessonID)
		}
	}

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(instructions) LIKE ? OR LOWER(content) LIKE ?",
			search, search, search)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count exercises: %w", err)
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
	var exercises []models.Exercise
	if err := query.Preload("Lesson").Preload("Creator").Find(&exercises).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch exercises: %w", err)
	}

	// Convert to response
	responses := make([]dto.ExerciseResponse, len(exercises))
	for i, exercise := range exercises {
		responses[i] = *s.toExerciseResponse(&exercise)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.ExerciseListResponse{
		Exercises:  responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetExerciseByID retrieves a single exercise by ID
func (s *ExerciseService) GetExerciseByID(id string) (*dto.ExerciseResponse, error) {
	exerciseID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid exercise ID")
	}

	var exercise models.Exercise
	if err := s.db.Where("id = ? AND deleted_at IS NULL", exerciseID).
		Preload("Lesson").
		Preload("Creator").
		First(&exercise).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("exercise not found")
		}
		return nil, errors.New("failed to fetch exercise: " + err.Error())
	}

	return s.toExerciseResponse(&exercise), nil
}

// GetExercisesByLesson retrieves all exercises for a specific lesson
func (s *ExerciseService) GetExercisesByLesson(lessonID string) ([]dto.ExerciseResponse, error) {
	lID, err := uuid.Parse(lessonID)
	if err != nil {
		return nil, errors.New("invalid lesson ID")
	}

	var exercises []models.Exercise
	if err := s.db.Where("lesson_id = ? AND deleted_at IS NULL", lID).
		Preload("Lesson").
		Preload("Creator").
		Order("created_at ASC").
		Find(&exercises).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch exercises: %w", err)
	}

	responses := make([]dto.ExerciseResponse, len(exercises))
	for i, exercise := range exercises {
		responses[i] = *s.toExerciseResponse(&exercise)
	}

	return responses, nil
}

// UpdateExercise updates an existing exercise
func (s *ExerciseService) UpdateExercise(id string, req *dto.UpdateExerciseRequest) (*dto.ExerciseResponse, error) {
	exerciseID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid exercise ID")
	}

	// Find existing exercise
	var exercise models.Exercise
	if err := s.db.Where("id = ? AND deleted_at IS NULL", exerciseID).First(&exercise).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("exercise not found")
		}
		return nil, errors.New("failed to fetch exercise: " + err.Error())
	}

	// Update fields
	if req.Title != "" {
		exercise.Title = req.Title
	}
	if req.Instructions != "" {
		exercise.Instructions = req.Instructions
	}
	if req.Content != "" {
		exercise.Content = req.Content
	}
	if req.TotalMarks > 0 {
		exercise.TotalMarks = req.TotalMarks
	}

	exercise.UpdatedAt = time.Now()

	if err := s.db.Save(&exercise).Error; err != nil {
		return nil, errors.New("failed to update exercise: " + err.Error())
	}

	// Preload relationships
	if err := s.db.Preload("Lesson").Preload("Creator").First(&exercise, exercise.ID).Error; err != nil {
		return nil, errors.New("failed to load exercise details: " + err.Error())
	}

	return s.toExerciseResponse(&exercise), nil
}

// DeleteExercise soft deletes an exercise
func (s *ExerciseService) DeleteExercise(id string) error {
	exerciseID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid exercise ID")
	}

	var exercise models.Exercise
	if err := s.db.Where("id = ? AND deleted_at IS NULL", exerciseID).First(&exercise).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("exercise not found")
		}
		return errors.New("failed to fetch exercise: " + err.Error())
	}

	if err := s.db.Delete(&exercise).Error; err != nil {
		return errors.New("failed to delete exercise: " + err.Error())
	}

	return nil
}

// validateExerciseRequest validates the exercise request
func (s *ExerciseService) validateExerciseRequest(req *dto.CreateExerciseRequest) error {
	if req.LessonID == "" {
		return errors.New("lesson ID is required")
	}
	if req.Title == "" {
		return errors.New("title is required")
	}
	if req.TotalMarks < 0 {
		return errors.New("total marks cannot be negative")
	}
	return nil
}

// toExerciseResponse converts model to response DTO
func (s *ExerciseService) toExerciseResponse(exercise *models.Exercise) *dto.ExerciseResponse {
	response := &dto.ExerciseResponse{
		ID:           exercise.ID.String(),
		LessonID:     exercise.LessonID.String(),
		Title:        exercise.Title,
		Instructions: exercise.Instructions,
		Content:      exercise.Content,
		TotalMarks:   exercise.TotalMarks,
		CreatedBy:    exercise.CreatedBy.String(),
		CreatedAt:    exercise.CreatedAt,
		UpdatedAt:    exercise.UpdatedAt,
	}

	// Add lesson details if preloaded
	if exercise.Lesson.ID != uuid.Nil {
		response.Lesson = &dto.LessonResponse{
			ID:    exercise.Lesson.ID.String(),
			Title: exercise.Lesson.Title,
		}
	}

	// Add creator details if preloaded
	if exercise.Creator.ID != uuid.Nil {
		response.Creator = &dto.UserResponse{
			ID:        exercise.Creator.ID.String(),
			FirstName: exercise.Creator.FirstName,
			LastName:  exercise.Creator.LastName,
			Email:     exercise.Creator.Email,
			Phone:     exercise.Creator.Phone,
			Role:      exercise.Creator.Role,
		}
	}

	return response
}