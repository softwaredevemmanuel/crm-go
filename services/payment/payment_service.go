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
	ErrPaymentNotFound = errors.New("payment not found")
	ErrInvoiceNotFound = errors.New("invoice not found")
)

type PaymentService struct {
	DB *gorm.DB
}

func NewPaymentService(db *gorm.DB) *PaymentService {
	return &PaymentService{DB: db}
}

// Create records a new payment against an invoice.
func (s *PaymentService) Create(invoiceID uuid.UUID, req dto.CreatePaymentRequest, recordedBy *uuid.UUID) (*models.Payment, error) {
	// 1. Verify invoice exists
	var inv models.Invoice
	if err := s.DB.First(&inv, "id = ?", invoiceID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvoiceNotFound
		}
		return nil, err
	}

	// 2. Setup status & timestamps
	status := models.PaymentSucceeded
	if req.Status != "" {
		status = models.PaymentStatus(req.Status)
	}

	var paidAt *time.Time
	if status == models.PaymentSucceeded {
		now := time.Now()
		paidAt = &now
	}
	if req.PaidAt != nil && *req.PaidAt != "" {
		t, err := time.Parse(time.RFC3339, *req.PaidAt)
		if err == nil {
			paidAt = &t
		}
	}

	// 3. Create Payment Record
	payment := &models.Payment{
		InvoiceID:       invoiceID,
		AmountCents:     req.AmountCents,
		Status:          status,
		Method:          models.PaymentMethod(req.Method),
		IsInstalment:    req.IsInstalment,
		InstalmentNo:    req.InstalmentNo,
		InstalmentTotal: req.InstalmentTotal,
		Notes:           req.Notes,
		PaidAt:          paidAt,
		RecordedBy:      recordedBy,
	}

	// 4. Save & Recalculate Invoice (in a transaction)
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(payment).Error; err != nil {
			return err
		}
		
		// Recalculate invoice totals and status
		return s.recalcInvoice(tx, invoiceID)
	})

	if err != nil {
		return nil, err
	}

	return payment, nil
}

// GetByID fetches a single payment.
func (s *PaymentService) GetByID(id uuid.UUID) (*models.Payment, error) {
	var p models.Payment
	if err := s.DB.First(&p, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}
	return &p, nil
}

// ListPayments paginates and filters payments.
func (s *PaymentService) ListPayments(q dto.ListPaymentsQuery) ([]models.Payment, int64, error) {
	var (
		payments []models.Payment
		total    int64
	)
	query := s.DB.Model(&models.Payment{})

	if q.InvoiceID != "" {
		if id, err := uuid.Parse(q.InvoiceID); err == nil {
			query = query.Where("invoice_id = ?", id)
		}
	}
	if q.Status != "" {
		query = query.Where("status = ?", q.Status)
	}
	if q.Method != "" {
		query = query.Where("method = ?", q.Method)
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
		"created_at": true, "updated_at": true, "amount_cents": true, "paid_at": true,
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
	if err := query.Offset(offset).Limit(q.Limit).Find(&payments).Error; err != nil {
		return nil, 0, err
	}
	return payments, total, nil
}

// Update updates the status/notes of a payment.
func (s *PaymentService) Update(id uuid.UUID, req dto.UpdatePaymentRequest) (*models.Payment, error) {
	p, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Status != nil {
		newStatus := models.PaymentStatus(*req.Status)
		// If transitioning to succeeded, ensure PaidAt is set
		if newStatus == models.PaymentSucceeded && p.PaidAt == nil {
			now := time.Now()
			p.PaidAt = &now
		}
		p.Status = newStatus
	}
	if req.Method != nil {
		p.Method = models.PaymentMethod(*req.Method)
	}
	if req.Notes != nil {
		p.Notes = *req.Notes
	}

	// Save & Recalculate
	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(p).Error; err != nil {
			return err
		}
		return s.recalcInvoice(tx, p.InvoiceID)
	})

	if err != nil {
		return nil, err
	}
	return p, nil
}

// Delete removes a payment and recalculates the invoice.
func (s *PaymentService) Delete(id uuid.UUID) error {
	p, err := s.GetByID(id)
	if err != nil {
		return err
	}

	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(p).Error; err != nil {
			return err
		}
		return s.recalcInvoice(tx, p.InvoiceID)
	})
}

// recalcInvoice is a helper that updates the Invoice's paid amounts.
func (s *PaymentService) recalcInvoice(tx *gorm.DB, invoiceID uuid.UUID) error {
	var inv models.Invoice
	if err := tx.Preload("Items").Preload("Payments").First(&inv, "id = ?", invoiceID).Error; err != nil {
		return err
	}

	// Your Invoice model has this method:
	inv.RecomputeTotals()
	
	// Update only the totals and status of the invoice
	return tx.Model(&inv).Updates(map[string]interface{}{
		"amount_paid_cents": inv.AmountPaidCents,
		"amount_due_cents":  inv.AmountDueCents,
		"status":            inv.Status,
	}).Error
}