// handlers/exam_scheme_item_handler.go
package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"crm-go/dto"
	"crm-go/services/exam_scheme_item"
)

type ExamSchemeItemHandler struct {
	examSchemeItemService *services.ExamSchemeItemService
}

func NewExamSchemeItemHandler(examSchemeItemService *services.ExamSchemeItemService) *ExamSchemeItemHandler {
	return &ExamSchemeItemHandler{
		examSchemeItemService: examSchemeItemService,
	}
}

// CreateExamSchemeItem handles the creation of a new exam scheme item
// @Summary Create an exam scheme item
// @Description Associate a scheme of work item with an exam
// @Tags Exam Scheme Items
// @Accept json
// @Produce json
// @Param request body dto.CreateExamSchemeItemRequest true "Exam scheme item request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/exam-scheme-items [post]
func (h *ExamSchemeItemHandler) CreateExamSchemeItem(c *gin.Context) {
	var req dto.CreateExamSchemeItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	item, err := h.examSchemeItemService.CreateExamSchemeItem(&req)
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
		"message": "Exam scheme item created successfully",
		"item":    item,
	})
}

// BulkCreateExamSchemeItems handles bulk creation of exam scheme items
// @Summary Bulk create exam scheme items
// @Description Associate multiple scheme of work items with an exam
// @Tags Exam Scheme Items
// @Accept json
// @Produce json
// @Param request body dto.BulkCreateExamSchemeItemsRequest true "Bulk exam scheme item request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/exam-scheme-items/bulk [post]
func (h *ExamSchemeItemHandler) BulkCreateExamSchemeItems(c *gin.Context) {
	var req dto.BulkCreateExamSchemeItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.examSchemeItemService.BulkCreateExamSchemeItems(&req)
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
		"message": "Exam scheme items created successfully",
		"data":    result,
	})
}

// GetAllExamSchemeItems handles fetching all exam scheme items
// @Summary Get all exam scheme items
// @Description Get a paginated list of all exam scheme items
// @Tags Exam Scheme Items
// @Accept json
// @Produce json
// @Param exam_id query string false "Filter by exam ID"
// @Param scheme_of_work_item_id query string false "Filter by scheme of work item ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.ExamSchemeItemListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/exam-scheme-items [get]
func (h *ExamSchemeItemHandler) GetAllExamSchemeItems(c *gin.Context) {
	var params dto.ExamSchemeItemQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.examSchemeItemService.GetAllExamSchemeItems(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Exam scheme items retrieved successfully",
		"data":    response,
	})
}

// GetExamSchemeItemsByExam handles fetching all scheme items for an exam
// @Summary Get exam scheme items by exam
// @Description Get all scheme of work items associated with a specific exam
// @Tags Exam Scheme Items
// @Accept json
// @Produce json
// @Param exam_id path string true "Exam ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/exam-scheme-items/exam/{exam_id} [get]
func (h *ExamSchemeItemHandler) GetExamSchemeItemsByExam(c *gin.Context) {
	examID := c.Param("exam_id")
	if examID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Exam ID is required",
		})
		return
	}

	items, err := h.examSchemeItemService.GetExamSchemeItemsByExam(examID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Exam scheme items retrieved successfully",
		"items":   items,
		"total":   len(items),
	})
}

// GetExamSchemeItemsBySchemeItem handles fetching all exams for a scheme item
// @Summary Get exam scheme items by scheme of work item
// @Description Get all exams associated with a specific scheme of work item
// @Tags Exam Scheme Items
// @Accept json
// @Produce json
// @Param scheme_of_work_item_id path string true "Scheme of Work Item ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/exam-scheme-items/scheme-item/{scheme_of_work_item_id} [get]
func (h *ExamSchemeItemHandler) GetExamSchemeItemsBySchemeItem(c *gin.Context) {
	schemeItemID := c.Param("scheme_of_work_item_id")
	if schemeItemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Scheme of work item ID is required",
		})
		return
	}

	items, err := h.examSchemeItemService.GetExamSchemeItemsBySchemeItem(schemeItemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Exam scheme items retrieved successfully",
		"items":   items,
		"total":   len(items),
	})
}

// DeleteExamSchemeItem handles deleting an exam scheme item
// @Summary Delete an exam scheme item
// @Description Remove an association between an exam and a scheme of work item
// @Tags Exam Scheme Items
// @Accept json
// @Produce json
// @Param exam_id path string true "Exam ID"
// @Param scheme_of_work_item_id path string true "Scheme of Work Item ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/exam-scheme-items/{exam_id}/{scheme_of_work_item_id} [delete]
func (h *ExamSchemeItemHandler) DeleteExamSchemeItem(c *gin.Context) {
	examID := c.Param("exam_id")
	schemeItemID := c.Param("scheme_of_work_item_id")

	if examID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Exam ID is required",
		})
		return
	}
	if schemeItemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Scheme of work item ID is required",
		})
		return
	}

	err := h.examSchemeItemService.DeleteExamSchemeItem(examID, schemeItemID)
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
		"message": "Exam scheme item deleted successfully",
	})
}

// DeleteAllExamSchemeItemsByExam handles deleting all scheme items for an exam
// @Summary Delete all exam scheme items for an exam
// @Description Remove all associations between an exam and scheme of work items
// @Tags Exam Scheme Items
// @Accept json
// @Produce json
// @Param exam_id path string true "Exam ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/exam-scheme-items/exam/{exam_id} [delete]
func (h *ExamSchemeItemHandler) DeleteAllExamSchemeItemsByExam(c *gin.Context) {
	examID := c.Param("exam_id")
	if examID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Exam ID is required",
		})
		return
	}

	err := h.examSchemeItemService.DeleteAllExamSchemeItemsByExam(examID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All exam scheme items deleted successfully",
	})
}