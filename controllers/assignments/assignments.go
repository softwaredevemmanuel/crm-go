package assignments

import (
	"crm-go/config"
	"crm-go/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// CreateAssignment handles the creation of a new assignment
// @Summary Create a new assignment
// @Description Create a new assignment
// @Tags Assignments
// @Accept json
// @Produce json
// @Param assignment body models.AssignmentInput true "Assignment"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.FailureResponse
// @Router /assignments [post]
func CreateAssignment(c *gin.Context) {
	var input models.AssignmentInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// 🔍 Ensure course exists
	var course models.Course
	if err := config.DB.First(&course, "id = ?", input.CourseID).Error; err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid course_id",
		})
		return
	}

	// 🔍 Optional chapter validation
	if input.ChapterID != nil {
		var chapter models.Chapter
		if err := config.DB.First(&chapter, "id = ?", *input.ChapterID).Error; err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "Invalid chapter_id",
			})
			return
		}
	}

	assignment := models.Assignment{
		ID:               uuid.New(),
		CourseID:         input.CourseID,
		ChapterID:        input.ChapterID,
		LessonID:         input.LessonID,
		PublisherID:      input.PublisherID,
		Title:            input.Title,
		Slug:             input.Slug,
		Description:      input.Description,
		Type:             input.Type,
		SubmissionType:   input.SubmissionType,
		Content:          input.Content,
		DueDate:          input.DueDate,
		Status:           input.Status,	
	}

	if err := config.DB.Create(&assignment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create assignment",
		})
		return
	}

	c.JSON(http.StatusCreated, "Message: Assignment created successfully")
}

// GetAssignments handles the retrieval of all assignments
// @Summary Get all assignments
// @Description Get all assignments
// @Tags Assignments
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 500 {object} models.FailureResponse
// @Router /assignments [get]
func GetAssignments(c *gin.Context) {
	var assignments []models.AssignmentListResponse

	if err := config.DB.
		Model(&models.Assignment{}).
		Select(`
			id,
			course_id,
			chapter_id,
			lesson_id,
			publisher_id,
			title,
			slug,
			description,
			content,
			type,
			submission_type,
			status,
			due_date,
			created_at,
			updated_at,
			approved_at,
			published_at,
			archived_at
		`).
		Order("created_at DESC").
		Scan(&assignments).Error; err != nil {

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to fetch assignments",
		})
		return
	}

	c.JSON(http.StatusOK, assignments)
}

// GetAssignmentByID handles the retrieval of a single assignment by ID
// @Summary Get a single assignment by ID
// @Description Get a single assignment by ID
// @Tags Assignments
// @Accept json
// @Produce json
// @Param id path string true "Assignment ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.NotFoundResponse
// @Failure 500 {object} models.FailureResponse
// @Router /assignments/{id} [get]
func GetAssignmentByID(c *gin.Context) {
	id := c.Param("id")

	assignmentID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid assignment ID",
		})
		return
	}

	var assignment models.Assignment

	err = config.DB.
		Preload("Course").
		Preload("Chapter").
		Preload("Lesson").
		// Preload("Submissions").
		Preload("Publisher").
		First(&assignment, "id = ?", assignmentID).Error

	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Assignment not found",
		})
		return
	}

	// 🔹 Map to response DTO
	response := models.AssignmentViewResponse{
		ID:               assignment.ID,
		Title:            assignment.Title,
		Slug:             assignment.Slug,
		Description:      assignment.Description,
		Content:          assignment.Content,
		Type:             assignment.Type,
		SubmissionType:   assignment.SubmissionType,
		Status:           assignment.Status,
		DueDate:          assignment.DueDate,
		CreatedAt:        assignment.CreatedAt,
		UpdatedAt:        assignment.UpdatedAt,

		Course: models.CourseResponse{
			ID:               assignment.Course.ID,
			Title:            assignment.Course.Title,
			Description:      assignment.Course.Description,
			Image:            assignment.Course.Image,
			VideoURL:         assignment.Course.VideoURL,
			TutorID:          assignment.Course.TutorID,
			LearningOutcomes: assignment.Course.LearningOutcomes,
			Requirements:     assignment.Course.Requirements,
		},

		Publisher: models.UserResponse{
			ID:        assignment.Publisher.ID,
			FirstName: assignment.Publisher.FirstName,
			LastName:  assignment.Publisher.LastName,
			Email:     assignment.Publisher.Email,
		},
	}

	if assignment.ChapterID != nil {
		response.Chapter = &models.ChapterResponse{
			ID:            assignment.Chapter.ID,
			CourseID:      assignment.Chapter.CourseID,
			Title:         assignment.Chapter.Title,
			Slug:          assignment.Chapter.Slug,
			Description:   assignment.Chapter.Description,
			ChapterNumber: assignment.Chapter.ChapterNumber,
			IsFree:        assignment.Chapter.IsFree,
			Status:        assignment.Chapter.Status,
			EstimatedTime: assignment.Chapter.EstimatedTime,
			TotalLessons:  assignment.Chapter.TotalLessons,
			TotalDuration: assignment.Chapter.TotalDuration,
		}
	}

	if assignment.LessonID != nil {
		response.Lesson = &models.LessonResponse{
			ID:          assignment.Lesson.ID,
			ChapterID:   assignment.Lesson.ChapterID,
			CourseID:    assignment.Lesson.CourseID,
			Title:       assignment.Lesson.Title,
			ContentType: assignment.Lesson.ContentType,
			ContentURL:  assignment.Lesson.ContentURL,
			CreatedAt:   assignment.Lesson.CreatedAt,
			UpdatedAt:   assignment.Lesson.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

func ToAssignmentResponse(a models.Assignment) models.AssignmentViewResponse {
	response := models.AssignmentViewResponse{
		ID:     a.ID,
		Title:  a.Title,
		Status: a.Status,
		Slug:   a.Slug,
		Description: a.Description,
		Type: a.Type,
		DueDate: a.DueDate,
		Course: models.CourseResponse{
			ID:    a.Course.ID,
			Title: a.Course.Title,
			Description: a.Course.Description,
			Image: a.Course.Image,
			VideoURL: a.Course.VideoURL,
			TutorID: a.Course.TutorID,
			LearningOutcomes: a.Course.LearningOutcomes,
			Requirements: a.Course.Requirements,
		},
		Publisher: models.UserResponse{
			ID:        a.Publisher.ID,
			FirstName: a.Publisher.FirstName,
			LastName:  a.Publisher.LastName,
			Email:     a.Publisher.Email,
			Role: 	a.Publisher.Role,
		},

	}

	if a.ChapterID != nil {
		response.Chapter = &models.ChapterResponse{
			ID:    a.Chapter.ID,
			Title: a.Chapter.Title,
			Slug:  a.Chapter.Slug,
			Description: a.Chapter.Description,
			ChapterNumber: a.Chapter.ChapterNumber,
			IsFree: a.Chapter.IsFree,
			Status: a.Chapter.Status,
			EstimatedTime: a.Chapter.EstimatedTime,
			TotalLessons: a.Chapter.TotalLessons,
			TotalDuration: a.Chapter.TotalDuration,
		}
	}

	if a.LessonID != nil {
		response.Lesson = &models.LessonResponse{
			ID:    a.Lesson.ID,
			ChapterID: a.Lesson.ChapterID,
			CourseID:  a.Lesson.CourseID,
			Title: a.Lesson.Title,
			ContentType: a.Lesson.ContentType,
			ContentURL:  a.Lesson.ContentURL,
			CreatedAt:   a.Lesson.CreatedAt,
			UpdatedAt:   a.Lesson.UpdatedAt,
		}
	}

	return response
}

// UpdateAssignment handles the updating of an existing assignment
// @Summary Update an existing assignment
// @Description Update an existing assignment
// @Tags Assignments
// @Accept json
// @Produce json
// @Param id path string true "Assignment ID"
// @Param assignment body models.AssignmentUpdateInput true "Assignment payload"
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.NotFoundResponse
// @Failure 500 {object} models.FailureResponse
// @Router /assignments/{id} [put]
func UpdateAssignment(c *gin.Context) {
	// 1️⃣ Parse assignment ID
	id := c.Param("id")

	assignmentID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid assignment ID",
		})
		return
	}

	// 2️⃣ Find assignment
	var assignment models.Assignment
	if err := config.DB.First(&assignment, "id = ?", assignmentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Assignment not found",
		})
		return
	}

	// 3️⃣ Bind update input
	var input models.AssignmentUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 4️⃣ Update assignment (only provided fields)
	if err := config.DB.Model(&assignment).
		Select(
			"Title",
			"Slug",
			"Description",
			"Content",
			"SubmissionType",
			"Status",
			"Type",
			"DueDate",
			"CourseID",
			"ChapterID",
			"LessonID",
		).
		Updates(input).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update assignment",
		})
		return
	}

	// 5️⃣ Reload assignment with relations
	if err := config.DB.
		Preload("Course").
		Preload("Chapter").
		Preload("Lesson").
		Preload("Publisher").
		First(&assignment, "id = ?", assignmentID).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load updated assignment",
		})
		return
	}

	// 6️⃣ Map to DTO
	response := ToAssignmentResponse(assignment)

	// 7️⃣ Return response
	c.JSON(http.StatusOK, response)
}



// DeleteAssignment handles the deletion of an existing assignment
// @Summary Delete an existing assignment
// @Description Delete an existing assignment
// @Tags Assignments
// @Accept json
// @Produce json
// @Param id path string true "Assignment ID"
// @Success 200 {object} models.DeleteSuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.NotFoundResponse
// @Failure 500 {object} models.FailureResponse
// @Router /assignments/{id} [delete]
func DeleteAssignment(c *gin.Context) {
	id := c.Param("id")

	assignmentID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid assignment ID",
		})
		return
	}

	result := config.DB.Delete(&models.Assignment{}, "id = ?", assignmentID)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to delete assignment",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Assignment not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Assignment deleted successfully",
	})
}
