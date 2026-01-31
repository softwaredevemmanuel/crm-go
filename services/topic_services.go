package services

import (
	"errors"
	"time"

	"crm-go/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TopicService struct {
	db *gorm.DB
}

func NewTopicService(db *gorm.DB) *TopicService {
	return &TopicService{db: db}
}

func (s *TopicService) CreateTopic(req models.TopicInput) (*models.Topic, error) {
	// Validate Chapter exists
	var chapter models.Chapter
	if err := s.db.First(&chapter, "id = ?", req.ChapterID).Error; err != nil {
		return nil, errors.New("chapter not found")
	}

	// Optional: Ensure chapter belongs to course
	if chapter.CourseID != req.CourseID {
		return nil, errors.New("chapter does not belong to this course")
	}

	topic := models.Topic{
		ID:          uuid.New(),
		CourseID:    req.CourseID,
		ChapterID:   req.ChapterID,
		Title:       req.Title,
		Description: req.Description,
		Order:       req.Order,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(&topic).Error; err != nil {
		return nil, err
	}

	return &topic, nil
}
