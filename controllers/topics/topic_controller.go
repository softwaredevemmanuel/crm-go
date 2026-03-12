package topic

import (
	"crm-go/config"
	"crm-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	// "gorm.io/gorm"
)

// CreateTopic godoc
// @Summary      Create a topic
// @Description  Create a new topic under a module
// @Tags         Topics
// @Accept       json
// @Produce      json
// @Param        topic body models.TopicInput true "Topic payload"
// @Success      201 {object} models.SuccessResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      409 {object} models.ConflictResponse
// @Failure      500 {object} models.FailureResponse
// @Router       /api/topics [post]
// @Security BearerAuth
func CreateTopic(c *gin.Context) {
	var input models.TopicInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// 🔒 Duplicate check
	var existing models.Topics
	if err := config.DB.
		Where("course_id = ? AND module_id = ? AND title = ?",
			input.CourseID, input.ModuleID, input.Title).
		First(&existing).Error; err == nil {

		c.JSON(http.StatusConflict, models.ErrorResponse{
			Error: "Topic already exists for this module",
		})
		return
	}

	topic := models.Topics{
		ID:          uuid.New(),
		CourseID:    input.CourseID,
		ModuleID:    input.ModuleID,
		Title:       input.Title,
		ContentType: input.ContentType,
		ContentURL:  input.ContentURL,
	}

	if err := config.DB.Create(&topic).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Topic created successfully",
	})
}

// GetTopics godoc
// @Summary      Get topics
// @Description  Get all topics, optionally filtered by course or module
// @Tags         Topics
// @Produce      json
// @Param        course_id   query   string  false  "Course ID"
// @Param        module_id  query   string  false  "Module ID"
// @Success      200 {array}  models.TopicResponse
// @Failure      400 {object} models.ErrorResponse
// @Router       /topics [get]
func GetAllTopics(c *gin.Context) {
	var topics []models.Topics

	query := config.DB

	if courseID := c.Query("course_id"); courseID != "" {
		if uid, err := uuid.Parse(courseID); err == nil {
			query = query.Where("course_id = ?", uid)
		}
	}

	if moduleID := c.Query("module_id"); moduleID != "" {
		if uid, err := uuid.Parse(moduleID); err == nil {
			query = query.Where("module_id = ?", uid)
		}
	}

	if err := query.
		Order("created_at DESC").
		Find(&topics).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch topics",
		})
		return
	}

	responses := make([]models.TopicResponse, 0, len(topics))
	for _, topic := range topics {
		responses = append(responses, models.TopicResponse{
			ID:          topic.ID,
			ModuleID:    topic.ModuleID,
			CourseID:    topic.CourseID,
			Title:       topic.Title,
			ContentType: topic.ContentType,
			ContentURL:  topic.ContentURL,
			CreatedAt:   topic.CreatedAt,
			UpdatedAt:   topic.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, responses)
}

// GetTopicByID godoc
// @Summary      Get topic details
// @Description  Get topic with module and course details
// @Tags         Topics
// @Produce      json
// @Param        id path string true "Topic ID"
// @Success      200 {object} models.TopicResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      404 {object} models.NotFoundResponse
// @Router       /topics/{id} [get]
func GetTopicByID(c *gin.Context) {
	id := c.Param("id")

	topicID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid topic ID",
		})
		return
	}

	var topic models.Topics

	// Fetch topic with relations
	if err := config.DB.
		Preload("Module").
		Preload("Module.Course").
		First(&topic, "id = ?", topicID).
		Error; err != nil {

		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Topic not found",
		})
		return
	}

	// Map to DTO
	response := models.TopicViewResponse{
		ID:          topic.ID,
		ModuleID:    topic.ModuleID,
		CourseID:    topic.CourseID,
		Title:       topic.Title,
		ContentType: topic.ContentType,
		ContentURL:  topic.ContentURL,
		CreatedAt:   topic.CreatedAt,
		UpdatedAt:   topic.UpdatedAt,
		Course: models.CourseMiniResponse{
			ID:    topic.Module.Course.ID,
			Title: topic.Module.Course.Title,
		},
		Module: &models.ModuleMiniResponse{
			ID:           topic.Module.ID,
			Title:        topic.Module.Title,
			ModuleNumber: topic.Module.ModuleNumber,
		},
	}

	c.JSON(http.StatusOK, response)
}

// UpdateTopic godoc
// @Summary      Update topic
// @Description  Update topic details
// @Tags         Topics
// @Accept       json
// @Produce      json
// @Param        id path string true "Topic ID"
// @Param        topic body models.TopicInput true "Topic payload"
// @Success      200 {object} models.TopicResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      404 {object} models.NotFoundResponse
// @Failure      500 {object} models.FailureResponse
// @Router       /api/topics/{id} [put]
// @Security BearerAuth
func UpdateTopic(c *gin.Context) {
	id := c.Param("id")

	topicID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid topic ID",
		})
		return
	}

	var topic models.Topics

	// 1️⃣ Check if topic exists
	if err := config.DB.First(&topic, "id = ?", topicID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Topic not found",
		})
		return
	}

	// 2️⃣ Bind update payload
	var input models.TopicUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// 3️⃣ Apply updates safely
	if err := config.DB.Model(&topic).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update topic",
		})
		return
	}

	// 4️⃣ Reload topic with relations
	if err := config.DB.
		Preload("Module").
		Preload("Module.Course").
		First(&topic, "id = ?", topicID).Error; err != nil {

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to load updated topic",
		})
		return
	}

	// 5️⃣ Map to TopicViewResponse DTO
	response := models.TopicViewResponse{
		ID:          topic.ID,
		ModuleID:    topic.ModuleID,
		CourseID:    topic.CourseID,
		Title:       topic.Title,
		ContentType: topic.ContentType,
		ContentURL:  topic.ContentURL,
		CreatedAt:   topic.CreatedAt,
		UpdatedAt:   topic.UpdatedAt,

		Course: models.CourseMiniResponse{
			ID:    topic.Module.Course.ID,
			Title: topic.Module.Course.Title,
		},

		Module: &models.ModuleMiniResponse{
			ID:           topic.Module.ID,
			Title:        topic.Module.Title,
			ModuleNumber: topic.Module.ModuleNumber,
		},
	}

	// 6️⃣ Return clean DTO
	c.JSON(http.StatusOK, response)
}

// DeleteTopic godoc
// @Summary      Delete topic
// @Description  Delete topic by ID
// @Tags         Topics
// @Produce      json
// @Param        id path string true "Topic ID"
// @Success      200 {object} models.SuccessResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      404 {object} models.NotFoundResponse
// @Failure      500 {object} models.FailureResponse
// @Router       /api/topics/{id} [delete]
// @Security BearerAuth
func DeleteTopic(c *gin.Context) {
	id := c.Param("id")

	topicID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid topic ID",
		})
		return
	}

	result := config.DB.Delete(&models.Topics{}, "id = ?", topicID)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to delete topic",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Topic not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Topic deleted successfully",
	})
}
