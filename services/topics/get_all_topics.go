package services

import (
	"errors"
	"fmt"

	"crm-go/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetLessonService struct {
	db *gorm.DB
}

func NewGetLessonService(db *gorm.DB) *GetLessonService {
	return &GetLessonService{db: db}
}

// GetAllLessons with filtering options
func (s *GetLessonService) GetAllLessons(filters models.LessonFilters) ([]models.LessonResponse, error) {
	var lessons []models.Lesson

	// Start building query
	query := s.db.Model(&models.Lesson{})

	// Apply filters if provided
	if filters.CourseID != uuid.Nil {
		query = query.Where("course_id = ?", filters.CourseID)
	}

	if filters.ModuleID != uuid.Nil {
		query = query.Where("module_id = ?", filters.ModuleID)
	}

	if filters.TutorID != uuid.Nil {
		// Assuming lessons don't have tutor_id directly
		// Join with modules or courses if needed
		query = query.Joins("JOIN modules ON modules.id = lessons.module_id").
			Where("modules.tutor_id = ?", filters.TutorID)
	}

	if filters.Search != "" {
		searchTerm := "%" + filters.Search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	// Apply sorting
	sortBy := "created_at"
	if filters.SortBy != "" {
		// Validate sort field to prevent SQL injection
		validSortFields := map[string]bool{
			"title": true, "order": true, "created_at": true, "updated_at": true,
		}
		if validSortFields[filters.SortBy] {
			sortBy = filters.SortBy
		}
	}

	sortOrder := "DESC"
	if filters.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	// For order field, need to quote it in PostgreSQL
	if sortBy == "order" {
		query = query.Order(fmt.Sprintf("\"order\" %s", sortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
	}

	// Apply pagination
	if filters.Page > 0 && filters.Limit > 0 {
		offset := (filters.Page - 1) * filters.Limit
		query = query.Offset(offset).Limit(filters.Limit)
	}

	// Execute query
	if err := query.Find(&lessons).Error; err != nil {
		return nil, errors.New("failed to fetch lessons: " + err.Error())
	}

	// Convert to response
	return s.lessonsToResponse(lessons, filters.TutorID), nil
}

// GetLessonsByModuleID - convenience method
func (s *GetLessonService) GetLessonsByModuleID(moduleID uuid.UUID, tutorID uuid.UUID) ([]models.LessonResponse, error) {
	filters := models.LessonFilters{
		ModuleID: moduleID,
		TutorID:  tutorID,
	}
	return s.GetAllLessons(filters)
}

// GetLessonsByCourseID - convenience method
func (s *GetLessonService) GetLessonsByCourseID(courseID uuid.UUID, tutorID uuid.UUID) ([]models.LessonResponse, error) {
	filters := models.LessonFilters{
		CourseID: courseID,
		TutorID:  tutorID,
	}
	return s.GetAllLessons(filters)
}

// Helper to convert slice of Lessons to slice of LessonResponses
func (s *GetLessonService) lessonsToResponse(lessons []models.Lesson, tutorID uuid.UUID) []models.LessonResponse {
	responses := make([]models.LessonResponse, len(lessons))
	for i, lesson := range lessons {
		responses[i] = models.LessonResponse{
			ID:          lesson.ID,
			CourseID:    lesson.CourseID,
			ModuleID:    lesson.ModuleID,
			TutorID:     lesson.TutorID,
			Title:       lesson.Title,
			Description: lesson.Description,
			Order:       lesson.Order,
			CreatedAt:   lesson.CreatedAt,
			UpdatedAt:   lesson.UpdatedAt,
		}
	}
	return responses
}

// GetLessonCount - get total count with filters
func (s *GetLessonService) GetLessonCount(filters models.LessonFilters) (int64, error) {
	var count int64

	query := s.db.Model(&models.Lesson{})

	// Apply the same filters as GetAllLessons
	if filters.CourseID != uuid.Nil {
		query = query.Where("course_id = ?", filters.CourseID)
	}

	if filters.ModuleID != uuid.Nil {
		query = query.Where("module_id = ?", filters.ModuleID)
	}

	if filters.TutorID != uuid.Nil {
		query = query.Joins("JOIN modules ON modules.id = lessons.module_id").
			Where("modules.tutor_id = ?", filters.TutorID)
	}

	if filters.Search != "" {
		searchTerm := "%" + filters.Search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, errors.New("failed to count lessons: " + err.Error())
	}

	return count, nil
}

// GetAllLessonsWithPagination - returns lessons with pagination metadata
func (s *GetLessonService) GetAllLessonsWithPagination(filters models.LessonFilters) (*models.PaginatedLessonsResponse, error) {
	// Get lessons
	lessons, err := s.GetAllLessons(filters)
	if err != nil {
		return nil, err
	}

	// Get total count
	totalCount, err := s.GetLessonCount(filters)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := 0
	if filters.Limit > 0 {
		totalPages = int((totalCount + int64(filters.Limit) - 1) / int64(filters.Limit))
	}

	return &models.PaginatedLessonsResponse{
		Data:       lessons,
		Total:      totalCount,
		Page:       filters.Page,
		Limit:      filters.Limit,
		TotalPages: totalPages,
	}, nil
}
