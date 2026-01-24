package chapters

import (
	"github.com/gin-gonic/gin"
	"crm-go/models"
	"crm-go/config"
	"github.com/google/uuid"
	"net/http"
	"errors"
	"gorm.io/gorm"
)

// CreateChapters handles the creation of a new chapter
//@Summary Create a new chapters
//@Description Create a new chapter
//@Tags Chapters
//@Accept json
//@Produce json
//@Param chapter body models.ChapterInput true "Chapter"
//@Success 201 {object} models.SuccessResponse
//@Failure 400 {object} models.ErrorResponse
//@Failure 409 {object} models.ConflictResponse
//@Failure 500 {object} models.FailureResponse
//@Router /api/chapters [post]
func CreateChapter(c *gin.Context) {
	var input models.ChapterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 🔍 Check for duplicate chapter number per course
	var existingChapter models.Chapter
	err := config.DB.
		Where("course_id = ? AND chapter_number = ?", input.CourseID, input.ChapterNumber).
		First(&existingChapter).Error

	if err == nil {
		// Record found → duplicate
		c.JSON(http.StatusConflict, models.ConflictResponse{Error: "Chapter number already exists for this course"})
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Real DB error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to validate chapter number",
		})
		return
	}

	chapter := models.Chapter{
		ID:            uuid.New(),
		CourseID:      input.CourseID,
		Title:         input.Title,
		Slug:          input.Slug,
		Description:   input.Description,
		ChapterNumber: input.ChapterNumber,
		IsFree:        input.IsFree,
		Status:        input.Status,
		EstimatedTime: input.EstimatedTime,
	}

	if err := config.DB.Create(&chapter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Chapter created successfully",
	})
}



// GetAllChapters handles the retrieval of all chapters
//@Summary Get all chapters
//@Description Get all chapters
//@Tags Chapters
//@Accept json
//@Produce json
//@Success 200 {object} models.SuccessResponse
//@Failure 500 {object} models.FailureResponse
//@Router /chapters [get]
func GetAllChapters(c *gin.Context) {
	var chapters []models.Chapter

	if err := config.DB.
		Preload("Course", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "title")
		}).
		Order("chapter_number ASC").
		Find(&chapters).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chapters"})
		return
	}

	var response []models.ChapterResponse

	for _, ch := range chapters {
		var r models.ChapterResponse
		r.ID = ch.ID
		r.CourseID = ch.CourseID
		r.Title = ch.Title
		r.Slug = ch.Slug
		r.Description = ch.Description
		r.ChapterNumber = ch.ChapterNumber
		r.IsFree = ch.IsFree
		r.Status = ch.Status
		r.EstimatedTime = ch.EstimatedTime
		r.TotalLessons = ch.TotalLessons
		r.TotalDuration = ch.TotalDuration
		r.CreatedAt = ch.CreatedAt
		r.UpdatedAt = ch.UpdatedAt

		response = append(response, r)
	}

	c.JSON(http.StatusOK, response)
}


// GetChapterByID handles the retrieval of a single chapter by ID
//@Summary Get a single chapter by ID
//@Description Get a single chapter by ID
//@Tags Chapters
//@Accept json
//@Produce json
//@Param id path string true "Chapter ID"
//@Success 200 {object} models.SuccessResponse
//@Failure 400 {object} models.ErrorResponse
//@Failure 404 {object} models.NotFoundResponse
//@Failure 500 {object} models.FailureResponse
//@Router /chapters/{id} [get]
func GetChapterByID(c *gin.Context) {
	id := c.Param("id")

	chapterID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter ID"})
		return
	}

	var chapter models.Chapter

	if err := config.DB.
		Preload("Course").
		Preload("Lessons").
		First(&chapter, "id = ?", chapterID).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "Chapter not found"})
		return
	}

	c.JSON(http.StatusOK, chapter)
}



// UpdateChapter handles the updating of an existing chapter
//@Summary Update an existing chapter
//@Description Update an existing chapter
//@Tags Chapters
//@Accept json
//@Produce json
//@Param id path string true "Chapter ID"
//@Param chapter body models.ChapterInput true "Chapter"
//@Success 200 {object} models.SuccessResponse
//@Failure 400 {object} models.ErrorResponse
//@Failure 404 {object} models.NotFoundResponse
//@Failure 500 {object} models.FailureResponse
//@Router /api/chapters/{id} [put]
func UpdateChapter(c *gin.Context) {
	id := c.Param("id")

	chapterID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter ID"})
		return
	}

	var chapter models.Chapter
	if err := config.DB.First(&chapter, "id = ?", chapterID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chapter not found"})
		return
	}

	var input models.ChapterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ Prevent duplicate chapter numbers per course
	var count int64
	config.DB.Model(&models.Chapter{}).
		Where("course_id = ? AND chapter_number = ? AND id != ?",
			input.CourseID,
			input.ChapterNumber,
			chapterID,
		).
		Count(&count)

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Chapter number already exists for this course",
		})
		return
	}

	updates := map[string]interface{}{
		"course_id":      input.CourseID,
		"title":          input.Title,
		"slug":           input.Slug,
		"description":    input.Description,
		"chapter_number": input.ChapterNumber,
		"is_free":        input.IsFree,
		"status":         input.Status,
		"estimated_time": input.EstimatedTime,
	}

	if err := config.DB.Model(&chapter).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload with relations
	if err := config.DB.
		Preload("Course").
		Preload("Lessons").
		First(&chapter, "id = ?", chapterID).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "Chapter not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chapter updated successfully", "details": chapter})
}


// DeleteChapter handles the deletion of an existing chapter
//@Summary Delete an existing chapter
//@Description Delete an existing chapter
//@Tags Chapters
//@Accept json
//@Produce json
//@Param id path string true "Chapter ID"
//@Success 200 {object} models.DeleteSuccessResponse
//@Failure 400 {object} models.ErrorResponse
//@Failure 404 {object} models.NotFoundResponse
//@Failure 500 {object} models.FailureResponse
//@Router /api/chapters/{id} [delete]
func DeleteChapter(c *gin.Context) {
	id := c.Param("id")

	chapterID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter ID"})
		return
	}

	result := config.DB.Delete(&models.Chapter{}, "id = ?", chapterID)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chapter"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chapter not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chapter deleted successfully"})
}


