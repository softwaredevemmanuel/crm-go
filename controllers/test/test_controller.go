// handlers/test_handler.go
package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/test"
)

type TestHandler struct {
	testService *services.TestService
}

func NewTestHandler(testService *services.TestService) *TestHandler {
	return &TestHandler{
		testService: testService,
	}
}

// CreateTest handles the creation of a new test
// @Summary Create a test
// @Description Create a new test
// @Tags Tests
// @Accept json
// @Produce json
// @Param request body dto.CreateTestRequest true "Test request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/tests [post]
func (h *TestHandler) CreateTest(c *gin.Context) {
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

	var req dto.CreateTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	test, err := h.testService.CreateTest(&req, userID)
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
		"message": "Test created successfully",
		"test":    test,
	})
}

// BulkCreateTests handles bulk creation of tests
// @Summary Bulk create tests
// @Description Create multiple tests
// @Tags Tests
// @Accept json
// @Produce json
// @Param request body dto.BulkCreateTestsRequest true "Bulk test request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/tests/bulk [post]
func (h *TestHandler) BulkCreateTests(c *gin.Context) {
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

	var req dto.BulkCreateTestsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.testService.BulkCreateTests(&req, userID)
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
		"message": "Tests created successfully",
		"data":    result,
	})
}

// GetAllTests handles fetching all tests
// @Summary Get all tests
// @Description Get a paginated list of all tests
// @Tags Tests
// @Accept json
// @Produce json
// @Param academic_session_id query string false "Filter by academic session ID"
// @Param term_id query string false "Filter by term ID"
// @Param subject_id query string false "Filter by subject ID"
// @Param class_id query string false "Filter by class ID"
// @Param status query string false "Filter by status"
// @Param search query string false "Search by title or test type"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.TestListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/tests [get]
func (h *TestHandler) GetAllTests(c *gin.Context) {
	var params dto.TestQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.testService.GetAllTests(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tests retrieved successfully",
		"data":    response,
	})
}

// GetTestByID handles fetching a single test by ID
// @Summary Get test by ID
// @Description Get a single test by its ID
// @Tags Tests
// @Accept json
// @Produce json
// @Param id path string true "Test ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/tests/{id} [get]
func (h *TestHandler) GetTestByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Test ID is required",
		})
		return
	}

	test, err := h.testService.GetTestByID(id)
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
		"message": "Test retrieved successfully",
		"test":    test,
	})
}

// GetTestsBySubject handles fetching all tests for a subject
// @Summary Get tests by subject
// @Description Get all tests for a specific subject
// @Tags Tests
// @Accept json
// @Produce json
// @Param subject_id path string true "Subject ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/tests/subject/{subject_id} [get]
func (h *TestHandler) GetTestsBySubject(c *gin.Context) {
	subjectID := c.Param("subject_id")
	if subjectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Subject ID is required",
		})
		return
	}

	tests, err := h.testService.GetTestsBySubject(subjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tests retrieved successfully",
		"tests":   tests,
		"total":   len(tests),
	})
}

// GetTestsByClass handles fetching all tests for a class
// @Summary Get tests by class
// @Description Get all tests for a specific class
// @Tags Tests
// @Accept json
// @Produce json
// @Param class_id path string true "Class ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/tests/class/{class_id} [get]
func (h *TestHandler) GetTestsByClass(c *gin.Context) {
	classID := c.Param("class_id")
	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Class ID is required",
		})
		return
	}

	tests, err := h.testService.GetTestsByClass(classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tests retrieved successfully",
		"tests":   tests,
		"total":   len(tests),
	})
}

// UpdateTest handles updating a test
// @Summary Update a test
// @Description Update an existing test
// @Tags Tests
// @Accept json
// @Produce json
// @Param id path string true "Test ID"
// @Param request body dto.UpdateTestRequest true "Update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/tests/{id} [put]
func (h *TestHandler) UpdateTest(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Test ID is required",
		})
		return
	}

	var req dto.UpdateTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	test, err := h.testService.UpdateTest(id, &req)
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
		"message": "Test updated successfully",
		"test":    test,
	})
}

// DeleteTest handles deleting a test (soft delete)
// @Summary Delete a test
// @Description Soft delete a test
// @Tags Tests
// @Accept json
// @Produce json
// @Param id path string true "Test ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/tests/{id} [delete]
func (h *TestHandler) DeleteTest(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Test ID is required",
		})
		return
	}

	err := h.testService.DeleteTest(id)
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
		"message": "Test deleted successfully",
	})
}