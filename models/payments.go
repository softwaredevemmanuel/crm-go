package models

import (
	"time"
	"github.com/google/uuid"

)
type Payment struct {
    ID                 uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    
    // Transaction Identification
    PaymentID          string         `gorm:"type:varchar(100);uniqueIndex;not null"` // Internal reference
    GatewayReference   string         `gorm:"type:varchar(200);index"` // Gateway transaction ID
    GatewaySessionID   string         `gorm:"type:varchar(200);index"` // Gateway session ID
    
    // Payer Information
    PayerID            uuid.UUID      `gorm:"type:uuid;not null;index"` // Student/user ID
    PayerEmail         string         `gorm:"type:varchar(255);index"`
    PayerName          string         `gorm:"type:varchar(255)"`
    PayerPhone         string         `gorm:"type:varchar(50)"`
    
    // Payment Details
    Amount             float64        `gorm:"type:decimal(12,2);not null"` // Gross amount
    Currency           string         `gorm:"type:varchar(3);default:'USD';not null"`

    
    // Payment Method
    PaymentMethod      string         `gorm:"type:varchar(50);not null;check:payment_method IN ('credit_card', 'debit_card', 'paypal', 'bank_transfer', 'wallet', 'crypto', 'cash', 'check')"`

    
    // Gateway & Processing
    Gateway            string         `gorm:"type:varchar(50);not null;check:gateway IN ('stripe', 'paypal', 'razorpay', 'paystack', 'flutterwave', 'bank', 'manual')"`

    
    // Status & Lifecycle
    Status             string         `gorm:"type:varchar(20);default:'pending';check:status IN ('pending', 'processing', 'completed', 'failed', 'refunded', 'partially_refunded', 'disputed', 'cancelled', 'expired')"`
    FailureReason      string         `gorm:"type:varchar(255)"` // Why payment failed
    FailureCode        string         `gorm:"type:varchar(100)"` // Error code
    
    // Timing
    InitiatedAt        time.Time      `gorm:"not null"`
    ProcessedAt        *time.Time     // When payment was processed
    ExpiresAt          *time.Time     // For pending payments
    RefundedAt         *time.Time     // When refund was issued
    
    // Related Entities
    InvoiceID          *uuid.UUID     `gorm:"type:uuid;index"` // Associated invoice
    OrderID            *uuid.UUID     `gorm:"type:uuid;index"` // Associated order
    SubscriptionID     *uuid.UUID     `gorm:"type:varchar(100);index"` // For recurring payments
    
    // Course/Product Context
    CourseID           *uuid.UUID     `gorm:"type:uuid;index"`
    CourseName         string         `gorm:"type:varchar(255)"`
    
 
   
    
    // Relationships
    Payer              User           `gorm:"foreignKey:PayerID"`
    Course             Course         `gorm:"foreignKey:CourseID"`
    
    // Timestamps
    CreatedAt          time.Time
    UpdatedAt          time.Time
}

// TableName specifies the table name
func (Payment) TableName() string {
    return "payments"
}