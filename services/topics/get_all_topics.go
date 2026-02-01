package services

import (
    "errors"
    "fmt"
    
    "crm-go/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type GetTopicService struct {
    db *gorm.DB
}

func NewGetTopicService(db *gorm.DB) *GetTopicService {
    return &GetTopicService{db: db}
}

// GetAllTopics with filtering options
func (s *GetTopicService) GetAllTopics(filters models.TopicFilters) ([]models.TopicResponse, error) {
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
func (s *GetTopicService) GetTopicsByChapterID(chapterID uuid.UUID, tutorID uuid.UUID) ([]models.TopicResponse, error) {
    filters := models.TopicFilters{
        ChapterID: chapterID,
        TutorID:   tutorID,
    }
    return s.GetAllTopics(filters)
}

// GetTopicsByCourseID - convenience method
func (s *GetTopicService) GetTopicsByCourseID(courseID uuid.UUID, tutorID uuid.UUID) ([]models.TopicResponse, error) {
    filters := models.TopicFilters{
        CourseID: courseID,
        TutorID:  tutorID,
    }
    return s.GetAllTopics(filters)
}

// Helper to convert slice of Topics to slice of TopicResponses
func (s *GetTopicService) topicsToResponse(topics []models.Topic, tutorID uuid.UUID) []models.TopicResponse {
    responses := make([]models.TopicResponse, len(topics))
    for i, topic := range topics {
        responses[i] = models.TopicResponse{
            ID:          topic.ID,
            CourseID:    topic.CourseID,
            ChapterID:   topic.ChapterID,
            TutorID:     topic.TutorID,
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
func (s *GetTopicService) GetTopicCount(filters models.TopicFilters) (int64, error) {
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
func (s *GetTopicService) GetAllTopicsWithPagination(filters models.TopicFilters) (*models.PaginatedTopicsResponse, error) {
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