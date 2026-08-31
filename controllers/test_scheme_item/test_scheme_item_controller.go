// handlers/test_scheme_item_handler.go
package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"crm-go/dto"
	"crm-go/services/test_scheme_item"
)

type TestSchemeItemHandler struct {
	testSchemeItemService *services.TestSchemeItemService
}

func NewTestSchemeItemHandler(testSchemeItemService *services.TestSchemeItemService) *TestSchemeItemHandler {
	return &TestSchemeItemHandler{
		testSchemeItemService: testSchemeItemService,
	}
}

// CreateTestSchemeItem handles the creation of a new test scheme item
// @Summary Create a test scheme item
// @Description Associate a scheme of work item with a test
// @Tags Test Scheme Items
// @Accept json
// @Produce json
// @Param request body dto.CreateTestSchemeItemRequest true "Test scheme item request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/test-scheme-items [post]
func (h *TestSchemeItemHandler) CreateTestSchemeItem(c *gin.Context) {
	var req dto.CreateTestSchemeItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	item, err := h.testSchemeItemService.CreateTestSchemeItem(&req)
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
		"message": "Test scheme item created successfully",
		"item":    item,
	})
}

// BulkCreateTestSchemeItems handles bulk creation of test scheme items
// @Summary Bulk create test scheme items
// @Description Associate multiple scheme of work items with a test
// @Tags Test Scheme Items
// @Accept json
// @Produce json
// @Param request body dto.BulkCreateTestSchemeItemsRequest true "Bulk test scheme item request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/test-scheme-items/bulk [post]
func (h *TestSchemeItemHandler) BulkCreateTestSchemeItems(c *gin.Context) {
	var req dto.BulkCreateTestSchemeItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.testSchemeItemService.BulkCreateTestSchemeItems(&req)
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
		"message": "Test scheme items created successfully",
		"data":    result,
	})
}

// GetAllTestSchemeItems handles fetching all test scheme items
// @Summary Get all test scheme items
// @Description Get a paginated list of all test scheme items
// @Tags Test Scheme Items
// @Accept json
// @Produce json
// @Param test_id query string false "Filter by test ID"
// @Param scheme_of_work_item_id query string false "Filter by scheme of work item ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.TestSchemeItemListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/test-scheme-items [get]
func (h *TestSchemeItemHandler) GetAllTestSchemeItems(c *gin.Context) {
	var params dto.TestSchemeItemQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.testSchemeItemService.GetAllTestSchemeItems(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Test scheme items retrieved successfully",
		"data":    response,
	})
}

// GetTestSchemeItemsByTest handles fetching all scheme items for a test
// @Summary Get test scheme items by test
// @Description Get all scheme of work items associated with a specific test
// @Tags Test Scheme Items
// @Accept json
// @Produce json
// @Param test_id path string true "Test ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/test-scheme-items/test/{test_id} [get]
func (h *TestSchemeItemHandler) GetTestSchemeItemsByTest(c *gin.Context) {
	testID := c.Param("test_id")
	if testID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Test ID is required",
		})
		return
	}

	items, err := h.testSchemeItemService.GetTestSchemeItemsByTest(testID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Test scheme items retrieved successfully",
		"items":   items,
		"total":   len(items),
	})
}

// GetTestSchemeItemsBySchemeItem handles fetching all tests for a scheme item
// @Summary Get test scheme items by scheme of work item
// @Description Get all tests associated with a specific scheme of work item
// @Tags Test Scheme Items
// @Accept json
// @Produce json
// @Param scheme_of_work_item_id path string true "Scheme of Work Item ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/test-scheme-items/scheme-item/{scheme_of_work_item_id} [get]
func (h *TestSchemeItemHandler) GetTestSchemeItemsBySchemeItem(c *gin.Context) {
	schemeItemID := c.Param("scheme_of_work_item_id")
	if schemeItemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Scheme of work item ID is required",
		})
		return
	}

	items, err := h.testSchemeItemService.GetTestSchemeItemsBySchemeItem(schemeItemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Test scheme items retrieved successfully",
		"items":   items,
		"total":   len(items),
	})
}

// DeleteTestSchemeItem handles deleting a test scheme item
// @Summary Delete a test scheme item
// @Description Remove an association between a test and a scheme of work item
// @Tags Test Scheme Items
// @Accept json
// @Produce json
// @Param test_id path string true "Test ID"
// @Param scheme_of_work_item_id path string true "Scheme of Work Item ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/test-scheme-items/{test_id}/{scheme_of_work_item_id} [delete]
func (h *TestSchemeItemHandler) DeleteTestSchemeItem(c *gin.Context) {
	testID := c.Param("test_id")
	schemeItemID := c.Param("scheme_of_work_item_id")

	if testID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Test ID is required",
		})
		return
	}
	if schemeItemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Scheme of work item ID is required",
		})
		return
	}

	err := h.testSchemeItemService.DeleteTestSchemeItem(testID, schemeItemID)
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
		"message": "Test scheme item deleted successfully",
	})
}

// DeleteAllTestSchemeItemsByTest handles deleting all scheme items for a test
// @Summary Delete all test scheme items for a test
// @Description Remove all associations between a test and scheme of work items
// @Tags Test Scheme Items
// @Accept json
// @Produce json
// @Param test_id path string true "Test ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/test-scheme-items/test/{test_id} [delete]
func (h *TestSchemeItemHandler) DeleteAllTestSchemeItemsByTest(c *gin.Context) {
	testID := c.Param("test_id")
	if testID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Test ID is required",
		})
		return
	}

	err := h.testSchemeItemService.DeleteAllTestSchemeItemsByTest(testID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All test scheme items deleted successfully",
	})
}