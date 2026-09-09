// services/objective_question_service.go
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

type ObjectiveQuestionService struct {
	db *gorm.DB
}

func NewObjectiveQuestionService(db *gorm.DB) *ObjectiveQuestionService {
	return &ObjectiveQuestionService{db: db}
}

// CreateQuestion creates a new objective question with options
func (s *ObjectiveQuestionService) CreateQuestion(req *dto.CreateObjectiveQuestionRequest) (*dto.ObjectiveQuestionResponse, error) {
	// Parse UUIDs
	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		return nil, errors.New("invalid subject ID")
	}

	schemeOfWorkID, err := uuid.Parse(req.SchemeOfWorkID)
	if err != nil {
		return nil, errors.New("invalid scheme of work ID")
	}

	var moduleID, topicID, lessonID uuid.UUID
	if req.ModuleID != "" {
		moduleID, err = uuid.Parse(req.ModuleID)
		if err != nil {
			return nil, errors.New("invalid module ID")
		}
	}

	if req.TopicID != "" {
		topicID, err = uuid.Parse(req.TopicID)
		if err != nil {
			return nil, errors.New("invalid topic ID")
		}
	}

	if req.LessonID != "" {
		lessonID, err = uuid.Parse(req.LessonID)
		if err != nil {
			return nil, errors.New("invalid lesson ID")
		}
	}

	// Verify subject exists
	var subject models.Subject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subject not found")
		}
		return nil, fmt.Errorf("failed to verify subject: %w", err)
	}

	// Verify scheme of work exists
	var scheme models.SchemeOfWork
	if err := s.db.Where("id = ? AND deleted_at IS NULL", schemeOfWorkID).First(&scheme).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("scheme of work not found")
		}
		return nil, fmt.Errorf("failed to verify scheme of work: %w", err)
	}

	// Validate options
	if len(req.Options) < 2 {
		return nil, errors.New("at least 2 options are required")
	}

	// Count correct answers
	correctCount := 0
	for _, opt := range req.Options {
		if opt.IsCorrect {
			correctCount++
		}
	}

	// Validate based on question type
	switch req.QuestionType {
	case "true_false":
		if len(req.Options) != 2 {
			return nil, errors.New("true/false questions must have exactly 2 options")
		}
		if correctCount != 1 {
			return nil, errors.New("true/false questions must have exactly 1 correct answer")
		}
	case "multiple_choice":
		if correctCount != 1 {
			return nil, errors.New("multiple choice questions must have exactly 1 correct answer")
		}
	case "multiple_response":
		if correctCount < 1 {
			return nil, errors.New("multiple response questions must have at least 1 correct answer")
		}
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	// Create question
	question := &models.ObjectiveQuestion{
		ID:                uuid.New(),
		QuestionText:      req.QuestionText,
		QuestionType:      req.QuestionType,
		DifficultyLevel:   req.DifficultyLevel,
		Points:            req.Points,
		ImageURL:          req.ImageURL,
		VideoURL:          req.VideoURL,
		SubjectID:         subjectID,
		SchemeOfWorkID:    schemeOfWorkID,
		ModuleID:          moduleID,
		TopicID:           topicID,
		LessonID:          lessonID,
		AnswerExplanation: req.AnswerExplanation,
		SolutionSteps:     req.SolutionSteps,
		Hint:              req.Hint,
		Status:            status,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.db.Create(question).Error; err != nil {
		return nil, fmt.Errorf("failed to create question: %w", err)
	}

	// Create options
	for i, optReq := range req.Options {
		option := &models.QuestionOption{
			ID:         uuid.New(),
			QuestionID: question.ID,
			OptionText: optReq.OptionText,
			IsCorrect:  optReq.IsCorrect,
			Order:      optReq.Order,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		if option.Order == 0 {
			option.Order = i + 1
		}
		if err := s.db.Create(option).Error; err != nil {
			return nil, fmt.Errorf("failed to create option: %w", err)
		}
	}

	// Load relationships for response
	if err := s.db.Preload("Subject").
		Preload("SchemeOfWork").
		Preload("Module").
		Preload("Topic").
		Preload("Lesson").
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("question_options.order ASC")
		}).
		First(question, question.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load question details: %w", err)
	}

	return s.toQuestionResponse(question), nil
}

// GetQuestionByID retrieves a question by ID
func (s *ObjectiveQuestionService) GetQuestionByID(id string) (*dto.ObjectiveQuestionResponse, error) {
	questionID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid question ID")
	}

	var question models.ObjectiveQuestion
	if err := s.db.Where("id = ? AND deleted_at IS NULL", questionID).
		Preload("Subject").
		Preload("SchemeOfWork").
		Preload("Module").
		Preload("Topic").
		Preload("Lesson").
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("question_options.order ASC")
		}).
		First(&question).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("question not found")
		}
		return nil, fmt.Errorf("failed to fetch question: %w", err)
	}

	return s.toQuestionResponse(&question), nil
}

// GetQuestions retrieves questions with filters
func (s *ObjectiveQuestionService) GetQuestions(params *dto.ObjectiveQuestionQueryParams) (*dto.ObjectiveQuestionListResponse, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}

	// Build query
	query := s.db.Model(&models.ObjectiveQuestion{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.SubjectID != "" {
		subjectID, err := uuid.Parse(params.SubjectID)
		if err == nil {
			query = query.Where("subject_id = ?", subjectID)
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

	if params.QuestionType != "" {
		query = query.Where("question_type = ?", params.QuestionType)
	}

	if params.DifficultyLevel != "" {
		query = query.Where("difficulty_level = ?", params.DifficultyLevel)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("question_text ILIKE ? OR answer_explanation ILIKE ?", search, search)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count questions: %w", err)
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
	var questions []models.ObjectiveQuestion
	if err := query.
		Preload("Subject").
		Preload("SchemeOfWork").
		Preload("Module").
		Preload("Topic").
		Preload("Lesson").
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("question_options.order ASC")
		}).
		Find(&questions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch questions: %w", err)
	}

	// Convert to response
	responses := make([]dto.ObjectiveQuestionResponse, len(questions))
	for i, question := range questions {
		responses[i] = *s.toQuestionResponse(&question)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.ObjectiveQuestionListResponse{
		Questions:  responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetQuestionsByLesson retrieves questions for a specific lesson
func (s *ObjectiveQuestionService) GetQuestionsByLesson(lessonID string) ([]dto.ObjectiveQuestionResponse, error) {
	lID, err := uuid.Parse(lessonID)
	if err != nil {
		return nil, errors.New("invalid lesson ID")
	}

	var questions []models.ObjectiveQuestion
	if err := s.db.Where("lesson_id = ? AND status = ? AND deleted_at IS NULL", lID, "active").
		Preload("Subject").
		Preload("SchemeOfWork").
		Preload("Module").
		Preload("Topic").
		Preload("Lesson").
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("question_options.order ASC")
		}).
		Find(&questions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch questions: %w", err)
	}

	responses := make([]dto.ObjectiveQuestionResponse, len(questions))
	for i, question := range questions {
		responses[i] = *s.toQuestionResponse(&question)
	}

	return responses, nil
}

// GetQuestionsByTopic retrieves questions for a specific topic
func (s *ObjectiveQuestionService) GetQuestionsByTopic(topicID string) ([]dto.ObjectiveQuestionResponse, error) {
	tID, err := uuid.Parse(topicID)
	if err != nil {
		return nil, errors.New("invalid topic ID")
	}

	var questions []models.ObjectiveQuestion
	if err := s.db.Where("topic_id = ? AND status = ? AND deleted_at IS NULL", tID, "active").
		Preload("Subject").
		Preload("SchemeOfWork").
		Preload("Module").
		Preload("Topic").
		Preload("Lesson").
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("question_options.order ASC")
		}).
		Find(&questions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch questions: %w", err)
	}

	responses := make([]dto.ObjectiveQuestionResponse, len(questions))
	for i, question := range questions {
		responses[i] = *s.toQuestionResponse(&question)
	}

	return responses, nil
}

// GetQuestionsByModule retrieves questions for a specific module
func (s *ObjectiveQuestionService) GetQuestionsByModule(moduleID string) ([]dto.ObjectiveQuestionResponse, error) {
	mID, err := uuid.Parse(moduleID)
	if err != nil {
		return nil, errors.New("invalid module ID")
	}

	var questions []models.ObjectiveQuestion
	if err := s.db.Where("module_id = ? AND status = ? AND deleted_at IS NULL", mID, "active").
		Preload("Subject").
		Preload("SchemeOfWork").
		Preload("Module").
		Preload("Topic").
		Preload("Lesson").
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("question_options.order ASC")
		}).
		Find(&questions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch questions: %w", err)
	}

	responses := make([]dto.ObjectiveQuestionResponse, len(questions))
	for i, question := range questions {
		responses[i] = *s.toQuestionResponse(&question)
	}

	return responses, nil
}

// UpdateQuestion updates an existing question
func (s *ObjectiveQuestionService) UpdateQuestion(id string, req *dto.UpdateObjectiveQuestionRequest) (*dto.ObjectiveQuestionResponse, error) {
	questionID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid question ID")
	}

	var question models.ObjectiveQuestion
	if err := s.db.Where("id = ? AND deleted_at IS NULL", questionID).First(&question).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("question not found")
		}
		return nil, fmt.Errorf("failed to fetch question: %w", err)
	}

	// Update fields
	if req.QuestionText != "" {
		question.QuestionText = req.QuestionText
	}
	if req.QuestionType != "" {
		question.QuestionType = req.QuestionType
	}
	if req.DifficultyLevel != "" {
		question.DifficultyLevel = req.DifficultyLevel
	}
	if req.Points > 0 {
		question.Points = req.Points
	}
	if req.ImageURL != "" {
		question.ImageURL = req.ImageURL
	}
	if req.VideoURL != "" {
		question.VideoURL = req.VideoURL
	}
	if req.AnswerExplanation != "" {
		question.AnswerExplanation = req.AnswerExplanation
	}
	if req.SolutionSteps != "" {
		question.SolutionSteps = req.SolutionSteps
	}
	if req.Hint != "" {
		question.Hint = req.Hint
	}
	if req.Status != "" {
		question.Status = req.Status
	}

	question.UpdatedAt = time.Now()

	if err := s.db.Save(&question).Error; err != nil {
		return nil, fmt.Errorf("failed to update question: %w", err)
	}

	// Update options if provided
	if req.Options != nil {
		// Delete existing options
		if err := s.db.Where("question_id = ?", question.ID).Delete(&models.QuestionOption{}).Error; err != nil {
			return nil, fmt.Errorf("failed to delete old options: %w", err)
		}

		// Create new options
		for i, optReq := range req.Options {
			option := &models.QuestionOption{
				ID:         uuid.New(),
				QuestionID: question.ID,
				OptionText: optReq.OptionText,
				IsCorrect:  optReq.IsCorrect,
				Order:      optReq.Order,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
			if option.Order == 0 {
				option.Order = i + 1
			}
			if err := s.db.Create(option).Error; err != nil {
				return nil, fmt.Errorf("failed to create option: %w", err)
			}
		}
	}

	// Load relationships for response
	if err := s.db.Preload("Subject").
		Preload("SchemeOfWork").
		Preload("Module").
		Preload("Topic").
		Preload("Lesson").
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("question_options.order ASC")
		}).
		First(&question, question.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load question details: %w", err)
	}

	return s.toQuestionResponse(&question), nil
}

// DeleteQuestion soft deletes a question
func (s *ObjectiveQuestionService) DeleteQuestion(id string) error {
	questionID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid question ID")
	}

	var question models.ObjectiveQuestion
	if err := s.db.Where("id = ? AND deleted_at IS NULL", questionID).First(&question).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("question not found")
		}
		return fmt.Errorf("failed to fetch question: %w", err)
	}

	if err := s.db.Delete(&question).Error; err != nil {
		return fmt.Errorf("failed to delete question: %w", err)
	}

	return nil
}

// toQuestionResponse converts model to response DTO
func (s *ObjectiveQuestionService) toQuestionResponse(question *models.ObjectiveQuestion) *dto.ObjectiveQuestionResponse {
	response := &dto.ObjectiveQuestionResponse{
		ID:                question.ID.String(),
		QuestionText:      question.QuestionText,
		QuestionType:      question.QuestionType,
		DifficultyLevel:   question.DifficultyLevel,
		Points:            question.Points,
		ImageURL:          question.ImageURL,
		VideoURL:          question.VideoURL,
		SubjectID:         question.SubjectID.String(),
		SchemeOfWorkID:    question.SchemeOfWorkID.String(),
		ModuleID:          question.ModuleID.String(),
		TopicID:           question.TopicID.String(),
		LessonID:          question.LessonID.String(),
		AnswerExplanation: question.AnswerExplanation,
		SolutionSteps:     question.SolutionSteps,
		Hint:              question.Hint,
		Status:            question.Status,
		CreatedAt:         question.CreatedAt,
		UpdatedAt:         question.UpdatedAt,
	}

	// Add options if they exist
	if len(question.Options) > 0 {
		options := make([]dto.OptionResponse, len(question.Options))
		for i, opt := range question.Options {
			options[i] = dto.OptionResponse{
				ID:         opt.ID.String(),
				OptionText: opt.OptionText,
				IsCorrect:  opt.IsCorrect,
				Order:      opt.Order,
			}
		}
		response.Options = options
	}

	return response
}