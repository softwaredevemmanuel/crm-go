// handlers/address_handler.go
package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/address"
)

type AddressHandler struct {
	addressService *services.AddressService
}

func NewAddressHandler(addressService *services.AddressService) *AddressHandler {
	return &AddressHandler{
		addressService: addressService,
	}
}

// CreateAddress handles the creation of a new address
// @Summary Create a new address
// @Description Create a new address for a user
// @Tags Addresses
// @Accept json
// @Produce json
// @Param request body dto.CreateAddressRequest true "Address creation request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/addresses [post]
func (h *AddressHandler) CreateAddress(c *gin.Context) {
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

	var req dto.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	address, err := h.addressService.CreateAddress(&req, userID)
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
		"message": "Address created successfully",
		"address": address,
	})
}

// GetAllAddresses handles fetching all addresses with pagination and filters
// @Summary Get all addresses
// @Description Get a paginated list of addresses with optional filtering
// @Tags Addresses
// @Accept json
// @Produce json
// @Param search query string false "Search by address, city, state, or user name"
// @Param user_id query string false "Filter by user ID"
// @Param address_type query string false "Filter by address type"
// @Param status query string false "Filter by status"
// @Param is_primary query bool false "Filter by primary status"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.AddressListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/addresses [get]
func (h *AddressHandler) GetAllAddresses(c *gin.Context) {
	var params dto.AddressQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.addressService.GetAllAddresses(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Addresses retrieved successfully",
		"data":    response,
	})
}

// GetAddressByID handles fetching a single address by ID
// @Summary Get address by ID
// @Description Get a single address by its ID
// @Tags Addresses
// @Accept json
// @Produce json
// @Param id path string true "Address ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/addresses/{id} [get]
func (h *AddressHandler) GetAddressByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Address ID is required",
		})
		return
	}

	address, err := h.addressService.GetAddressByID(id)
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
		"message": "Address retrieved successfully",
		"address": address,
	})
}

// GetAddressesByUser handles fetching all addresses for a user
// @Summary Get addresses by user
// @Description Get all addresses associated with a specific user
// @Tags Addresses
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/addresses/user/{user_id} [get]
func (h *AddressHandler) GetAddressesByUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	addresses, err := h.addressService.GetAddressesByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Addresses retrieved successfully",
		"addresses": addresses,
	})
}

// GetPrimaryAddressByUser handles fetching the primary address for a user
// @Summary Get primary address by user
// @Description Get the primary address associated with a specific user
// @Tags Addresses
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/addresses/user/{user_id}/primary [get]
func (h *AddressHandler) GetPrimaryAddressByUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	address, err := h.addressService.GetPrimaryAddressByUser(userID)
	if err != nil {
		if strings.Contains(err.Error(), "no primary address found") {
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
		"message": "Primary address retrieved successfully",
		"address": address,
	})
}

// UpdateAddress handles updating an existing address
// @Summary Update an address
// @Description Update an existing address by ID
// @Tags Addresses
// @Accept json
// @Produce json
// @Param id path string true "Address ID"
// @Param request body dto.UpdateAddressRequest true "Address update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/addresses/{id} [put]
func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Address ID is required",
		})
		return
	}

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

	var req dto.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	address, err := h.addressService.UpdateAddress(id, &req, userID)
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
		"message": "Address updated successfully",
		"address": address,
	})
}

// DeleteAddress handles deleting an address (soft delete)
// @Summary Delete an address
// @Description Soft delete an address by ID
// @Tags Addresses
// @Accept json
// @Produce json
// @Param id path string true "Address ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/addresses/{id} [delete]
func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Address ID is required",
		})
		return
	}

	err := h.addressService.DeleteAddress(id)
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
		"message": "Address deleted successfully",
	})
}