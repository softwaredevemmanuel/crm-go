// services/daily_report_service.go
package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"crm-go/dto"
	"crm-go/models"
)

type DailyReportService struct {
	db *gorm.DB
}

func NewDailyReportService(db *gorm.DB) *DailyReportService {
	return &DailyReportService{db: db}
}

// CreateReport creates a new daily report
// CreateReport creates a new daily report
func (s *DailyReportService) CreateReport(teacherID uuid.UUID, req *dto.CreateDailyReportRequest) (*dto.DailyReportResponse, error) {
	// Parse UUIDs
	schemeOfWorkID, err := uuid.Parse(req.SchemeOfWorkID)
	if err != nil {
		return nil, errors.New("invalid scheme of work ID")
	}

	moduleID, err := uuid.Parse(req.ModuleID)
	if err != nil {
		return nil, errors.New("invalid module ID")
	}

	topicID, err := uuid.Parse(req.TopicID)
	if err != nil {
		return nil, errors.New("invalid topic ID")
	}

	lessonID, err := uuid.Parse(req.LessonID)
	if err != nil {
		return nil, errors.New("invalid lesson ID")
	}

	academicSessionID, err := uuid.Parse(req.AcademicSessionID)
	if err != nil {
		return nil, errors.New("invalid academic session ID")
	}

	// Verify that the teacher exists
	var teacher models.User
	if err := s.db.Where("id = ? AND deleted_at IS NULL", teacherID).First(&teacher).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("teacher not found")
		}
		return nil, errors.New("failed to verify teacher")
	}

	// Verify that the teacher is a teacher or admin
	if teacher.Position != "teacher" && teacher.Position != "admin" {
		return nil, errors.New("user is not a teacher or admin")
	}

	// Verify that the scheme of work exists
	var scheme models.SchemeOfWork
	if err := s.db.Where("id = ? AND deleted_at IS NULL", schemeOfWorkID).First(&scheme).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("scheme of work not found")
		}
		return nil, errors.New("failed to verify scheme of work")
	}

	// Verify that the module exists
	var module models.Module
	if err := s.db.Where("id = ? AND deleted_at IS NULL", moduleID).First(&module).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("module not found")
		}
		return nil, errors.New("failed to verify module")
	}

	// Verify that the topic exists
	var topic models.Topic
	if err := s.db.Where("id = ? AND deleted_at IS NULL", topicID).First(&topic).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("topic not found")
		}
		return nil, errors.New("failed to verify topic")
	}

	// Verify that the lesson exists
	var lesson models.Lesson
	if err := s.db.Where("id = ? AND deleted_at IS NULL", lessonID).First(&lesson).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lesson not found")
		}
		return nil, errors.New("failed to verify lesson")
	}

	// Verify that the academic session exists and is current
	var academicSession models.AcademicSession
	if err := s.db.Where("id = ? AND deleted_at IS NULL", academicSessionID).First(&academicSession).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("academic session not found")
		}
		return nil, errors.New("failed to verify academic session")
	}

	// Check if the academic session is current
	if !academicSession.IsCurrent {
		return nil, errors.New("reports can only be created for the current academic session")
	}

	// Check if the academic session is active
	if academicSession.Status != "active" {
		return nil, errors.New("academic session is not active")
	}

	// Check if a report already exists for this lesson in the same academic session
	// This checks for any report with the same lesson_id and academic_session_id
	var existingReport models.DailyReport
	err = s.db.Where(
		"lesson_id = ? AND academic_session_id = ? AND deleted_at IS NULL",
		lessonID, academicSessionID,
	).First(&existingReport).Error

	if err == nil {
		return nil, errors.New("a report already exists for this lesson in the current academic session")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("failed to check for existing report")
	}

	// Additional check: If report_date is provided, check for duplicate on the same date
	// This is a secondary check for the same day
	if !req.ReportDate.IsZero() {
		var dateReport models.DailyReport
		startOfDay := time.Date(req.ReportDate.Year(), req.ReportDate.Month(), req.ReportDate.Day(), 0, 0, 0, 0, req.ReportDate.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)
		
		err = s.db.Where(
			"teacher_id = ? AND lesson_id = ? AND report_date >= ? AND report_date < ? AND academic_session_id = ? AND deleted_at IS NULL",
			teacherID, lessonID, startOfDay, endOfDay, academicSessionID,
		).First(&dateReport).Error

		if err == nil {
			return nil, errors.New("a report already exists for this lesson on the selected date")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("failed to check for existing report on date")
		}
	}

	// Create the report
	report := &models.DailyReport{
		ID:                uuid.New(),
		TeacherID:         teacherID,
		SchemeOfWorkID:    schemeOfWorkID,
		ModuleID:          moduleID,
		TopicID:           topicID,
		LessonID:          lessonID,
		AcademicSessionID: academicSessionID,
		ReportDate:        req.ReportDate,
		Status:            req.Status,
		StudentsCovered:   req.StudentsCovered,
		Report:            req.Report,
		Challenges:        req.Challenges,
		Remarks:           req.Remarks,
		StudentsPresent:   req.StudentsPresent,
		StudentsAbsent:    req.StudentsAbsent,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.db.Create(report).Error; err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	// Preload relationships for response
	if err := s.db.Preload("Teacher").
		Preload("SchemeOfWork").
		Preload("Module").
		Preload("Topic").
		Preload("Lesson").
		Preload("AcademicSession").
		First(report, report.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load report details: %w", err)
	}

	return s.toReportResponse(report), nil
}

// GetReportByID retrieves a report by ID
func (s *DailyReportService) GetReportByID(id string) (*dto.DailyReportResponse, error) {
	reportID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid report ID")
	}

	var report models.DailyReport
	if err := s.db.Where("id = ? AND deleted_at IS NULL", reportID).
		Preload("Teacher").
		Preload("SchemeOfWork").
		Preload("Module").
		Preload("Topic").
		Preload("Lesson").
		First(&report).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("report not found")
		}
		return nil, fmt.Errorf("failed to fetch report: %w", err)
	}

	return s.toReportResponse(&report), nil
}

// GetReports retrieves reports with filters
func (s *DailyReportService) GetReports(params *dto.DailyReportQueryParams) (*dto.DailyReportListResponse, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = "report_date"
	}
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}

	// Build query
	query := s.db.Model(&models.DailyReport{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.TeacherID != "" {
		teacherID, err := uuid.Parse(params.TeacherID)
		if err == nil {
			query = query.Where("teacher_id = ?", teacherID)
		}
	}

	if params.SchemeOfWorkID != "" {
		schemeID, err := uuid.Parse(params.SchemeOfWorkID)
		if err == nil {
			query = query.Where("scheme_of_work_id = ?", schemeID)
		}
	}

	if params.ModuleID != "" {
		moduleID, err := uuid.Parse(params.ModuleID)
		if err == nil {
			query = query.Where("module_id = ?", moduleID)
		}
	}

	if params.TopicID != "" {
		topicID, err := uuid.Parse(params.TopicID)
		if err == nil {
			query = query.Where("topic_id = ?", topicID)
		}
	}

	if params.LessonID != "" {
		lessonID, err := uuid.Parse(params.LessonID)
		if err == nil {
			query = query.Where("lesson_id = ?", lessonID)
		}
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.DateFrom != "" {
		dateFrom, err := time.Parse("2006-01-02", params.DateFrom)
		if err == nil {
			query = query.Where("report_date >= ?", dateFrom)
		}
	}

	if params.DateTo != "" {
		dateTo, err := time.Parse("2006-01-02", params.DateTo)
		if err == nil {
			query = query.Where("report_date <= ?", dateTo)
		}
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count reports: %w", err)
	}

	// Apply sorting
	sortDirection := "DESC"
	if strings.ToLower(params.SortOrder) == "asc" {
		sortDirection = "ASC"
	}
	query = query.Order(fmt.Sprintf("%s %s", params.SortBy, sortDirection))

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	// Execute with preloads
	var reports []models.DailyReport
	if err := query.Preload("Teacher").
		Preload("SchemeOfWork").
		Preload("Module").
		Preload("Topic").
		Preload("Lesson").
		Find(&reports).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch reports: %w", err)
	}

	// Convert to response
	responses := make([]dto.DailyReportResponse, len(reports))
	for i, report := range reports {
		responses[i] = *s.toReportResponse(&report)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.DailyReportListResponse{
		Reports:    responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// UpdateReport updates an existing report
func (s *DailyReportService) UpdateReport(id string, req *dto.UpdateDailyReportRequest) (*dto.DailyReportResponse, error) {
	reportID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid report ID")
	}

	var report models.DailyReport
	if err := s.db.Where("id = ? AND deleted_at IS NULL", reportID).First(&report).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("report not found")
		}
		return nil, fmt.Errorf("failed to fetch report: %w", err)
	}

	// Update fields
	if req.Status != "" {
		report.Status = req.Status
	}
	if req.StudentsCovered != "" {
		report.StudentsCovered = req.StudentsCovered
	}
	if req.Report != "" {
		report.Report = req.Report
	}
	if req.Challenges != "" {
		report.Challenges = req.Challenges
	}
	if req.Remarks != "" {
		report.Remarks = req.Remarks
	}
	if req.StudentsPresent > 0 {
		report.StudentsPresent = req.StudentsPresent
	}
	if req.StudentsAbsent > 0 {
		report.StudentsAbsent = req.StudentsAbsent
	}

	report.UpdatedAt = time.Now()

	if err := s.db.Save(&report).Error; err != nil {
		return nil, fmt.Errorf("failed to update report: %w", err)
	}

	// Preload relationships
	if err := s.db.Preload("Teacher").
		Preload("SchemeOfWork").
		Preload("Module").
		Preload("Topic").
		Preload("Lesson").
		First(&report, report.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load report details: %w", err)
	}

	return s.toReportResponse(&report), nil
}

// DeleteReport soft deletes a report
func (s *DailyReportService) DeleteReport(id string) error {
	reportID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid report ID")
	}

	var report models.DailyReport
	if err := s.db.Where("id = ? AND deleted_at IS NULL", reportID).First(&report).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("report not found")
		}
		return fmt.Errorf("failed to fetch report: %w", err)
	}

	if err := s.db.Delete(&report).Error; err != nil {
		return fmt.Errorf("failed to delete report: %w", err)
	}

	return nil
}

// GetReportStats retrieves statistics for a teacher
func (s *DailyReportService) GetReportStats(teacherID string) (*dto.DailyReportStats, error) {
	tID, err := uuid.Parse(teacherID)
	if err != nil {
		return nil, errors.New("invalid teacher ID")
	}

	// Verify teacher exists
	var teacher models.User
	if err := s.db.Where("id = ? AND deleted_at IS NULL", tID).First(&teacher).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("teacher not found")
		}
		return nil, fmt.Errorf("failed to verify teacher: %w", err)
	}

	var stats dto.DailyReportStats

	// Get total reports
	if err := s.db.Model(&models.DailyReport{}).
		Where("teacher_id = ? AND deleted_at IS NULL", tID).
		Count(&stats.TotalReports).Error; err != nil {
		return nil, fmt.Errorf("failed to count total reports: %w", err)
	}

	// Get reports by status
	var statusCounts []struct {
		Status string
		Count  int64
	}
	if err := s.db.Model(&models.DailyReport{}).
		Select("status, count(*) as count").
		Where("teacher_id = ? AND deleted_at IS NULL", tID).
		Group("status").
		Scan(&statusCounts).Error; err != nil {
		return nil, fmt.Errorf("failed to count status reports: %w", err)
	}

	for _, sc := range statusCounts {
		switch sc.Status {
		case "not_taught":
			stats.NotTaught = sc.Count
		case "in_progress":
			stats.InProgress = sc.Count
		case "completed":
			stats.Completed = sc.Count
		}
	}

	// Calculate average attendance
	var result struct {
		TotalPresent int
		TotalAbsent  int
		Count        int
	}
	if err := s.db.Model(&models.DailyReport{}).
		Select("COALESCE(SUM(students_present), 0) as total_present, COALESCE(SUM(students_absent), 0) as total_absent, COUNT(*) as count").
		Where("teacher_id = ? AND deleted_at IS NULL", tID).
		Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate attendance: %w", err)
	}

	stats.TotalStudents = result.TotalPresent + result.TotalAbsent
	if result.Count > 0 {
		stats.AverageAttendance = (result.TotalPresent * 100) / (result.TotalPresent + result.TotalAbsent)
	}

	return &stats, nil
}

// GetReportsByLesson retrieves reports for a specific lesson
func (s *DailyReportService) GetReportsByLesson(lessonID string, teacherID string) ([]dto.DailyReportResponse, error) {
	lID, err := uuid.Parse(lessonID)
	if err != nil {
		return nil, errors.New("invalid lesson ID")
	}

	var reports []models.DailyReport
	query := s.db.Where("lesson_id = ? AND deleted_at IS NULL", lID)

	if teacherID != "" {
		tID, err := uuid.Parse(teacherID)
		if err == nil {
			query = query.Where("teacher_id = ?", tID)
		}
	}

	if err := query.Preload("Teacher").
		Preload("SchemeOfWork").
		Preload("Module").
		Preload("Topic").
		Preload("Lesson").
		Order("report_date DESC").
		Find(&reports).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch reports: %w", err)
	}

	responses := make([]dto.DailyReportResponse, len(reports))
	for i, report := range reports {
		responses[i] = *s.toReportResponse(&report)
	}

	return responses, nil
}

// toReportResponse converts model to response DTO
func (s *DailyReportService) toReportResponse(report *models.DailyReport) *dto.DailyReportResponse {
	response := &dto.DailyReportResponse{
		ID:                report.ID.String(),
		TeacherID:         report.TeacherID.String(),
		SchemeOfWorkID:    report.SchemeOfWorkID.String(),
		ModuleID:          report.ModuleID.String(),
		TopicID:           report.TopicID.String(),
		LessonID:          report.LessonID.String(),
		AcademicSessionID: report.AcademicSessionID.String(),
		ReportDate:        report.ReportDate,
		Status:            report.Status,
		StudentsCovered:   report.StudentsCovered,
		Report:            report.Report,
		Challenges:        report.Challenges,
		Remarks:           report.Remarks,
		StudentsPresent:   report.StudentsPresent,
		StudentsAbsent:    report.StudentsAbsent,
		CreatedAt:         report.CreatedAt,
		UpdatedAt:         report.UpdatedAt,
	}

	// Add teacher details
	if report.Teacher.ID != uuid.Nil {
		response.Teacher = &dto.UserResponse{
			ID:        report.Teacher.ID.String(),
			FirstName: report.Teacher.FirstName,
			LastName:  report.Teacher.LastName,
			Email:     report.Teacher.Email,
		}
	}

	// Add academic session details
	if report.AcademicSession.ID != uuid.Nil {
		response.AcademicSession = &dto.AcademicSessionResponse{
			ID:          report.AcademicSession.ID.String(),
			AcademicYear:        report.AcademicSession.AcademicYear,
			StartDate:   report.AcademicSession.StartDate,
			EndDate:     report.AcademicSession.EndDate,
			IsCurrent:   report.AcademicSession.IsCurrent,
			Status:      report.AcademicSession.Status,
			Description: report.AcademicSession.Description,
			CreatedAt:   report.AcademicSession.CreatedAt,
			UpdatedAt:   report.AcademicSession.UpdatedAt,
		}
	}

	return response
}