// controllers/inventory_request_controller.go
package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/inventory_request"
)

type InventoryRequestController struct {
	requestService *services.InventoryRequestService
}

func NewInventoryRequestController(requestService *services.InventoryRequestService) *InventoryRequestController {
	return &InventoryRequestController{
		requestService: requestService,
	}
}

// CreateInventoryRequest creates a new inventory request
// @Summary Create inventory request
// @Description Create a new request for an inventory item
// @Tags Inventory Requests
// @Accept json
// @Produce json
// @Param request body dto.CreateInventoryRequestRequest true "Request data"
// @Success 201 {object} map[string]interface{} "Request created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or insufficient stock"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Section or item not found"
// @Failure 409 {object} map[string]interface{} "Duplicate pending request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-requests [post]
func (c *InventoryRequestController) CreateInventoryRequest(ctx *gin.Context) {
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

	var req dto.CreateInventoryRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	request, err := c.requestService.CreateInventoryRequest(userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "already have a pending") {
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
		"message": "Inventory request created successfully",
		"data":    request,
	})
}

// GetAllInventoryRequests retrieves all inventory requests
// @Summary Get all inventory requests
// @Description Get a paginated list of inventory requests with optional filters
// @Tags Inventory Requests
// @Accept json
// @Produce json
// @Param requested_by_id query string false "Filter by requester ID"
// @Param inventory_section_id query string false "Filter by section ID"
// @Param inventory_item_id query string false "Filter by item ID"
// @Param status query string false "Filter by status" Enums(pending, approved, rejected, fulfilled, cancelled)
// @Param search query string false "Search by request code, purpose, or notes"
// @Param date_from query string false "Filter by date from (YYYY-MM-DD)"
// @Param date_to query string false "Filter by date to (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param sort_by query string false "Sort field" default(created_at)
// @Param sort_order query string false "Sort direction" default(desc) Enums(asc, desc)
// @Success 200 {object} map[string]interface{} "Requests retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-requests [get]
func (c *InventoryRequestController) GetAllInventoryRequests(ctx *gin.Context) {
	var params dto.InventoryRequestQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := c.requestService.GetInventoryRequests(&params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Requests retrieved successfully",
		"data":    response,
	})
}

// GetMyInventoryRequests retrieves the current user's requests
// @Summary Get my inventory requests
// @Description Get all inventory requests created by the current user
// @Tags Inventory Requests
// @Accept json
// @Produce json
// @Param status query string false "Filter by status" Enums(pending, approved, rejected, fulfilled, cancelled)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{} "Requests retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-requests/my-requests [get]
func (c *InventoryRequestController) GetMyInventoryRequests(ctx *gin.Context) {
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

	var params dto.InventoryRequestQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := c.requestService.GetMyRequests(userID, &params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Requests retrieved successfully",
		"data":    response,
	})
}

// GetInventoryRequestByID retrieves a request by ID
// @Summary Get inventory request by ID
// @Description Get a specific inventory request by its ID
// @Tags Inventory Requests
// @Accept json
// @Produce json
// @Param id path string true "Request ID"
// @Success 200 {object} map[string]interface{} "Request retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Request ID is required"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Request not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-requests/{id} [get]
func (c *InventoryRequestController) GetInventoryRequestByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request ID is required"})
		return
	}

	request, err := c.requestService.GetInventoryRequestByID(id)
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
		"message": "Request retrieved successfully",
		"data":    request,
	})
}

// UpdateInventoryRequest updates a pending request
// @Summary Update inventory request
// @Description Update a pending inventory request (only by requester)
// @Tags Inventory Requests
// @Accept json
// @Produce json
// @Param id path string true "Request ID"
// @Param request body dto.UpdateInventoryRequestRequest true "Update data"
// @Success 200 {object} map[string]interface{} "Request updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or not pending"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "You can only update your own requests"
// @Failure 404 {object} map[string]interface{} "Request not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-requests/{id} [put]
func (c *InventoryRequestController) UpdateInventoryRequest(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request ID is required"})
		return
	}

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

	var req dto.UpdateInventoryRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	request, err := c.requestService.UpdateInventoryRequest(id, userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "only update your own") {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "only pending") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Request updated successfully",
		"data":    request,
	})
}

// ApproveInventoryRequest approves a request
// @Summary Approve inventory request
// @Description Approve a pending inventory request (admin only)
// @Tags Inventory Requests
// @Accept json
// @Produce json
// @Param id path string true "Request ID"
// @Param request body dto.ApproveInventoryRequestRequest true "Approval data"
// @Success 200 {object} map[string]interface{} "Request approved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or insufficient stock"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Request not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-requests/{id}/approve [post]
func (c *InventoryRequestController) ApproveInventoryRequest(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request ID is required"})
		return
	}

	approverIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	approverID, err := uuid.Parse(approverIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dto.ApproveInventoryRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	request, err := c.requestService.ApproveInventoryRequest(id, approverID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "only pending") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Request approved successfully",
		"data":    request,
	})
}

// RejectInventoryRequest rejects a request
// @Summary Reject inventory request
// @Description Reject a pending inventory request (admin only)
// @Tags Inventory Requests
// @Accept json
// @Produce json
// @Param id path string true "Request ID"
// @Param request body dto.RejectInventoryRequestRequest true "Rejection data"
// @Success 200 {object} map[string]interface{} "Request rejected successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or not pending"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Request not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-requests/{id}/reject [post]
func (c *InventoryRequestController) RejectInventoryRequest(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request ID is required"})
		return
	}

	rejectorIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	rejectorID, err := uuid.Parse(rejectorIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dto.RejectInventoryRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	request, err := c.requestService.RejectInventoryRequest(id, rejectorID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "only pending") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Request rejected successfully",
		"data":    request,
	})
}

// FulfillInventoryRequest fulfills a request
// @Summary Fulfill inventory request
// @Description Fulfill an approved request and reduce stock
// @Tags Inventory Requests
// @Accept json
// @Produce json
// @Param id path string true "Request ID"
// @Param request body dto.FulfillInventoryRequestRequest true "Fulfillment data"
// @Success 200 {object} map[string]interface{} "Request fulfilled successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or insufficient stock"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Request not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-requests/{id}/fulfill [post]
func (c *InventoryRequestController) FulfillInventoryRequest(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request ID is required"})
		return
	}

	fulfillerIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	fulfillerID, err := uuid.Parse(fulfillerIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dto.FulfillInventoryRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	request, err := c.requestService.FulfillInventoryRequest(id, fulfillerID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "only approved") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Request fulfilled successfully",
		"data":    request,
	})
}

// CancelInventoryRequest cancels a request
// @Summary Cancel inventory request
// @Description Cancel a pending or approved request (by requester)
// @Tags Inventory Requests
// @Accept json
// @Produce json
// @Param id path string true "Request ID"
// @Success 200 {object} map[string]interface{} "Request cancelled successfully"
// @Failure 400 {object} map[string]interface{} "Cannot cancel this request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "You can only cancel your own requests"
// @Failure 404 {object} map[string]interface{} "Request not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-requests/{id}/cancel [post]
func (c *InventoryRequestController) CancelInventoryRequest(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request ID is required"})
		return
	}

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

	if err := c.requestService.CancelInventoryRequest(id, userID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "only cancel your own") {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Request cancelled successfully",
	})
}

// DeleteInventoryRequest deletes a request
// @Summary Delete inventory request
// @Description Soft delete an inventory request
// @Tags Inventory Requests
// @Accept json
// @Produce json
// @Param id path string true "Request ID"
// @Success 200 {object} map[string]interface{} "Request deleted successfully"
// @Failure 400 {object} map[string]interface{} "Request ID is required"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Request not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-requests/{id} [delete]
func (c *InventoryRequestController) DeleteInventoryRequest(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request ID is required"})
		return
	}

	if err := c.requestService.DeleteInventoryRequest(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Request deleted successfully",
	})
}

// GetInventoryRequestStats retrieves statistics for inventory requests
// @Summary Get inventory request statistics
// @Description Get statistics for inventory requests
// @Tags Inventory Requests
// @Accept json
// @Produce json
// @Param section_id query string false "Filter by section ID"
// @Success 200 {object} map[string]interface{} "Statistics retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid section ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-requests/stats [get]
func (c *InventoryRequestController) GetInventoryRequestStats(ctx *gin.Context) {
	sectionID := ctx.Query("section_id")

	stats, err := c.requestService.GetInventoryRequestStats(sectionID)
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