package dto

import (
	"time"

	"github.com/google/uuid"

	"crm-go/models"
)

// CreateInvoiceItemRequest — payload for creating a line item.
type CreateInvoiceItemRequest struct {
	TermID          uuid.UUID `json:"term_id" binding:"required"`
	ItemType        string    `json:"item_type" binding:"required,oneof=tuition course material exam registration hostel transport library uniform fine discount custom"`
	Title           string    `json:"title" binding:"required,min=1,max=255"`
	Description     string    `json:"description" binding:"max=5000"`
	SKU             string    `json:"sku" binding:"omitempty,max=100"`
	UnitPriceCents  int64     `json:"unit_price_cents"`
	Quantity        int       `json:"quantity" binding:"required,gte=1"`
	DiscountCents   int64     `json:"discount_cents" binding:"gte=0"`
	DiscountPercent float64   `json:"discount_percent" binding:"gte=0,lte=100"`
	DisplayOrder    int       `json:"display_order"`
	Notes           string    `json:"notes" binding:"max=2000"`
}

// UpdateInvoiceItemRequest — partial update.
// NOTE: TermID is intentionally immutable — move an item by deleting + recreating.
type UpdateInvoiceItemRequest struct {
	ItemType        *string  `json:"item_type" binding:"omitempty,oneof=tuition course material exam registration hostel transport library uniform fine discount custom"`
	Title           *string  `json:"title" binding:"omitempty,min=1,max=255"`
	Description     *string  `json:"description" binding:"omitempty,max=5000"`
	SKU             *string  `json:"sku" binding:"omitempty,max=100"`
	UnitPriceCents  *int64   `json:"unit_price_cents"`
	Quantity        *int     `json:"quantity" binding:"omitempty,gte=1"`
	DiscountCents   *int64   `json:"discount_cents" binding:"omitempty,gte=0"`
	DiscountPercent *float64 `json:"discount_percent" binding:"omitempty,gte=0,lte=100"`
	DisplayOrder    *int     `json:"display_order"`
	Notes           *string  `json:"notes" binding:"omitempty,max=2000"`
}

// TermMiniResponse — lightweight Term info embedded in the item response.
type TermMiniResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name,omitempty"`
	StartDate *string   `json:"start_date,omitempty"`
	EndDate   *string   `json:"end_date,omitempty"`
}

// InvoiceItemResponse — API shape returned to the client.
type InvoiceItemResponse struct {
	ID              uuid.UUID         `json:"id"`
	TermID          uuid.UUID         `json:"term_id"`
	Term            *TermMiniResponse `json:"term,omitempty"` // populated when preloaded

	ItemType        string `json:"item_type"`
	Title           string `json:"title"`
	Description     string `json:"description,omitempty"`
	SKU             string `json:"sku,omitempty"`
	UnitPriceCents  int64  `json:"unit_price_cents"`
	Quantity        int    `json:"quantity"`
	DiscountCents   int64  `json:"discount_cents"`
	DiscountPercent float64 `json:"discount_percent"`
	LineTotalCents  int64  `json:"line_total_cents"` // derived, never stored
	DisplayOrder    int    `json:"display_order"`
	Notes           string `json:"notes,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FromInvoiceItemModel maps the model to the response.
// Set includeTerm=true to embed the preloaded Term.
func FromInvoiceItemModel(it *models.InvoiceItem, includeTerm bool) InvoiceItemResponse {
	lineTotal := it.UnitPriceCents*int64(it.Quantity) - it.DiscountCents
	if lineTotal < 0 {
		lineTotal = 0
	}

	resp := InvoiceItemResponse{
		ID:              it.ID,
		TermID:          it.TermID,
		ItemType:        string(it.ItemType),
		Title:           it.Title,
		Description:     it.Description,
		SKU:             it.SKU,
		UnitPriceCents:  it.UnitPriceCents,
		Quantity:        it.Quantity,
		DiscountCents:   it.DiscountCents,
		DiscountPercent: it.DiscountPercent,
		LineTotalCents:  lineTotal,
		DisplayOrder:    it.DisplayOrder,
		Notes:           it.Notes,
		CreatedAt:       it.CreatedAt,
		UpdatedAt:       it.UpdatedAt,
	}

	if includeTerm && it.TermID != uuid.Nil {
		// InvoiceItem does not carry a loaded Term on the model itself.
		// The caller can populate resp.Term separately if needed.
		_ = includeTerm
	}

	return resp
}

// ListInvoiceItemsQuery — filters + pagination.
type ListInvoiceItemsQuery struct {
	Page      int    `form:"page,default=1" binding:"min=1"`
	Limit     int    `form:"limit,default=50" binding:"min=1,max=200"`
	Search    string `form:"search"`
	ItemType  string `form:"item_type"`
	TermID    string `form:"term_id"`
	SortBy    string `form:"sort_by,default=display_order"`
	SortOrder string `form:"sort_order,default=asc"`
}

type SingleInvoiceItemResponse struct {
	Data InvoiceItemResponse `json:"data"`
}

type PaginatedInvoiceItemsResponse struct {
	Data       []InvoiceItemResponse `json:"data"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalPages int                   `json:"total_pages"`
}

type ReorderInvoiceItemRequest struct {
	ID           uuid.UUID `json:"id" binding:"required"`
	DisplayOrder int       `json:"display_order" binding:"required"`
}

type ReorderInvoiceItemsRequest struct {
	Items []ReorderInvoiceItemRequest `json:"items" binding:"required,min=1,dive"`
}

