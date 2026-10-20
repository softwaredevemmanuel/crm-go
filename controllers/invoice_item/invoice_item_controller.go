package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/invoice_item"
)

type InvoiceItemController struct {
	Service *services.InvoiceItemService
}

func NewInvoiceItemController(service *services.InvoiceItemService) *InvoiceItemController {
	return &InvoiceItemController{Service: service}
}

// CreateInvoiceItem godoc
//
//	@Summary		Create an invoice item
//	@Description	Adds a single charge (tuition, hostel, exam, discount, etc.) under a Term.
//	@Tags			Invoice Items
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.CreateInvoiceItemRequest	true	"Item payload"
//	@Success		201		{object}	dto.SingleInvoiceItemResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Validation failed"
//	@Failure		404		{object}	dto.ErrorResponse	"Term not found"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Security		BearerAuth
//	@Router			/api/invoice-items [post]
func (ctrl *InvoiceItemController) CreateInvoiceItem(c *gin.Context) {
	var req dto.CreateInvoiceItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:  "validation failed",
			Fields: map[string]string{"body": err.Error()},
		})
		return
	}

	item, err := ctrl.Service.Create(req)
	if err != nil {
		if errors.Is(err, services.ErrTermNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.SingleInvoiceItemResponse{
		Data: dto.FromInvoiceItemModel(item, true),
	})
}

// ListInvoiceItems godoc
//
//	@Summary		List invoice items (paginated)
//	@Tags			Invoice Items
//	@Produce		json
//	@Param			page		query		int		false	"Page number (default 1)"				minimum(1)
//	@Param			limit		query		int		false	"Items per page (default 50, max 200)"
//	@Param			search		query		string	false	"Search in title, description, SKU"
//	@Param			item_type	query		string	false	"Filter by item type"	Enums(tuition, course, material, exam, registration, hostel, transport, library, uniform, fine, discount, custom)
//	@Param			term_id		query		string	false	"Filter by term UUID"
//	@Param			sort_by		query		string	false	"Sort field"	Enums(display_order, created_at, updated_at, title, unit_price_cents, item_type)
//	@Param			sort_order	query		string	false	"Sort direction"	Enums(asc, desc)
//	@Success		200			{object}	dto.PaginatedInvoiceItemsResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/invoice-items [get]
func (ctrl *InvoiceItemController) ListInvoiceItems(c *gin.Context) {
	var q dto.ListInvoiceItemsQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid query", Message: err.Error()})
		return
	}

	items, total, err := ctrl.Service.List(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	resp := make([]dto.InvoiceItemResponse, 0, len(items))
	for i := range items {
		resp = append(resp, dto.FromInvoiceItemModel(&items[i], true))
	}
	totalPages := int((total + int64(q.Limit) - 1) / int64(q.Limit))

	c.JSON(http.StatusOK, dto.PaginatedInvoiceItemsResponse{
		Data:       resp,
		Total:      total,
		Page:       q.Page,
		Limit:      q.Limit,
		TotalPages: totalPages,
	})
}

// GetInvoiceItem godoc
//
//	@Summary		Get an invoice item by ID
//	@Tags			Invoice Items
//	@Produce		json
//	@Param			id	path		string	true	"Invoice Item UUID"	format(uuid)
//	@Success		200	{object}	dto.SingleInvoiceItemResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/invoice-items/{id} [get]
func (ctrl *InvoiceItemController) GetInvoiceItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid item id"})
		return
	}

	item, err := ctrl.Service.GetByID(id)
	if err != nil {
		if errors.Is(err, services.ErrInvoiceItemNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SingleInvoiceItemResponse{
		Data: dto.FromInvoiceItemModel(item, true),
	})
}

// ListInvoiceItemsByTerm godoc
//
//	@Summary		List items belonging to a specific term
//	@Tags			Invoice Items
//	@Produce		json
//	@Param			term_id	path		string	true	"Term UUID"	format(uuid)
//	@Success		200		{array}		dto.InvoiceItemResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/term/{term_id}/invoice-items [get]
func (ctrl *InvoiceItemController) ListInvoiceItemsByTerm(c *gin.Context) {
	termID, err := uuid.Parse(c.Param("term_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid term id"})
		return
	}

	items, err := ctrl.Service.ListByTerm(termID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	resp := make([]dto.InvoiceItemResponse, 0, len(items))
	for i := range items {
		resp = append(resp, dto.FromInvoiceItemModel(&items[i], true))
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// UpdateInvoiceItem godoc
//
//	@Summary		Update an invoice item (partial)
//	@Description	TermID is immutable. To move an item to a different term, delete + recreate.
//	@Tags			Invoice Items
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string							true	"Invoice Item UUID"	format(uuid)
//	@Param			payload	body		dto.UpdateInvoiceItemRequest	true	"Fields to update"
//	@Success		200		{object}	dto.SingleInvoiceItemResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		404		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/invoice-items/{id} [patch]
func (ctrl *InvoiceItemController) UpdateInvoiceItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid item id"})
		return
	}

	var req dto.UpdateInvoiceItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:  "validation failed",
			Fields: map[string]string{"body": err.Error()},
		})
		return
	}

	item, err := ctrl.Service.Update(id, req)
	if err != nil {
		if errors.Is(err, services.ErrInvoiceItemNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SingleInvoiceItemResponse{
		Data: dto.FromInvoiceItemModel(item, true),
	})
}

// DeleteInvoiceItem godoc
//
//	@Summary		Delete an invoice item
//	@Description	Soft-deletes the line item.
//	@Tags			Invoice Items
//	@Produce		json
//	@Param			id	path	string	true	"Invoice Item UUID"	format(uuid)
//	@Success		204	"No Content"
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/invoice-items/{id} [delete]
func (ctrl *InvoiceItemController) DeleteInvoiceItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid item id"})
		return
	}

	if err := ctrl.Service.Delete(id); err != nil {
		if errors.Is(err, services.ErrInvoiceItemNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

