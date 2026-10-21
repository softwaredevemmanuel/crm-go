package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/upskill_program"
)

type UpskillProgramController struct {
	Service *services.UpskillProgramService
}

func NewUpskillProgramController(service *services.UpskillProgramService) *UpskillProgramController {
	return &UpskillProgramController{Service: service}
}

// CreateProgram godoc
//
//	@Summary		Create a new upskill program
//	@Description	Creates an upskill program with an optional initial set of courses.
//	@Tags			Upskill Programs
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.CreateUpskillProgramRequest	true	"Program payload"
//	@Success		201		{object}	dto.SingleUpskillProgramResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		409		{object}	dto.ErrorResponse	"Name already exists"
//	@Security		BearerAuth
//	@Router			/api/upskill-programs [post]
func (ctrl *UpskillProgramController) CreateProgram(c *gin.Context) {
	var req dto.CreateUpskillProgramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:  "validation failed",
			Fields: map[string]string{"body": err.Error()},
		})
		return
	}
	p, err := ctrl.Service.Create(req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrProgramNameTaken):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, dto.SingleUpskillProgramResponse{Data: dto.FromUpskillProgramModel(p, true)})
}

// ListPrograms godoc
//
//	@Summary		List upskill programs (paginated)
//	@Tags			Upskill Programs
//	@Produce		json
//	@Param			page		query		int		false	"Page number"
//	@Param			limit		query		int		false	"Items per page"
//	@Param			search		query		string	false	"Search by name or description"
//	@Param			is_active	query		string	false	"Filter by active status"	Enums(true, false)
//	@Param			sort_by		query		string	false	"Sort field"		Enums(created_at, updated_at, name, is_active)
//	@Param			sort_order	query		string	false	"Sort direction"	Enums(asc, desc)
//	@Success		200			{object}	dto.PaginatedUpskillProgramsResponse
//	@Security		BearerAuth
//	@Router			/api/upskill-programs [get]
func (ctrl *UpskillProgramController) ListPrograms(c *gin.Context) {
	var q dto.ListUpskillProgramsQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid query", Message: err.Error()})
		return
	}
	programs, total, err := ctrl.Service.List(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	items := make([]dto.UpskillProgramResponse, 0, len(programs))
	for i := range programs {
		items = append(items, dto.FromUpskillProgramModel(&programs[i], true))
	}
	totalPages := int((total + int64(q.Limit) - 1) / int64(q.Limit))

	c.JSON(http.StatusOK, dto.PaginatedUpskillProgramsResponse{
		Data:       items,
		Total:      total,
		Page:       q.Page,
		Limit:      q.Limit,
		TotalPages: totalPages,
	})
}

// GetProgram godoc
//
//	@Summary		Get an upskill program by ID
//	@Tags			Upskill Programs
//	@Produce		json
//	@Param			id	path		string	true	"Program UUID"
//	@Success		200	{object}	dto.SingleUpskillProgramResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/upskill-programs/{id} [get]
func (ctrl *UpskillProgramController) GetProgram(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid program id"})
		return
	}
	p, err := ctrl.Service.GetByID(id)
	if err != nil {
		if errors.Is(err, services.ErrProgramNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.SingleUpskillProgramResponse{Data: dto.FromUpskillProgramModel(p, true)})
}

// GetProgramBySlug godoc
//
//	@Summary		Get an upskill program by slug
//	@Tags			Upskill Programs
//	@Produce		json
//	@Param			slug	path		string	true	"Program slug"
//	@Success		200		{object}	dto.SingleUpskillProgramResponse
//	@Failure		404		{object}	dto.ErrorResponse
//	@Router			/api/upskill-programs/slug/{slug} [get]
func (ctrl *UpskillProgramController) GetProgramBySlug(c *gin.Context) {
	slug := c.Param("slug")
	p, err := ctrl.Service.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, services.ErrProgramNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.SingleUpskillProgramResponse{Data: dto.FromUpskillProgramModel(p, true)})
}

// UpdateProgram godoc
//
//	@Summary		Update an upskill program (partial)
//	@Tags			Upskill Programs
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string							true	"Program UUID"
//	@Param			payload	body		dto.UpdateUpskillProgramRequest	true	"Fields to update"
//	@Success		200		{object}	dto.SingleUpskillProgramResponse
//	@Security		BearerAuth
//	@Router			/api/upskill-programs/{id} [patch]
func (ctrl *UpskillProgramController) UpdateProgram(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid program id"})
		return
	}
	var req dto.UpdateUpskillProgramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "validation failed", Message: err.Error()})
		return
	}
	p, err := ctrl.Service.Update(id, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrProgramNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, services.ErrProgramNameTaken), errors.Is(err, services.ErrProgramSlugTaken):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, dto.SingleUpskillProgramResponse{Data: dto.FromUpskillProgramModel(p, true)})
}

// DeleteProgram godoc
//
//	@Summary		Delete an upskill program
//	@Tags			Upskill Programs
//	@Produce		json
//	@Param			id	path	string	true	"Program UUID"
//	@Success		204
//	@Failure		404	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/upskill-programs/{id} [delete]
func (ctrl *UpskillProgramController) DeleteProgram(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid program id"})
		return
	}
	if err := ctrl.Service.Delete(id); err != nil {
		if errors.Is(err, services.ErrProgramNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ToggleProgramActive godoc
//
//	@Summary		Toggle active status
//	@Tags			Upskill Programs
//	@Produce		json
//	@Param			id		path		string	true	"Program UUID"
//	@Param			active	query		bool	true	"Active state"
//	@Success		200		{object}	dto.SingleUpskillProgramResponse
//	@Security		BearerAuth
//	@Router			/api/upskill-programs/{id}/active [patch]
func (ctrl *UpskillProgramController) ToggleProgramActive(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid program id"})
		return
	}
	active := c.Query("active") == "true"
	p, err := ctrl.Service.ToggleActive(id, active)
	if err != nil {
		if errors.Is(err, services.ErrProgramNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.SingleUpskillProgramResponse{Data: dto.FromUpskillProgramModel(p, true)})
}

// AddCourseToProgram godoc
//
//	@Summary		Add a course to a program
//	@Tags			Upskill Programs
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Program UUID"
//	@Param			payload	body		dto.ProgramCourseInput	true	"Course to add"
//	@Success		200		{object}	dto.SingleUpskillProgramResponse
//	@Security		BearerAuth
//	@Router			/api/upskill-programs/{id}/courses [post]
func (ctrl *UpskillProgramController) AddCourseToProgram(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid program id"})
		return
	}
	var req dto.ProgramCourseInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "validation failed", Message: err.Error()})
		return
	}
	p, err := ctrl.Service.AddCourse(id, req)
	if err != nil {
		if errors.Is(err, services.ErrProgramCourseNoDup) {
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.SingleUpskillProgramResponse{Data: dto.FromUpskillProgramModel(p, true)})
}

// RemoveCourseFromProgram godoc
//
//	@Summary		Remove a course from a program
//	@Tags			Upskill Programs
//	@Produce		json
//	@Param			id			path	string	true	"Program UUID"
//	@Param			course_id	path	string	true	"Course UUID"
//	@Success		200			{object}	dto.SingleUpskillProgramResponse
//	@Security		BearerAuth
//	@Router			/api/upskill-programs/{id}/courses/{course_id} [delete]
func (ctrl *UpskillProgramController) RemoveCourseFromProgram(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid program id"})
		return
	}
	courseID, err := uuid.Parse(c.Param("course_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid course id"})
		return
	}
	p, err := ctrl.Service.RemoveCourse(id, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.SingleUpskillProgramResponse{Data: dto.FromUpskillProgramModel(p, true)})
}

// ReorderProgramCourses godoc
//
//	@Summary		Reorder the courses within a program
//	@Description	Provide a JSON body with the ordered list of course IDs.
//	@Tags			Upskill Programs
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"Program UUID"
//	@Param			payload	body		map[string][]string	true	"Ordered course IDs"
//	@Success		200		{object}	dto.SingleUpskillProgramResponse
//	@Security		BearerAuth
//	@Router			/api/upskill-programs/{id}/courses/reorder [post]
func (ctrl *UpskillProgramController) ReorderProgramCourses(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid program id"})
		return
	}

	var body struct {
		CourseIDs []string `json:"course_ids" binding:"required,min=1,dive,uuid4"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "validation failed", Message: err.Error()})
		return
	}

	ids := make([]uuid.UUID, 0, len(body.CourseIDs))
	for _, s := range body.CourseIDs {
		uid, _ := uuid.Parse(s)
		ids = append(ids, uid)
	}

	p, err := ctrl.Service.ReorderCourses(id, ids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.SingleUpskillProgramResponse{Data: dto.FromUpskillProgramModel(p, true)})
}