package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"


)
type Product struct {
	ID               uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name             string     `gorm:"type:varchar(255);not null" json:"name"`
	Description      string     `gorm:"type:text" json:"description"` // Use text for longer descriptions
	Price            float64    `gorm:"type:decimal(12,2);not null" json:"price"` // Use float64 or better: create a Price type
	CompareAtPrice   float64    `gorm:"type:decimal(12,2)" json:"compare_at_price"`
	Image            string     `gorm:"type:varchar(500)" json:"image"` // URLs can be long
	RequiresShipping bool       `gorm:"default:true" json:"requires_shipping"`
	Status           string     `gorm:"type:varchar(20);default:'draft';check:status IN ('draft', 'active', 'inactive', 'discontinued')" json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

}

type CourseProductTable struct {
    ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
   	CourseID   uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE;" `
	ProductID uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE;" `
    CreatedAt  time.Time
}

// Add these methods to your CourseCategory model for proper JSON marshaling/unmarshaling
func (c *CourseProductTable) BeforeCreate(tx *gorm.DB) (err error) {
    if c.ID == uuid.Nil {
        c.ID = uuid.New()
    }
    return
}