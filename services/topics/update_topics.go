package services

import (
    "errors"
    "time"

    "crm-go/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type UpdateTopicService struct {
    db *gorm.DB
}

func NewUpdateTopicService(db *gorm.DB) *UpdateTopicService {
    return &UpdateTopicService{db: db}
}

// UpdateTopicWithTx - for use with transactions
func (s *UpdateTopicService) UpdateTopicWithTx(tx *gorm.DB, topicID uuid.UUID, req models.TopicInput) (*models.TopicResponse, error) {
    // Fetch existing topic
    var topic models.Topic
    if err := tx.First(&topic, "id = ?", topicID).Error; err != nil {
        return nil, errors.New("topic not found")
    }
    
    // Validate chapter exists
    var chapter models.Chapter
    if err := tx.First(&chapter, "id = ?", req.ChapterID).Error; err != nil {
        return nil, errors.New("chapter not found")
    }
    
    // Ensure chapter belongs to course
    if chapter.CourseID != req.CourseID {
        return nil, errors.New("chapter does not belong to this course")
    }
    
    // Update fields
    topic.CourseID = req.CourseID
    topic.ChapterID = req.ChapterID
    topic.TutorID = req.TutorID
    topic.Title = req.Title
    topic.Description = req.Description
    topic.Order = req.Order
    topic.UpdatedAt = time.Now()
    
    // Save changes
    if err := tx.Save(&topic).Error; err != nil {
        return nil, err
    }
    
    // Convert to response
    response := s.topicToResponse(&topic, req.TutorID)
    return response, nil
}


// Helper function to convert Topic to TopicResponse
func (s *UpdateTopicService) topicToResponse(topic *models.Topic, tutorID uuid.UUID) *models.TopicResponse {
    return &models.TopicResponse{
        ID:          topic.ID,
        CourseID:    topic.CourseID,
        ChapterID:   topic.ChapterID,
        TutorID:     tutorID,
        Title:       topic.Title,
        Description: topic.Description,
        Order:       topic.Order,
        CreatedAt:   topic.CreatedAt,
        UpdatedAt:   topic.UpdatedAt,
    }
}