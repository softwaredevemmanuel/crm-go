// services/lesson_plan_service.go
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

type LessonPlanService struct {
	db *gorm.DB
}

func NewLessonPlanService(db *gorm.DB) *LessonPlanService {
	return &LessonPlanService{db: db}
}

// CreateLessonPlan creates a new lesson plan
func (s *LessonPlanService) CreateLessonPlan(req *dto.CreateLessonPlanRequest, userID uuid.UUID) (*dto.LessonPlanResponse, error) {
	// Validate input
	if err := s.validateLessonPlanRequest(req); err != nil {
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

	// Check if lesson plan already exists for this lesson
	var existing models.LessonPlan
	if err := s.db.Where("lesson_id = ? AND deleted_at IS NULL", lessonID).First(&existing).Error; err == nil {
		return nil, errors.New("lesson plan already exists for this lesson")
	}

	// Create lesson plan
	lessonPlan := &models.LessonPlan{
		ID:                    uuid.New(),
		LessonID:              lessonID,
		PreviousKnowledge:     req.PreviousKnowledge,
		BehaviouralObjectives: req.BehaviouralObjectives,
		TeachingAids:          req.TeachingAids,
		Introduction:          req.Introduction,
		LessonContent:         req.LessonContent,
		TeacherActivities:     req.TeacherActivities,
		StudentActivities:     req.StudentActivities,
		Conclusion:            req.Conclusion,
		Evaluation:            req.Evaluation,
		CreatedBy:             userID,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	if err := s.db.Create(lessonPlan).Error; err != nil {
		return nil, errors.New("failed to create lesson plan: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("Lesson").Preload("Creator").First(lessonPlan, lessonPlan.ID).Error; err != nil {
		return nil, errors.New("failed to load lesson plan details: " + err.Error())
	}

	return s.toLessonPlanResponse(lessonPlan), nil
}

// GetAllLessonPlans retrieves all lesson plans with pagination and filters
func (s *LessonPlanService) GetAllLessonPlans(params *dto.LessonPlanQueryParams) (*dto.LessonPlanListResponse, error) {
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
	query := s.db.Model(&models.LessonPlan{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.LessonID != "" {
		lessonID, err := uuid.Parse(params.LessonID)
		if err == nil {
			query = query.Where("lesson_id = ?", lessonID)
		}
	}

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where(
			"LOWER(previous_knowledge) LIKE ? OR LOWER(behavioural_objectives) LIKE ? OR LOWER(teaching_aids) LIKE ? OR LOWER(introduction) LIKE ? OR LOWER(lesson_content) LIKE ?",
			search, search, search, search, search,
		)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count lesson plans: %w", err)
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
	var lessonPlans []models.LessonPlan
	if err := query.Preload("Lesson").Preload("Creator").Find(&lessonPlans).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch lesson plans: %w", err)
	}

	// Convert to response
	responses := make([]dto.LessonPlanResponse, len(lessonPlans))
	for i, lessonPlan := range lessonPlans {
		responses[i] = *s.toLessonPlanResponse(&lessonPlan)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.LessonPlanListResponse{
		LessonPlans: responses,
		Total:       total,
		Page:        params.Page,
		Limit:       params.Limit,
		TotalPages:  totalPages,
	}, nil
}

// GetLessonPlanByID retrieves a single lesson plan by ID
func (s *LessonPlanService) GetLessonPlanByID(id string) (*dto.LessonPlanResponse, error) {
	lessonPlanID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid lesson plan ID")
	}

	var lessonPlan models.LessonPlan
	if err := s.db.Where("id = ? AND deleted_at IS NULL", lessonPlanID).
		Preload("Lesson").
		Preload("Creator").
		First(&lessonPlan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lesson plan not found")
		}
		return nil, errors.New("failed to fetch lesson plan: " + err.Error())
	}

	return s.toLessonPlanResponse(&lessonPlan), nil
}

// GetLessonPlanByLesson retrieves a lesson plan by lesson ID
func (s *LessonPlanService) GetLessonPlanByLesson(lessonID string) (*dto.LessonPlanResponse, error) {
	lID, err := uuid.Parse(lessonID)
	if err != nil {
		return nil, errors.New("invalid lesson ID")
	}

	var lessonPlan models.LessonPlan
	if err := s.db.Where("lesson_id = ? AND deleted_at IS NULL", lID).
		Preload("Lesson").
		Preload("Creator").
		First(&lessonPlan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lesson plan not found for this lesson")
		}
		return nil, errors.New("failed to fetch lesson plan: " + err.Error())
	}

	return s.toLessonPlanResponse(&lessonPlan), nil
}

// UpdateLessonPlan updates an existing lesson plan
func (s *LessonPlanService) UpdateLessonPlan(id string, req *dto.UpdateLessonPlanRequest) (*dto.LessonPlanResponse, error) {
	lessonPlanID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid lesson plan ID")
	}

	// Find existing lesson plan
	var lessonPlan models.LessonPlan
	if err := s.db.Where("id = ? AND deleted_at IS NULL", lessonPlanID).First(&lessonPlan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lesson plan not found")
		}
		return nil, errors.New("failed to fetch lesson plan: " + err.Error())
	}

	// Update fields
	if req.PreviousKnowledge != "" {
		lessonPlan.PreviousKnowledge = req.PreviousKnowledge
	}
	if req.BehaviouralObjectives != "" {
		lessonPlan.BehaviouralObjectives = req.BehaviouralObjectives
	}
	if req.TeachingAids != "" {
		lessonPlan.TeachingAids = req.TeachingAids
	}
	if req.Introduction != "" {
		lessonPlan.Introduction = req.Introduction
	}
	if req.LessonContent != "" {
		lessonPlan.LessonContent = req.LessonContent
	}
	if req.TeacherActivities != "" {
		lessonPlan.TeacherActivities = req.TeacherActivities
	}
	if req.StudentActivities != "" {
		lessonPlan.StudentActivities = req.StudentActivities
	}
	if req.Conclusion != "" {
		lessonPlan.Conclusion = req.Conclusion
	}
	if req.Evaluation != "" {
		lessonPlan.Evaluation = req.Evaluation
	}

	lessonPlan.UpdatedAt = time.Now()

	if err := s.db.Save(&lessonPlan).Error; err != nil {
		return nil, errors.New("failed to update lesson plan: " + err.Error())
	}

	// Preload relationships
	if err := s.db.Preload("Lesson").Preload("Creator").First(&lessonPlan, lessonPlan.ID).Error; err != nil {
		return nil, errors.New("failed to load lesson plan details: " + err.Error())
	}

	return s.toLessonPlanResponse(&lessonPlan), nil
}

// DeleteLessonPlan soft deletes a lesson plan
func (s *LessonPlanService) DeleteLessonPlan(id string) error {
	lessonPlanID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid lesson plan ID")
	}

	var lessonPlan models.LessonPlan
	if err := s.db.Where("id = ? AND deleted_at IS NULL", lessonPlanID).First(&lessonPlan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("lesson plan not found")
		}
		return errors.New("failed to fetch lesson plan: " + err.Error())
	}

	if err := s.db.Delete(&lessonPlan).Error; err != nil {
		return errors.New("failed to delete lesson plan: " + err.Error())
	}

	return nil
}

// validateLessonPlanRequest validates the lesson plan request
func (s *LessonPlanService) validateLessonPlanRequest(req *dto.CreateLessonPlanRequest) error {
	if req.LessonID == "" {
		return errors.New("lesson ID is required")
	}
	return nil
}

// toLessonPlanResponse converts model to response DTO
func (s *LessonPlanService) toLessonPlanResponse(lessonPlan *models.LessonPlan) *dto.LessonPlanResponse {
	response := &dto.LessonPlanResponse{
		ID:                    lessonPlan.ID.String(),
		LessonID:              lessonPlan.LessonID.String(),
		PreviousKnowledge:     lessonPlan.PreviousKnowledge,
		BehaviouralObjectives: lessonPlan.BehaviouralObjectives,
		TeachingAids:          lessonPlan.TeachingAids,
		Introduction:          lessonPlan.Introduction,
		LessonContent:         lessonPlan.LessonContent,
		TeacherActivities:     lessonPlan.TeacherActivities,
		StudentActivities:     lessonPlan.StudentActivities,
		Conclusion:            lessonPlan.Conclusion,
		Evaluation:            lessonPlan.Evaluation,
		CreatedBy:             lessonPlan.CreatedBy.String(),
		CreatedAt:             lessonPlan.CreatedAt,
		UpdatedAt:             lessonPlan.UpdatedAt,
	}

	// Add lesson details if preloaded
	if lessonPlan.Lesson.ID != uuid.Nil {
		response.Lesson = &dto.LessonResponse{
			ID:    lessonPlan.Lesson.ID.String(),
			Title: lessonPlan.Lesson.Title,
		}
	}

	// Add creator details if preloaded
	if lessonPlan.Creator.ID != uuid.Nil {
		response.Creator = &dto.UserResponse{
			ID:        lessonPlan.Creator.ID.String(),
			FirstName: lessonPlan.Creator.FirstName,
			LastName:  lessonPlan.Creator.LastName,
			Email:     lessonPlan.Creator.Email,
			Phone:     lessonPlan.Creator.Phone,
			Role:      lessonPlan.Creator.Role,
		}
	}

	return response
}