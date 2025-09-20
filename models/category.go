package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"

)

type Category struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name       string    `gorm:"size:255;not null" json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CategoryInput struct {
	Name    string `json:"name" binding:"required"`
}

type CourseCategoryTable struct {
    ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
    CourseID   uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE;" `
    CategoryID uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE;" `
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

// Add these methods to your CourseCategory model for proper JSON marshaling/unmarshaling
func (c *CourseCategoryTable) BeforeCreate(tx *gorm.DB) (err error) {
    if c.ID == uuid.Nil {
        c.ID = uuid.New()
    }
    return
}
