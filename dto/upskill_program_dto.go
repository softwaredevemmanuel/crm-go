package dto

import (
	"time"

	"github.com/google/uuid"

	"crm-go/models"
)

// ============================================================
// Program Course (pivot)
// ============================================================
type ProgramCourseInput struct {
	CourseID    uuid.UUID `json:"course_id" binding:"required"`
	CourseOrder int       `json:"course_order" binding:"omitempty,gte=1"`
	IsRequired  bool      `json:"is_required"`
}

type ProgramCourseResponse struct {
	ID          uuid.UUID         `json:"id"`
	CourseID    uuid.UUID         `json:"course_id"`
	CourseOrder int               `json:"course_order"`
	IsRequired  bool              `json:"is_required"`
	Course      *CourseMiniResponse `json:"course,omitempty"`
}

// CourseMiniResponse — lightweight course info embedded in a program
type CourseMiniResponse struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Code         string    `json:"code"`
	Slug         string    `json:"slug"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	Level        string    `json:"level,omitempty"`
}

// ============================================================
// Program
// ============================================================
type CreateUpskillProgramRequest struct {
	Name        string                `json:"name" binding:"required,min=2,max=150"`
	Slug        string                `json:"slug" binding:"omitempty,max=150"`
	Description string                `json:"description" binding:"max=5000"`
	IsActive    *bool                 `json:"is_active"`
	Courses     []ProgramCourseInput  `json:"courses" binding:"omitempty,dive"`
}

type UpdateUpskillProgramRequest struct {
	Name        *string               `json:"name" binding:"omitempty,min=2,max=150"`
	Slug        *string               `json:"slug" binding:"omitempty,max=150"`
	Description *string               `json:"description" binding:"omitempty,max=5000"`
	IsActive    *bool                 `json:"is_active"`
	Courses     *[]ProgramCourseInput `json:"courses" binding:"omitempty,dive"`
}

type UpskillProgramResponse struct {
	ID          uuid.UUID               `json:"id"`
	Name        string                  `json:"name"`
	Slug        string                  `json:"slug"`
	Description string                  `json:"description"`
	IsActive    bool                    `json:"is_active"`
	Courses     []ProgramCourseResponse `json:"courses,omitempty"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

// ============================================================
// List / envelopes
// ============================================================
type ListUpskillProgramsQuery struct {
	Page      int    `form:"page,default=1" binding:"min=1"`
	Limit     int    `form:"limit,default=20" binding:"min=1,max=100"`
	Search    string `form:"search"`
	IsActive  string `form:"is_active"`
	SortBy    string `form:"sort_by,default=created_at"`
	SortOrder string `form:"sort_order,default=desc"`
}

type SingleUpskillProgramResponse struct {
	Data UpskillProgramResponse `json:"data"`
}

type PaginatedUpskillProgramsResponse struct {
	Data       []UpskillProgramResponse `json:"data"`
	Total      int64                    `json:"total"`
	Page       int                      `json:"page"`
	Limit      int                      `json:"limit"`
	TotalPages int                      `json:"total_pages"`
}


// ============================================================
// Mapping helpers
// ============================================================
func FromUpskillProgramCourseModel(pc *models.UpskillProgramCourse) ProgramCourseResponse {
	resp := ProgramCourseResponse{
		ID:          pc.ID,
		CourseID:    pc.CourseID,
		CourseOrder: pc.CourseOrder,
		IsRequired:  pc.IsRequired,
	}
	if pc.Course.ID != uuid.Nil {
		resp.Course = &CourseMiniResponse{
			ID:           pc.Course.ID,
			Title:        pc.Course.Title,
			Code:         pc.Course.Code,
			Slug:         pc.Course.Slug,
			ThumbnailURL: pc.Course.ThumbnailURL,
			Level:        string(pc.Course.Level),
		}
	}
	return resp
}

func FromUpskillProgramModel(p *models.UpskillProgram, includeCourses bool) UpskillProgramResponse {
	resp := UpskillProgramResponse{
		ID:          p.ID,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
	if includeCourses {
		resp.Courses = make([]ProgramCourseResponse, 0, len(p.Courses))
		for i := range p.Courses {
			resp.Courses = append(resp.Courses, FromUpskillProgramCourseModel(&p.Courses[i]))
		}
	}
	return resp
}