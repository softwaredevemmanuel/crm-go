// controllers/inventory_item_controller.go
package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/inventory_item"
)

type InventoryItemController struct {
	itemService *services.InventoryItemService
}

func NewInventoryItemController(itemService *services.InventoryItemService) *InventoryItemController {
	return &InventoryItemController{
		itemService: itemService,
	}
}

// CreateInventoryItem creates a new inventory item
// @Summary Create inventory item
// @Description Create a new inventory item with images
// @Tags Inventory Items
// @Accept json
// @Produce json
// @Param request body dto.CreateInventoryItemRequest true "Item data"
// @Success 201 {object} map[string]interface{} "Item created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Inventory section not found"
// @Failure 409 {object} map[string]interface{} "Item code already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-items [post]
func (c *InventoryItemController) CreateInventoryItem(ctx *gin.Context) {
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dto.CreateInventoryItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	item, err := c.itemService.CreateInventoryItem(userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Inventory item created successfully",
		"data":    item,
	})
}

// GetAllInventoryItems retrieves all inventory items
// @Summary Get all inventory items
// @Description Get a paginated list of inventory items with optional filters
// @Tags Inventory Items
// @Accept json
// @Produce json
// @Param inventory_section_id query string false "Filter by Inventory Section ID"
// @Param search query string false "Search by name, code, description, brand, or model"
// @Param condition query string false "Filter by condition" Enums(new, good, fair, poor, damaged, obsolete, expired)
// @Param status query string false "Filter by status" Enums(active, inactive, archived, disposed, out_of_stock)
// @Param is_chemical query bool false "Filter chemical items only"
// @Param hazard_class query string false "Filter by hazard class" Enums(none, flammable, toxic, corrosive, explosive, radioactive, biohazard, oxidizer)
// @Param location query string false "Filter by location"
// @Param low_stock query bool false "Filter items with low stock"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param sort_by query string false "Sort field" default(created_at)
// @Param sort_order query string false "Sort direction" default(desc) Enums(asc, desc)
// @Success 200 {object} map[string]interface{} "Items retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-items [get]
func (c *InventoryItemController) GetAllInventoryItems(ctx *gin.Context) {
	var params dto.InventoryItemQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := c.itemService.GetInventoryItems(&params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Items retrieved successfully",
		"data":    response,
	})
}

// GetInventoryItemByID retrieves an item by ID
// @Summary Get inventory item by ID
// @Description Get a specific inventory item by its ID with all details
// @Tags Inventory Items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} map[string]interface{} "Item retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Item ID is required"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Item not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-items/{id} [get]
func (c *InventoryItemController) GetInventoryItemByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Item ID is required"})
		return
	}

	item, err := c.itemService.GetInventoryItemByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "invalid") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Item retrieved successfully",
		"data":    item,
	})
}

// GetInventoryItemsBySection retrieves all items in a section
// @Summary Get inventory items by section
// @Description Get all inventory items in a specific section
// @Tags Inventory Items
// @Accept json
// @Produce json
// @Param section_id path string true "Inventory Section ID"
// @Success 200 {object} map[string]interface{} "Items retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Section ID is required"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-items/section/{section_id} [get]
func (c *InventoryItemController) GetInventoryItemsBySection(ctx *gin.Context) {
	sectionID := ctx.Param("section_id")
	if sectionID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Section ID is required"})
		return
	}

	items, err := c.itemService.GetInventoryItemsBySection(sectionID)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Items retrieved successfully",
		"data":    items,
	})
}

// UpdateInventoryItem updates an existing item
// @Summary Update inventory item
// @Description Update an existing inventory item
// @Tags Inventory Items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param request body dto.UpdateInventoryItemRequest true "Update data"
// @Success 200 {object} map[string]interface{} "Item updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Item not found"
// @Failure 409 {object} map[string]interface{} "Item code already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-items/{id} [put]
func (c *InventoryItemController) UpdateInventoryItem(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Item ID is required"})
		return
	}

	var req dto.UpdateInventoryItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	item, err := c.itemService.UpdateInventoryItem(id, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Item updated successfully",
		"data":    item,
	})
}

// DeleteInventoryItem deletes an item
// @Summary Delete inventory item
// @Description Soft delete an inventory item
// @Tags Inventory Items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} map[string]interface{} "Item deleted successfully"
// @Failure 400 {object} map[string]interface{} "Item ID is required"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Item not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-items/{id} [delete]
func (c *InventoryItemController) DeleteInventoryItem(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Item ID is required"})
		return
	}

	if err := c.itemService.DeleteInventoryItem(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Item deleted successfully",
	})
}

// GetInventoryItemStats retrieves statistics for inventory items
// @Summary Get inventory item statistics
// @Description Get statistics for inventory items, optionally filtered by section
// @Tags Inventory Items
// @Accept json
// @Produce json
// @Param section_id query string false "Filter by Section ID"
// @Success 200 {object} map[string]interface{} "Statistics retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid section ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-items/stats [get]
func (c *InventoryItemController) GetInventoryItemStats(ctx *gin.Context) {
	sectionID := ctx.Query("section_id")

	stats, err := c.itemService.GetInventoryItemStats(sectionID)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Statistics retrieved successfully",
		"data":    stats,
	})
}