// handlers/user_handler.go
package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"crm-go/dto"
	"crm-go/services/users"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetAllUsers handles fetching all users with pagination and filters
// @Summary Get all users
// @Description Get a paginated list of users with optional filtering
// @Tags Users
// @Accept json
// @Produce json
// @Param search query string false "Search by name, email, or phone"
// @Param role query string false "Filter by role (admin, staff, student, parent, user)"
// @Param status query string false "Filter by status (active, inactive)"
// @Param is_active query bool false "Filter by active status"
// @Param is_verified query bool false "Filter by verification status"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param sort_by query string false "Sort by field" default(created_at)
// @Param sort_order query string false "Sort order" default(desc)
// @Success 200 {object} dto.UserListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	var params dto.UserQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.userService.GetAllUsers(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved successfully",
		"data":    response,
	})
}



// GetUsersByRole handles fetching all users with a specific role
// @Summary Get users by role
// @Description Get all users with a specific role
// @Tags Users
// @Accept json
// @Produce json
// @Param role path string true "Role (admin, staff, student, parent, user)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/users/role/{role} [get]
func (h *UserHandler) GetUsersByRole(c *gin.Context) {
	role := c.Param("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Role is required",
		})
		return
	}

	users, err := h.userService.GetUsersByRole(role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved successfully",
		"users":   users,
	})
}



// DeleteUser handles deleting a user (soft delete)
// @Summary Delete a user
// @Description Soft delete a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	err := h.userService.DeleteUser(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "cannot delete") {
			c.JSON(http.StatusConflict, gin.H{
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
		"message": "User deleted successfully",
	})
}