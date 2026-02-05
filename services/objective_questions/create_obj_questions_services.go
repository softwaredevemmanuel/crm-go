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

type ObjectiveQuestionService struct {
    db *gorm.DB
}

func NewObjectiveQuestionService(db *gorm.DB) *ObjectiveQuestionService {
    return &ObjectiveQuestionService{db: db}
}

// ValidateQuestionInput validates the question input
func (s *ObjectiveQuestionService) validateQuestionInput(req models.ObjectiveQuestionInput) error {
    // Validate course exists
    var course models.Course
    if err := s.db.First(&course, "id = ?", req.CourseID).Error; err != nil {
        return errors.New("course not found")
    }
    
    // Validate creator exists
    var creator models.User
    if err := s.db.First(&creator, "id = ?", req.CreatedBy).Error; err != nil {
        return errors.New("creator not found")
    }
    
    // Validate chapter if provided
    if req.ChapterID != nil && *req.ChapterID != uuid.Nil {
        var chapter models.Chapter
        if err := s.db.First(&chapter, "id = ?", req.ChapterID).Error; err != nil {
            return errors.New("chapter not found")
        }
        if chapter.CourseID != req.CourseID {
            return errors.New("chapter does not belong to this course")
        }
    }
    
    // Validate topic if provided
    if req.TopicID != nil && *req.TopicID != uuid.Nil {
        var topic models.Topic
        if err := s.db.First(&topic, "id = ?", req.TopicID).Error; err != nil {
            return errors.New("topic not found")
        }
        if topic.CourseID != req.CourseID {
            return errors.New("topic does not belong to this course")
        }
        
        // If chapter is provided, ensure topic belongs to chapter
        if req.ChapterID != nil && *req.ChapterID != uuid.Nil && topic.ChapterID != *req.ChapterID {
            return errors.New("topic does not belong to the specified chapter")
        }
    }
    
    // Validate question type specific rules
    switch req.QuestionType {
    case "multiple_choice", "multiple_response":
        if len(req.Options) < 2 {
            return errors.New("multiple choice questions require at least 2 options")
        }
        if len(req.Options) > 10 {
            return errors.New("cannot have more than 10 options")
        }
        
        // Check for at least one correct answer
        hasCorrect := false
        for _, opt := range req.Options {
            if opt.IsCorrect {
                hasCorrect = true
                break
            }
        }
        if !hasCorrect {
            return errors.New("at least one option must be marked as correct")
        }
        
        // For multiple choice (single answer), ensure only one correct
        if req.QuestionType == "multiple_choice" {
            correctCount := 0
            for _, opt := range req.Options {
                if opt.IsCorrect {
                    correctCount++
                }
            }
            if correctCount > 1 {
                return errors.New("multiple choice questions can only have one correct answer")
            }
        }
        
    case "true_false":
        // True/false questions don't need options
        if len(req.Options) > 0 {
            return errors.New("true/false questions should not have options")
        }
        
    case "matching", "ordering":
        if len(req.Options) < 3 {
            return errors.New("matching/ordering questions require at least 3 options")
        }
    }
    
    return nil
}

// CreateObjectiveQuestion - main create function
func (s *ObjectiveQuestionService) CreateObjectiveQuestion(req models.ObjectiveQuestionInput) (*models.ObjectiveQuestionResponse, error) {
    // Validate input
    if err := s.validateQuestionInput(req); err != nil {
        return nil, err
    }
    
    // Set defaults if not provided
    if req.QuestionType == "" {
        req.QuestionType = "multiple_choice"
    }
    if req.DifficultyLevel == "" {
        req.DifficultyLevel = "medium"
    }
    if req.Points == 0 {
        req.Points = 1
    }
    
    // Start transaction
    tx := s.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // Create question
    question := models.ObjectiveQuestion{
        ID:                uuid.New(),
        QuestionText:      strings.TrimSpace(req.QuestionText),
        QuestionType:      req.QuestionType,
        DifficultyLevel:   req.DifficultyLevel,
        Points:            req.Points,
        ImageURL:          strings.TrimSpace(req.ImageURL),
        VideoURL:          strings.TrimSpace(req.VideoURL),
        CourseID:          req.CourseID,
        ChapterID:         req.ChapterID,
        TopicID:           req.TopicID,
        AnswerExplanation: strings.TrimSpace(req.AnswerExplanation),
        SolutionSteps:     strings.TrimSpace(req.SolutionSteps),
        Hint:              strings.TrimSpace(req.Hint),
        CreatedBy:         req.CreatedBy,
        CreatedAt:         time.Now(),
        UpdatedAt:         time.Now(),
        IsApproved:        req.IsApproved,
    }
    
    // Save question
    if err := tx.Create(&question).Error; err != nil {
        tx.Rollback()
        return nil, errors.New("failed to save question: " + err.Error())
    }
    
    // Create options if provided
    var options []models.QuestionOption
    if len(req.Options) > 0 {
        for i, optInput := range req.Options {
            option := models.QuestionOption{
                ID:         uuid.New(),
                QuestionID: question.ID,
                OptionText: strings.TrimSpace(optInput.OptionText),
                IsCorrect:  optInput.IsCorrect,
                Explanation: strings.TrimSpace(optInput.Explanation),
                SortOrder:      optInput.SortOrder,
                CreatedAt:  time.Now(),
                UpdatedAt:  time.Now(),
            }

            // If sort order not specified, use index
            if option.SortOrder == 0 {
                option.SortOrder = i + 1
            }
            
            if err := tx.Create(&option).Error; err != nil {
                tx.Rollback()
                return nil, fmt.Errorf("failed to save option %d: %v", i+1, err)
            }
            
            options = append(options, option)
        }
    }
    
    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        return nil, errors.New("failed to save question: " + err.Error())
    }
    
    // Convert to response
    return s.questionToResponse(&question, options, true), nil
}

// CreateObjectiveQuestionWithTx - for use with transactions
func (s *ObjectiveQuestionService) CreateObjectiveQuestionWithTx(tx *gorm.DB, req models.ObjectiveQuestionInput) (*models.ObjectiveQuestionResponse, error) {
    // Validate input
    if err := s.validateQuestionInput(req); err != nil {
        return nil, err
    }
    
    // Set defaults
    if req.QuestionType == "" { req.QuestionType = "multiple_choice" }
    if req.DifficultyLevel == "" { req.DifficultyLevel = "medium" }
    if req.Points == 0 { req.Points = 1 }
    
    // Create question
    question := models.ObjectiveQuestion{
        ID:                uuid.New(),
        QuestionText:      strings.TrimSpace(req.QuestionText),
        QuestionType:      req.QuestionType,
        DifficultyLevel:   req.DifficultyLevel,
        Points:            req.Points,
        ImageURL:          strings.TrimSpace(req.ImageURL),
        VideoURL:          strings.TrimSpace(req.VideoURL),
        CourseID:          req.CourseID,
        ChapterID:         req.ChapterID,
        TopicID:           req.TopicID,
        AnswerExplanation: strings.TrimSpace(req.AnswerExplanation),
        SolutionSteps:     strings.TrimSpace(req.SolutionSteps),
        Hint:              strings.TrimSpace(req.Hint),
        CreatedBy:         req.CreatedBy,
        CreatedAt:         time.Now(),
        UpdatedAt:         time.Now(),
        IsApproved:        req.IsApproved,
    }
    
    // Save question
    if err := tx.Create(&question).Error; err != nil {
        return nil, errors.New("failed to save question: " + err.Error())
    }
    
    // Create options
    var options []models.QuestionOption
    for i, optInput := range req.Options {
        option := models.QuestionOption{
            ID:         uuid.New(),
            QuestionID: question.ID,
            OptionText: strings.TrimSpace(optInput.OptionText),
            IsCorrect:  optInput.IsCorrect,
            Explanation: strings.TrimSpace(optInput.Explanation),
            SortOrder:      optInput.SortOrder,
            CreatedAt:  time.Now(),
            UpdatedAt:  time.Now(),
        }

        if option.SortOrder == 0 {
            option.SortOrder = i + 1
        }
        
        if err := tx.Create(&option).Error; err != nil {
            return nil, fmt.Errorf("failed to save option: %v", err)
        }
        
        options = append(options, option)
    }
    
    return s.questionToResponse(&question, options, false), nil
}

// CreateBulkQuestions - create multiple questions at once
func (s *ObjectiveQuestionService) CreateBulkQuestions(requests []models.ObjectiveQuestionInput) ([]models.ObjectiveQuestionResponse, error) {
    tx := s.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    responses := make([]models.ObjectiveQuestionResponse, 0, len(requests))
    
    for i, req := range requests {
        response, err := s.CreateObjectiveQuestionWithTx(tx, req)
        if err != nil {
            tx.Rollback()
            return nil, fmt.Errorf("failed to create question %d: %v", i+1, err)
        }
        responses = append(responses, *response)
    }
    
    if err := tx.Commit().Error; err != nil {
        return nil, errors.New("failed to save questions: " + err.Error())
    }
    
    return responses, nil
}

// Helper to convert question to response
func (s *ObjectiveQuestionService) questionToResponse(question *models.ObjectiveQuestion, options []models.QuestionOption, withDetails bool) *models.ObjectiveQuestionResponse {
    response := &models.ObjectiveQuestionResponse{
        ID:                question.ID,
        QuestionText:      question.QuestionText,
        QuestionType:      question.QuestionType,
        DifficultyLevel:   question.DifficultyLevel,
        Points:            question.Points,
        ImageURL:          question.ImageURL,
        VideoURL:          question.VideoURL,
        CourseID:          question.CourseID,
        ChapterID:         question.ChapterID,
        TopicID:           question.TopicID,
        AnswerExplanation: question.AnswerExplanation,
        SolutionSteps:     question.SolutionSteps,
        Hint:              question.Hint,
        CreatedBy:         question.CreatedBy,
        CreatedAt:         question.CreatedAt,
        UpdatedAt:         question.UpdatedAt,
        IsApproved:        question.IsApproved,
    }
    
    // Add options
    for _, option := range options {
        response.Options = append(response.Options, models.QuestionOptionResponse{
            ID:          option.ID,
            QuestionID:  option.QuestionID,
            OptionText:  option.OptionText,
            IsCorrect:   option.IsCorrect,
            Explanation: option.Explanation,
            SortOrder:       option.SortOrder,
        })
    }
    
    // Add details if requested
    if withDetails {
        // Get course name
        var course models.Course
        if err := s.db.First(&course, "id = ?", question.CourseID).Error; err == nil {
            response.CourseName = course.Title
        }
        
        // Get chapter name if exists
        if question.ChapterID != nil && *question.ChapterID != uuid.Nil {
            var chapter models.Chapter
            if err := s.db.First(&chapter, "id = ?", question.ChapterID).Error; err == nil {
                response.ChapterName = chapter.Title
            }
        }
        
        // Get topic name if exists
        if question.TopicID != nil && *question.TopicID != uuid.Nil {
            var topic models.Topic
            if err := s.db.First(&topic, "id = ?", question.TopicID).Error; err == nil {
                response.TopicName = topic.Title
            }
        }
        
        // Get creator name
        var creator models.User
        if err := s.db.First(&creator, "id = ?", question.CreatedBy).Error; err == nil {
            response.CreatorName = fmt.Sprintf("%s %s", creator.FirstName, creator.LastName)
        }
        
        // Get statistics (optional)
        var stats struct {
            TotalAttempts   int64
            CorrectAttempts int64
        }
        
        s.db.Raw(`
            SELECT 
                COUNT(*) as total_attempts,
                SUM(CASE WHEN is_correct THEN 1 ELSE 0 END) as correct_attempts
            FROM question_attempts
            WHERE question_id = ?
        `, question.ID).Scan(&stats)
        
        response.TotalAttempts = int(stats.TotalAttempts)
        response.CorrectAttempts = int(stats.CorrectAttempts)
        
        if stats.TotalAttempts > 0 {
            response.SuccessRate = float64(stats.CorrectAttempts) / float64(stats.TotalAttempts) * 100
        }
    }
    
    return response
}

// GenerateOptionsForTrueFalse - auto-generate options for true/false questions
func (s *ObjectiveQuestionService) generateTrueFalseOptions(questionID uuid.UUID) []models.QuestionOption {
    return []models.QuestionOption{
        {
            ID:         uuid.New(),
            QuestionID: questionID,
            OptionText: "True",
            IsCorrect:  true, // Assuming answer is True
            SortOrder:      1,
            CreatedAt:  time.Now(),
            UpdatedAt:  time.Now(),
        },
        {
            ID:         uuid.New(),
            QuestionID: questionID,
            OptionText: "False",
            IsCorrect:  false,
            SortOrder:      2,
            CreatedAt:  time.Now(),
            UpdatedAt:  time.Now(),
        },
    }
}

// ValidateQuestionAnswer - validate answer for a question
func (s *ObjectiveQuestionService) ValidateQuestionAnswer(questionID uuid.UUID, answer interface{}) (bool, error) {
    var question models.ObjectiveQuestion
    if err := s.db.First(&question, "id = ?", questionID).Error; err != nil {
        return false, errors.New("question not found")
    }
    
    switch question.QuestionType {
    case "multiple_choice":
        // Answer should be option ID or index
        return s.validateMultipleChoiceAnswer(questionID, answer)
    case "true_false":
        // Answer should be boolean
        return s.validateTrueFalseAnswer(questionID, answer)
    case "multiple_response":
        // Answer should be array of option IDs
        return s.validateMultipleResponseAnswer(questionID, answer)
    default:
        return false, errors.New("question type not supported for auto-validation")
    }
}

func (s *ObjectiveQuestionService) validateMultipleChoiceAnswer(questionID uuid.UUID, answer interface{}) (bool, error) {
    // Implementation for multiple choice validation
    return false, nil
}

func (s *ObjectiveQuestionService) validateTrueFalseAnswer(questionID uuid.UUID, answer interface{}) (bool, error) {
    // Implementation for true/false validation
    return false, nil
}

func (s *ObjectiveQuestionService) validateMultipleResponseAnswer(questionID uuid.UUID, answer interface{}) (bool, error) {
    // Implementation for multiple response validation
    return false, nil
}