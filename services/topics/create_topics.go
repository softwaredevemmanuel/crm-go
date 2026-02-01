package services

import (
    "errors"
    "time"
    
    "crm-go/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
    // "fmt"
)

type CreateTopicService struct {
    db *gorm.DB
}

func NewCreateTopicService(db *gorm.DB) *CreateTopicService {
    return &CreateTopicService{db: db}
}

// CreateTopicWithTx - for use with transactions
func (s *CreateTopicService) CreateTopicWithTx(tx *gorm.DB, req models.TopicInput) (*models.TopicResponse, error) {
    // Validate Chapter exists
    var chapter models.Chapter
    if err := tx.First(&chapter, "id = ?", req.ChapterID).Error; err != nil {
        return nil, errors.New("chapter not found")
    }
    
    // Ensure chapter belongs to course
    if chapter.CourseID != req.CourseID {
        return nil, errors.New("chapter does not belong to this course")
    }
    
	 // Check if topic already exists with same title in same chapter
    var existingTopic models.Topic
    err := tx.Where("chapter_id = ? AND LOWER(title) = LOWER(?)", req.ChapterID, req.Title).
        First(&existingTopic).Error
    
    if err == nil {
        return nil, errors.New("topic with this title already exists in this chapter")
    } else if err != gorm.ErrRecordNotFound {
        // Some other database error
        return nil, err
    }
    
    // Also check if order number is already used in same chapter
    var topicWithSameOrder models.Topic
    err = tx.Where("chapter_id = ? AND \"order\" = ?", req.ChapterID, req.Order).
        First(&topicWithSameOrder).Error
    
    if err == nil {
        return nil, errors.New("a topic with this order number already exists in this chapter")
    } else if err != gorm.ErrRecordNotFound {
        return nil, err
    }

	   // Validate Chapter exists
    if err := tx.First(&chapter, "id = ?", req.ChapterID).Error; err != nil {
        return nil, errors.New("chapter not found")
    }
    
    // Ensure chapter belongs to course
    if chapter.CourseID != req.CourseID {
        return nil, errors.New("chapter does not belong to this course")
    }
    
    // Create topic
    topic := models.Topic{
        ID:          uuid.New(),
        CourseID:    req.CourseID,
        ChapterID:   req.ChapterID,
        TutorID:   	 req.TutorID,
        Title:       req.Title,
        Description: req.Description,
        Order:       req.Order,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    
    if err := tx.Create(&topic).Error; err != nil {
        return nil, err
    }
    
    // Convert to response
    response := s.topicToResponse(&topic, req.TutorID)
    return response, nil
}



// Helper function to convert Topic to TopicResponse
func (s *CreateTopicService) topicToResponse(topic *models.Topic, tutorID uuid.UUID) *models.TopicResponse {
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


