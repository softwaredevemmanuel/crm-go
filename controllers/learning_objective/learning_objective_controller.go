// handlers/learning_objective_handler.go
package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"crm-go/dto"
	"crm-go/services/learning_objective"
)

type LearningObjectiveHandler struct {
	objectiveService *services.LearningObjectiveService
}

func NewLearningObjectiveHandler(objectiveService *services.LearningObjectiveService) *LearningObjectiveHandler {
	return &LearningObjectiveHandler{
		objectiveService: objectiveService,
	}
}

// CreateLearningObjective handles the creation of a new learning objective
// @Summary Create a learning objective
// @Description Create a new learning objective for a scheme of work item
// @Tags Learning Objectives
// @Accept json
// @Produce json
// @Param request body dto.CreateLearningObjectiveRequest true "Objective request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/learning-objectives [post]
func (h *LearningObjectiveHandler) CreateLearningObjective(c *gin.Context) {
	var req dto.CreateLearningObjectiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	objective, err := h.objectiveService.CreateLearningObjective(&req)
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
		"message":   "Learning objective created successfully",
		"objective": objective,
	})
}

// BulkCreateLearningObjectives handles bulk creation of learning objectives
// @Summary Bulk create learning objectives
// @Description Create multiple learning objectives for a scheme of work item
// @Tags Learning Objectives
// @Accept json
// @Produce json
// @Param request body dto.BulkCreateLearningObjectivesRequest true "Bulk objective request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/learning-objectives/bulk [post]
func (h *LearningObjectiveHandler) BulkCreateLearningObjectives(c *gin.Context) {
	var req dto.BulkCreateLearningObjectivesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.objectiveService.BulkCreateLearningObjectives(&req)
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
		"message": "Learning objectives created successfully",
		"data":    result,
	})
}

// GetAllObjectives handles fetching all learning objectives
// @Summary Get all learning objectives
// @Description Get a paginated list of all learning objectives
// @Tags Learning Objectives
// @Accept json
// @Produce json
// @Param scheme_of_work_item_id query string false "Filter by scheme of work item ID"
// @Param search query string false "Search by objective text"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.LearningObjectiveListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/learning-objectives [get]
func (h *LearningObjectiveHandler) GetAllObjectives(c *gin.Context) {
	var params dto.LearningObjectiveQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.objectiveService.GetAllObjectives(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Learning objectives retrieved successfully",
		"data":    response,
	})
}

// GetObjectiveByID handles fetching a single objective by ID
// @Summary Get objective by ID
// @Description Get a single learning objective by its ID
// @Tags Learning Objectives
// @Accept json
// @Produce json
// @Param id path string true "Objective ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/learning-objectives/{id} [get]
func (h *LearningObjectiveHandler) GetObjectiveByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Objective ID is required",
		})
		return
	}

	objective, err := h.objectiveService.GetObjectiveByID(id)
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
		"message":   "Learning objective retrieved successfully",
		"objective": objective,
	})
}

// GetObjectivesBySchemeItem handles fetching all objectives for a scheme of work item
// @Summary Get objectives by scheme of work item
// @Description Get all learning objectives for a specific scheme of work item
// @Tags Learning Objectives
// @Accept json
// @Produce json
// @Param scheme_of_work_item_id path string true "Scheme of Work Item ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/learning-objectives/scheme-item/{scheme_of_work_item_id} [get]
func (h *LearningObjectiveHandler) GetObjectivesBySchemeItem(c *gin.Context) {
	schemeItemID := c.Param("scheme_of_work_item_id")
	if schemeItemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Scheme of work item ID is required",
		})
		return
	}

	objectives, err := h.objectiveService.GetObjectivesBySchemeItem(schemeItemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Learning objectives retrieved successfully",
		"objectives": objectives,
		"total":      len(objectives),
	})
}

// UpdateLearningObjective handles updating a learning objective
// @Summary Update a learning objective
// @Description Update an existing learning objective
// @Tags Learning Objectives
// @Accept json
// @Produce json
// @Param id path string true "Objective ID"
// @Param request body dto.UpdateLearningObjectiveRequest true "Update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/learning-objectives/{id} [put]
func (h *LearningObjectiveHandler) UpdateLearningObjective(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Objective ID is required",
		})
		return
	}

	var req dto.UpdateLearningObjectiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	objective, err := h.objectiveService.UpdateLearningObjective(id, &req)
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

	c.JSON(http.StatusOK, gin.H{
		"message":   "Learning objective updated successfully",
		"objective": objective,
	})
}

// DeleteLearningObjective handles deleting a learning objective (soft delete)
// @Summary Delete a learning objective
// @Description Soft delete a learning objective
// @Tags Learning Objectives
// @Accept json
// @Produce json
// @Param id path string true "Objective ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/learning-objectives/{id} [delete]
func (h *LearningObjectiveHandler) DeleteLearningObjective(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Objective ID is required",
		})
		return
	}

	err := h.objectiveService.DeleteLearningObjective(id)
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
		"message": "Learning objective deleted successfully",
	})
}