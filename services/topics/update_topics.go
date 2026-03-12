package services

import (
	"errors"
	"time"

	"crm-go/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateLessonService struct {
	db *gorm.DB
}

func NewUpdateLessonService(db *gorm.DB) *UpdateLessonService {
	return &UpdateLessonService{db: db}
}

// UpdateLessonWithTx - for use with transactions
func (s *UpdateLessonService) UpdateLessonWithTx(tx *gorm.DB, lessonID uuid.UUID, req models.LessonInput) (*models.LessonResponse, error) {
	// Fetch existing lesson
	var lesson models.Lesson
	if err := tx.First(&lesson, "id = ?", lessonID).Error; err != nil {
		return nil, errors.New("lesson not found")
	}

	// Validate module exists
	var module models.Module
	if err := tx.First(&module, "id = ?", req.ModuleID).Error; err != nil {
		return nil, errors.New("module not found")
	}

	// Ensure module belongs to course
	if module.CourseID != req.CourseID {
		return nil, errors.New("module does not belong to this course")
	}

	// Update fields
	lesson.CourseID = req.CourseID
	lesson.ModuleID = req.ModuleID
	lesson.TutorID = req.TutorID
	lesson.Title = req.Title
	lesson.Description = req.Description
	lesson.Order = req.Order
	lesson.UpdatedAt = time.Now()

	// Save changes
	if err := tx.Save(&lesson).Error; err != nil {
		return nil, err
	}

	// Convert to response
	response := s.lessonToResponse(&lesson, req.TutorID)
	return response, nil
}

// Helper function to convert Lesson to LessonResponse
func (s *UpdateLessonService) lessonToResponse(lesson *models.Lesson, tutorID uuid.UUID) *models.LessonResponse {
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
