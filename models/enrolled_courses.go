package models

import (
	"time"
	"github.com/google/uuid"


)

type Enrollment struct {
    ID                 uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
   	StudentID uuid.UUID `gorm:"type:uuid;not null;index:idx_student_course,unique" json:"student_id"`
	CourseID  uuid.UUID `gorm:"type:uuid;not null;index:idx_student_course,unique" json:"course_id"`
    Status             string     `gorm:"type:varchar(20);default:'pending';check:status IN ('pending', 'active', 'completed', 'cancelled', 'expired', 'suspended')" json:"status"`
    EnrollmentDate     time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"enrollment_date"`
    StartDate          *time.Time `json:"start_date,omitempty"`
    CompletionDate     *time.Time `json:"completion_date,omitempty"`
    ExpirationDate     *time.Time `json:"expiration_date,omitempty"`
    LastAccessed       *time.Time `json:"last_accessed,omitempty"`
    ProgressPercentage int        `gorm:"default:0;check:progress_percentage >= 0 AND progress_percentage <= 100" json:"progress_percentage"`
    TotalTimeSpent     int        `gorm:"default:0" json:"total_time_spent"` // seconds
    FinalGrade         *float64   `gorm:"type:decimal(5,2)" json:"final_grade,omitempty"`
	
    CertificateIssued  bool       `gorm:"default:false" json:"certificate_issued"`
    CertificateID      *uuid.UUID `gorm:"type:uuid" json:"certificate_id,omitempty"`
    PricePaid          float64    `gorm:"type:decimal(12,2);default:0" json:"price_paid"`
    Currency           string     `gorm:"type:varchar(3);default:'USD'" json:"currency"`
    PaymentMethod      string     `gorm:"type:varchar(50)" json:"payment_method,omitempty"`
    PaymentStatus      string     `gorm:"type:varchar(20);default:'pending';check:payment_status IN ('pending', 'paid', 'failed', 'refunded', 'free')" json:"payment_status"`
    TransactionID      string     `gorm:"type:varchar(255)" json:"transaction_id,omitempty"`
    DiscountApplied    float64    `gorm:"type:decimal(12,2);default:0" json:"discount_applied"`
    CouponCode         string     `gorm:"type:varchar(100)" json:"coupon_code,omitempty"`
    AccessLevel        string     `gorm:"type:varchar(20);default:'full';check:access_level IN ('full', 'trial', 'preview', 'limited')" json:"access_level"`
    TrialEndsAt        *time.Time `json:"trial_ends_at,omitempty"`
 

    // Relationships
    Student    User   `gorm:"foreignKey:StudentID" json:"student,omitempty"`
    Course     Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
    Certificate Certificate `gorm:"foreignKey:CertificateID" json:"certificate,omitempty"`
}

// Add this index in your migration
// CREATE UNIQUE INDEX idx_enrollments_student_course ON enrollments(student_id, course_id) WHERE deleted_at IS NULL;