package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/invoice"
)

type InvoiceController struct {
	Service *services.InvoiceService
}

func NewInvoiceController(service *services.InvoiceService) *InvoiceController {
	return &InvoiceController{Service: service}
}

// CreateInvoice godoc
//
//	@Summary		Create an invoice
//	@Description	Admin creates a new draft invoice.
//	@Tags			Invoices
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.CreateInvoiceRequest	true	"Invoice payload"
//	@Success		201		{object}	dto.SingleInvoiceResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Validation failed"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Security		BearerAuth
//	@Router			/api/invoices [post]
func (ctrl *InvoiceController) CreateInvoice(c *gin.Context) {
	var req dto.CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:  "validation failed",
			Fields: map[string]string{"body": err.Error()},
		})
		return
	}

	var createdBy *uuid.UUID
	if uidVal, ok := c.Get("user_id"); ok {
		if s, ok := uidVal.(string); ok {
			if id, err := uuid.Parse(s); err == nil {
				createdBy = &id
			}
		}
	}

	inv, err := ctrl.Service.Create(req, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.SingleInvoiceResponse{
		Data: dto.FromInvoiceModel(inv),
	})
}

// ListInvoices godoc
//
//	@Summary		List invoices (paginated)
//	@Description	Admin lists invoices with search, filters, and pagination.
//	@Tags			Invoices
//	@Produce		json
//	@Param			page			query		int		false	"Page number (default 1)"				minimum(1)
//	@Param			limit			query		int		false	"Items per page (default 20, max 100)"
//	@Param			search			query		string	false	"Search by invoice number"
//	@Param			status			query		string	false	"Filter by status"	Enums(draft, pending, partial, paid, overdue, cancelled, refunded)
//	@Param			student_id		query		string	false	"Filter by student UUID"
//	@Param			from			query		string	false	"Created after (RFC3339)"
//	@Param			to				query		string	false	"Created before (RFC3339)"
//	@Param			sort_by			query		string	false	"Sort field"		Enums(created_at, updated_at, invoice_number, total_cents, status, issued_at)
//	@Param			sort_order		query		string	false	"Sort direction"	Enums(asc, desc)
//	@Success		200				{object}	dto.PaginatedInvoicesResponse
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/invoices [get]
func (ctrl *InvoiceController) ListInvoices(c *gin.Context) {
	var q dto.ListInvoicesQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid query", Message: err.Error()})
		return
	}

	invoices, total, err := ctrl.Service.List(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	items := make([]dto.InvoiceResponse, 0, len(invoices))
	for i := range invoices {
		items = append(items, dto.FromInvoiceModel(&invoices[i]))
	}
	totalPages := int((total + int64(q.Limit) - 1) / int64(q.Limit))

	c.JSON(http.StatusOK, dto.PaginatedInvoicesResponse{
		Data:       items,
		Total:      total,
		Page:       q.Page,
		Limit:      q.Limit,
		TotalPages: totalPages,
	})
}

// GetInvoice godoc
//
//	@Summary		Get an invoice by ID
//	@Description	Returns the invoice with line items and payment history.
//	@Tags			Invoices
//	@Produce		json
//	@Param			id	path		string	true	"Invoice UUID"	format(uuid)
//	@Success		200	{object}	dto.SingleInvoiceResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/invoices/{id} [get]
func (ctrl *InvoiceController) GetInvoice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid invoice id"})
		return
	}
	inv, err := ctrl.Service.GetByID(id)
	if err != nil {
		if errors.Is(err, services.ErrInvoiceNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	// Here we return the full model so Items and Payments are included.
	c.JSON(http.StatusOK, dto.SingleInvoiceResponse{Data: inv})
}

// UpdateInvoice godoc
//
//	@Summary		Update an invoice (partial)
//	@Tags			Invoices
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Invoice UUID"	format(uuid)
//	@Param			payload	body		dto.UpdateInvoiceRequest	true	"Fields to update"
//	@Success		200		{object}	dto.SingleInvoiceResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		404		{object}	dto.ErrorResponse
//	@Failure		409		{object}	dto.ErrorResponse	"Invoice not editable"
//	@Security		BearerAuth
//	@Router			/api/invoices/{id} [patch]
func (ctrl *InvoiceController) UpdateInvoice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid invoice id"})
		return
	}
	var req dto.UpdateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:  "validation failed",
			Fields: map[string]string{"body": err.Error()},
		})
		return
	}
	inv, err := ctrl.Service.Update(id, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvoiceNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, services.ErrInvoiceNotEditable):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		}
		return
	}
	// Return updated full invoice (with items + payments)
	// Optionally reload if not preloaded
	c.JSON(http.StatusOK, dto.SingleInvoiceResponse{Data: inv})
}

// DeleteInvoice godoc
//
//	@Summary		Delete an invoice (soft)
//	@Description	Only draft or cancelled invoices can be deleted.
//	@Tags			Invoices
//	@Produce		json
//	@Param			id	path	string	true	"Invoice UUID"	format(uuid)
//	@Success		204	"No Content"
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		409	{object}	dto.ErrorResponse	"Invoice not deletable"
//	@Security		BearerAuth
//	@Router			/api/invoices/{id} [delete]
func (ctrl *InvoiceController) DeleteInvoice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid invoice id"})
		return
	}
	if err := ctrl.Service.Delete(id); err != nil {
		switch {
		case errors.Is(err, services.ErrInvoiceNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

// RecalculateInvoice godoc
//
//	@Summary		Recalculate invoice totals
//	@Description	Recomputes subtotal/tax/total/due from items and payments. Useful after bulk changes.
//	@Tags			Invoices
//	@Produce		json
//	@Param			id	path		string	true	"Invoice UUID"	format(uuid)
//	@Success		200	{object}	dto.SingleInvoiceResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/invoices/{id}/recalculate [post]
func (ctrl *InvoiceController) RecalculateInvoice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid invoice id"})
		return
	}
	if err := ctrl.Service.Recalculate(id); err != nil {
		if errors.Is(err, services.ErrInvoiceNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	inv, _ := ctrl.Service.GetByID(id)
	c.JSON(http.StatusOK, dto.SingleInvoiceResponse{
		Data: inv, // returns with items + payments
	})
}

// InvoiceSummary godoc
//
//	@Summary		Invoice dashboard summary
//	@Description	Returns totals for the admin dashboard.
//	@Tags			Invoices
//	@Produce		json
//	@Success		200	{object}	dto.InvoiceSummaryResponse
//	@Security		BearerAuth
//	@Router			/api/invoices/summary [get]
func (ctrl *InvoiceController) InvoiceSummary(c *gin.Context) {
	summary, err := ctrl.Service.Summary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}