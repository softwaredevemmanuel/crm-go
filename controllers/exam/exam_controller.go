// handlers/exam_handler.go
package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/exam"
)

type ExamHandler struct {
	examService *services.ExamService
}

func NewExamHandler(examService *services.ExamService) *ExamHandler {
	return &ExamHandler{
		examService: examService,
	}
}

// CreateExam handles the creation of a new exam
// @Summary Create an exam
// @Description Create a new exam
// @Tags Exams
// @Accept json
// @Produce json
// @Param request body dto.CreateExamRequest true "Exam request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/exams [post]
func (h *ExamHandler) CreateExam(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: user ID not found",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req dto.CreateExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	exam, err := h.examService.CreateExam(&req, userID)
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
		"message": "Exam created successfully",
		"exam":    exam,
	})
}

// BulkCreateExams handles bulk creation of exams
// @Summary Bulk create exams
// @Description Create multiple exams
// @Tags Exams
// @Accept json
// @Produce json
// @Param request body dto.BulkCreateExamsRequest true "Bulk exam request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/exams/bulk [post]
func (h *ExamHandler) BulkCreateExams(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: user ID not found",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req dto.BulkCreateExamsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.examService.BulkCreateExams(&req, userID)
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
		"message": "Exams created successfully",
		"data":    result,
	})
}

// GetAllExams handles fetching all exams
// @Summary Get all exams
// @Description Get a paginated list of all exams
// @Tags Exams
// @Accept json
// @Produce json
// @Param academic_session_id query string false "Filter by academic session ID"
// @Param term_id query string false "Filter by term ID"
// @Param subject_id query string false "Filter by subject ID"
// @Param class_id query string false "Filter by class ID"
// @Param status query string false "Filter by status"
// @Param search query string false "Search by title or exam type"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.ExamListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/exams [get]
func (h *ExamHandler) GetAllExams(c *gin.Context) {
	var params dto.ExamQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.examService.GetAllExams(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Exams retrieved successfully",
		"data":    response,
	})
}

// GetExamByID handles fetching a single exam by ID
// @Summary Get exam by ID
// @Description Get a single exam by its ID
// @Tags Exams
// @Accept json
// @Produce json
// @Param id path string true "Exam ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/exams/{id} [get]
func (h *ExamHandler) GetExamByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Exam ID is required",
		})
		return
	}

	exam, err := h.examService.GetExamByID(id)
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
		"message": "Exam retrieved successfully",
		"exam":    exam,
	})
}

// GetExamsBySubject handles fetching all exams for a subject
// @Summary Get exams by subject
// @Description Get all exams for a specific subject
// @Tags Exams
// @Accept json
// @Produce json
// @Param subject_id path string true "Subject ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/exams/subject/{subject_id} [get]
func (h *ExamHandler) GetExamsBySubject(c *gin.Context) {
	subjectID := c.Param("subject_id")
	if subjectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Subject ID is required",
		})
		return
	}

	exams, err := h.examService.GetExamsBySubject(subjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Exams retrieved successfully",
		"exams":   exams,
		"total":   len(exams),
	})
}

// GetExamsByClass handles fetching all exams for a class
// @Summary Get exams by class
// @Description Get all exams for a specific class
// @Tags Exams
// @Accept json
// @Produce json
// @Param class_id path string true "Class ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/exams/class/{class_id} [get]
func (h *ExamHandler) GetExamsByClass(c *gin.Context) {
	classID := c.Param("class_id")
	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Class ID is required",
		})
		return
	}

	exams, err := h.examService.GetExamsByClass(classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Exams retrieved successfully",
		"exams":   exams,
		"total":   len(exams),
	})
}

// UpdateExam handles updating an exam
// @Summary Update an exam
// @Description Update an existing exam
// @Tags Exams
// @Accept json
// @Produce json
// @Param id path string true "Exam ID"
// @Param request body dto.UpdateExamRequest true "Update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/exams/{id} [put]
func (h *ExamHandler) UpdateExam(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Exam ID is required",
		})
		return
	}

	var req dto.UpdateExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	exam, err := h.examService.UpdateExam(id, &req)
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
		"message": "Exam updated successfully",
		"exam":    exam,
	})
}

// DeleteExam handles deleting an exam (soft delete)
// @Summary Delete an exam
// @Description Soft delete an exam
// @Tags Exams
// @Accept json
// @Produce json
// @Param id path string true "Exam ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/exams/{id} [delete]
func (h *ExamHandler) DeleteExam(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Exam ID is required",
		})
		return
	}

	err := h.examService.DeleteExam(id)
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
		"message": "Exam deleted successfully",
	})
}