package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentSucceeded PaymentStatus = "succeeded"
	PaymentFailed    PaymentStatus = "failed"
	PaymentRefunded  PaymentStatus = "refunded"
	PaymentPartial   PaymentStatus = "partial_refund"
)

type PaymentMethod string

const (
	PaymentMethodCard        PaymentMethod = "card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodCash        PaymentMethod = "cash"
	PaymentMethodWallet      PaymentMethod = "wallet"
	PaymentMethodPaystack    PaymentMethod = "paystack"
	PaymentMethodFlutterwave PaymentMethod = "flutterwave"
	PaymentMethodStripe      PaymentMethod = "stripe"
	PaymentMethodCheque      PaymentMethod = "cheque"
)

type Payment struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	InvoiceID uuid.UUID `gorm:"type:uuid;not null;index" json:"invoice_id"`

	// Amounts (minor units)
	AmountCents         int64 `gorm:"not null" json:"amount_cents"`          // amount of this payment

	// Payment metadata
	Status        PaymentStatus `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	Method        PaymentMethod `gorm:"type:varchar(30);index" json:"method"`


	// Instalment tracking
	IsInstalment    bool `gorm:"default:false" json:"is_instalment"`
	InstalmentNo    int  `gorm:"default:0" json:"instalment_no,omitempty"`
	InstalmentTotal int  `gorm:"default:0" json:"instalment_total,omitempty"`


	// Timestamps
	PaidAt      *time.Time `json:"paid_at,omitempty"`

	// Notes / audit
	Notes       string     `gorm:"type:text" json:"notes,omitempty"`
	RecordedBy  *uuid.UUID `gorm:"type:uuid" json:"recorded_by,omitempty"` // admin who recorded it

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Invoice Invoice `gorm:"foreignKey:InvoiceID" json:"-"`
}

func (Payment) TableName() string {
	return "payments"
}
