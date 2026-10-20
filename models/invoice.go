package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvoiceStatus string

const (
	InvoiceDraft     InvoiceStatus = "draft"
	InvoicePending   InvoiceStatus = "pending"
	InvoicePartial   InvoiceStatus = "partial"
	InvoicePaid      InvoiceStatus = "paid"
	InvoiceOverdue   InvoiceStatus = "overdue"
	InvoiceCancelled InvoiceStatus = "cancelled"
	InvoiceRefunded  InvoiceStatus = "refunded"
)

type Invoice struct {
	ID            uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	InvoiceNumber string        `gorm:"type:varchar(50);uniqueIndex;not null" json:"invoice_number" example:"INV-2026-00042"`
	Status        InvoiceStatus `gorm:"type:varchar(20);not null;default:'draft';index" json:"status" example:"pending"`

	// Student / customer
	StudentID     *uuid.UUID `gorm:"type:uuid;index" json:"student_id,omitempty"`

	// Currency — all amounts in minor units (kobo/cents)
	Currency string `gorm:"type:varchar(3);not null;default:'NGN'" json:"currency" example:"NGN"`

	// Totals (denormalized — recomputed from items + payments)
	SubtotalCents   int64 `gorm:"not null;default:0" json:"subtotal_cents"`
	DiscountCents   int64 `gorm:"not null;default:0" json:"discount_cents"`
	TaxCents        int64 `gorm:"not null;default:0" json:"tax_cents"`
	TotalCents      int64 `gorm:"not null;default:0" json:"total_cents"`
	AmountPaidCents int64 `gorm:"not null;default:0" json:"amount_paid_cents"`
	AmountDueCents  int64 `gorm:"not null;default:0" json:"amount_due_cents"`

	TaxRatePercent float64 `gorm:"type:numeric(5,2);default:0" json:"tax_rate_percent"`
	TaxLabel       string  `gorm:"type:varchar(50)" json:"tax_label,omitempty"`

	// Discount metadata
	DiscountCode   string `gorm:"type:varchar(50)" json:"discount_code,omitempty"`
	DiscountReason string `gorm:"type:varchar(255)" json:"discount_reason,omitempty"`

	// Notes & documents
	Notes  string `gorm:"type:text" json:"notes,omitempty"`
	Terms  string `gorm:"type:text" json:"terms,omitempty"`
	PdfURL string `gorm:"type:text" json:"pdf_url,omitempty"`

	// Lifecycle timestamps
	IssuedAt    time.Time  `json:"issued_at"`
	DueAt       *time.Time `json:"due_at,omitempty"`
	SentAt      *time.Time `json:"sent_at,omitempty"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`
	RefundedAt  *time.Time `json:"refunded_at,omitempty"`

	// Ownership
	CreatedBy *uuid.UUID `gorm:"type:uuid" json:"created_by,omitempty"`

	// Relationships
	Items []InvoiceItem `gorm:"foreignKey:InvoiceID;references:ID" json:"items,omitempty"`
	Payments []Payment     `gorm:"foreignKey:InvoiceID;constraint:OnDelete:CASCADE;" json:"payments,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Invoice) TableName() string {
	return "invoices"
}

// RecomputeTotals recalculates subtotal/tax/total/due from the current items and payments slices.
// Call this after any change to items or payments.
func (inv *Invoice) RecomputeTotals() {
    var subtotal int64

    for i := range inv.Items {
        // Calculate each item's total first.
        inv.Items[i].RecomputeLineTotal()
        subtotal += inv.Items[i].LineTotalCents
    }

    inv.SubtotalCents = subtotal

    discount := inv.DiscountCents
    if discount < 0 {
        discount = 0
    }
    if discount > subtotal {
        discount = subtotal
    }

    inv.TotalCents = subtotal - discount

    var paid int64
    for i := range inv.Payments {
        if inv.Payments[i].Status == PaymentSucceeded {
            paid += inv.Payments[i].AmountCents
        }
    }

    inv.AmountPaidCents = paid
    inv.AmountDueCents = inv.TotalCents - paid

    if inv.AmountDueCents < 0 {
        inv.AmountDueCents = 0
    }
}

// RecomputeStatusIfOverdue flips the status to overdue when due_at has passed.
func (inv *Invoice) RecomputeStatusIfOverdue() {
	if inv.DueAt == nil {
		return
	}
	if inv.Status != InvoicePending && inv.Status != InvoicePartial {
		return
	}
	if time.Now().After(*inv.DueAt) {
		inv.Status = InvoiceOverdue
	}
}