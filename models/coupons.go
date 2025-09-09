package models

import (
	"time"
	"github.com/google/uuid"

)

type Coupon struct {
    ID                 uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    
    // Coupon Identification
    Code               string         `gorm:"type:varchar(100);uniqueIndex;not null"` // Discount code
    Name               string         `gorm:"type:varchar(255);not null"`             // Display name
    Description        string         `gorm:"type:text"`                              // Detailed description
    
    // Discount Configuration
    DiscountType       string         `gorm:"type:varchar(20);default:'percentage';check:discount_type IN ('percentage', 'fixed', 'free', 'bogo')"` // BOGO = Buy One Get One
    DiscountValue      float64        `gorm:"type:decimal(10,2);not null"`            // 10.00 for 10% or $10
    Currency           string         `gorm:"type:varchar(3);default:'USD'"`          // For fixed amounts
    
    // Usage Limits
    UsageLimit         int            `gorm:"default:0"`          // Total usage limit (0 = unlimited)

    
    // Validity Period
    ValidFrom          *time.Time     `gorm:"index"`              // Start date/time
    ValidUntil         *time.Time     `gorm:"index"`              // Expiration date/time
    IsActive           bool           `gorm:"default:true;index"` // Master activation switch
    

    
    // User Targeting
    UserIDs            []uuid.UUID    `gorm:"type:uuid[]"`        // Specific users only

    
    
    // Status & Audit
    Status             string         `gorm:"type:varchar(20);default:'active';check:status IN ('draft', 'active', 'paused', 'expired', 'archived')"`

    // Relationships
    Creator            User           `gorm:"foreignKey:CreatedBy"`
    Courses            []Course       `gorm:"many2many:coupon_courses;"`
    Categories         []Category     `gorm:"many2many:coupon_categories;"`
    
    // Timestamps
    CreatedAt          time.Time
    UpdatedAt          time.Time
    ApprovedAt         *time.Time
    ArchivedAt         *time.Time
}

// TableName specifies the table name
func (Coupon) TableName() string {
    return "coupons"
}