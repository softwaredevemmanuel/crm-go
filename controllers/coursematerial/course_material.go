package coursematerial

import (
	"crm-go/config"
	"crm-go/models"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"strings"
)

// CreateCourseMaterial handles the creation of a new course material
// @Summary Create a new course material
// @Description Create a new course material
// @Tags Course Materials
// @Accept json
// @Produce json
// @Param chapter body models.CreateCourseMaterialRequest true "Course Material"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ConflictResponse
// @Failure 500 {object} models.FailureResponse
// @Router /api/course-materials [post]
// @Security BearerAuth
func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = regexp.MustCompile(`[^a-z0-9\-]`).ReplaceAllString(slug, "")
	return slug
}

func CreateCourseMaterial(c *gin.Context) {
	var req models.CreateCourseMaterialRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var existing models.CourseMaterial

	err := config.DB.
		Where("course_id = ? AND title = ?", req.CourseID, req.Title).
		First(&existing).Error

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "A course material with this title already exists for this course",
		})
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error while checking duplicates",
		})
		return
	}

	// Optional: validate material type
	validTypes := map[string]bool{
		"document": true, "video": true, "audio": true, "image": true,
		"code": true, "presentation": true, "spreadsheet": true,
		"archive": true, "link": true, "external": true,
		"exercise": true, "quiz": true, "template": true,
	}

	if !validTypes[req.Type] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid material type",
		})
		return
	}

	material := models.CourseMaterial{
		CourseID:    req.CourseID,
		ChapterID:   req.ChapterID,
		LessonID:    req.LessonID,
		Title:       req.Title,
		Description: req.Description,
		Slug:        generateSlug(req.Title),
		Type:        req.Type,
		FileURL:     req.FileURL,
		Status:      req.Status,
	}

	// Default status if not provided
	if material.Status == "" {
		material.Status = "draft"
	}

	if err := config.DB.Create(&material).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create course material",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Course material created successfully"})
}

// GetCourseMaterials godoc
// @Summary      Get course materials
// @Description  Get all course materials, optionally filtered by course, chapter, or lesson
// @Tags         Course Materials
// @Produce      json
// @Param        course_id   query   string  false  "Course ID"
// @Param        chapter_id  query   string  false  "Chapter ID"
// @Param        lesson_id   query   string  false  "Lesson ID"
// @Param        status      query   string  false  "Status"
// @Param        type        query   string  false  "Material Type"
// @Success      200 {array}  models.SuccessResponse
// @Failure      400 {object} models.ErrorResponse
// @Router       /course-materials [get]
func GetCourseMaterials(c *gin.Context) {
	var materials []models.CourseMaterial
	var responses []models.CourseMaterialResponse

	// Query params
	courseID := c.Query("course_id")
	chapterID := c.Query("chapter_id")
	lessonID := c.Query("lesson_id")
	status := c.Query("status")
	materialType := c.Query("type")

	query := config.DB.Model(&models.CourseMaterial{})

	// Filters
	if courseID != "" {
		id, err := uuid.Parse(courseID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course_id"})
			return
		}
		query = query.Where("course_id = ?", id)
	}

	if chapterID != "" {
		id, err := uuid.Parse(chapterID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter_id"})
			return
		}
		query = query.Where("chapter_id = ?", id)
	}

	if lessonID != "" {
		id, err := uuid.Parse(lessonID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson_id"})
			return
		}
		query = query.Where("lesson_id = ?", id)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if materialType != "" {
		query = query.Where("type = ?", materialType)
	}

	// Execute
	if err := query.
		Order("created_at DESC").
		Find(&materials).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch course materials",
		})
		return
	}

	// Map to response
	for _, m := range materials {
		responses = append(responses, models.CourseMaterialResponse{
			ID:          m.ID,
			CourseID:    m.CourseID,
			ChapterID:   m.ChapterID,
			LessonID:    m.LessonID,
			Title:       m.Title,
			Description: m.Description,
			Slug:        m.Slug,
			Type:        m.Type,
			FileURL:     m.FileURL,
			Status:      m.Status,
			CreatedAt:   m.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, responses)
}

// GetCourseMaterialByID godoc
// @Summary      Get course material details
// @Description  Get course material with chapter and lesson details
// @Tags         Course Materials
// @Produce      json
// @Param        id path string true "Course Material ID"
// @Success      200 {object} models.CourseMaterialViewResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      404 {object} models.NotFoundResponse
// @Router       /course-materials/{id} [get]
func GetCourseMaterialByID(c *gin.Context) {
	id := c.Param("id")

	materialID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid course material ID",
		})
		return
	}

	var material models.CourseMaterial

	if err := config.DB.
		Preload("Course").
		Preload("Chapter").
		Preload("Lesson").
		First(&material, "id = ?", materialID).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Course material not found",
		})
		return
	}

	// Build response
	response := models.CourseMaterialViewResponse{
		ID:          material.ID,
		CourseID:    material.CourseID,
		ChapterID:   material.ChapterID,
		LessonID:    material.LessonID,
		Title:       material.Title,
		Description: material.Description,
		Slug:        material.Slug,
		Type:        material.Type,
		FileURL:     material.FileURL,
		Status:      material.Status,
		CreatedAt:   material.CreatedAt,

		Course: models.CourseMiniResponse{
			ID:    material.Course.ID,
			Title: material.Course.Title,
		},
	}

	// Optional Chapter
	if material.ChapterID != nil {
		response.Chapter = &models.ChapterMiniResponse{
			ID:            material.Chapter.ID,
			Title:         material.Chapter.Title,
			ChapterNumber: material.Chapter.ChapterNumber,
		}
	}

	// Optional Lesson
	if material.LessonID != nil {
		response.Lesson = &models.LessonMiniResponse{
			ID:          material.Lesson.ID,
			Title:       material.Lesson.Title,
			ContentType: material.Lesson.ContentType,
			ContentURL:  material.Lesson.ContentURL,
		}
	}

	c.JSON(http.StatusOK, response)
}

// UpdateCourseMaterial godoc
// @Summary      Update course material
// @Description  Update course material details
// @Tags         Course Materials
// @Accept       json
// @Produce      json
// @Param        id path string true "Course Material ID"
// @Param        material body models.UpdateCourseMaterialRequest true "Course Material payload"
// @Success      200 {object} models.CourseMaterialResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      404 {object} models.NotFoundResponse
// @Failure      500 {object} models.FailureResponse
// @Router       /api/course-materials/{id} [put]
// @Security BearerAuth
func UpdateCourseMaterial(c *gin.Context) {
	id := c.Param("id")

	materialID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid course material ID",
		})
		return
	}

	var req models.UpdateCourseMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var material models.CourseMaterial

	// Fetch existing material
	if err := config.DB.
		Preload("Course").
		First(&material, "id = ?", materialID).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Course material not found",
		})
		return
	}

	// 🔒 Prevent duplicate title within the same course
	var count int64
	config.DB.Model(&models.CourseMaterial{}).
		Where("course_id = ? AND title = ? AND id != ?",
			material.CourseID,
			req.Title,
			material.ID,
		).
		Count(&count)

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Course material with this title already exists in the course",
		})
		return
	}

	// Update fields
	material.Title = req.Title
	material.Description = req.Description
	material.Type = req.Type
	material.FileURL = req.FileURL
	material.Status = req.Status
	material.ChapterID = req.ChapterID
	material.LessonID = req.LessonID

	if err := config.DB.Save(&material).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update course material",
		})
		return
	}

	// Reload relations
	if err := config.DB.
		Preload("Course").
		Preload("Chapter").
		Preload("Lesson").
		First(&material, "id = ?", material.ID).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load updated material",
		})
		return
	}

	// Build response
	response := models.CourseMaterialViewResponse{
		ID:          material.ID,
		CourseID:    material.CourseID,
		ChapterID:   material.ChapterID,
		LessonID:    material.LessonID,
		Title:       material.Title,
		Description: material.Description,
		Slug:        material.Slug,
		Type:        material.Type,
		FileURL:     material.FileURL,
		Status:      material.Status,
		CreatedAt:   material.CreatedAt,

		Course: models.CourseMiniResponse{
			ID:    material.Course.ID,
			Title: material.Course.Title,
		},
	}

	// Optional Chapter
	if material.ChapterID != nil {
		response.Chapter = &models.ChapterMiniResponse{
			ID:            material.Chapter.ID,
			Title:         material.Chapter.Title,
			ChapterNumber: material.Chapter.ChapterNumber,
		}
	}

	// Optional Lesson
	if material.LessonID != nil {
		response.Lesson = &models.LessonMiniResponse{
			ID:          material.Lesson.ID,
			Title:       material.Lesson.Title,
			ContentType: material.Lesson.ContentType,
			ContentURL:  material.Lesson.ContentURL,
		}
	}

	c.JSON(http.StatusOK, response)
}

// DeleteWithArchive moved to services/delete_service.go to avoid duplicate declarations
// (implementation exists later in this file under the services/delete_service.go section)

// services/delete_service.go
func DeleteWithArchive(
    tx *gorm.DB,
    entityType string,
    entityID uuid.UUID,
    data any,        // This is for archiving only, NOT for deletion
    deletedBy *uuid.UUID,
    reason string,
    model interface{}, // Add this parameter - the actual model to delete
) error {
    
    // First, marshal the data for archiving
    jsonData, err := json.Marshal(data)
    if err != nil {
        return err
    }

    // Create the archive record
    archive := models.DeletedRecord{
        EntityType: entityType,
        EntityID:   entityID,
        Data:       jsonData,
        DeletedBy:  deletedBy,
        Reason:     reason,
    }

    if err := tx.Create(&archive).Error; err != nil {
        return err
    }

    // Delete the actual model, NOT the data parameter
    if err := tx.Delete(model, "id = ?", entityID).Error; err != nil {
        return err
    }

    return nil
}

func GetUserIDFromContext(c *gin.Context) *uuid.UUID {
	v, exists := c.Get("user_id")
	if !exists {
		return nil
	}

	idStr, ok := v.(string)
	if !ok {
		return nil
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil
	}

	return &id
}

// DeleteCourseMaterial godoc
// @Summary      Delete course material
// @Description  Delete course material by ID
// @Tags         Course Materials
// @Produce      json
// @Param        id path string true "Course Material ID"
// @Success      200 {object} models.SuccessResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      404 {object} models.NotFoundResponse
// @Failure      500 {object} models.FailureResponse
// @Router       /api/course-materials/{id} [delete]
// @Security BearerAuth
func DeleteCourseMaterial(c *gin.Context) {
    id := c.Param("id")

    materialID, err := uuid.Parse(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    // Get the material to archive
    var material models.CourseMaterial
    if err := config.DB.
        First(&material, "id = ?", materialID).
        Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Material not found"})
        return
    }

    userID := GetUserIDFromContext(c)

    // Convert to response struct for archiving
    materialResponse := models.CourseMaterialResponse{
        ID:          material.ID,
        CourseID:    material.CourseID,
        ChapterID:   material.ChapterID,
        LessonID:    material.LessonID,
        Title:       material.Title,
        Description: material.Description,
        Slug:        material.Slug,
        Type:        material.Type,
        FileURL:     material.FileURL,
        Status:      material.Status,
        CreatedAt:   material.CreatedAt,
    }

	// ✅ TRANSACTION STARTS HERE
	tx := config.DB.Begin()

	if err := DeleteWithArchive(
		tx,
		"course_materials",
		material.ID,
		&materialResponse, // Data to archive
		userID,
		"Admin deleted course material",
		&models.CourseMaterial{}, // The actual model to delete
	); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Delete failed: " + err.Error(),
		})
		return
	}

	tx.Commit()

    c.JSON(http.StatusOK, gin.H{
        "message": "Material deleted successfully",
    })
}