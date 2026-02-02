// controllers/grade_controller.go
package controllers

import (
    "net/http"
    "strconv"
    
    "crm-go/models"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
	"time"
)

// GetAllGrades handler
// @Summary Get all grades
// @Description Get grades with filtering, sorting, and pagination
// @Tags grades
// @Accept json
// @Produce json
// @Param student_id query string false "Filter by student ID"
// @Param course_id query string false "Filter by course ID"
// @Param assignment_id query string false "Filter by assignment ID"
// @Param tutor_id query string false "Filter by tutor ID"
// @Param min_score query number false "Minimum score" minimum(0) maximum(100)
// @Param max_score query number false "Maximum score" minimum(0) maximum(100)
// @Param grade query string false "Filter by grade letter (A, B, C, D, F)"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param search query string false "Search in remarks"
// @Param with_details query boolean false "Include student/course details"
// @Param page query int false "Page number" default(1) minimum(1)
// @Param limit query int false "Items per page" default(20) minimum(1) maximum(100)
// @Param sort_by query string false "Sort field" Enums(created_at, updated_at, score, student, course)
// @Param sort_order query string false "Sort order" Enums(asc, desc, ASC, DESC)
// @Success 200 {object} models.PaginatedGradesResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /grades [get]
// @Security BearerAuth
func (ctl *GradeController) GetAllGrades(c *gin.Context) {
    var filters models.GradeFilters
    
    // Parse UUID filters
    if studentID := c.Query("student_id"); studentID != "" {
        if id, err := uuid.Parse(studentID); err == nil {
            filters.StudentID = id
        }
    }
    
    if courseID := c.Query("course_id"); courseID != "" {
        if id, err := uuid.Parse(courseID); err == nil {
            filters.CourseID = id
        }
    }
    
    if assignmentID := c.Query("assignment_id"); assignmentID != "" {
        if id, err := uuid.Parse(assignmentID); err == nil {
            filters.AssignmentID = id
        }
    }
    
    if tutorID := c.Query("tutor_id"); tutorID != "" {
        if id, err := uuid.Parse(tutorID); err == nil {
            filters.TutorID = id
        }
    }
    
    // Parse numeric filters
    if minScore := c.Query("min_score"); minScore != "" {
        if val, err := strconv.ParseFloat(minScore, 64); err == nil {
            filters.MinScore = val
        }
    }
    
    if maxScore := c.Query("max_score"); maxScore != "" {
        if val, err := strconv.ParseFloat(maxScore, 64); err == nil {
            filters.MaxScore = val
        }
    }
    
    // Parse string filters
    filters.GradeLetter = c.Query("grade")
    filters.Search = c.Query("search")
    
    // Parse date filters
    if startDate := c.Query("start_date"); startDate != "" {
        if t, err := time.Parse("2006-01-02", startDate); err == nil {
            filters.StartDate = t
        }
    }
    
    if endDate := c.Query("end_date"); endDate != "" {
        if t, err := time.Parse("2006-01-02", endDate); err == nil {
            filters.EndDate = t
        }
    }
    
    // Parse boolean
    if withDetails := c.Query("with_details"); withDetails != "" {
        if val, err := strconv.ParseBool(withDetails); err == nil {
            filters.WithDetails = val
        }
    }
    
    // Parse pagination
    if page := c.Query("page"); page != "" {
        if val, err := strconv.Atoi(page); err == nil && val > 0 {
            filters.Page = val
        }
    }
    
    if limit := c.Query("limit"); limit != "" {
        if val, err := strconv.Atoi(limit); err == nil && val > 0 {
            if val > 100 {
                val = 100 // Enforce maximum
            }
            filters.Limit = val
        }
    }
    
    // Parse sorting
    filters.SortBy = c.Query("sort_by")
    filters.SortOrder = c.Query("sort_order")
    
    // If user is tutor/teacher, filter by their courses automatically
    if filters.TutorID == uuid.Nil {
        if tutorID, exists := c.Get("user_id"); exists {
            if id, ok := tutorID.(uuid.UUID); ok {
                filters.TutorID = id
            }
        }
    }
    
    // Get paginated results
    result, err := ctl.gradeService.GetAllGradesWithPagination(filters)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, result)
}

// GetGradeStats handler
// @Summary Get grade statistics
// @Description Get statistics for filtered grades
// @Tags grades
// @Accept json
// @Produce json
// @Param student_id query string false "Filter by student ID"
// @Param course_id query string false "Filter by course ID"
// @Success 200 {object} models.GradeStats
// @Router /grades/stats [get]
// @Security BearerAuth
func (ctl *GradeController) GetGradeStats(c *gin.Context) {
    var filters models.GradeFilters
    
    // Parse filters (similar to GetAllGrades)
    // ... parsing code ...
    
    stats, err := ctl.gradeService.GetGradeStats(filters)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, stats)
}

// GetStudentGrades handler
// @Summary Get grades for a specific student
// @Description Get all grades for a student
// @Tags grades
// @Accept json
// @Produce json
// @Param student_id path string true "Student ID"
// @Param with_details query boolean false "Include course details"
// @Success 200 {array} models.GradeResponse
// @Router /grades/student/{student_id}/grades [get]
// @Security BearerAuth
func (ctl *GradeController) GetStudentGrades(c *gin.Context) {
    studentID, err := uuid.Parse(c.Param("student_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid student ID",
        })
        return
    }
    
    withDetails := false
    if details := c.Query("with_details"); details != "" {
        if val, err := strconv.ParseBool(details); err == nil {
            withDetails = val
        }
    }
    
    grades, err := ctl.gradeService.GetGradesByStudent(studentID, withDetails)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "data": grades,
        "count": len(grades),
    })
}

// GetCourseGrades handler
// @Summary Get grades for a specific course
// @Description Get all grades for a course
// @Tags grades
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Param with_details query boolean false "Include student details"
// @Success 200 {array} models.GradeResponse
// @Router /grades/courses/{course_id}/grades [get]
// @Security BearerAuth
func (ctl *GradeController) GetCourseGrades(c *gin.Context) {
    courseID, err := uuid.Parse(c.Param("course_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid course ID",
        })
        return
    }
    
    withDetails := false
    if details := c.Query("with_details"); details != "" {
        if val, err := strconv.ParseBool(details); err == nil {
            withDetails = val
        }
    }
    
    grades, err := ctl.gradeService.GetGradesByCourse(courseID, withDetails)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    // Get statistics for this course
    stats, _ := ctl.gradeService.GetGradeStats(models.GradeFilters{
        CourseID: courseID,
    })
    
    c.JSON(http.StatusOK, gin.H{
        "data":     grades,
        "count":    len(grades),
        "stats":    stats,
    })
}