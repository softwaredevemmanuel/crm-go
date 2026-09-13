// controllers/inventory_section_controller.go
package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"crm-go/dto"
	"crm-go/services/inventory_section"
)

type InventorySectionController struct {
	sectionService *services.InventorySectionService
}

func NewInventorySectionController(sectionService *services.InventorySectionService) *InventorySectionController {
	return &InventorySectionController{
		sectionService: sectionService,
	}
}

// CreateInventorySection creates a new inventory section
// @Summary Create inventory section
// @Description Create a new inventory section (e.g., Laboratory, ICT Room)
// @Tags Inventory Sections
// @Accept json
// @Produce json
// @Param request body dto.CreateInventorySectionRequest true "Section data"
// @Success 201 {object} map[string]interface{} "Section created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Head of department user not found"
// @Failure 409 {object} map[string]interface{} "Section with this code already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-sections [post]
func (c *InventorySectionController) CreateInventorySection(ctx *gin.Context) {
	var req dto.CreateInventorySectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	section, err := c.sectionService.CreateInventorySection(&req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Inventory section created successfully",
		"data":    section,
	})
}

// GetAllInventorySections retrieves all inventory sections
// @Summary Get all inventory sections
// @Description Get a paginated list of inventory sections with optional filters
// @Tags Inventory Sections
// @Accept json
// @Produce json
// @Param search query string false "Search by name, code, or description"
// @Param status query string false "Filter by status" Enums(active, inactive, closed, renovation)
// @Param location query string false "Filter by location"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param sort_by query string false "Sort field" default(created_at)
// @Param sort_order query string false "Sort direction" default(desc) Enums(asc, desc)
// @Success 200 {object} map[string]interface{} "Sections retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-sections [get]
func (c *InventorySectionController) GetAllInventorySections(ctx *gin.Context) {
	var params dto.InventorySectionQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	
	response, err := c.sectionService.GetInventorySections(&params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sections retrieved successfully",
		"data":    response,
	})
}

// GetInventorySectionByID retrieves an inventory section by ID
// @Summary Get inventory section by ID
// @Description Get a specific inventory section by its ID with all details
// @Tags Inventory Sections
// @Accept json
// @Produce json
// @Param id path string true "Section ID"
// @Success 200 {object} map[string]interface{} "Section retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Section ID is required"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Section not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-sections/{id} [get]
func (c *InventorySectionController) GetInventorySectionByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Section ID is required",
		})
		return
	}

	section, err := c.sectionService.GetInventorySectionByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "invalid") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Section retrieved successfully",
		"data":    section,
	})
}

// UpdateInventorySection updates an existing inventory section
// @Summary Update inventory section
// @Description Update an existing inventory section's details
// @Tags Inventory Sections
// @Accept json
// @Produce json
// @Param id path string true "Section ID"
// @Param request body dto.UpdateInventorySectionRequest true "Update data"
// @Success 200 {object} map[string]interface{} "Section updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Section not found"
// @Failure 409 {object} map[string]interface{} "Section code already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-sections/{id} [put]
func (c *InventorySectionController) UpdateInventorySection(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Section ID is required",
		})
		return
	}

	var req dto.UpdateInventorySectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	section, err := c.sectionService.UpdateInventorySection(id, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Section updated successfully",
		"data":    section,
	})
}

// DeleteInventorySection deletes an inventory section
// @Summary Delete inventory section
// @Description Soft delete an inventory section. Cannot delete if it has items.
// @Tags Inventory Sections
// @Accept json
// @Produce json
// @Param id path string true "Section ID"
// @Success 200 {object} map[string]interface{} "Section deleted successfully"
// @Failure 400 {object} map[string]interface{} "Section ID is required or section has items"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Section not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-sections/{id} [delete]
func (c *InventorySectionController) DeleteInventorySection(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Section ID is required",
		})
		return
	}

	if err := c.sectionService.DeleteInventorySection(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "cannot delete") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Section deleted successfully",
	})
}