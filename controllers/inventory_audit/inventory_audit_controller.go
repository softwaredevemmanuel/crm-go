// controllers/inventory_audit_controller.go
package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/inventory_audit"
)

type InventoryAuditController struct {
	auditService *services.InventoryAuditService
}

func NewInventoryAuditController(auditService *services.InventoryAuditService) *InventoryAuditController {
	return &InventoryAuditController{
		auditService: auditService,
	}
}

// CreateInventoryAudit creates a new audit
// @Summary Create inventory audit
// @Description Create a new inventory audit with optional items
// @Tags Inventory Audits
// @Accept json
// @Produce json
// @Param request body dto.CreateInventoryAuditRequest true "Audit data"
// @Success 201 {object} map[string]interface{} "Audit created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Section or item not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-audits [post]
func (c *InventoryAuditController) CreateInventoryAudit(ctx *gin.Context) {
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

	var req dto.CreateInventoryAuditRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	audit, err := c.auditService.CreateInventoryAudit(userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Inventory audit created successfully",
		"data":    audit,
	})
}

// GetAllInventoryAudits retrieves all audits
// @Summary Get all inventory audits
// @Description Get a paginated list of inventory audits with filters
// @Tags Inventory Audits
// @Accept json
// @Produce json
// @Param inventory_section_id query string false "Filter by Section ID"
// @Param conducted_by_id query string false "Filter by conductor ID"
// @Param status query string false "Filter by status" Enums(in_progress, completed, approved, rejected, cancelled)
// @Param search query string false "Search by code or notes"
// @Param date_from query string false "Filter by date from (YYYY-MM-DD)"
// @Param date_to query string false "Filter by date to (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param sort_by query string false "Sort field" default(created_at)
// @Param sort_order query string false "Sort direction" default(desc) Enums(asc, desc)
// @Success 200 {object} map[string]interface{} "Audits retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-audits [get]
func (c *InventoryAuditController) GetAllInventoryAudits(ctx *gin.Context) {
	var params dto.InventoryAuditQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := c.auditService.GetInventoryAudits(&params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Audits retrieved successfully",
		"data":    response,
	})
}

// GetInventoryAuditByID retrieves an audit by ID
// @Summary Get inventory audit by ID
// @Description Get a specific audit by its ID with all items
// @Tags Inventory Audits
// @Accept json
// @Produce json
// @Param id path string true "Audit ID"
// @Success 200 {object} map[string]interface{} "Audit retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Audit ID is required"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Audit not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-audits/{id} [get]
func (c *InventoryAuditController) GetInventoryAuditByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Audit ID is required"})
		return
	}

	audit, err := c.auditService.GetInventoryAuditByID(id)
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
		"message": "Audit retrieved successfully",
		"data":    audit,
	})
}

// UpdateInventoryAudit updates an in-progress audit
// @Summary Update inventory audit
// @Description Update an in-progress audit
// @Tags Inventory Audits
// @Accept json
// @Produce json
// @Param id path string true "Audit ID"
// @Param request body dto.UpdateInventoryAuditRequest true "Update data"
// @Success 200 {object} map[string]interface{} "Audit updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or not in-progress"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Not authorized"
// @Failure 404 {object} map[string]interface{} "Audit not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-audits/{id} [put]
func (c *InventoryAuditController) UpdateInventoryAudit(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Audit ID is required"})
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

	var req dto.UpdateInventoryAuditRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	audit, err := c.auditService.UpdateInventoryAudit(id, userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "only update your own") {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Audit updated successfully",
		"data":    audit,
	})
}

// AddAuditItem adds an item to an in-progress audit
// @Summary Add item to audit
// @Description Add an item to an in-progress audit
// @Tags Inventory Audits
// @Accept json
// @Produce json
// @Param id path string true "Audit ID"
// @Param request body dto.AuditItemInput true "Item data"
// @Success 200 {object} map[string]interface{} "Item added successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or duplicate"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Audit or item not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-audits/{id}/items [post]
func (c *InventoryAuditController) AddAuditItem(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Audit ID is required"})
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

	var req dto.AuditItemInput
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	audit, err := c.auditService.AddAuditItem(id, userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "already been added") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Item added to audit successfully",
		"data":    audit,
	})
}

// UpdateAuditItem updates an audit item
// @Summary Update audit item
// @Description Update a specific item within an in-progress audit
// @Tags Inventory Audits
// @Accept json
// @Produce json
// @Param id path string true "Audit ID"
// @Param item_id path string true "Audit Item ID"
// @Param request body dto.UpdateAuditItemRequest true "Update data"
// @Success 200 {object} map[string]interface{} "Item updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Audit or item not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-audits/{id}/items/{item_id} [put]
func (c *InventoryAuditController) UpdateAuditItem(ctx *gin.Context) {
	id := ctx.Param("id")
	itemID := ctx.Param("item_id")

	if id == "" || itemID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Audit ID and Item ID are required"})
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

	var req dto.UpdateAuditItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	audit, err := c.auditService.UpdateAuditItem(id, itemID, userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Audit item updated successfully",
		"data":    audit,
	})
}

// RemoveAuditItem removes an item from an audit
// @Summary Remove audit item
// @Description Remove an item from an in-progress audit
// @Tags Inventory Audits
// @Accept json
// @Produce json
// @Param id path string true "Audit ID"
// @Param item_id path string true "Audit Item ID"
// @Success 200 {object} map[string]interface{} "Item removed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Audit or item not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-audits/{id}/items/{item_id} [delete]
func (c *InventoryAuditController) RemoveAuditItem(ctx *gin.Context) {
	id := ctx.Param("id")
	itemID := ctx.Param("item_id")

	if id == "" || itemID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Audit ID and Item ID are required"})
		return
	}

	if err := c.auditService.RemoveAuditItem(id, itemID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Item removed from audit successfully",
	})
}

// CompleteInventoryAudit completes an audit
// @Summary Complete inventory audit
// @Description Mark an in-progress audit as completed
// @Tags Inventory Audits
// @Accept json
// @Produce json
// @Param id path string true "Audit ID"
// @Success 200 {object} map[string]interface{} "Audit completed successfully"
// @Failure 400 {object} map[string]interface{} "Cannot complete audit"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Not authorized"
// @Failure 404 {object} map[string]interface{} "Audit not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-audits/{id}/complete [post]
func (c *InventoryAuditController) CompleteInventoryAudit(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Audit ID is required"})
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

	audit, err := c.auditService.CompleteInventoryAudit(id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "only complete your own") {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Audit completed successfully",
		"data":    audit,
	})
}

// ApproveInventoryAudit approves a completed audit
// @Summary Approve inventory audit
// @Description Approve a completed audit and apply stock adjustments
// @Tags Inventory Audits
// @Accept json
// @Produce json
// @Param id path string true "Audit ID"
// @Param request body dto.ApproveInventoryAuditRequest true "Approval data"
// @Success 200 {object} map[string]interface{} "Audit approved successfully"
// @Failure 400 {object} map[string]interface{} "Cannot approve audit"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Audit not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-audits/{id}/approve [post]
func (c *InventoryAuditController) ApproveInventoryAudit(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Audit ID is required"})
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

	var req dto.ApproveInventoryAuditRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	audit, err := c.auditService.ApproveInventoryAudit(id, approverID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Audit approved successfully",
		"data":    audit,
	})
}

// RejectInventoryAudit rejects a completed audit
// @Summary Reject inventory audit
// @Description Reject a completed audit with reason
// @Tags Inventory Audits
// @Accept json
// @Produce json
// @Param id path string true "Audit ID"
// @Param request body dto.RejectInventoryAuditRequest true "Rejection data"
// @Success 200 {object} map[string]interface{} "Audit rejected successfully"
// @Failure 400 {object} map[string]interface{} "Cannot reject audit"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Audit not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-audits/{id}/reject [post]
func (c *InventoryAuditController) RejectInventoryAudit(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Audit ID is required"})
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

	var req dto.RejectInventoryAuditRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	audit, err := c.auditService.RejectInventoryAudit(id, rejectorID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Audit rejected successfully",
		"data":    audit,
	})
}

// CancelInventoryAudit cancels an in-progress audit
// @Summary Cancel inventory audit
// @Description Cancel an in-progress audit
// @Tags Inventory Audits
// @Accept json
// @Produce json
// @Param id path string true "Audit ID"
// @Success 200 {object} map[string]interface{} "Audit cancelled successfully"
// @Failure 400 {object} map[string]interface{} "Cannot cancel audit"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Not authorized"
// @Failure 404 {object} map[string]interface{} "Audit not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-audits/{id}/cancel [post]
func (c *InventoryAuditController) CancelInventoryAudit(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Audit ID is required"})
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

	if err := c.auditService.CancelInventoryAudit(id, userID); err != nil {
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
		"message": "Audit cancelled successfully",
	})
}

// DeleteInventoryAudit deletes an audit
// @Summary Delete inventory audit
// @Description Soft delete an inventory audit
// @Tags Inventory Audits
// @Accept json
// @Produce json
// @Param id path string true "Audit ID"
// @Success 200 {object} map[string]interface{} "Audit deleted successfully"
// @Failure 400 {object} map[string]interface{} "Audit ID is required"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Audit not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-audits/{id} [delete]
func (c *InventoryAuditController) DeleteInventoryAudit(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Audit ID is required"})
		return
	}

	if err := c.auditService.DeleteInventoryAudit(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Audit deleted successfully",
	})
}

// GetInventoryAuditStats retrieves statistics for audits
// @Summary Get inventory audit statistics
// @Description Get statistics for inventory audits
// @Tags Inventory Audits
// @Accept json
// @Produce json
// @Param section_id query string false "Filter by section ID"
// @Success 200 {object} map[string]interface{} "Statistics retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid section ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/inventory-audits/stats [get]
func (c *InventoryAuditController) GetInventoryAuditStats(ctx *gin.Context) {
	sectionID := ctx.Query("section_id")

	stats, err := c.auditService.GetInventoryAuditStats(sectionID)
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