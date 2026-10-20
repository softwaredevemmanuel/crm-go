package controllers

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"crm-go/config"
	"crm-go/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"encoding/json"
	"strconv"
	"time"
)

// ============================================================
// Input DTOs
// ============================================================

type CreateCourseInput struct {
	Title            string   `json:"title" binding:"required,min=3,max=255"`
	Code             string   `json:"code" binding:"omitempty,min=2,max=50"`
	Slug             string   `json:"slug" binding:"omitempty,max=255"`
	Description      string   `json:"description" binding:"max=5000"`
	Image            string   `json:"image" binding:"omitempty,url"`
	ThumbnailURL     string   `json:"thumbnail_url" binding:"omitempty,url"`
	CoverImageURL    string   `json:"cover_image_url" binding:"omitempty,url"`
	VideoURL         string   `json:"video_url" binding:"omitempty,url"`
	IntroVideoURL    string   `json:"intro_video_url" binding:"omitempty,url"`
	Level            string   `json:"level" binding:"omitempty,oneof=beginner intermediate advanced"`
	Status           string   `json:"status" binding:"omitempty,oneof=draft published archived"`
	LearningOutcomes []string `json:"learning_outcomes"`
	Requirements     []string `json:"requirements"`
	IsFree           *bool    `json:"is_free"`
	PriceCents       int      `json:"price_cents" binding:"omitempty,gte=0"`
	Currency         string   `json:"currency" binding:"omitempty,len=3"`
	DurationMinutes  int      `json:"duration_minutes" binding:"omitempty,gte=0"`
	RequiresEnrollment *bool  `json:"requires_enrollment"`
	IsRestricted     *bool    `json:"is_restricted"`
}

type UpdateCourseInput struct {
	Title            *string   `json:"title" binding:"omitempty,min=3,max=255"`
	Code             *string   `json:"code" binding:"omitempty,min=2,max=50"`
	Slug             *string   `json:"slug" binding:"omitempty,max=255"`
	Description      *string   `json:"description" binding:"omitempty,max=5000"`
	Image            *string   `json:"image" binding:"omitempty,url"`
	ThumbnailURL     *string   `json:"thumbnail_url" binding:"omitempty,url"`
	CoverImageURL    *string   `json:"cover_image_url" binding:"omitempty,url"`
	VideoURL         *string   `json:"video_url" binding:"omitempty,url"`
	IntroVideoURL    *string   `json:"intro_video_url" binding:"omitempty,url"`
	Level            *string   `json:"level" binding:"omitempty,oneof=beginner intermediate advanced"`
	Status           *string   `json:"status" binding:"omitempty,oneof=draft published archived"`
	LearningOutcomes *[]string `json:"learning_outcomes"`
	Requirements     *[]string `json:"requirements"`
	IsFree           *bool     `json:"is_free"`
	PriceCents       *int      `json:"price_cents" binding:"omitempty,gte=0"`
	Currency         *string   `json:"currency" binding:"omitempty,len=3"`
	DurationMinutes  *int      `json:"duration_minutes" binding:"omitempty,gte=0"`
	RequiresEnrollment *bool   `json:"requires_enrollment"`
	IsRestricted     *bool     `json:"is_restricted"`
}

// ============================================================
// Response DTOs
// ============================================================

type CourseResponse struct {
	ID               string         `json:"id"`
	Title            string         `json:"title"`
	Code             string         `json:"code"`
	Slug             string         `json:"slug"`
	Description      string         `json:"description"`
	Image            string         `json:"image"`
	ThumbnailURL     string         `json:"thumbnail_url"`
	CoverImageURL    string         `json:"cover_image_url"`
	VideoURL         string         `json:"video_url"`
	IntroVideoURL    string         `json:"intro_video_url"`
	Level            string         `json:"level"`
	Status           string         `json:"status"`
	LearningOutcomes []string       `json:"learning_outcomes"`
	Requirements     []string       `json:"requirements"`
	IsFree           bool           `json:"is_free"`
	PriceCents       int            `json:"price_cents"`
	Currency         string         `json:"currency"`
	DurationMinutes  int            `json:"duration_minutes"`
	TotalLessons     int            `json:"total_lessons"`
	EnrollmentCount  int            `json:"enrollment_count"`
	Rating           float64        `json:"rating"`
	RatingCount      int            `json:"rating_count"`
	RequiresEnrollment bool         `json:"requires_enrollment"`
	IsRestricted     bool           `json:"is_restricted"`
	PublishedAt      *string        `json:"published_at,omitempty"`
	CreatedAt        string         `json:"created_at"`
	UpdatedAt        string         `json:"updated_at"`
}

type CourseMiniResponse struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Slug         string `json:"slug"`
	ThumbnailURL string `json:"thumbnail_url"`
	Level        string `json:"level"`
}

type CourseWithCategoryMatesResponse struct {
	Course         CourseResponse   `json:"course"`
	RelatedCourses []CourseResponse `json:"related_courses"`
}

// ============================================================
// Mapping helpers
// ============================================================

func toCourseResponse(m *models.Course) CourseResponse {
	var publishedAt *string
	if m.PublishedAt != nil {
		s := m.PublishedAt.Format("2006-01-02T15:04:05Z07:00")
		publishedAt = &s
	}

	return CourseResponse{
		ID:                 m.ID.String(),
		Title:              m.Title,
		Code:               m.Code,
		Slug:               m.Slug,
		Description:        m.Description,
		Image:              m.Image,
		ThumbnailURL:       m.ThumbnailURL,
		CoverImageURL:      m.CoverImageURL,
		VideoURL:           m.VideoURL,
		IntroVideoURL:      m.IntroVideoURL,
		Level:              string(m.Level),
		Status:             string(m.Status),
		LearningOutcomes:   decodeStringArray(m.LearningOutcomes),
		Requirements:       decodeStringArray(m.Requirements),
		IsFree:             m.IsFree,
		PriceCents:         m.PriceCents,
		Currency:           m.Currency,
		DurationMinutes:    m.DurationMinutes,
		TotalLessons:       m.TotalLessons,
		EnrollmentCount:    m.EnrollmentCount,
		Rating:             m.Rating,
		RatingCount:        m.RatingCount,
		RequiresEnrollment: m.RequiresEnrollment,
		IsRestricted:       m.IsRestricted,
		PublishedAt:        publishedAt,
		CreatedAt:          m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:          m.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}



// decodeStringArray converts a JSONB column into a []string.
func decodeStringArray(raw datatypes.JSON) []string {
	if len(raw) == 0 {
		return []string{}
	}
	var out []string
	if err := json.Unmarshal(raw, &out); err != nil {
		return []string{}
	}
	if out == nil {
		return []string{}
	}
	return out
}

// encodeStringArray marshals a []string into datatypes.JSON.
func encodeStringArray(values []string) datatypes.JSON {
	if values == nil {
		values = []string{}
	}
	raw, err := datatypes.JSON([]byte("[]")).MarshalJSON()
	_ = raw
	_ = err
	// Use JSON marshal directly
	b, err := jsonMarshal(values)
	if err != nil {
		return datatypes.JSON([]byte("[]"))
	}
	return datatypes.JSON(b)
}

// jsonMarshal is a tiny wrapper to avoid importing encoding/json in many places.
func jsonMarshal(v interface{}) ([]byte, error) {
	return jsonMarshalImpl(v)
}

// ============================================================
// Slug helper (inline, no utils package)
// ============================================================

var (
	nonAlnum = regexp.MustCompile(`[^a-z0-9\s-]`)
	spaces   = regexp.MustCompile(`[\s-]+`)
)

func slugify(input string) string {
	s := strings.ToLower(strings.TrimSpace(input))
	s = nonAlnum.ReplaceAllString(s, "")
	s = spaces.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

// ============================================================
// CreateCourse
// ============================================================

// CreateCourse godoc
//
//	@Summary		Create a new course
//	@Description	Admin creates a new course. Prevents duplicate title, code, or slug.
//	@Tags			Courses
//	@Accept			json
//	@Produce		json
//	@Param			course	body		CreateCourseInput	true	"Course details"
//	@Success		201		{object}	CourseResponse
//	@Failure		400		{object}	map[string]string	"Invalid request body"
//	@Failure		409		{object}	map[string]string	"Course with the same title, code, or slug already exists"
//	@Failure		500		{object}	map[string]string	"Failed to create course"
//	@Router			/api/courses [post]
//	@Security		BearerAuth
func CreateCourse(c *gin.Context) {
	var input CreateCourseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB


	// Defaults
	level := input.Level
	if level == "" {
		level = string(models.LevelBeginner)
	}
	status := input.Status
	if status == "" {
		status = string(models.StatusDraft)
	}

	// Code fallback
	code := strings.TrimSpace(input.Code)
	if code == "" {
		// Derive from title + short random suffix
		base := strings.ToUpper(strings.ReplaceAll(slugify(input.Title), "-", ""))
		if len(base) > 6 {
			base = base[:6]
		}
		code = base + "-" + strings.ToUpper(uuid.NewString()[:4])
	}
	code = strings.ToUpper(code)

	// Slug fallback
	slug := input.Slug
	if slug == "" {
		slug = slugify(input.Title)
	}

	// Uniqueness: title
	var existing models.Course
	if err := db.Where("LOWER(title) = ?", strings.ToLower(input.Title)).
		First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Course with the same title already exists"})
		return
	}

	// Uniqueness: code
	if err := db.Where("code = ?", code).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Course with the same code already exists"})
		return
	}

	// Uniqueness: slug
	if err := db.Where("slug = ?", slug).First(&existing).Error; err == nil {
		// Auto-disambiguate slug instead of failing
		slug = slug + "-" + uuid.NewString()[:6]
	}

	// Defaults for booleans
	isFree := true
	if input.IsFree != nil {
		isFree = *input.IsFree
	}
	requiresEnrollment := true
	if input.RequiresEnrollment != nil {
		requiresEnrollment = *input.RequiresEnrollment
	}
	isRestricted := false
	if input.IsRestricted != nil {
		isRestricted = *input.IsRestricted
	}
	currency := input.Currency
	if currency == "" {
		currency = "NGN"
	}

	course := models.Course{
		Title:              strings.TrimSpace(input.Title),
		Code:               code,
		Slug:               slug,
		Description:        input.Description,
		Image:              input.Image,
		ThumbnailURL:       input.ThumbnailURL,
		CoverImageURL:      input.CoverImageURL,
		VideoURL:           input.VideoURL,
		IntroVideoURL:      input.IntroVideoURL,
		Level:              models.CourseLevel(level),
		Status:             models.CourseStatus(status),
		LearningOutcomes:   encodeStringArray(input.LearningOutcomes),
		Requirements:       encodeStringArray(input.Requirements),
		IsFree:             isFree,
		PriceCents:         input.PriceCents,
		Currency:           currency,
		DurationMinutes:    input.DurationMinutes,
		RequiresEnrollment: requiresEnrollment,
		IsRestricted:       isRestricted,
	}

	if err := db.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, toCourseResponse(&course))
}

// ============================================================
// GetCourses
// ============================================================

// GetCourses godoc
//
//	@Summary		List all courses
//	@Description	Retrieve all courses available in the system. Supports pagination, search, and filters.
//	@Tags			Courses
//	@Produce		json
//	@Param			page		query		int		false	"Page number (default 1)"
//	@Param			limit		query		int		false	"Items per page (default 20, max 100)"
//	@Param			search		query		string	false	"Search in title, code, description"
//	@Param			status		query		string	false	"Filter by status"	Enums(draft, published, archived)
//	@Param			level		query		string	false	"Filter by level"	Enums(beginner, intermediate, advanced)
//	@Success		200			{object}	map[string]interface{}
//	@Failure		500			{object}	map[string]string	"Failed to fetch courses"
//	@Router			/courses [get]
func GetCourses(c *gin.Context) {
	db := config.DB

	page := 1
	limit := 20
	if p := c.Query("page"); p != "" {
		if v, err := strconvAtoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if l := c.Query("limit"); l != "" {
		if v, err := strconvAtoi(l); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	query := db.Model(&models.Course{})

	if search := strings.TrimSpace(c.Query("search")); search != "" {
		s := "%" + strings.ToLower(search) + "%"
		query = query.Where(
			"LOWER(title) LIKE ? OR LOWER(code) LIKE ? OR LOWER(description) LIKE ?",
			s, s, s,
		)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if level := c.Query("level"); level != "" {
		query = query.Where("level = ?", level)
	}


	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count courses"})
		return
	}

	var courses []models.Course
	offset := (page - 1) * limit
	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch courses"})
		return
	}

	items := make([]CourseResponse, 0, len(courses))
	for i := range courses {
		items = append(items, toCourseResponse(&courses[i]))
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{
		"data":        items,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}

// ============================================================
// GetCourseByID
// ============================================================

// GetCourseByID godoc
//
//	@Summary		Get a course by ID
//	@Description	Retrieve details of a specific course using its ID
//	@Tags			Courses
//	@Produce		json
//	@Param			id	path		string	true	"Course ID"	format(uuid)
//	@Success		200	{object}	CourseResponse
//	@Failure		400	{object}	map[string]string	"Invalid course ID"
//	@Failure		404	{object}	map[string]string	"Course not found"
//	@Router			/courses/{id} [get]
func GetCourseByID(c *gin.Context) {
	id := c.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	db := config.DB
	var course models.Course
	if err := db.First(&course, "id = ?", uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch course"})
		return
	}

	c.JSON(http.StatusOK, toCourseResponse(&course))
}

// ============================================================
// GetCourseBySlug
// ============================================================

// GetCourseBySlug godoc
//
//	@Summary		Get a course by slug
//	@Description	Retrieve a course by its URL-friendly slug.
//	@Tags			Courses
//	@Produce		json
//	@Param			slug	path		string	true	"Course slug"
//	@Success		200		{object}	CourseResponse
//	@Failure		404		{object}	map[string]string	"Course not found"
//	@Router			/courses/slug/{slug} [get]
func GetCourseBySlug(c *gin.Context) {
	slug := c.Param("slug")

	db := config.DB
	var course models.Course
	if err := db.Where("slug = ?", slug).First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch course"})
		return
	}

	c.JSON(http.StatusOK, toCourseResponse(&course))
}

// ============================================================
// GetProductsWithRelatedCourseID
// ============================================================

// GetProductsWithRalatedCourseID godoc
//
//	@Summary		Get products related to a course
//	@Description	Retrieve a course along with products linked to it via the pivot table
//	@Tags			Courses
//	@Produce		json
//	@Param			id	path		string	true	"Course ID"	format(uuid)
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	map[string]string	"Invalid course ID"
//	@Failure		404	{object}	map[string]string	"Course not found"
//	@Failure		500	{object}	map[string]string	"Failed to fetch related products"
//	@Router			/courses/{id}/products [get]
func GetProductsWithRalatedCourseID(c *gin.Context) {
	courseID := c.Param("id")
	db := config.DB

	uid, err := uuid.Parse(courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var course models.Course
	if err := db.First(&course, "id = ?", uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch course"})
		return
	}

	var relatedProducts []models.Product
	if err := db.
		Joins("JOIN course_product_tables cp ON cp.product_id = products.id").
		Where("cp.course_id = ?", uid).
		Find(&relatedProducts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch related products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"course":           toCourseResponse(&course),
		"related_products": relatedProducts,
	})
}

// ============================================================
// UpdateCourse (partial)
// ============================================================

// UpdateCourse godoc
//
//	@Summary		Update an existing course (partial)
//	@Description	Admin updates course details by ID. Only fields present in the payload are updated.
//	@Tags			Courses
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Course ID"	format(uuid)
//	@Param			course	body		UpdateCourseInput	true	"Updated course fields"
//	@Success		200		{object}	CourseResponse
//	@Failure		400		{object}	map[string]string	"Invalid course ID or request body"
//	@Failure		404		{object}	map[string]string	"Course not found"
//	@Failure		409		{object}	map[string]string	"Code or slug already in use"
//	@Failure		500		{object}	map[string]string	"Failed to update course"
//	@Router			/api/courses/{id} [put]
//	@Security		BearerAuth
func UpdateCourse(c *gin.Context) {
	id := c.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	db := config.DB
	var course models.Course
	if err := db.First(&course, "id = ?", uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch course"})
		return
	}

	var input UpdateCourseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Code uniqueness
	if input.Code != nil && strings.ToUpper(*input.Code) != course.Code {
		newCode := strings.ToUpper(strings.TrimSpace(*input.Code))
		var existing models.Course
		if err := db.Where("code = ? AND id <> ?", newCode, uid).
			First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Course with the same code already exists"})
			return
		}
		course.Code = newCode
	}

	// Slug uniqueness
	if input.Slug != nil && *input.Slug != course.Slug {
		newSlug := *input.Slug
		var existing models.Course
		if err := db.Where("slug = ? AND id <> ?", newSlug, uid).
			First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Course with the same slug already exists"})
			return
		}
		course.Slug = newSlug
	}

	// Title uniqueness
	if input.Title != nil && *input.Title != course.Title {
		var existing models.Course
		if err := db.Where("LOWER(title) = ? AND id <> ?", strings.ToLower(*input.Title), uid).
			First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Course with the same title already exists"})
			return
		}
		course.Title = strings.TrimSpace(*input.Title)
	}

	// Simple scalar fields
	if input.Description != nil {
		course.Description = *input.Description
	}
	if input.Image != nil {
		course.Image = *input.Image
	}
	if input.ThumbnailURL != nil {
		course.ThumbnailURL = *input.ThumbnailURL
	}
	if input.CoverImageURL != nil {
		course.CoverImageURL = *input.CoverImageURL
	}
	if input.VideoURL != nil {
		course.VideoURL = *input.VideoURL
	}
	if input.IntroVideoURL != nil {
		course.IntroVideoURL = *input.IntroVideoURL
	}

	if input.Level != nil {
		course.Level = models.CourseLevel(*input.Level)
	}
	if input.Status != nil {
		course.Status = models.CourseStatus(*input.Status)
		// Stamp published_at when transitioning into "published"
		if course.Status == models.StatusPublished && course.PublishedAt == nil {
			now := timeNow()
			course.PublishedAt = &now
		}
	}
	if input.LearningOutcomes != nil {
		course.LearningOutcomes = encodeStringArray(*input.LearningOutcomes)
	}
	if input.Requirements != nil {
		course.Requirements = encodeStringArray(*input.Requirements)
	}
	if input.IsFree != nil {
		course.IsFree = *input.IsFree
	}
	if input.PriceCents != nil {
		course.PriceCents = *input.PriceCents
	}
	if input.Currency != nil {
		course.Currency = *input.Currency
	}
	if input.DurationMinutes != nil {
		course.DurationMinutes = *input.DurationMinutes
	}
	if input.RequiresEnrollment != nil {
		course.RequiresEnrollment = *input.RequiresEnrollment
	}
	if input.IsRestricted != nil {
		course.IsRestricted = *input.IsRestricted
	}

	if err := db.Save(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		return
	}

	c.JSON(http.StatusOK, toCourseResponse(&course))
}

// ============================================================
// DeleteCourse (soft delete)
// ============================================================

// DeleteCourse godoc
//
//	@Summary		Delete a course
//	@Description	Admin soft-deletes a course by ID. The record is preserved with a deleted_at timestamp.
//	@Tags			Courses
//	@Produce		json
//	@Param			id	path		string	true	"Course ID"	format(uuid)
//	@Success		200	{object}	map[string]string	"Course deleted successfully"
//	@Failure		400	{object}	map[string]string	"Invalid course ID"
//	@Failure		404	{object}	map[string]string	"Course not found"
//	@Failure		500	{object}	map[string]string	"Failed to delete course"
//	@Router			/api/courses/{id} [delete]
//	@Security		BearerAuth
func DeleteCourse(c *gin.Context) {
	id := c.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	db := config.DB
	var course models.Course
	if err := db.First(&course, "id = ?", uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch course"})
		return
	}

	if err := db.Delete(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}


// jsonMarshalImpl is a thin wrapper so we can keep `encoding/json` out of the main controller file.
func jsonMarshalImpl(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// strconvAtoi wraps strconv.Atoi.
func strconvAtoi(s string) (int, error) {
	return strconv.Atoi(s)
}

// timeNow returns the current time — kept as a helper so tests can stub it.
func timeNow() time.Time {
	return time.Now()
}