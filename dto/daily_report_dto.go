// dto/daily_report_dto.go
package dto

import (
	"time"
)

// CreateDailyReportRequest represents the request to create a daily report
type CreateDailyReportRequest struct {
	SchemeOfWorkID    string    `json:"scheme_of_work_id" binding:"required"`
	ModuleID          string    `json:"module_id" binding:"required"`
	TopicID           string    `json:"topic_id" binding:"required"`
	LessonID          string    `json:"lesson_id" binding:"required"`
	AcademicSessionID string    `json:"academic_session_id" binding:"required"`
	ReportDate        time.Time `json:"report_date" binding:"required"`
	Status            string    `json:"status" binding:"required,oneof=not_taught in_progress completed"`
	StudentsCovered   string    `json:"students_covered"`
	Report            string    `json:"report"`
	Challenges        string    `json:"challenges"`
	Remarks           string    `json:"remarks"`
	StudentsPresent   int       `json:"students_present"`
	StudentsAbsent    int       `json:"students_absent"`
}

// UpdateDailyReportRequest represents the request to update a daily report
type UpdateDailyReportRequest struct {
	Status          string `json:"status" binding:"omitempty,oneof=not_taught in_progress completed"`
	StudentsCovered string `json:"students_covered"`
	Report          string `json:"report"`
	Challenges      string `json:"challenges"`
	Remarks         string `json:"remarks"`
	StudentsPresent int    `json:"students_present"`
	StudentsAbsent  int    `json:"students_absent"`
}

// DailyReportResponse represents the response for a daily report
type DailyReportResponse struct {
	ID                string    `json:"id"`
	TeacherID         string    `json:"teacher_id"`
	SchemeOfWorkID    string    `json:"scheme_of_work_id"`
	ModuleID          string    `json:"module_id"`
	TopicID           string    `json:"topic_id"`
	LessonID          string    `json:"lesson_id"`
	AcademicSessionID string    `json:"academic_session_id"`
	ReportDate        time.Time `json:"report_date"`
	Status            string    `json:"status"`
	StudentsCovered   string    `json:"students_covered"`
	Report            string    `json:"report"`
	Challenges        string    `json:"challenges"`
	Remarks           string    `json:"remarks"`
	StudentsPresent   int       `json:"students_present"`
	StudentsAbsent    int       `json:"students_absent"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// Nested relationships
	Teacher         *UserResponse            `json:"teacher,omitempty"`
	SchemeOfWork    *SchemeOfWorkResponse    `json:"scheme_of_work,omitempty"`
	Module          *ModuleResponse          `json:"module,omitempty"`
	Topic           *TopicResponse           `json:"topic,omitempty"`
	Lesson          *LessonResponse          `json:"lesson,omitempty"`
	AcademicSession *AcademicSessionResponse `json:"academic_session,omitempty"`
}

// DailyReportListResponse represents a paginated list of daily reports
type DailyReportListResponse struct {
	Reports    []DailyReportResponse `json:"reports"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalPages int                   `json:"total_pages"`
}

// DailyReportQueryParams represents query parameters for filtering reports
type DailyReportQueryParams struct {
	TeacherID         string `form:"teacher_id"`
	SchemeOfWorkID    string `form:"scheme_of_work_id"`
	ModuleID          string `form:"module_id"`
	TopicID           string `form:"topic_id"`
	LessonID          string `form:"lesson_id"`
	AcademicSessionID string `form:"academic_session_id"`
	Status            string `form:"status" binding:"omitempty,oneof=not_taught in_progress completed"`
	DateFrom          string `form:"date_from"`
	DateTo            string `form:"date_to"`
	Page              int    `form:"page" default:"1"`
	Limit             int    `form:"limit" default:"20"`
	SortBy            string `form:"sort_by" default:"report_date"`
	SortOrder         string `form:"sort_order" default:"desc"`
}

// DailyReportStats represents statistics for a teacher's reports
type DailyReportStats struct {
	TotalReports      int64 `json:"total_reports"`
	NotTaught         int64 `json:"not_taught"`
	InProgress        int64 `json:"in_progress"`
	Completed         int64 `json:"completed"`
	TotalStudents     int   `json:"total_students"`
	AverageAttendance int   `json:"average_attendance"`
}
