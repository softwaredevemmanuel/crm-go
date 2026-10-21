package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"crm-go/dto"
	"crm-go/models"
)

var (
	ErrInvoiceNotFound    = errors.New("invoice not found")
	ErrInvoiceNotEditable = errors.New("only draft or pending invoices can be edited")
)

type InvoiceService struct {
	DB *gorm.DB
}

func NewInvoiceService(db *gorm.DB) *InvoiceService {
	return &InvoiceService{DB: db}
}

// ---------- Invoice number generator ----------
func (s *InvoiceService) nextInvoiceNumber() (string, error) {
	year := time.Now().Year()
	prefix := fmt.Sprintf("INV-%d-", year)

	var last models.Invoice
	err := s.DB.Unscoped().
		Where("invoice_number LIKE ?", prefix+"%").
		Order("invoice_number DESC").
		First(&last).Error

	seq := 1
	if err == nil {
		var n int
		if _, e := fmt.Sscanf(last.InvoiceNumber, prefix+"%d", &n); e == nil {
			seq = n + 1
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}
	return fmt.Sprintf("%s%05d", prefix, seq), nil
}

// ---------- Create ----------
func (s *InvoiceService) Create(req dto.CreateInvoiceRequest, createdBy *uuid.UUID) (*models.Invoice, error) {
	var inv *models.Invoice

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		num, err := s.nextInvoiceNumber()
		if err != nil {
			return err
		}

		// 1. Find the requested catalog items
		var catalogItems []models.InvoiceItem
		if err := tx.Where("id IN ?", req.ItemIDs).Find(&catalogItems).Error; err != nil {
			return err
		}

		if len(catalogItems) != len(req.ItemIDs) {
			return errors.New("one or more invoice items not found")
		}

		// 2. Create the Invoice Header
		inv = &models.Invoice{
			InvoiceNumber:  num,
			Status:         models.InvoiceStatus(req.Status),
			StudentID:      req.StudentID,
			Currency:       req.Currency,
			DiscountCents:  req.DiscountCents,
			TaxRatePercent: req.TaxRatePercent,
			TaxLabel:       req.TaxLabel,
			Notes:          req.Notes,
			IssuedAt:       time.Now(),
			CreatedBy:      createdBy,
			Items:          catalogItems, // <--- Attach the products here!
		}
		
		if req.DueAt != nil && *req.DueAt != "" {
			t, _ := time.Parse(time.RFC3339, *req.DueAt)
			inv.DueAt = &t
		}

		// 3. Save the Invoice. 
		// Because we added `Items`, GORM will automatically:
		// - Create the Invoice
		// - Insert rows into `invoice_line_items` linking this invoice to the catalog items
		if err := tx.Create(inv).Error; err != nil {
			return err
		}

		// 4. Recalculate totals
		inv.RecomputeTotals()

		if err := tx.Save(inv).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return inv, nil
}


// ---------- Get ----------
func (s *InvoiceService) GetByID(id uuid.UUID) (*models.Invoice, error) {
	var inv models.Invoice
	err := s.DB.
		Preload("Student"). // <-- ADD THIS
		Preload("Items").
		Preload("Payments").
		First(&inv, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInvoiceNotFound
	}
	return &inv, err
}

// ---------- List ----------
func (s *InvoiceService) List(q dto.ListInvoicesQuery) ([]models.Invoice, int64, error) {
	var (
		invoices []models.Invoice
		total    int64
	)
	query := s.DB.Model(&models.Invoice{})

	if q.Search != "" {
		s := "%" + strings.ToLower(q.Search) + "%"
		query = query.Where("LOWER(invoice_number) LIKE ?", s)
	}
	if q.Status != "" {
		query = query.Where("status = ?", q.Status)
	}
	if q.StudentID != "" {
		if uid, err := uuid.Parse(q.StudentID); err == nil {
			query = query.Where("student_id = ?", uid)
		}
	}
	if q.From != "" {
		if t, err := time.Parse(time.RFC3339, q.From); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if q.To != "" {
		if t, err := time.Parse(time.RFC3339, q.To); err == nil {
			query = query.Where("created_at <= ?", t)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	allowedSort := map[string]bool{
		"created_at": true, "updated_at": true, "invoice_number": true,
		"total_cents": true, "status": true, "issued_at": true,
	}
	sortBy := q.SortBy
	if !allowedSort[sortBy] {
		sortBy = "created_at"
	}
	order := strings.ToLower(q.SortOrder)
	if order != "asc" && order != "desc" {
		order = "desc"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, order))

	offset := (q.Page - 1) * q.Limit
	
	// --- HERE IS THE FIX ---
	// Preload the Student relationship, and optionally only select the columns we need
	if err := query.
		Preload("Student", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "first_name", "last_name", "email") // Only fetch what we need (optional but good for performance)
		}).
		Offset(offset).Limit(q.Limit).
		Find(&invoices).Error; err != nil {
		return nil, 0, err
	}
	
	return invoices, total, nil
}

// ---------- Update ----------
func (s *InvoiceService) Update(id uuid.UUID, req dto.UpdateInvoiceRequest) (*models.Invoice, error) {
	inv, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Prevent editing critical fields on paid/cancelled/refunded invoices
	if inv.Status == models.InvoicePaid || inv.Status == models.InvoiceCancelled || inv.Status == models.InvoiceRefunded {
		return nil, ErrInvoiceNotEditable
	}

	// Start Transaction
	err = s.DB.Transaction(func(tx *gorm.DB) error {
		
		// 1. Update scalar fields
		if req.StudentID != nil {
			inv.StudentID = req.StudentID
		}
		if req.Currency != nil {
			inv.Currency = strings.ToUpper(*req.Currency)
		}
		if req.DiscountCents != nil {
			inv.DiscountCents = *req.DiscountCents
		}
		if req.DiscountCode != nil {
			inv.DiscountCode = *req.DiscountCode
		}
		if req.DiscountReason != nil {
			inv.DiscountReason = *req.DiscountReason
		}
		if req.TaxRatePercent != nil {
			inv.TaxRatePercent = *req.TaxRatePercent
		}
		if req.TaxLabel != nil {
			inv.TaxLabel = *req.TaxLabel
		}
		if req.Notes != nil {
			inv.Notes = *req.Notes
		}
		if req.Terms != nil {
			inv.Terms = *req.Terms
		}
		if req.PdfURL != nil {
			inv.PdfURL = *req.PdfURL
		}
		if req.DueAt != nil {
			if *req.DueAt == "" {
				inv.DueAt = nil
			} else {
				t, err := time.Parse(time.RFC3339, *req.DueAt)
				if err != nil {
					return fmt.Errorf("invalid due_at: %w", err)
				}
				inv.DueAt = &t
			}
		}
		if req.Status != nil {
			newStatus := models.InvoiceStatus(*req.Status)
			// Auto-manage timestamps
			if newStatus == models.InvoicePaid && inv.PaidAt == nil {
				now := time.Now()
				inv.PaidAt = &now
			}
			if newStatus == models.InvoiceCancelled && inv.CancelledAt == nil {
				now := time.Now()
				inv.CancelledAt = &now
			}
			if newStatus == models.InvoiceRefunded && inv.RefundedAt == nil {
				now := time.Now()
				inv.RefundedAt = &now
			}
			inv.Status = newStatus
		}

		// 2. Handle Items Update (Many-to-Many)
		// If ItemIDs are provided, we swap the current line items
		if len(req.ItemIDs) > 0 {
			// Fetch the new catalog items to attach
			var newItems []models.InvoiceItem
			if err := tx.Where("id IN ?", req.ItemIDs).Find(&newItems).Error; err != nil {
				return err
			}
			
			// Safety check
			if len(newItems) != len(req.ItemIDs) {
				return errors.New("one or more invoice items not found")
			}

			// Replace the associations. 
			// GORM will automatically remove old links in the pivot table 
			// and insert the new ones.
			if err := tx.Model(inv).Association("Items").Replace(newItems); err != nil {
				return err
			}
			
			// Set the local struct so RecomputeTotals has the right data
			inv.Items = newItems
		}

		// 3. Recompute totals based on the (possibly new) items
		inv.RecomputeTotals()
		inv.RecomputeStatusIfOverdue()

		// 4. Save the invoice header (totals, status, etc.)
		if err := tx.Save(inv).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Reload to ensure all relationships are fresh before returning to the controller
	return s.GetByID(id)
}

// ---------- Delete ----------
func (s *InvoiceService) Delete(id uuid.UUID) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Unlink all items first
		if err := tx.Model(&models.InvoiceItem{}).Where("invoice_id = ?", id).Update("invoice_id", nil).Error; err != nil {
			return err
		}

		// 2. Delete the invoice
		if err := tx.Delete(&models.Invoice{}, "id = ?", id).Error; err != nil {
			return err
		}

		return nil
	})
}

// ---------- Recalculate ----------
// Loads items + payments from DB and recomputes totals + status.
func (s *InvoiceService) Recalculate(invoiceID uuid.UUID) error {
	inv, err := s.GetByID(invoiceID)
	if err != nil {
		return err
	}
	inv.RecomputeTotals()
	inv.RecomputeStatusIfOverdue()
	return s.DB.Save(inv).Error
}

// ---------- Summary ----------
func (s *InvoiceService) Summary() (*dto.InvoiceSummaryResponse, error) {
	summary := &dto.InvoiceSummaryResponse{}

	if err := s.DB.Model(&models.Invoice{}).Count(&summary.TotalInvoices).Error; err != nil {
		return nil, err
	}
	_ = s.DB.Model(&models.Invoice{}).Select("COALESCE(SUM(total_cents),0)").Scan(&summary.TotalBilledCents).Error
	_ = s.DB.Model(&models.Invoice{}).Select("COALESCE(SUM(amount_paid_cents),0)").Scan(&summary.TotalPaidCents).Error
	_ = s.DB.Model(&models.Invoice{}).Select("COALESCE(SUM(amount_due_cents),0)").Scan(&summary.TotalDueCents).Error
	_ = s.DB.Model(&models.Invoice{}).Where("status = ?", models.InvoiceOverdue).Count(&summary.OverdueCount).Error

	return summary, nil
}
