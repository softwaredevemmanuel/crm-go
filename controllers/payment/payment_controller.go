package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/payment"
)

type PaymentController struct {
	Service *services.PaymentService
}

func NewPaymentController(service *services.PaymentService) *PaymentController {
	return &PaymentController{Service: service}
}

// CreatePayment godoc
//
//	@Summary		Record a payment against an invoice
//	@Description	Supports partial and instalment payments. Automatically recomputes the invoice totals and status.
//	@Tags			Payments
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Invoice UUID"	format(uuid)
//	@Param			payload	body		dto.CreatePaymentRequest	true	"Payment payload"
//	@Success		201		{object}	dto.SinglePaymentResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		404		{object}	dto.ErrorResponse	"Invoice not found"
//	@Security		BearerAuth
//	@Router			/api/invoices/{id}/payments [post]
func (ctrl *PaymentController) CreatePayment(c *gin.Context) {
	invoiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid invoice id"})
		return
	}

	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:  "validation failed",
			Fields: map[string]string{"body": err.Error()},
		})
		return
	}

	var recordedBy *uuid.UUID
	if uidVal, ok := c.Get("user_id"); ok {
		if s, ok := uidVal.(string); ok {
			if id, err := uuid.Parse(s); err == nil {
				recordedBy = &id
			}
		}
	}

	payment, err := ctrl.Service.Create(invoiceID, req, recordedBy)
	if err != nil {
		if errors.Is(err, services.ErrInvoiceNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.SinglePaymentResponse{
		Data: dto.FromPaymentModel(payment),
	})
}

// ListPayments godoc
//
//	@Summary		List payments (paginated)
//	@Tags			Payments
//	@Produce		json
//	@Param			page		query		int		false	"Page number"
//	@Param			limit		query		int		false	"Items per page"
//	@Param			status		query		string	false	"Filter by status"
//	@Param			method		query		string	false	"Filter by method"
//	@Param			invoice_id	query		string	false	"Filter by invoice UUID"
//	@Success		200			{object}	dto.PaginatedPaymentsResponse
//	@Security		BearerAuth
//	@Router			/api/payments [get]
func (ctrl *PaymentController) ListPayments(c *gin.Context) {
	var q dto.ListPaymentsQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid query", Message: err.Error()})
		return
	}
	payments, total, err := ctrl.Service.ListPayments(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	items := make([]dto.PaymentResponse, 0, len(payments))
	for i := range payments {
		items = append(items, dto.FromPaymentModel(&payments[i]))
	}
	totalPages := int((total + int64(q.Limit) - 1) / int64(q.Limit))

	c.JSON(http.StatusOK, dto.PaginatedPaymentsResponse{
		Data:       items,
		Total:      total,
		Page:       q.Page,
		Limit:      q.Limit,
		TotalPages: totalPages,
	})
}

// ListInvoicePayments godoc
//
//	@Summary		List payments for a specific invoice
//	@Tags			Payments
//	@Produce		json
//	@Param			id	path		string	true	"Invoice UUID"	format(uuid)
//	@Success		200	{object}	map[string]interface{}
//	@Security		BearerAuth
//	@Router			/api/invoices/{id}/payments [get]
func (ctrl *PaymentController) ListInvoicePayments(c *gin.Context) {
	invoiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid invoice id"})
		return
	}
	
	q := dto.ListPaymentsQuery{
		Page:      1,
		Limit:     100,
		InvoiceID: invoiceID.String(),
		SortBy:    "created_at",
		SortOrder: "desc",
	}
	payments, _, err := ctrl.Service.ListPayments(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	resp := make([]dto.PaymentResponse, 0, len(payments))
	for i := range payments {
		resp = append(resp, dto.FromPaymentModel(&payments[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// GetPayment godoc
//
//	@Summary		Get a payment by ID
//	@Tags			Payments
//	@Produce		json
//	@Param			id	path		string	true	"Payment UUID"	format(uuid)
//	@Success		200	{object}	dto.SinglePaymentResponse
//	@Security		BearerAuth
//	@Router			/api/payments/{id} [get]
func (ctrl *PaymentController) GetPayment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid payment id"})
		return
	}
	p, err := ctrl.Service.GetByID(id)
	if err != nil {
		if errors.Is(err, services.ErrPaymentNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.SinglePaymentResponse{
		Data: dto.FromPaymentModel(p),
	})
}

// UpdatePayment godoc
//
//	@Summary		Update a payment (status/notes)
//	@Tags			Payments
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Payment UUID"	format(uuid)
//	@Param			payload	body		dto.UpdatePaymentRequest	true	"Fields to update"
//	@Success		200		{object}	dto.SinglePaymentResponse
//	@Security		BearerAuth
//	@Router			/api/payments/{id} [patch]
func (ctrl *PaymentController) UpdatePayment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid payment id"})
		return
	}
	var req dto.UpdatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:  "validation failed",
			Fields: map[string]string{"body": err.Error()},
		})
		return
	}
	p, err := ctrl.Service.Update(id, req)
	if err != nil {
		if errors.Is(err, services.ErrPaymentNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.SinglePaymentResponse{
		Data: dto.FromPaymentModel(p),
	})
}

// DeletePayment godoc
//
//	@Summary		Delete a payment
//	@Description	Hard deletes the payment and recomputes the parent invoice.
//	@Tags			Payments
//	@Produce		json
//	@Param			id	path	string	true	"Payment UUID"	format(uuid)
//	@Success		204	"No Content"
//	@Security		BearerAuth
//	@Router			/api/payments/{id} [delete]
func (ctrl *PaymentController) DeletePayment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid payment id"})
		return
	}
	if err := ctrl.Service.Delete(id); err != nil {
		if errors.Is(err, services.ErrPaymentNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}