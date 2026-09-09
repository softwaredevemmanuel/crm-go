// controllers/objective_question_controller.go
package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"crm-go/dto"
	"crm-go/services/objective_question"
)

type ObjectiveQuestionController struct {
	questionService *services.ObjectiveQuestionService
}

func NewObjectiveQuestionController(questionService *services.ObjectiveQuestionService) *ObjectiveQuestionController {
	return &ObjectiveQuestionController{
		questionService: questionService,
	}
}

// CreateQuestion creates a new objective question
// @Summary Create objective question
// @Description Create a new objective question with options
// @Tags Objective Questions
// @Accept json
// @Produce json
// @Param request body dto.CreateObjectiveQuestionRequest true "Question data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/objective-questions [post]
func (c *ObjectiveQuestionController) CreateQuestion(ctx *gin.Context) {
	var req dto.CreateObjectiveQuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	question, err := c.questionService.CreateQuestion(&req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Question created successfully",
		"data":    question,
	})
}

// GetQuestionByID retrieves a question by ID
// @Summary Get question by ID
// @Description Get a specific objective question by ID
// @Tags Objective Questions
// @Accept json
// @Produce json
// @Param id path string true "Question ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/objective-questions/{id} [get]
func (c *ObjectiveQuestionController) GetQuestionByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Question ID is required",
		})
		return
	}

	question, err := c.questionService.GetQuestionByID(id)
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
		"message": "Question retrieved successfully",
		"data":    question,
	})
}

// GetQuestions retrieves questions with filters
// @Summary Get questions
// @Description Get a list of objective questions with filters
// @Tags Objective Questions
// @Accept json
// @Produce json
// @Param subject_id query string false "Subject ID"
// @Param scheme_of_work_id query string false "Scheme of Work ID"
// @Param module_id query string false "Module ID"
// @Param topic_id query string false "Topic ID"
// @Param lesson_id query string false "Lesson ID"
// @Param question_type query string false "Question Type"
// @Param difficulty_level query string false "Difficulty Level"
// @Param status query string false "Status"
// @Param search query string false "Search"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/objective-questions [get]
func (c *ObjectiveQuestionController) GetQuestions(ctx *gin.Context) {
	var params dto.ObjectiveQuestionQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := c.questionService.GetQuestions(&params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Questions retrieved successfully",
		"data":    response,
	})
}

// GetQuestionsByLesson retrieves questions for a specific lesson
// @Summary Get questions by lesson
// @Description Get all objective questions for a specific lesson
// @Tags Objective Questions
// @Accept json
// @Produce json
// @Param lesson_id path string true "Lesson ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/objective-questions/lesson/{lesson_id} [get]
func (c *ObjectiveQuestionController) GetQuestionsByLesson(ctx *gin.Context) {
	lessonID := ctx.Param("lesson_id")
	if lessonID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Lesson ID is required",
		})
		return
	}

	questions, err := c.questionService.GetQuestionsByLesson(lessonID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Questions retrieved successfully",
		"data":    questions,
	})
}

// GetQuestionsByTopic retrieves questions for a specific topic
// @Summary Get questions by topic
// @Description Get all objective questions for a specific topic
// @Tags Objective Questions
// @Accept json
// @Produce json
// @Param topic_id path string true "Topic ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/objective-questions/topic/{topic_id} [get]
func (c *ObjectiveQuestionController) GetQuestionsByTopic(ctx *gin.Context) {
	topicID := ctx.Param("topic_id")
	if topicID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Topic ID is required",
		})
		return
	}

	questions, err := c.questionService.GetQuestionsByTopic(topicID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Questions retrieved successfully",
		"data":    questions,
	})
}

// GetQuestionsByModule retrieves questions for a specific module
// @Summary Get questions by module
// @Description Get all objective questions for a specific module
// @Tags Objective Questions
// @Accept json
// @Produce json
// @Param module_id path string true "Module ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/objective-questions/module/{module_id} [get]
func (c *ObjectiveQuestionController) GetQuestionsByModule(ctx *gin.Context) {
	moduleID := ctx.Param("module_id")
	if moduleID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Module ID is required",
		})
		return
	}

	questions, err := c.questionService.GetQuestionsByModule(moduleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Questions retrieved successfully",
		"data":    questions,
	})
}

// UpdateQuestion updates an existing question
// @Summary Update question
// @Description Update an existing objective question
// @Tags Objective Questions
// @Accept json
// @Produce json
// @Param id path string true "Question ID"
// @Param request body dto.UpdateObjectiveQuestionRequest true "Update data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/objective-questions/{id} [put]
func (c *ObjectiveQuestionController) UpdateQuestion(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Question ID is required",
		})
		return
	}

	var req dto.UpdateObjectiveQuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	question, err := c.questionService.UpdateQuestion(id, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Question updated successfully",
		"data":    question,
	})
}

// DeleteQuestion deletes a question
// @Summary Delete question
// @Description Soft delete an objective question
// @Tags Objective Questions
// @Accept json
// @Produce json
// @Param id path string true "Question ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/objective-questions/{id} [delete]
func (c *ObjectiveQuestionController) DeleteQuestion(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Question ID is required",
		})
		return
	}

	if err := c.questionService.DeleteQuestion(id); err != nil {
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
		"message": "Question deleted successfully",
	})
}