package annoncements

import (
	"github.com/gin-gonic/gin"
	"crm-go/models"
	"crm-go/config"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Error string `json:"error" example:"Invalid announcement ID"`
}
type FailureResponse struct {
	Error string `json:"error" example:"Failed to create announcement"`
}
type NotFoundResponse struct {
	Error string `json:"error" example:"Announcement not found"`
}
type SuccessResponse struct {
	Message string `json:"message" example:"Announcement created successfully"`
}
type DeleteSuccessResponse struct {
	Message string `json:"message" example:"Announcement deleted successfully"`
}

type AnnouncementInput struct {
	Title     string    `json:"title" example:"System Maintenance"`
	Message   string    `json:"message" example:"The platform will be unavailable from 2AM to 4AM."`
	Type      string    `json:"type" example:"maintenance"` // general, update, maintenance, urgent
	Audience  string    `json:"audience" example:"all"`      // all, students, tutors, admins
	CreatedBy uuid.UUID `json:"created_by" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`

	StartDate *time.Time `json:"start_date,omitempty" example:"2026-01-22T02:00:00Z"`
	EndDate   *time.Time `json:"end_date,omitempty" example:"2026-01-22T04:00:00Z"`

	IsPinned bool `json:"is_pinned" example:"true"`
}

// CreateAnnouncement godoc
// @Summary Create a new announcement
// @Description Create a new announcement
// @Tags Announcements
// @Accept json
// @Produce json
// @Param announcement body AnnouncementInput true "Announcement"
// @Success 201 {object} models.Announcement
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/announcements [post]
// @Security BearerAuth
func CreateAnnouncement(c *gin.Context) {
	var input AnnouncementInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	announcements := models.Announcement{
		Title:     input.Title,
		Message:   input.Message,
		Type:      input.Type,
		Audience:  input.Audience,
		CreatedBy: input.CreatedBy,
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
		IsPinned:  input.IsPinned,
	}

	if err := config.DB.Create(&announcements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create announcement"})
		return
	}

	c.JSON(http.StatusCreated, announcements)
}


// GetAnnouncements godoc
// @Summary Get all announcements
// @Description Get all announcements
// @Tags Announcements
// @Accept json
// @Produce json
// @Success 200 {array} models.Announcement
// @Failure 500 {object} FailureResponse
// @Router /announcements [get]
func GetAnnouncements(c *gin.Context) {
	var announcements []models.Announcement

	if err := config.DB.
		Preload("UserDetails").
		Order("is_pinned DESC, created_at DESC").
		Find(&announcements).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch announcements",
		})
		return
	}

	c.JSON(http.StatusOK, announcements)
}



// GetAnnouncementByID godoc
// @Summary Get an announcement by ID
// @Description Get an announcement by ID
// @Tags Announcements
// @Accept json
// @Produce json
// @Param id path string true "Announcement ID"
// @Success 200 {object} models.Announcement
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 500 {object} FailureResponse
// @Router /announcements/{id} [get]
func GetAnnouncementByID(c *gin.Context) {
	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid announcement ID"})
		return
	}

	var announcement models.Announcement

	err = config.DB.
		Preload("UserDetails").
		First(&announcement, "id = ?", uid).Error

	if err := config.DB.First(&announcement, "id = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Announcement not found"})
		return
	}

	c.JSON(http.StatusOK, announcement)
}

// UpdateAnnouncement godoc
// @Summary Update an announcement by ID
// @Description Update an announcement by ID
// @Tags Announcements
// @Accept json
// @Produce json
// @Param id path string true "Announcement ID"
// @Param announcement body AnnouncementInput true "Announcement"
// @Success 200 {object} models.Announcement
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 500 {object} FailureResponse
// @Router /api/announcements/{id} [put]
// @Security BearerAuth
func UpdateAnnouncement(c *gin.Context) {
	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid announcement ID"})
		return
	}

	var announcement models.Announcement

	if err := config.DB.First(&announcement, "id = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Announcement not found"})
		return
	}

	var input models.Announcement
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&announcement).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update announcement"})
		return
	}

	c.JSON(http.StatusOK, announcement)
}

// DeleteAnnouncement godoc
// @Summary Delete an announcement by ID
// @Description Delete an announcement by ID
// @Tags Announcements
// @Accept json
// @Produce json
// @Param id path string true "Announcement ID"
// @Success 200 {object} DeleteSuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 500 {object} FailureResponse
// @Router /api/announcements/{id} [delete]
// @Security BearerAuth
func DeleteAnnouncement(c *gin.Context) {
	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid announcement ID"})
		return
	}

	var announcement models.Announcement

	if err := config.DB.First(&announcement, "id = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Announcement not found"})
		return
	}

	if err := config.DB.Delete(&announcement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete announcement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Announcement deleted successfully"})
}
