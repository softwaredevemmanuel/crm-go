package services

import (
    "errors"
    "time"
    
    "crm-go/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "fmt"
)

type TopicService struct {
    db *gorm.DB
}

func NewTopicService(db *gorm.DB) *TopicService {
    return &TopicService{db: db}
}

// CreateTopicWithTx - for use with transactions
func (s *TopicService) CreateTopicWithTx(tx *gorm.DB, req models.TopicInput) (*models.TopicResponse, error) {
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

// CreateTopic - regular version without transaction
func (s *TopicService) CreateTopic(req models.TopicInput) (*models.TopicResponse, error) {
    return s.CreateTopicWithTx(s.db, req)
}

// UpdateTopicWithTx - for use with transactions
func (s *TopicService) UpdateTopicWithTx(tx *gorm.DB, topicID uuid.UUID, req models.TopicInput) (*models.TopicResponse, error) {
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

// UpdateTopic - regular version without transaction
func (s *TopicService) UpdateTopic(topicID uuid.UUID, req models.TopicInput) (*models.TopicResponse, error) {
    return s.UpdateTopicWithTx(s.db, topicID, req)
}

// Helper function to convert Topic to TopicResponse
func (s *TopicService) topicToResponse(topic *models.Topic, tutorID uuid.UUID) *models.TopicResponse {
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




// GetAllTopics with filtering options
func (s *TopicService) GetAllTopics(filters models.TopicFilters) ([]models.TopicResponse, error) {
    var topics []models.Topic
    
    // Start building query
    query := s.db.Model(&models.Topic{})
    
    // Apply filters if provided
    if filters.CourseID != uuid.Nil {
        query = query.Where("course_id = ?", filters.CourseID)
    }
    
    if filters.ChapterID != uuid.Nil {
        query = query.Where("chapter_id = ?", filters.ChapterID)
    }
    
    if filters.TutorID != uuid.Nil {
        // Assuming topics don't have tutor_id directly
        // Join with chapters or courses if needed
        query = query.Joins("JOIN chapters ON chapters.id = topics.chapter_id").
            Where("chapters.tutor_id = ?", filters.TutorID)
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
    if err := query.Find(&topics).Error; err != nil {
        return nil, errors.New("failed to fetch topics: " + err.Error())
    }
    
    // Convert to response
    return s.topicsToResponse(topics, filters.TutorID), nil
}

// GetTopicsByChapterID - convenience method
func (s *TopicService) GetTopicsByChapterID(chapterID uuid.UUID, tutorID uuid.UUID) ([]models.TopicResponse, error) {
    filters := models.TopicFilters{
        ChapterID: chapterID,
        TutorID:   tutorID,
    }
    return s.GetAllTopics(filters)
}

// GetTopicsByCourseID - convenience method
func (s *TopicService) GetTopicsByCourseID(courseID uuid.UUID, tutorID uuid.UUID) ([]models.TopicResponse, error) {
    filters := models.TopicFilters{
        CourseID: courseID,
        TutorID:  tutorID,
    }
    return s.GetAllTopics(filters)
}

// Helper to convert slice of Topics to slice of TopicResponses
func (s *TopicService) topicsToResponse(topics []models.Topic, tutorID uuid.UUID) []models.TopicResponse {
    responses := make([]models.TopicResponse, len(topics))
    for i, topic := range topics {
        responses[i] = models.TopicResponse{
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
    return responses
}

// GetTopicCount - get total count with filters
func (s *TopicService) GetTopicCount(filters models.TopicFilters) (int64, error) {
    var count int64
    
    query := s.db.Model(&models.Topic{})
    
    // Apply the same filters as GetAllTopics
    if filters.CourseID != uuid.Nil {
        query = query.Where("course_id = ?", filters.CourseID)
    }
    
    if filters.ChapterID != uuid.Nil {
        query = query.Where("chapter_id = ?", filters.ChapterID)
    }
    
    if filters.TutorID != uuid.Nil {
        query = query.Joins("JOIN chapters ON chapters.id = topics.chapter_id").
            Where("chapters.tutor_id = ?", filters.TutorID)
    }
    
    if filters.Search != "" {
        searchTerm := "%" + filters.Search + "%"
        query = query.Where("title ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
    }
    
    if err := query.Count(&count).Error; err != nil {
        return 0, errors.New("failed to count topics: " + err.Error())
    }
    
    return count, nil
}

// GetAllTopicsWithPagination - returns topics with pagination metadata
func (s *TopicService) GetAllTopicsWithPagination(filters models.TopicFilters) (*models.PaginatedTopicsResponse, error) {
    // Get topics
    topics, err := s.GetAllTopics(filters)
    if err != nil {
        return nil, err
    }
    
    // Get total count
    totalCount, err := s.GetTopicCount(filters)
    if err != nil {
        return nil, err
    }
    
    // Calculate total pages
    totalPages := 0
    if filters.Limit > 0 {
        totalPages = int((totalCount + int64(filters.Limit) - 1) / int64(filters.Limit))
    }
    
    return &models.PaginatedTopicsResponse{
        Data:       topics,
        Total:      totalCount,
        Page:       filters.Page,
        Limit:      filters.Limit,
        TotalPages: totalPages,
    }, nil
}


// Optional: Get topic by ID
func (s *TopicService) GetTopicByID(topicID uuid.UUID, tutorID uuid.UUID) (*models.TopicResponse, error) {
    var topic models.Topic
    if err := s.db.First(&topic, "id = ?", topicID).Error; err != nil {
        return nil, errors.New("topic not found")
    }
    
    return s.topicToResponse(&topic, tutorID), nil
}