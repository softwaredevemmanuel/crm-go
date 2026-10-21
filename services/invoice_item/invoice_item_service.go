package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"crm-go/dto"
	"crm-go/models"
)

var (
	ErrInvoiceItemNotFound = errors.New("invoice item not found")
	ErrTermNotFound        = errors.New("term not found")
)

type InvoiceItemService struct {
	DB *gorm.DB
}

func NewInvoiceItemService(db *gorm.DB) *InvoiceItemService {
	return &InvoiceItemService{DB: db}
}

// ---------- Create ----------
func (s *InvoiceItemService) Create(req dto.CreateInvoiceItemRequest,) (*models.InvoiceItem, error) {
	// Validate parent term exists.
	var count int64
	if err := s.DB.Model(&models.Term{}).
		Where("id = ?", req.TermID).
		Count(&count).Error; err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, ErrTermNotFound
	}

	title := strings.TrimSpace(req.Title)
	itemType := models.InvoiceItemType(req.ItemType)

	// Check for duplicate item type, term, and title.
	var existing models.InvoiceItem
	err := s.DB.
		Where(
			"item_type = ? AND term_id = ? AND LOWER(TRIM(title)) = LOWER(?) AND LOWER(TRIM(description)) = LOWER(?)",
			itemType,
			req.TermID,
			title,
			strings.TrimSpace(req.Description),
		).
		First(&existing).Error

	if err == nil {
		return nil, errors.New(
			"an invoice item with this type and title already exists for this term",
		)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	item := &models.InvoiceItem{
		TermID:          req.TermID,
		ItemType:        itemType,
		Title:           title,
		Description:     strings.TrimSpace(req.Description),
		SKU:             strings.TrimSpace(req.SKU),
		UnitPriceCents:  req.UnitPriceCents,
		Quantity:        req.Quantity,
		DiscountCents:   req.DiscountCents,
		DiscountPercent: req.DiscountPercent,
		DisplayOrder:    req.DisplayOrder,
		Notes:           req.Notes,
	}

	item.RecomputeLineTotal()

	if err := s.DB.Create(item).Error; err != nil {
		return nil, err
	}

	// Reload with Term preloaded for response.
	return s.GetByID(item.ID)
}

// ---------- Get ----------
func (s *InvoiceItemService) GetByID(id uuid.UUID) (*models.InvoiceItem, error) {
	var item models.InvoiceItem
	err := s.DB.
		Preload("Term").
		First(&item, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInvoiceItemNotFound
	}
	return &item, err
}

// ---------- List by term ----------
func (s *InvoiceItemService) ListByTerm(termID uuid.UUID) ([]models.InvoiceItem, error) {
	var items []models.InvoiceItem
	err := s.DB.
		Preload("Term").
		Where("term_id = ?", termID).
		Order("display_order ASC, created_at ASC").
		Find(&items).Error
	return items, err
}

// ---------- List (paginated, filterable) ----------
func (s *InvoiceItemService) List(q dto.ListInvoiceItemsQuery) ([]models.InvoiceItem, int64, error) {
	var (
		items []models.InvoiceItem
		total int64
	)

	query := s.DB.Model(&models.InvoiceItem{})

	if q.Search != "" {
		s := "%" + strings.ToLower(q.Search) + "%"
		query = query.Where(
			"LOWER(title) LIKE ? OR LOWER(description) LIKE ? OR LOWER(sku) LIKE ?",
			s, s, s,
		)
	}
	if q.ItemType != "" {
		query = query.Where("item_type = ?", q.ItemType)
	}
	if q.TermID != "" {
		if tid, err := uuid.Parse(q.TermID); err == nil {
			query = query.Where("term_id = ?", tid)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	allowedSort := map[string]bool{
		"display_order": true, "created_at": true, "updated_at": true,
		"title": true, "unit_price_cents": true, "item_type": true,
	}
	sortBy := q.SortBy
	if !allowedSort[sortBy] {
		sortBy = "display_order"
	}
	sortOrder := strings.ToLower(q.SortOrder)
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

	offset := (q.Page - 1) * q.Limit
	
	// Execute the query with Preload
	if err := query.
		Preload("Term", func(db *gorm.DB) *gorm.DB {
			// Optional: Only select the ID and Name columns from the Terms table
			return db.Select("id", "name", "description") 
		}).
		Offset(offset).Limit(q.Limit).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// ---------- Update ----------
func (s *InvoiceItemService) Update(id uuid.UUID, req dto.UpdateInvoiceItemRequest) (*models.InvoiceItem, error) {
	item, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.ItemType != nil {
		item.ItemType = models.InvoiceItemType(*req.ItemType)
	}
	if req.Title != nil {
		item.Title = strings.TrimSpace(*req.Title)
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	if req.SKU != nil {
		item.SKU = strings.TrimSpace(*req.SKU)
	}
	if req.UnitPriceCents != nil {
		item.UnitPriceCents = *req.UnitPriceCents
	}
	if req.Quantity != nil {
		item.Quantity = *req.Quantity
	}
	if req.DiscountCents != nil {
		item.DiscountCents = *req.DiscountCents
	}
	if req.DiscountPercent != nil {
		item.DiscountPercent = *req.DiscountPercent
		// Reset fixed discount so the percentage drives the amount again
		item.DiscountCents = 0
	}
	if req.DisplayOrder != nil {
		item.DisplayOrder = *req.DisplayOrder
	}
	if req.Notes != nil {
		item.Notes = *req.Notes
	}

	item.RecomputeLineTotal()

	if err := s.DB.Save(item).Error; err != nil {
		return nil, err
	}

	return s.GetByID(item.ID)
}

// ---------- Delete ----------
func (s *InvoiceItemService) Delete(id uuid.UUID) error {
	item, err := s.GetByID(id)
	if err != nil {
		return err
	}
	return s.DB.Delete(item).Error
}

// ---------- Delete by term ----------
func (s *InvoiceItemService) DeleteByTerm(termID uuid.UUID) error {
	return s.DB.
		Where("term_id = ?", termID).
		Delete(&models.InvoiceItem{}).Error
}

