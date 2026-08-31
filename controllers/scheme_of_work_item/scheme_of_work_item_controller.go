// handlers/scheme_of_work_item_handler.go
package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"crm-go/dto"
	"crm-go/services/scheme_of_work_item"
)

type SchemeOfWorkItemHandler struct {
	itemService *services.SchemeOfWorkItemService
}

func NewSchemeOfWorkItemHandler(itemService *services.SchemeOfWorkItemService) *SchemeOfWorkItemHandler {
	return &SchemeOfWorkItemHandler{
		itemService: itemService,
	}
}

// CreateSchemeOfWorkItem handles the creation of a new scheme of work item
// @Summary Create a scheme of work item
// @Description Create a new scheme of work item
// @Tags Scheme of Work Items
// @Accept json
// @Produce json
// @Param request body dto.CreateSchemeOfWorkItemRequest true "Item request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/scheme-items [post]
func (h *SchemeOfWorkItemHandler) CreateSchemeOfWorkItem(c *gin.Context) {
	var req dto.CreateSchemeOfWorkItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	item, err := h.itemService.CreateSchemeOfWorkItem(&req)
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
		"message": "Scheme of work item created successfully",
		"item":    item,
	})
}

// BulkCreateSchemeItems handles bulk creation of scheme of work items
// @Summary Bulk create scheme of work items
// @Description Create multiple scheme of work items
// @Tags Scheme of Work Items
// @Accept json
// @Produce json
// @Param request body dto.BulkCreateSchemeItemsRequest true "Bulk item request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/scheme-items/bulk [post]
func (h *SchemeOfWorkItemHandler) BulkCreateSchemeItems(c *gin.Context) {
	var req dto.BulkCreateSchemeItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.itemService.BulkCreateSchemeItems(&req)
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
		"message": "Scheme items created successfully",
		"data":    result,
	})
}

// GetAllItems handles fetching all scheme of work items
// @Summary Get all scheme of work items
// @Description Get a paginated list of all scheme of work items
// @Tags Scheme of Work Items
// @Accept json
// @Produce json
// @Param scheme_of_work_id query string false "Filter by scheme of work ID"
// @Param module_id query string false "Filter by module ID"
// @Param search query string false "Search by topic, subtopic, or content"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.SchemeOfWorkItemListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/scheme-items [get]
func (h *SchemeOfWorkItemHandler) GetAllItems(c *gin.Context) {
	var params dto.SchemeOfWorkItemQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.itemService.GetAllItems(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Items retrieved successfully",
		"data":    response,
	})
}

// GetItemByID handles fetching a single item by ID
// @Summary Get item by ID
// @Description Get a single scheme of work item by its ID
// @Tags Scheme of Work Items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/scheme-items/{id} [get]
func (h *SchemeOfWorkItemHandler) GetItemByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Item ID is required",
		})
		return
	}

	item, err := h.itemService.GetItemByID(id)
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
		"message": "Item retrieved successfully",
		"item":    item,
	})
}

// GetItemsBySchemeOfWork handles fetching all items for a scheme of work
// @Summary Get items by scheme of work
// @Description Get all items for a specific scheme of work
// @Tags Scheme of Work Items
// @Accept json
// @Produce json
// @Param scheme_of_work_id path string true "Scheme of Work ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/scheme-items/scheme/{scheme_of_work_id} [get]
func (h *SchemeOfWorkItemHandler) GetItemsBySchemeOfWork(c *gin.Context) {
	schemeOfWorkID := c.Param("scheme_of_work_id")
	if schemeOfWorkID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Scheme of work ID is required",
		})
		return
	}

	items, err := h.itemService.GetItemsBySchemeOfWork(schemeOfWorkID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Items retrieved successfully",
		"items":   items,
		"total":   len(items),
	})
}

// GetItemsByModule handles fetching all items for a module
// @Summary Get items by module
// @Description Get all items for a specific module
// @Tags Scheme of Work Items
// @Accept json
// @Produce json
// @Param module_id path string true "Module ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/scheme-items/module/{module_id} [get]
func (h *SchemeOfWorkItemHandler) GetItemsByModule(c *gin.Context) {
	moduleID := c.Param("module_id")
	if moduleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Module ID is required",
		})
		return
	}

	items, err := h.itemService.GetItemsByModule(moduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Items retrieved successfully",
		"items":   items,
		"total":   len(items),
	})
}

// UpdateSchemeOfWorkItem handles updating an item
// @Summary Update an item
// @Description Update an existing scheme of work item
// @Tags Scheme of Work Items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param request body dto.UpdateSchemeOfWorkItemRequest true "Update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/scheme-items/{id} [put]
func (h *SchemeOfWorkItemHandler) UpdateSchemeOfWorkItem(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Item ID is required",
		})
		return
	}

	var req dto.UpdateSchemeOfWorkItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	item, err := h.itemService.UpdateSchemeOfWorkItem(id, &req)
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
		"message": "Item updated successfully",
		"item":    item,
	})
}

// DeleteSchemeOfWorkItem handles deleting an item (soft delete)
// @Summary Delete an item
// @Description Soft delete a scheme of work item
// @Tags Scheme of Work Items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/scheme-items/{id} [delete]
func (h *SchemeOfWorkItemHandler) DeleteSchemeOfWorkItem(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Item ID is required",
		})
		return
	}

	err := h.itemService.DeleteSchemeOfWorkItem(id)
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
		"message": "Item deleted successfully",
	})
}