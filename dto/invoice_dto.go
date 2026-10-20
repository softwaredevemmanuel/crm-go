package dto

import (
	"time"

	"github.com/google/uuid"

	"crm-go/models"
)

// ============================================================
// Invoice DTOs
// ============================================================

type CreateInvoiceRequest struct {
	StudentID      *uuid.UUID `json:"student_id"`
	ItemIDs []uuid.UUID `json:"item_ids" binding:"required,min=1,dive"`
	Currency       string     `json:"currency" binding:"omitempty,len=3"`
	DiscountCents  int64      `json:"discount_cents" binding:"gte=0"`
	DiscountCode   string     `json:"discount_code" binding:"omitempty,max=50"`
	DiscountReason string     `json:"discount_reason" binding:"omitempty,max=255"`
	TaxRatePercent float64    `json:"tax_rate_percent" binding:"gte=0,lte=100"`
	TaxLabel       string     `json:"tax_label" binding:"omitempty,max=50"`
	Notes          string     `json:"notes"`
	Terms          string     `json:"terms"`
	DueAt          *string    `json:"due_at"` // RFC3339
}

type UpdateInvoiceRequest struct {
	StudentID      *uuid.UUID `json:"student_id"`
	Currency       *string    `json:"currency" binding:"omitempty,len=3"`
	DiscountCents  *int64     `json:"discount_cents" binding:"omitempty,gte=0"`
	DiscountCode   *string    `json:"discount_code" binding:"omitempty,max=50"`
	DiscountReason *string    `json:"discount_reason" binding:"omitempty,max=255"`
	TaxRatePercent *float64   `json:"tax_rate_percent" binding:"omitempty,gte=0,lte=100"`
	TaxLabel       *string    `json:"tax_label" binding:"omitempty,max=50"`
	Notes          *string    `json:"notes"`
	Terms          *string    `json:"terms"`
	PdfURL         *string    `json:"pdf_url"`
	Status         *string    `json:"status" binding:"omitempty,oneof=draft pending partial paid overdue cancelled refunded"`
	DueAt          *string    `json:"due_at"`
}

type InvoiceResponse struct {
	ID            uuid.UUID  `json:"id"`
	InvoiceNumber string     `json:"invoice_number"`
	Status        string     `json:"status"`
	StudentID     *uuid.UUID `json:"student_id,omitempty"`

	Currency string `json:"currency"`

	SubtotalCents   int64 `json:"subtotal_cents"`
	DiscountCents   int64 `json:"discount_cents"`
	TaxCents        int64 `json:"tax_cents"`
	TotalCents      int64 `json:"total_cents"`
	AmountPaidCents int64 `json:"amount_paid_cents"`
	AmountDueCents  int64 `json:"amount_due_cents"`

	TaxRatePercent float64 `json:"tax_rate_percent"`
	TaxLabel       string  `json:"tax_label,omitempty"`

	DiscountCode   string `json:"discount_code,omitempty"`
	DiscountReason string `json:"discount_reason,omitempty"`

	Notes  string `json:"notes,omitempty"`
	Terms  string `json:"terms,omitempty"`
	PdfURL string `json:"pdf_url,omitempty"`

	IssuedAt    time.Time  `json:"issued_at"`
	DueAt       *time.Time `json:"due_at,omitempty"`
	SentAt      *time.Time `json:"sent_at,omitempty"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`
	RefundedAt  *time.Time `json:"refunded_at,omitempty"`

	CreatedBy *uuid.UUID `json:"created_by,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Embedded relations (populated on single-fetch)
	// Items/Payments are omitted from the DTO here to avoid circular imports if InvoiceItem DTO is elsewhere.
	// Include their IDs if you prefer, or use dedicated endpoints. 
	// For a full invoice, we'll return the model directly in the controller if needed.
}

func FromInvoiceModel(inv *models.Invoice) InvoiceResponse {
	return InvoiceResponse{
		ID:              inv.ID,
		InvoiceNumber:   inv.InvoiceNumber,
		Status:          string(inv.Status),
		StudentID:       inv.StudentID,
		Currency:        inv.Currency,
		SubtotalCents:   inv.SubtotalCents,
		DiscountCents:   inv.DiscountCents,
		TaxCents:        inv.TaxCents,
		TotalCents:      inv.TotalCents,
		AmountPaidCents: inv.AmountPaidCents,
		AmountDueCents:  inv.AmountDueCents,
		TaxRatePercent:  inv.TaxRatePercent,
		TaxLabel:        inv.TaxLabel,
		DiscountCode:    inv.DiscountCode,
		DiscountReason:  inv.DiscountReason,
		Notes:           inv.Notes,
		Terms:           inv.Terms,
		PdfURL:          inv.PdfURL,
		IssuedAt:        inv.IssuedAt,
		DueAt:           inv.DueAt,
		SentAt:          inv.SentAt,
		PaidAt:          inv.PaidAt,
		CancelledAt:     inv.CancelledAt,
		RefundedAt:      inv.RefundedAt,
		CreatedBy:       inv.CreatedBy,
		CreatedAt:       inv.CreatedAt,
		UpdatedAt:       inv.UpdatedAt,
	}
}

// ListInvoicesQuery — filters + pagination.
type ListInvoicesQuery struct {
	Page          int    `form:"page,default=1" binding:"min=1"`
	Limit         int    `form:"limit,default=20" binding:"min=1,max=100"`
	Search        string `form:"search"`
	Status        string `form:"status"`
	StudentID     string `form:"student_id"`
	From          string `form:"from"` // RFC3339
	To            string `form:"to"`   // RFC3339
	SortBy        string `form:"sort_by,default=created_at"`
	SortOrder     string `form:"sort_order,default=desc"`
}

type SingleInvoiceResponse struct {
	Data interface{} `json:"data"` // Use interface to allow model with preloaded items/payments if needed
}

type PaginatedInvoicesResponse struct {
	Data       []InvoiceResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

type InvoiceSummaryResponse struct {
	TotalInvoices    int64 `json:"total_invoices"`
	TotalBilledCents int64 `json:"total_billed_cents"`
	TotalPaidCents   int64 `json:"total_paid_cents"`
	TotalDueCents    int64 `json:"total_due_cents"`
	OverdueCount     int64 `json:"overdue_count"`
}

