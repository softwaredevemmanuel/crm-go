// controllers/modules/module_controller.go
package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"crm-go/dto"
	"crm-go/services/modules"
)

type ModuleHandler struct {
	moduleService *services.ModuleService
}

func NewModuleHandler(moduleService *services.ModuleService) *ModuleHandler {
	return &ModuleHandler{
		moduleService: moduleService,
	}
}

// CreateModule handles the creation of a new module
// @Summary Create a module
// @Description Create a new module for a scheme of work
// @Tags Modules
// @Accept json
// @Produce json
// @Param request body dto.CreateModuleRequest true "Module request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/modules [post]
func (h *ModuleHandler) CreateModule(c *gin.Context) {
	var req dto.CreateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	module, err := h.moduleService.CreateModule(&req)
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
		"message": "Module created successfully",
		"module":  module,
	})
}

// BulkCreateModules handles bulk creation of modules
// @Summary Bulk create modules
// @Description Create multiple modules for a scheme of work
// @Tags Modules
// @Accept json
// @Produce json
// @Param request body dto.BulkCreateModulesRequest true "Bulk module request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/modules/bulk [post]
func (h *ModuleHandler) BulkCreateModules(c *gin.Context) {
	var req dto.BulkCreateModulesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.moduleService.BulkCreateModules(&req)
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
		"message": "Modules created successfully",
		"data":    result,
	})
}

// GetAllModules handles fetching all modules
// @Summary Get all modules
// @Description Get a paginated list of all modules
// @Tags Modules
// @Accept json
// @Produce json
// @Param scheme_of_work_id query string false "Filter by scheme of work ID"
// @Param search query string false "Search by title or description"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.ModuleListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/modules [get]
func (h *ModuleHandler) GetAllModules(c *gin.Context) {
	var params dto.ModuleQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.moduleService.GetAllModules(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Modules retrieved successfully",
		"data":    response,
	})
}

// GetModuleByID handles fetching a single module by ID
// @Summary Get module by ID
// @Description Get a single module by its ID
// @Tags Modules
// @Accept json
// @Produce json
// @Param id path string true "Module ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/modules/{id} [get]
func (h *ModuleHandler) GetModuleByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Module ID is required",
		})
		return
	}

	module, err := h.moduleService.GetModuleByID(id)
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
		"message": "Module retrieved successfully",
		"module":  module,
	})
}

// GetModulesBySchemeOfWork handles fetching all modules for a scheme of work
// @Summary Get modules by scheme of work
// @Description Get all modules for a specific scheme of work
// @Tags Modules
// @Accept json
// @Produce json
// @Param scheme_of_work_id path string true "Scheme of Work ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/modules/scheme/{scheme_of_work_id} [get]
func (h *ModuleHandler) GetModulesBySchemeOfWork(c *gin.Context) {
	schemeOfWorkID := c.Param("scheme_of_work_id")
	if schemeOfWorkID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Scheme of work ID is required",
		})
		return
	}

	modules, err := h.moduleService.GetModulesBySchemeOfWork(schemeOfWorkID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Modules retrieved successfully",
		"modules": modules,
		"total":   len(modules),
	})
}

// UpdateModule handles updating a module
// @Summary Update a module
// @Description Update an existing module
// @Tags Modules
// @Accept json
// @Produce json
// @Param id path string true "Module ID"
// @Param request body dto.UpdateModuleRequest true "Update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/modules/{id} [put]
func (h *ModuleHandler) UpdateModule(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Module ID is required",
		})
		return
	}

	var req dto.UpdateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	module, err := h.moduleService.UpdateModule(id, &req)
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

	c.JSON(http.StatusOK, gin.H{
		"message": "Module updated successfully",
		"module":  module,
	})
}

// ReorderModules handles reordering modules
// @Summary Reorder modules
// @Description Reorder modules within a scheme of work
// @Tags Modules
// @Accept json
// @Produce json
// @Param request body dto.ReorderModulesRequest true "Reorder request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/modules/reorder [put]
func (h *ModuleHandler) ReorderModules(c *gin.Context) {
	var req dto.ReorderModulesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := h.moduleService.ReorderModules(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Modules reordered successfully",
	})
}

// DeleteModule handles deleting a module (soft delete)
// @Summary Delete a module
// @Description Soft delete a module
// @Tags Modules
// @Accept json
// @Produce json
// @Param id path string true "Module ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/modules/{id} [delete]
func (h *ModuleHandler) DeleteModule(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Module ID is required",
		})
		return
	}

	err := h.moduleService.DeleteModule(id)
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
		"message": "Module deleted successfully",
	})
}