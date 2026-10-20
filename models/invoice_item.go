package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvoiceItemType string

const (
	ItemTypeTuition      InvoiceItemType = "tuition"
	ItemTypeCourse       InvoiceItemType = "course"
	ItemTypeMaterial     InvoiceItemType = "material"
	ItemTypeExam         InvoiceItemType = "exam"
	ItemTypeRegistration InvoiceItemType = "registration"
	ItemTypeHostel       InvoiceItemType = "hostel"
	ItemTypeTransport    InvoiceItemType = "transport"
	ItemTypeLibrary      InvoiceItemType = "library"
	ItemTypeUniform      InvoiceItemType = "uniform"
	ItemTypeFine         InvoiceItemType = "fine"
	ItemTypeDiscount     InvoiceItemType = "discount" // negative line
	ItemTypeCustom       InvoiceItemType = "custom"
)

type InvoiceItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TermID    uuid.UUID `gorm:"type:uuid;not null;index" json:"term_id"`
	InvoiceID *uuid.UUID `gorm:"type:uuid;index" json:"invoice_id,omitempty"`

	// Classification
	ItemType    InvoiceItemType `gorm:"type:varchar(30);not null;default:'custom';index" json:"item_type"`

	// Display snapshot — never re-join for historical accuracy
	Title       string `gorm:"type:varchar(255);not null" json:"title"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	SKU         string `gorm:"type:varchar(100)" json:"sku,omitempty"`

	// Pricing (minor units). UnitPriceCents may be negative for discount/credit lines.
	UnitPriceCents int64 `gorm:"not null;default:0" json:"unit_price_cents"`
	Quantity       int   `gorm:"not null;default:1" json:"quantity"`

	// Per-line discount & tax
	DiscountCents   int64   `gorm:"not null;default:0" json:"discount_cents"`
	DiscountPercent float64 `gorm:"type:numeric(5,2);default:0" json:"discount_percent"`
	TaxRatePercent  float64 `gorm:"type:numeric(5,2);default:0" json:"tax_rate_percent"`
	TaxCents        int64   `gorm:"not null;default:0" json:"tax_cents"`

	// Computed line total = (unit * qty) - discount
	LineTotalCents int64 `gorm:"not null;default:0" json:"line_total_cents"`

	// Display
	DisplayOrder int    `gorm:"not null;default:0" json:"display_order"`
	Notes        string `gorm:"type:text" json:"notes,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Term Term `gorm:"foreignKey:TermID;references:ID" json:"term,omitempty"`
}

func (InvoiceItem) TableName() string {
	return "invoice_items"
}

// RecomputeLineTotal recalcs line_total and per-line tax.
func (it *InvoiceItem) RecomputeLineTotal() {
	gross := it.UnitPriceCents * int64(it.Quantity)

	if it.DiscountPercent > 0 && it.DiscountCents == 0 {
		it.DiscountCents = int64(float64(gross) * it.DiscountPercent / 100.0)
	}

	it.LineTotalCents = gross - it.DiscountCents

	if it.TaxRatePercent > 0 {
		it.TaxCents = int64(float64(it.LineTotalCents) * it.TaxRatePercent / 100.0)
	}
}