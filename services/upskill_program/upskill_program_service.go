package services

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"crm-go/dto"
	"crm-go/models"
)

var (
	ErrProgramNotFound     = errors.New("upskill program not found")
	ErrProgramNameTaken    = errors.New("a program with this name already exists")
	ErrProgramSlugTaken    = errors.New("a program with this slug already exists")
	ErrProgramCourseNoDup  = errors.New("duplicate course in program (same course added twice)")
)

// Slugify helper (inline, no utils package)
var (
	nonAlnum = regexp.MustCompile(`[^a-z0-9\s-]`)
	spaces   = regexp.MustCompile(`[\s-]+`)
)

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = nonAlnum.ReplaceAllString(s, "")
	s = spaces.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

type UpskillProgramService struct {
	DB *gorm.DB
}

func NewUpskillProgramService(db *gorm.DB) *UpskillProgramService {
	return &UpskillProgramService{DB: db}
}

// ---------- Create ----------
func (s *UpskillProgramService) Create(req dto.CreateUpskillProgramRequest) (*models.UpskillProgram, error) {
	// Name uniqueness
	var count int64
	if err := s.DB.Model(&models.UpskillProgram{}).Where("name = ?", req.Name).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrProgramNameTaken
	}

	// Slug
	slug := req.Slug
	if slug == "" {
		slug = slugify(req.Name)
	}
	if err := s.DB.Model(&models.UpskillProgram{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		// Disambiguate
		slug = fmt.Sprintf("%s-%s", slug, uuid.NewString()[:6])
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// Check duplicate courses in payload
	seen := map[uuid.UUID]bool{}
	for _, pc := range req.Courses {
		if seen[pc.CourseID] {
			return nil, ErrProgramCourseNoDup
		}
		seen[pc.CourseID] = true
	}

	// Validate that all requested courses exist
	if len(req.Courses) > 0 {
		ids := make([]uuid.UUID, 0, len(req.Courses))
		for _, pc := range req.Courses {
			ids = append(ids, pc.CourseID)
		}
		var found int64
		if err := s.DB.Model(&models.Course{}).Where("id IN ?", ids).Count(&found).Error; err != nil {
			return nil, err
		}
		if int(found) != len(ids) {
			return nil, fmt.Errorf("one or more courses were not found")
		}
	}

	program := &models.UpskillProgram{
		Name:        strings.TrimSpace(req.Name),
		Slug:        slug,
		Description: req.Description,
		IsActive:    isActive,
	}

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(program).Error; err != nil {
			return err
		}

		// Attach courses via pivot
		for i, pc := range req.Courses {
			order := pc.CourseOrder
			if order == 0 {
				order = i + 1
			}
			pivot := models.UpskillProgramCourse{
				ProgramID:   program.ID,
				CourseID:    pc.CourseID,
				CourseOrder: order,
				IsRequired:  pc.IsRequired,
			}
			if err := tx.Create(&pivot).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Reload with courses
	return s.GetByID(program.ID)
}

// ---------- Get ----------
func (s *UpskillProgramService) GetByID(id uuid.UUID) (*models.UpskillProgram, error) {
	var p models.UpskillProgram
	err := s.DB.
		Preload("Courses", func(db *gorm.DB) *gorm.DB {
			return db.Order("course_order ASC")
		}).
		Preload("Courses.Course").
		First(&p, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrProgramNotFound
	}
	return &p, err
}

// ---------- List ----------
func (s *UpskillProgramService) List(q dto.ListUpskillProgramsQuery) ([]models.UpskillProgram, int64, error) {
	var (
		programs []models.UpskillProgram
		total    int64
	)
	query := s.DB.Model(&models.UpskillProgram{})

	if q.Search != "" {
		s := "%" + strings.ToLower(q.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", s, s)
	}
	if q.IsActive == "true" {
		query = query.Where("is_active = ?", true)
	} else if q.IsActive == "false" {
		query = query.Where("is_active = ?", false)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	allowedSort := map[string]bool{
		"created_at": true, "updated_at": true, "name": true, "is_active": true,
	}
	sortBy := q.SortBy
	if !allowedSort[sortBy] {
		sortBy = "created_at"
	}
	order := strings.ToLower(q.SortOrder)
	if order != "asc" && order != "desc" {
		order = "desc"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, order))

	offset := (q.Page - 1) * q.Limit
	if err := query.
		Preload("Courses", func(db *gorm.DB) *gorm.DB {
			return db.Order("course_order ASC")
		}).
		Preload("Courses.Course").
		Offset(offset).Limit(q.Limit).
		Find(&programs).Error; err != nil {
		return nil, 0, err
	}
	return programs, total, nil
}

func (s *UpskillProgramService) GetBySlug(slug string) (*models.UpskillProgram, error) {
	var p models.UpskillProgram
	err := s.DB.
		Preload("Courses", func(db *gorm.DB) *gorm.DB {
			return db.Order("course_order ASC")
		}).
		Preload("Courses.Course").
		Where("slug = ?", slug).
		First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrProgramNotFound
	}
	return &p, err
}

// ---------- Update ----------
func (s *UpskillProgramService) Update(id uuid.UUID, req dto.UpdateUpskillProgramRequest) (*models.UpskillProgram, error) {
	program, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil && *req.Name != program.Name {
		var count int64
		if err := s.DB.Model(&models.UpskillProgram{}).
			Where("name = ? AND id <> ?", *req.Name, id).
			Count(&count).Error; err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, ErrProgramNameTaken
		}
		program.Name = strings.TrimSpace(*req.Name)
	}

	if req.Slug != nil && *req.Slug != program.Slug {
		var count int64
		if err := s.DB.Model(&models.UpskillProgram{}).
			Where("slug = ? AND id <> ?", *req.Slug, id).
			Count(&count).Error; err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, ErrProgramSlugTaken
		}
		program.Slug = *req.Slug
	}

	if req.Description != nil {
		program.Description = *req.Description
	}
	if req.IsActive != nil {
		program.IsActive = *req.IsActive
	}

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(program).Error; err != nil {
			return err
		}

		// Replace course associations if provided
		if req.Courses != nil {
			// Delete all existing pivots for this program
			if err := tx.Where("program_id = ?", program.ID).Delete(&models.UpskillProgramCourse{}).Error; err != nil {
				return err
			}

			// Validate the new course IDs
			seen := map[uuid.UUID]bool{}
			ids := make([]uuid.UUID, 0, len(*req.Courses))
			for _, pc := range *req.Courses {
				if seen[pc.CourseID] {
					return ErrProgramCourseNoDup
				}
				seen[pc.CourseID] = true
				ids = append(ids, pc.CourseID)
			}

			if len(ids) > 0 {
				var found int64
				if err := tx.Model(&models.Course{}).Where("id IN ?", ids).Count(&found).Error; err != nil {
					return err
				}
				if int(found) != len(ids) {
					return fmt.Errorf("one or more courses were not found")
				}
			}

			// Insert new pivots
			for i, pc := range *req.Courses {
				order := pc.CourseOrder
				if order == 0 {
					order = i + 1
				}
				pivot := models.UpskillProgramCourse{
					ProgramID:   program.ID,
					CourseID:    pc.CourseID,
					CourseOrder: order,
					IsRequired:  pc.IsRequired,
				}
				if err := tx.Create(&pivot).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

// ---------- Delete ----------
func (s *UpskillProgramService) Delete(id uuid.UUID) error {
	_, err := s.GetByID(id)
	if err != nil {
		return err
	}
	// Cascades to upskill_program_courses
	return s.DB.Delete(&models.UpskillProgram{}, "id = ?", id).Error
}

// ---------- Toggle active ----------
func (s *UpskillProgramService) ToggleActive(id uuid.UUID, active bool) (*models.UpskillProgram, error) {
	program, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	program.IsActive = active
	if err := s.DB.Save(program).Error; err != nil {
		return nil, err
	}
	return program, nil
}

// ---------- Add a single course to a program ----------
func (s *UpskillProgramService) AddCourse(programID uuid.UUID, req dto.ProgramCourseInput) (*models.UpskillProgram, error) {
	// Ensure program exists
	if _, err := s.GetByID(programID); err != nil {
		return nil, err
	}

	// Check for duplicate
	var count int64
	if err := s.DB.Model(&models.UpskillProgramCourse{}).
		Where("program_id = ? AND course_id = ?", programID, req.CourseID).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrProgramCourseNoDup
	}

	// Ensure course exists
	if err := s.DB.Model(&models.Course{}).
		Where("id = ?", req.CourseID).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("course not found")
	}

	order := req.CourseOrder
	if order == 0 {
		// Append at end
		var maxOrder int
		s.DB.Model(&models.UpskillProgramCourse{}).
			Where("program_id = ?", programID).
			Select("COALESCE(MAX(course_order), 0)").
			Scan(&maxOrder)
		order = maxOrder + 1
	}

	pivot := models.UpskillProgramCourse{
		ProgramID:   programID,
		CourseID:    req.CourseID,
		CourseOrder: order,
		IsRequired:  req.IsRequired,
	}
	if err := s.DB.Create(&pivot).Error; err != nil {
		return nil, err
	}
	return s.GetByID(programID)
}

// ---------- Remove a single course from a program ----------
func (s *UpskillProgramService) RemoveCourse(programID uuid.UUID, courseID uuid.UUID) (*models.UpskillProgram, error) {
	if _, err := s.GetByID(programID); err != nil {
		return nil, err
	}

	var count int64
	if err := s.DB.Model(&models.UpskillProgramCourse{}).
		Where("program_id = ? AND course_id = ?", programID, courseID).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("course not found in program")
	}

	if err := s.DB.
		Where("program_id = ? AND course_id = ?", programID, courseID).
		Delete(&models.UpskillProgramCourse{}).Error; err != nil {
		return nil, err
	}

	return s.GetByID(programID)
}

// ---------- Reorder courses ----------
func (s *UpskillProgramService) ReorderCourses(programID uuid.UUID, orderedCourseIDs []uuid.UUID) (*models.UpskillProgram, error) {
	if _, err := s.GetByID(programID); err != nil {
		return nil, err
	}

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		for i, courseID := range orderedCourseIDs {
			// Update course_order (1-based)
			if err := tx.Model(&models.UpskillProgramCourse{}).
				Where("program_id = ? AND course_id = ?", programID, courseID).
				Update("course_order", i+1).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.GetByID(programID)
}

// Keep import used
var _ = time.Now