// controllers/lesson_report_controller.go
package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/lesson_report"
)

type LessonReportController struct {
	reportService *services.LessonReportService
}

func NewLessonReportController(reportService *services.LessonReportService) *LessonReportController {
	return &LessonReportController{
		reportService: reportService,
	}
}

// CreateReport creates a new lesson report
// @Summary Create lesson report
// @Description Create a new lesson teaching report
// @Tags lesson Reports
// @Accept json
// @Produce json
// @Param request body dto.CreateLessonReportRequest true "Report data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lesson-reports [post]
// controllers/lesson_report_controller.go - Updated CreateReport
func (c *LessonReportController) CreateReport(ctx *gin.Context) {
	// Get teacher ID from context (set by auth middleware)
	teacherIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	teacherID, err := uuid.Parse(teacherIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req dto.CreateLessonReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Validate academic session ID is provided
	if req.AcademicSessionID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Academic session ID is required",
		})
		return
	}

	report, err := c.reportService.CreateReport(teacherID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "lesson report created successfully",
		"data":    report,
	})
}

// GetReportByID retrieves a report by ID
// @Summary Get report by ID
// @Description Get a specific lesson report by ID
// @Tags lesson Reports
// @Accept json
// @Produce json
// @Param id path string true "Report ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lesson-reports/{id} [get]
func (c *LessonReportController) GetReportByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Report ID is required",
		})
		return
	}

	report, err := c.reportService.GetReportByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Report retrieved successfully",
		"data":    report,
	})
}

// GetReports retrieves reports with filters
// @Summary Get reports
// @Description Get a list of lesson reports with filters
// @Tags lesson Reports
// @Accept json
// @Produce json
// @Param teacher_id query string false "Teacher ID"
// @Param scheme_of_work_id query string false "Scheme of Work ID"
// @Param module_id query string false "Module ID"
// @Param topic_id query string false "Topic ID"
// @Param lesson_id query string false "Lesson ID"
// @Param status query string false "Status (not_taught, in_progress, completed)"
// @Param date_from query string false "Date from (YYYY-MM-DD)"
// @Param date_to query string false "Date to (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/fetch-lesson-reports [get]
func (c *LessonReportController) GetReports(ctx *gin.Context) {
	var params dto.LessonReportQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}



	response, err := c.reportService.GetReports(&params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Reports retrieved successfully",
		"data":    response,
	})
}

// UpdateReport updates an existing report
// @Summary Update report
// @Description Update an existing lesson report
// @Tags lesson Reports
// @Accept json
// @Produce json
// @Param id path string true "Report ID"
// @Param request body dto.UpdateLessonReportRequest true "Update data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lesson-reports/{id} [put]
func (c *LessonReportController) UpdateReport(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Report ID is required",
		})
		return
	}

	var req dto.UpdateLessonReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	report, err := c.reportService.UpdateReport(id, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Report updated successfully",
		"data":    report,
	})
}

// DeleteReport deletes a report
// @Summary Delete report
// @Description Soft delete a lesson report
// @Tags lesson Reports
// @Accept json
// @Produce json
// @Param id path string true "Report ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lesson-reports/{id} [delete]
func (c *LessonReportController) DeleteReport(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Report ID is required",
		})
		return
	}

	if err := c.reportService.DeleteReport(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Report deleted successfully",
	})
}

// GetReportStats retrieves statistics for a teacher
// @Summary Get report statistics
// @Description Get statistics for a teacher's lesson reports
// @Tags lesson Reports
// @Accept json
// @Produce json
// @Param teacher_id query string false "Teacher ID (defaults to authenticated user)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lesson-reports/stats [get]
func (c *LessonReportController) GetReportStats(ctx *gin.Context) {
	teacherID := ctx.Query("teacher_id")
	if teacherID == "" {
		// Use authenticated user
		teacherIDStr, exists := ctx.Get("user_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}
		teacherID = teacherIDStr.(string)
	}

	stats, err := c.reportService.GetReportStats(teacherID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Statistics retrieved successfully",
		"data":    stats,
	})
}

// GetReportsByLesson retrieves reports for a specific lesson
// @Summary Get reports by lesson
// @Description Get all reports for a specific lesson
// @Tags lesson Reports
// @Accept json
// @Produce json
// @Param lesson_id path string true "Lesson ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lesson-reports/lesson/{lesson_id} [get]
func (c *LessonReportController) GetReportsByLesson(ctx *gin.Context) {
	lessonID := ctx.Param("lesson_id")
	if lessonID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Lesson ID is required",
		})
		return
	}

	// Get teacher ID from context
	teacherIDStr, _ := ctx.Get("user_id")
	teacherID := ""
	if teacherIDStr != nil {
		teacherID = teacherIDStr.(string)
	}

	reports, err := c.reportService.GetReportsByLesson(lessonID, teacherID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Reports retrieved successfully",
		"data":    reports,
	})
}