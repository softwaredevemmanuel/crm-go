package controllers

import (
    "net/http"
    "strconv"
    
    "crm-go/models"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)



// GetAllTopics retrieves all topics with optional filtering
// @Summary Get all topics
// @Description Get all topics with optional filtering
// @Tags topics
// @Accept json
// @Produce json
// @Param course_id query string false "Course ID"
// @Param chapter_id query string false "Chapter ID"
// @Param search query string false "Search term"
// @Param sort_by query string false "Sort field (title, order, created_at, updated_at)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedTopicsResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /topics [get]
// @Security BearerAuth
func (ctl *TopicController) GetAllTopics(c *gin.Context) {
    // Parse query parameters
    var filters models.TopicFilters
    
    // Parse UUID parameters
    if courseIDStr := c.Query("course_id"); courseIDStr != "" {
        if courseID, err := uuid.Parse(courseIDStr); err == nil {
            filters.CourseID = courseID
        }
    }
    
    if chapterIDStr := c.Query("chapter_id"); chapterIDStr != "" {
        if chapterID, err := uuid.Parse(chapterIDStr); err == nil {
            filters.ChapterID = chapterID
        }
    }
    
    // Parse other parameters
    filters.Search = c.Query("search")
    filters.SortBy = c.Query("sort_by")
    filters.SortOrder = c.Query("sort_order")
    
    // Parse pagination parameters
    if pageStr := c.Query("page"); pageStr != "" {
        if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
            filters.Page = page
        }
    }
    
    if limitStr := c.Query("limit"); limitStr != "" {
        if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
            filters.Limit = limit
        }
    }
    
    // Get tutor ID from context (assuming it's set by authentication middleware)
    if tutorID, exists := c.Get("tutor_id"); exists {
        if id, ok := tutorID.(uuid.UUID); ok {
            filters.TutorID = id
        }
    }
    
    // Get paginated results
    result, err := ctl.topicService.GetAllTopicsWithPagination(filters)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, result)
}

// GetTopicsByChapter retrieves topics by chapter ID
// @Summary Get topics by chapter
// @Description Get all topics for a specific chapter
// @Tags topics
// @Accept json
// @Produce json
// @Param chapter_id path string true "Chapter ID"
// @Success 200 {array} models.TopicResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/chapters/{chapter_id}/topics [get]
// @Security BearerAuth
func (ctl *TopicController) GetTopicsByChapter(c *gin.Context) {
    chapterID, err := uuid.Parse(c.Param("chapter_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid chapter ID",
        })
        return
    }
    
    // Get tutor ID from context
    tutorID := uuid.Nil
    if tutorIDVal, exists := c.Get("tutor_id"); exists {
        if id, ok := tutorIDVal.(uuid.UUID); ok {
            tutorID = id
        }
    }
    
    topics, err := ctl.topicService.GetTopicsByChapterID(chapterID, tutorID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "data": topics,
    })
}

// GetTopicsByCourse retrieves topics by course ID
// @Summary Get topics by course
// @Description Get all topics for a specific course
// @Tags topics
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Success 200 {array} models.TopicResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/courses/{course_id}/topics [get]
// @Security BearerAuth
func (ctl *TopicController) GetTopicsByCourse(c *gin.Context) {
    courseID, err := uuid.Parse(c.Param("course_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid course ID",
        })
        return
    }
    
    // Get tutor ID from context
    tutorID := uuid.Nil
    if tutorIDVal, exists := c.Get("tutor_id"); exists {
        if id, ok := tutorIDVal.(uuid.UUID); ok {
            tutorID = id
        }
    }
    
    topics, err := ctl.topicService.GetTopicsByCourseID(courseID, tutorID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "data": topics,
    })
}