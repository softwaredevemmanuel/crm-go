// services/grade_service.go
package services

import (
    "errors"
    "fmt"
    "strings"
    "time"
    
    "crm-go/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

// GetAllGrades - main function to get grades with filters
func (s *GradeService) GetAllGrades(filters models.GradeFilters) ([]models.GradeResponse, error) {
    var grades []models.Grade
    
    // Start building query
    query := s.db.Model(&models.Grade{})
    
    // Apply filters
    query = s.applyGradeFilters(query, filters)
    
    // Apply sorting
    query = s.applyGradeSorting(query, filters)
    
    // Apply pagination
    if filters.Page > 0 && filters.Limit > 0 {
        offset := (filters.Page - 1) * filters.Limit
        query = query.Offset(offset).Limit(filters.Limit)
    }
    
    // Execute query
    if err := query.Find(&grades).Error; err != nil {
        return nil, errors.New("failed to fetch grades: " + err.Error())
    }
    
    // Convert to response
    return s.gradesToResponse(grades, filters.WithDetails), nil
}

// applyGradeFilters - helper to apply all filters
func (s *GradeService) applyGradeFilters(query *gorm.DB, filters models.GradeFilters) *gorm.DB {
    // Filter by student
    if filters.StudentID != uuid.Nil {
        query = query.Where("student_id = ?", filters.StudentID)
    }
    
    // Filter by course
    if filters.CourseID != uuid.Nil {
        query = query.Where("course_id = ?", filters.CourseID)
    }
    
    // Filter by assignment
    if filters.AssignmentID != uuid.Nil {
        query = query.Where("assignment_id = ?", filters.AssignmentID)
    }
    
    // Filter by assignment is null
    if filters.AssignmentID == uuid.Nil && 
       strings.Contains(filters.Search, "assignment:null") {
        query = query.Where("assignment_id IS NULL")
    }
    
    // Filter by score range
    if filters.MinScore > 0 {
        query = query.Where("score >= ?", filters.MinScore)
    }
    
    if filters.MaxScore > 0 && filters.MaxScore <= 100 {
        query = query.Where("score <= ?", filters.MaxScore)
    }
    
    // Filter by grade letter
    if filters.GradeLetter != "" {
        query = query.Where("grade = ?", filters.GradeLetter)
    }
    
    // Filter by date range
    if !filters.StartDate.IsZero() {
        query = query.Where("created_at >= ?", filters.StartDate)
    }
    
    if !filters.EndDate.IsZero() {
        // Add one day to include the entire end date
        endDate := filters.EndDate.Add(24 * time.Hour)
        query = query.Where("created_at < ?", endDate)
    }
    
    // Search in remarks
    if filters.Search != "" && !strings.Contains(filters.Search, "assignment:null") {
        searchTerm := "%" + strings.ToLower(filters.Search) + "%"
        query = query.Where("LOWER(remarks) LIKE ?", searchTerm)
    }
    
    // Filter by tutor (teacher) - only show grades for courses they teach
    if filters.TutorID != uuid.Nil {
        query = query.Joins("JOIN courses ON courses.id = grades.course_id").
            Where("courses.tutor_id = ?", filters.TutorID)
    }
    
    return query
}

// applyGradeSorting - helper to apply sorting
func (s *GradeService) applyGradeSorting(query *gorm.DB, filters models.GradeFilters) *gorm.DB {
    // Default sorting
    sortBy := "created_at"
    sortOrder := "DESC"
    
    // Override with user preferences
    if filters.SortBy != "" {
        switch filters.SortBy {
        case "score":
            sortBy = "score"
        case "student":
            sortBy = "student_id" // Would need join for name sorting
        case "course":
            sortBy = "course_id"  // Would need join for title sorting
        case "updated_at":
            sortBy = "updated_at"
        default:
            sortBy = filters.SortBy
        }
    }
    
    if filters.SortOrder != "" {
        sortOrder = strings.ToUpper(filters.SortOrder)
    }
    
    return query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
}

// GetGradesCount - get total count with filters
func (s *GradeService) GetGradesCount(filters models.GradeFilters) (int64, error) {
    var count int64
    
    query := s.db.Model(&models.Grade{})
    query = s.applyGradeFilters(query, filters)
    
    if err := query.Count(&count).Error; err != nil {
        return 0, errors.New("failed to count grades: " + err.Error())
    }
    
    return count, nil
}

// GetAllGradesWithPagination - returns grades with pagination metadata
func (s *GradeService) GetAllGradesWithPagination(filters models.GradeFilters) (*models.PaginatedGradesResponse, error) {
    // Get grades
    grades, err := s.GetAllGrades(filters)
    if err != nil {
        return nil, err
    }
    
    // Get total count
    totalCount, err := s.GetGradesCount(filters)
    if err != nil {
        return nil, err
    }
    
    // Calculate total pages
    totalPages := 0
    if filters.Limit > 0 {
        totalPages = int((totalCount + int64(filters.Limit) - 1) / int64(filters.Limit))
    }
    
    return &models.PaginatedGradesResponse{
        Data:       grades,
        Total:      totalCount,
        Page:       filters.Page,
        Limit:      filters.Limit,
        TotalPages: totalPages,
        Filters:    filters,
    }, nil
}

// GetGradeStats - get statistics for filtered grades
func (s *GradeService) GetGradeStats(filters models.GradeFilters) (*models.GradeStats, error) {
    var stats models.GradeStats
    stats.GradeDistribution = make(map[string]int)
    
    query := s.db.Model(&models.Grade{})
    query = s.applyGradeFilters(query, filters)
    
    // Get total count
    if err := query.Count(&stats.TotalCount).Error; err != nil {
        return nil, err
    }
    
    // Get average score
    if err := query.Select("COALESCE(AVG(score), 0)").Scan(&stats.AverageScore).Error; err != nil {
        return nil, err
    }
    
    // Get highest score
    if err := query.Select("COALESCE(MAX(score), 0)").Scan(&stats.HighestScore).Error; err != nil {
        return nil, err
    }
    
    // Get lowest score
    if err := query.Select("COALESCE(MIN(score), 0)").Scan(&stats.LowestScore).Error; err != nil {
        return nil, err
    }
    
    // Get grade distribution
    var gradeDist []struct {
        Grade string
        Count int
    }
    
    if err := query.Select("grade, COUNT(*) as count").
        Group("grade").
        Order("grade").
        Scan(&gradeDist).Error; err != nil {
        return nil, err
    }
    
    for _, dist := range gradeDist {
        stats.GradeDistribution[dist.Grade] = dist.Count
    }
    
    return &stats, nil
}

// GetGradesByStudent - convenience method
func (s *GradeService) GetGradesByStudent(studentID uuid.UUID, withDetails bool) ([]models.GradeResponse, error) {
    filters := models.GradeFilters{
        StudentID:   studentID,
        WithDetails: withDetails,
        SortBy:      "course_id",
        SortOrder:   "asc",
    }
    return s.GetAllGrades(filters)
}

// GetGradesByCourse - convenience method
func (s *GradeService) GetGradesByCourse(courseID uuid.UUID, withDetails bool) ([]models.GradeResponse, error) {
    filters := models.GradeFilters{
        CourseID:    courseID,
        WithDetails: withDetails,
        SortBy:      "score",
        SortOrder:   "desc",
    }
    return s.GetAllGrades(filters)
}

// GetStudentCourseGrade - get specific student's grade for a course
func (s *GradeService) GetStudentCourseGrade(studentID, courseID uuid.UUID) (*models.GradeResponse, error) {
    filters := models.GradeFilters{
        StudentID: studentID,
        CourseID:  courseID,
        Limit:     1,
    }
    
    grades, err := s.GetAllGrades(filters)
    if err != nil {
        return nil, err
    }
    
    if len(grades) == 0 {
        return nil, errors.New("grade not found for this student and course")
    }
    
    return &grades[0], nil
}

// Helper to convert multiple grades to responses
func (s *GradeService) gradesToResponse(grades []models.Grade, withDetails bool) []models.GradeResponse {
    responses := make([]models.GradeResponse, len(grades))
    
    // If we need details, preload related data efficiently
    if withDetails && len(grades) > 0 {
        // Batch load students and courses
        studentIDs := make([]uuid.UUID, 0, len(grades))
        courseIDs := make([]uuid.UUID, 0, len(grades))
        
        for _, grade := range grades {
            studentIDs = append(studentIDs, grade.StudentID)
            courseIDs = append(courseIDs, grade.CourseID)
        }
        
        // Get students
        var students []models.User
        s.db.Where("id IN ?", studentIDs).Find(&students)
        studentMap := make(map[uuid.UUID]models.User)
        for _, student := range students {
            studentMap[student.ID] = student
        }
        
        // Get courses
        var courses []models.Course
        s.db.Where("id IN ?", courseIDs).Find(&courses)
        courseMap := make(map[uuid.UUID]models.Course)
        for _, course := range courses {
            courseMap[course.ID] = course
        }
        
        // Build responses with details
        for i, grade := range grades {
            var studentName, courseName string
            
            if student, exists := studentMap[grade.StudentID]; exists {
                studentName = fmt.Sprintf("%s %s", student.FirstName, student.LastName)
            }
            
            if course, exists := courseMap[grade.CourseID]; exists {
                courseName = course.Title
            }
            
            responses[i] = models.GradeResponse{
                ID:           grade.ID,
                StudentID:    grade.StudentID,
                StudentName:  studentName,
                CourseID:     grade.CourseID,
                TutorID:      grade.TutorID,
                CourseName:   courseName,
                AssignmentID: grade.AssignmentID,
                Score:        grade.Score,
                Grade:        grade.Grade,
                Percentage:   grade.Score,
                Remarks:      grade.Remarks,
                CreatedAt:    grade.CreatedAt,
                UpdatedAt:    grade.UpdatedAt,
            }
        }
    } else {
        // Without details
        for i, grade := range grades {
            responses[i] = models.GradeResponse{
                ID:           grade.ID,
                StudentID:    grade.StudentID,
                CourseID:     grade.CourseID,
                TutorID:      grade.TutorID,
                AssignmentID: grade.AssignmentID,
                Score:        grade.Score,
                Grade:        grade.Grade,
                Percentage:   grade.Score,
                Remarks:      grade.Remarks,
                CreatedAt:    grade.CreatedAt,
                UpdatedAt:    grade.UpdatedAt,
            }
        }
    }
    
    return responses
}