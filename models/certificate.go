package models

import (
	"time"
	"github.com/google/uuid"


)

type Certificate struct {
    ID                   uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
    StudentID            uuid.UUID      `gorm:"type:uuid;not null" json:"student_id"`
    CourseID             uuid.UUID      `gorm:"type:uuid;not null" json:"course_id"`
    
    // Certificate Identification
    CertificateNumber    string         `gorm:"type:varchar(100);unique;not null" json:"certificate_number"`
    CertificateType      string         `gorm:"type:varchar(30);default:'completion';check:certificate_type IN ('completion', 'achievement', 'participation', 'excellence', 'professional')" json:"certificate_type"`
    
    // Certificate Content
    Title                string         `gorm:"type:varchar(255);not null" json:"title"`
    Description          string         `gorm:"type:text" json:"description,omitempty"`
    IssuingOrganization  string         `gorm:"type:varchar(255);default:'Your Organization Name'" json:"issuing_organization"`
    IssuerName           string         `gorm:"type:varchar(255)" json:"issuer_name,omitempty"`
    IssuerTitle          string         `gorm:"type:varchar(255)" json:"issuer_title,omitempty"`
    IssuerSignatureURL   string         `gorm:"type:varchar(500)" json:"issuer_signature_url,omitempty"`
    OrganizationLogoURL  string         `gorm:"type:varchar(500)" json:"organization_logo_url,omitempty"`
    

    // Student Performance
    FinalGrade           *float64       `gorm:"type:decimal(5,2)" json:"final_grade,omitempty"`
    GradeScale           string         `gorm:"type:varchar(50);default:'percentage'" json:"grade_scale,omitempty"`
    PerformanceLevel     string         `gorm:"type:varchar(50)" json:"performance_level,omitempty"`
    CreditsEarned        *float64       `gorm:"type:decimal(4,1)" json:"credits_earned,omitempty"`
    TotalHours           *float64       `gorm:"type:decimal(6,2)" json:"total_hours,omitempty"`
    
    // Certificate Details
    TemplateID           *uuid.UUID     `gorm:"type:uuid" json:"template_id,omitempty"`
    BackgroundURL        string         `gorm:"type:varchar(500)" json:"background_url,omitempty"`
    CertificateURL       string         `gorm:"type:varchar(500);not null" json:"certificate_url"`
    ThumbnailURL         string         `gorm:"type:varchar(500)" json:"thumbnail_url,omitempty"`
    DigitalBadgeURL      string         `gorm:"type:varchar(500)" json:"digital_badge_url,omitempty"`
    
    // Verification
    VerificationHash     string         `gorm:"type:varchar(255);unique;not null" json:"verification_hash"`
    VerificationURL      string         `gorm:"type:varchar(500);not null" json:"verification_url"`
    QRCodeURL            string         `gorm:"type:varchar(500)" json:"qr_code_url,omitempty"`
    IsVerifiable         bool           `gorm:"default:true" json:"is_verifiable"`
    VerificationCount    int            `gorm:"default:0" json:"verification_count"`
    LastVerifiedAt       *time.Time     `json:"last_verified_at,omitempty"`
    
    
    // Accreditation
    AccreditationBody    string         `gorm:"type:varchar(255)" json:"accreditation_body,omitempty"`
    AccreditationID      string         `gorm:"type:varchar(100)" json:"accreditation_id,omitempty"`
    CEUCredits           *float64       `gorm:"type:decimal(4,1)" json:"ceu_credits,omitempty"`
    
  
    
    // Relationships
    Student              User           `gorm:"foreignKey:StudentID" json:"student,omitempty"`
    Course               Course         `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}