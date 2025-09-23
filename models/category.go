package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"

)

type Category struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name       string    `gorm:"size:255;not null" json:"name"`
	Description       string    `gorm:"size:255; default:No Description" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

    // Relationships
    Courses []CourseCategoryTable `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
}

type CategoryInput struct {
	Name    string `json:"name" binding:"required"`
    Description       string    `gorm:"size:255;not null" json:"description"`

}

type CourseCategoryTable struct {
    ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
    CourseID   uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE;" `
    CategoryID uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE;" `
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

type CreateCourseCategoryRequest struct {
    CourseID   string `json:"course_id" binding:"required,uuid4"`
    CategoryID string `json:"category_id" binding:"required,uuid4"`
}

type CreateCourseProductRequest struct {
    CourseID   string `json:"course_id" binding:"required,uuid4"`
    ProductID string `json:"product_id" binding:"required,uuid4"`
}
// Add these methods to your CourseCategory model for proper JSON marshaling/unmarshaling
func (c *CourseCategoryTable) BeforeCreate(tx *gorm.DB) (err error) {
    if c.ID == uuid.Nil {
        c.ID = uuid.New()
    }
    return
}
