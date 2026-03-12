package services

import (
	"errors"
	"time"

	"crm-go/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	// "fmt"
)

type CreateLessonService struct {
	db *gorm.DB
}

func NewCreateLessonService(db *gorm.DB) *CreateLessonService {
	return &CreateLessonService{db: db}
}

// CreateLessonWithTx - for use with transactions
func (s *CreateLessonService) CreateLessonWithTx(tx *gorm.DB, req models.LessonInput) (*models.LessonResponse, error) {
	// Validate Module exists
	var module models.Module
	if err := tx.First(&module, "id = ?", req.ModuleID).Error; err != nil {
		return nil, errors.New("module not found")
	}

	// Ensure module belongs to course
	if module.CourseID != req.CourseID {
		return nil, errors.New("module does not belong to this course")
	}

	// Check if lesson already exists with same title in same module
	var existingLesson models.Lesson
	err := tx.Where("module_id = ? AND LOWER(title) = LOWER(?)", req.ModuleID, req.Title).
		First(&existingLesson).Error

	if err == nil {
		return nil, errors.New("lesson with this title already exists in this module")
	} else if err != gorm.ErrRecordNotFound {
		// Some other database error
		return nil, err
	}

	// Also check if order number is already used in same module
	var lessonWithSameOrder models.Lesson
	err = tx.Where("module_id = ? AND \"order\" = ?", req.ModuleID, req.Order).
		First(&lessonWithSameOrder).Error

	if err == nil {
		return nil, errors.New("a lesson with this order number already exists in this module")
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Validate Module exists
	if err := tx.First(&module, "id = ?", req.ModuleID).Error; err != nil {
		return nil, errors.New("module not found")
	}

	// Ensure module belongs to course
	if module.CourseID != req.CourseID {
		return nil, errors.New("module does not belong to this course")
	}

	// Create lesson
	lesson := models.Lesson{
		ID:          uuid.New(),
		CourseID:    req.CourseID,
		ModuleID:    req.ModuleID,
		TutorID:     req.TutorID,
		Title:       req.Title,
		Description: req.Description,
		Order:       req.Order,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := tx.Create(&lesson).Error; err != nil {
		return nil, err
	}

	// Convert to response
	response := s.lessonToResponse(&lesson, req.TutorID)
	return response, nil
}

// Helper function to convert Lesson to LessonResponse
func (s *CreateLessonService) lessonToResponse(lesson *models.Lesson, tutorID uuid.UUID) *models.LessonResponse {
	return &models.LessonResponse{
		ID:          lesson.ID,
		CourseID:    lesson.CourseID,
		ModuleID:    lesson.ModuleID,
		TutorID:     tutorID,
		Title:       lesson.Title,
		Description: lesson.Description,
		Order:       lesson.Order,
		CreatedAt:   lesson.CreatedAt,
		UpdatedAt:   lesson.UpdatedAt,
	}
}
