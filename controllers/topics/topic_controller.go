// handlers/topic_handler.go
package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"crm-go/dto"
	"crm-go/services/topics"
)

type TopicHandler struct {
	topicService *services.TopicService
}

func NewTopicHandler(topicService *services.TopicService) *TopicHandler {
	return &TopicHandler{
		topicService: topicService,
	}
}

// CreateTopic handles the creation of a new topic
// @Summary Create a topic
// @Description Create a new topic for a module
// @Tags Topics
// @Accept json
// @Produce json
// @Param request body dto.CreateTopicRequest true "Topic request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/topics [post]
func (h *TopicHandler) CreateTopic(c *gin.Context) {
	var req dto.CreateTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	topic, err := h.topicService.CreateTopic(&req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Topic created successfully",
		"topic":   topic,
	})
}

// BulkCreateTopics handles bulk creation of topics
// @Summary Bulk create topics
// @Description Create multiple topics for a module
// @Tags Topics
// @Accept json
// @Produce json
// @Param request body dto.BulkCreateTopicsRequest true "Bulk topic request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/topics/bulk [post]
func (h *TopicHandler) BulkCreateTopics(c *gin.Context) {
	var req dto.BulkCreateTopicsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.topicService.BulkCreateTopics(&req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Topics created successfully",
		"data":    result,
	})
}

// GetAllTopics handles fetching all topics
// @Summary Get all topics
// @Description Get a paginated list of all topics
// @Tags Topics
// @Accept json
// @Produce json
// @Param module_id query string false "Filter by module ID"
// @Param search query string false "Search by title or description"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.TopicListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/topics [get]
func (h *TopicHandler) GetAllTopics(c *gin.Context) {
	var params dto.TopicQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.topicService.GetAllTopics(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Topics retrieved successfully",
		"data":    response,
	})
}

// GetTopicByID handles fetching a single topic by ID
// @Summary Get topic by ID
// @Description Get a single topic by its ID
// @Tags Topics
// @Accept json
// @Produce json
// @Param id path string true "Topic ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/topics/{id} [get]
func (h *TopicHandler) GetTopicByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Topic ID is required",
		})
		return
	}

	topic, err := h.topicService.GetTopicByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Topic retrieved successfully",
		"topic":   topic,
	})
}

// GetTopicsByModule handles fetching all topics for a module
// @Summary Get topics by module
// @Description Get all topics for a specific module
// @Tags Topics
// @Accept json
// @Produce json
// @Param module_id path string true "Module ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/topics/module/{module_id} [get]
func (h *TopicHandler) GetTopicsByModule(c *gin.Context) {
	moduleID := c.Param("module_id")
	if moduleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Module ID is required",
		})
		return
	}

	topics, err := h.topicService.GetTopicsByModule(moduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Topics retrieved successfully",
		"topics":  topics,
		"total":   len(topics),
	})
}

// UpdateTopic handles updating a topic
// @Summary Update a topic
// @Description Update an existing topic
// @Tags Topics
// @Accept json
// @Produce json
// @Param id path string true "Topic ID"
// @Param request body dto.UpdateTopicRequest true "Update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/topics/{id} [put]
func (h *TopicHandler) UpdateTopic(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Topic ID is required",
		})
		return
	}

	var req dto.UpdateTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	topic, err := h.topicService.UpdateTopic(id, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Topic updated successfully",
		"topic":   topic,
	})
}

// ReorderTopics handles reordering topics
// @Summary Reorder topics
// @Description Reorder topics within a module
// @Tags Topics
// @Accept json
// @Produce json
// @Param request body dto.ReorderTopicsRequest true "Reorder request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/topics/reorder [put]
func (h *TopicHandler) ReorderTopics(c *gin.Context) {
	var req dto.ReorderTopicsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := h.topicService.ReorderTopics(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Topics reordered successfully",
	})
}

// DeleteTopic handles deleting a topic (soft delete)
// @Summary Delete a topic
// @Description Soft delete a topic
// @Tags Topics
// @Accept json
// @Produce json
// @Param id path string true "Topic ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/topics/{id} [delete]
func (h *TopicHandler) DeleteTopic(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Topic ID is required",
		})
		return
	}

	err := h.topicService.DeleteTopic(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Topic deleted successfully",
	})
}