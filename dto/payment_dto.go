package dto

import (
	"time"

	"github.com/google/uuid"

	"crm-go/models"
)

type CreatePaymentRequest struct {
	AmountCents     int64  `json:"amount_cents" binding:"required,gt=0"`
	Method          string `json:"method" binding:"required,oneof=card bank_transfer cash wallet paystack flutterwave stripe cheque"`
	Reference       string `json:"reference" binding:"omitempty,max=255"`
	GatewayTxnID    string `json:"gateway_txn_id" binding:"omitempty,max=255"`
	IsInstalment    bool   `json:"is_instalment"`
	InstalmentNo    int    `json:"instalment_no" binding:"omitempty,gte=0"`
	InstalmentTotal int    `json:"instalment_total" binding:"omitempty,gte=0"`
	Status          string `json:"status" binding:"omitempty,oneof=pending succeeded failed"`
	Notes           string `json:"notes"`
	PaidAt          *string `json:"paid_at"` // optional RFC3339 — defaults to now for succeeded
}

type UpdatePaymentRequest struct {
	Status          *string `json:"status" binding:"omitempty,oneof=pending succeeded failed refunded partial_refund"`
	Method          *string `json:"method" binding:"omitempty,oneof=card bank_transfer cash wallet paystack flutterwave stripe cheque"`
	Reference       *string `json:"reference" binding:"omitempty,max=255"`
	GatewayTxnID    *string `json:"gateway_txn_id" binding:"omitempty,max=255"`
	FailureReason   *string `json:"failure_reason" binding:"omitempty,max=500"`
	Notes           *string `json:"notes"`
}

type RefundPaymentRequest struct {
	AmountCents int64  `json:"amount_cents" binding:"required,gt=0"`
	Reason      string `json:"reason" binding:"omitempty,max=500"`
}

type PaymentResponse struct {
	ID              uuid.UUID  `json:"id"`
	InvoiceID       uuid.UUID  `json:"invoice_id"`
	AmountCents     int64      `json:"amount_cents"`
	RefundedCents   int64      `json:"refunded_cents"`
	NetAmountCents  int64      `json:"net_amount_cents"`
	Status          string     `json:"status"`
	Method          string     `json:"method"`
	Reference       string     `json:"reference,omitempty"`
	GatewayTxnID    string     `json:"gateway_txn_id,omitempty"`
	IsInstalment    bool       `json:"is_instalment"`
	InstalmentNo    int        `json:"instalment_no,omitempty"`
	InstalmentTotal int        `json:"instalment_total,omitempty"`
	FailureReason   string     `json:"failure_reason,omitempty"`
	RefundReason    string     `json:"refund_reason,omitempty"`
	PaidAt          *time.Time `json:"paid_at,omitempty"`
	RefundedAt      *time.Time `json:"refunded_at,omitempty"`
	FailedAt        *time.Time `json:"failed_at,omitempty"`
	Notes           string     `json:"notes,omitempty"`
	RecordedBy      *uuid.UUID `json:"recorded_by,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func FromPaymentModel(p *models.Payment) PaymentResponse {
	return PaymentResponse{
		ID:              p.ID,
		InvoiceID:       p.InvoiceID,
		AmountCents:     p.AmountCents,
		Status:          string(p.Status),
		Method:          string(p.Method),
		IsInstalment:    p.IsInstalment,
		InstalmentNo:    p.InstalmentNo,
		InstalmentTotal: p.InstalmentTotal,
		PaidAt:          p.PaidAt,
		Notes:           p.Notes,
		RecordedBy:      p.RecordedBy,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

type SinglePaymentResponse struct {
	Data PaymentResponse `json:"data"`
}

type PaginatedPaymentsResponse struct {
	Data       []PaymentResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

type ListPaymentsQuery struct {
	Page      int    `form:"page,default=1" binding:"min=1"`
	Limit     int    `form:"limit,default=20" binding:"min=1,max=100"`
	Status    string `form:"status"`
	Method    string `form:"method"`
	InvoiceID string `form:"invoice_id"`
	Search    string `form:"search"`
	From      string `form:"from"`
	To        string `form:"to"`
	SortBy    string `form:"sort_by,default=created_at"`
	SortOrder string `form:"sort_order,default=desc"`
}