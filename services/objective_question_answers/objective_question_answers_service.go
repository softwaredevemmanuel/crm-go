// services/student_answer_service.go
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

type StudentAnswerService struct {
	db *gorm.DB
}

func NewStudentAnswerService(db *gorm.DB) *StudentAnswerService {
	return &StudentAnswerService{db: db}
}

// SubmitAnswer submits a student's answer to a question
func (s *StudentAnswerService) SubmitAnswer(studentID uuid.UUID, req *dto.SubmitAnswerRequest) (*dto.SubmitAnswerResponse, error) {
	// Parse question ID
	questionID, err := uuid.Parse(req.QuestionID)
	if err != nil {
		return nil, errors.New("invalid question ID")
	}

	// Fetch the question with options
	var question models.ObjectiveQuestion
	if err := s.db.Where("id = ? AND deleted_at IS NULL", questionID).
		Preload("Options").
		First(&question).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("question not found")
		}
		return nil, fmt.Errorf("failed to fetch question: %w", err)
	}

	// Validate question is active
	if question.Status != "active" {
		return nil, errors.New("question is not active")
	}

	// Check if student has already answered this question
	var existingAnswers []models.ObjectiveQuestionAnswer
	if err := s.db.Where("student_id = ? AND question_id = ? AND deleted_at IS NULL", studentID, questionID).
		Order("attempt_number DESC").
		Find(&existingAnswers).Error; err != nil {
		return nil, fmt.Errorf("failed to check existing answers: %w", err)
	}

	previousAttempts := len(existingAnswers)
	attemptNumber := previousAttempts + 1

	// Validate answer based on question type
	var isCorrect bool
	var score int
	var selectedOptionIDsStr string
	var selectedOptionID uuid.UUID

	switch question.QuestionType {
	case "multiple_choice":
		// Single selection - validate exactly one option
		if req.SelectedOptionID == "" {
			// No option selected - mark as wrong
			isCorrect = false
			score = 0
			selectedOptionID = uuid.Nil
			break
		}
		optionID, err := uuid.Parse(req.SelectedOptionID)
		if err != nil {
			return nil, errors.New("invalid option ID")
		}
		selectedOptionID = optionID

		// Find the selected option
		var selectedOption models.QuestionOption
		for _, opt := range question.Options {
			if opt.ID == optionID {
				selectedOption = opt
				break
			}
		}
		if selectedOption.ID == uuid.Nil {
			// Option not found - mark as wrong
			isCorrect = false
			score = 0
			break
		}

		isCorrect = selectedOption.IsCorrect
		score = question.Points
		if !isCorrect {
			score = 0
		}

	case "true_false":
		// Single selection - validate exactly one option
		if req.SelectedOptionID == "" {
			// No option selected - mark as wrong
			isCorrect = false
			score = 0
			selectedOptionID = uuid.Nil
			break
		}
		optionID, err := uuid.Parse(req.SelectedOptionID)
		if err != nil {
			return nil, errors.New("invalid option ID")
		}
		selectedOptionID = optionID

		var selectedOption models.QuestionOption
		for _, opt := range question.Options {
			if opt.ID == optionID {
				selectedOption = opt
				break
			}
		}
		if selectedOption.ID == uuid.Nil {
			// Option not found - mark as wrong
			isCorrect = false
			score = 0
			break
		}

		isCorrect = selectedOption.IsCorrect
		score = question.Points
		if !isCorrect {
			score = 0
		}

	case "multiple_response":
		// Multiple selection - validate at least one option
		if len(req.SelectedOptionIDs) == 0 {
			// No options selected - mark as wrong
			isCorrect = false
			score = 0
			selectedOptionIDsStr = ""
			break
		}

		// Parse all selected option IDs
		var selectedOptionIDs []uuid.UUID
		for _, idStr := range req.SelectedOptionIDs {
			id, err := uuid.Parse(idStr)
			if err != nil {
				return nil, errors.New("invalid option ID format")
			}
			selectedOptionIDs = append(selectedOptionIDs, id)
		}

		// Check if all selected options are valid and count correct selections
		correctCount := 0
		totalCorrect := 0
		allCorrect := true

		for _, opt := range question.Options {
			if opt.IsCorrect {
				totalCorrect++
				// Check if this correct option was selected
				found := false
				for _, selID := range selectedOptionIDs {
					if selID == opt.ID {
						found = true
						break
					}
				}
				if found {
					correctCount++
				} else {
					allCorrect = false
				}
			} else {
				// Check if any incorrect option was selected
				for _, selID := range selectedOptionIDs {
					if selID == opt.ID {
						allCorrect = false
						break
					}
				}
			}
		}

		// All correct options must be selected and no incorrect options
		isCorrect = allCorrect && correctCount == totalCorrect && len(selectedOptionIDs) == totalCorrect
		score = question.Points
		if !isCorrect {
			// Partial scoring: give points based on correct selections
			if correctCount > 0 {
				score = (correctCount * question.Points) / totalCorrect
			} else {
				score = 0
			}
		}

		// Store selected option IDs as comma-separated string
		var idStrs []string
		for _, id := range selectedOptionIDs {
			idStrs = append(idStrs, id.String())
		}
		selectedOptionIDsStr = strings.Join(idStrs, ",")

	case "matching":
		// For matching questions, store as comma-separated pairs
		if len(req.SelectedOptionIDs) == 0 {
			// No options selected - mark as wrong
			isCorrect = false
			score = 0
			selectedOptionIDsStr = ""
			break
		}
		selectedOptionIDsStr = strings.Join(req.SelectedOptionIDs, ",")
		// For matching, we need to check each pair
		// This is a simplified version - implement proper matching logic
		isCorrect = false
		score = 0

	case "ordering":
		// For ordering questions, store as comma-separated ordered list
		if len(req.SelectedOptionIDs) == 0 {
			// No options selected - mark as wrong
			isCorrect = false
			score = 0
			selectedOptionIDsStr = ""
			break
		}
		selectedOptionIDsStr = strings.Join(req.SelectedOptionIDs, ",")
		// For ordering, we need to check if the order is correct
		// This is a simplified version - implement proper ordering logic
		isCorrect = false
		score = 0

	default:
		return nil, errors.New("unsupported question type")
	}

	// Create the answer record
	answer := &models.ObjectiveQuestionAnswer{
		ID:                uuid.New(),
		StudentID:         studentID,
		QuestionID:        questionID,
		SelectedOptionIDs: selectedOptionIDsStr,
		SelectedOptionID:  selectedOptionID,
		IsCorrect:         isCorrect,
		Score:             score,
		TimeSpent:         req.TimeSpent,
		AttemptNumber:     attemptNumber,
		Status:            "active",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.db.Create(answer).Error; err != nil {
		return nil, fmt.Errorf("failed to save answer: %w", err)
	}

	// Get correct answer for response
	var correctAnswer string
	var explanation string
	for _, opt := range question.Options {
		if opt.IsCorrect {
			if correctAnswer != "" {
				correctAnswer += ", "
			}
			correctAnswer += opt.OptionText
		}
	}
	explanation = question.AnswerExplanation

	// If no option was selected, set correct answer message
	if req.SelectedOptionID == "" && len(req.SelectedOptionIDs) == 0 {
		correctAnswer = "You did not select any option"
		explanation = "Please select an option to answer this question"
	}

	response := &dto.SubmitAnswerResponse{
		IsCorrect:        isCorrect,
		Score:            score,
		CorrectAnswer:    correctAnswer,
		Explanation:      explanation,
		TotalAttempts:    attemptNumber,
		PreviousAttempts: previousAttempts,
		TimeSpent:        req.TimeSpent,
	}

	return response, nil
}

// GetStudentAnswers retrieves all answers for a student
func (s *StudentAnswerService) GetStudentAnswers(
	studentID uuid.UUID,
	params *dto.StudentAnswerQueryParams,
) ([]dto.StudentAnswerResponse, error) {

	var answers []models.ObjectiveQuestionAnswer

	query := s.db.
		Where("student_id = ? AND deleted_at IS NULL", studentID)

	// Filter by question
	if params.QuestionID != "" {
		questionID, err := uuid.Parse(params.QuestionID)
		if err == nil {
			query = query.Where("question_id = ?", questionID)
		}
	}

	// Filter by correctness
	if params.IsCorrect != nil {
		query = query.Where("is_correct = ?", *params.IsCorrect)
	}

	// Filter by status
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// Apply sorting
	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}

	sortDirection := "DESC"
	if params.SortOrder == "asc" {
		sortDirection = "ASC"
	}

	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortDirection))

	// Apply pagination
	if params.Limit > 0 {
		page := params.Page
		if page < 1 {
			page = 1
		}

		offset := (page - 1) * params.Limit

		query = query.
			Offset(offset).
			Limit(params.Limit)
	}

	// Execute query
	if err := query.
		Preload("Question").
		Find(&answers).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch answers: %w", err)
	}

	// Convert to response
	responses := make([]dto.StudentAnswerResponse, len(answers))

	for i, answer := range answers {
		responses[i] = s.toAnswerResponse(&answer)
	}

	return responses, nil
}

// GetStudentStats retrieves statistics for a student
func (s *StudentAnswerService) GetStudentStats(studentID uuid.UUID, lessonID string) (*dto.StudentAnswerStats, error) {
	// Parse lesson ID if provided
	var lessonIDPtr *uuid.UUID
	if lessonID != "" {
		id, err := uuid.Parse(lessonID)
		if err != nil {
			return nil, errors.New("invalid lesson ID")
		}
		lessonIDPtr = &id
	}

	// Build query for answers
	query := s.db.Model(&models.ObjectiveQuestionAnswer{}).
		Where("student_id = ? AND deleted_at IS NULL", studentID)

	// If lesson ID is provided, join with questions to filter by lesson
	if lessonIDPtr != nil {
		query = query.Joins("JOIN objective_questions ON objective_question_answers.question_id = objective_questions.id").
			Where("objective_questions.lesson_id = ?", lessonIDPtr)
	}

	// Get total questions attempted
	var totalQuestions int64
	if err := query.Count(&totalQuestions).Error; err != nil {
		return nil, fmt.Errorf("failed to count questions: %w", err)
	}

	// Get correct answers count
	var correctAnswers int64
	if err := query.Where("is_correct = ?", true).Count(&correctAnswers).Error; err != nil {
		return nil, fmt.Errorf("failed to count correct answers: %w", err)
	}

	// Get total score
	var totalScore int
	if err := query.Select("COALESCE(SUM(score), 0)").Scan(&totalScore).Error; err != nil {
		return nil, fmt.Errorf("failed to get total score: %w", err)
	}

	// Get average time
	var avgTime float64
	if err := query.Select("COALESCE(AVG(time_spent), 0)").Scan(&avgTime).Error; err != nil {
		return nil, fmt.Errorf("failed to get average time: %w", err)
	}

	scorePercentage := 0.0
	if totalQuestions > 0 {
		scorePercentage = (float64(correctAnswers) / float64(totalQuestions)) * 100
	}

	stats := &dto.StudentAnswerStats{
		TotalQuestions:   int(totalQuestions),
		CorrectAnswers:   int(correctAnswers),
		IncorrectAnswers: int(totalQuestions) - int(correctAnswers),
		ScorePercentage:  scorePercentage,
		AverageTime:      avgTime,
		TotalScore:       totalScore,
	}

	return stats, nil
}

// GetQuestionAttempts retrieves all attempts for a specific question by a student
func (s *StudentAnswerService) GetQuestionAttempts(studentID uuid.UUID, questionID string) ([]dto.StudentAnswerResponse, error) {
	qID, err := uuid.Parse(questionID)
	if err != nil {
		return nil, errors.New("invalid question ID")
	}

	var answers []models.ObjectiveQuestionAnswer
	if err := s.db.Where("student_id = ? AND question_id = ? AND deleted_at IS NULL", studentID, qID).
		Order("attempt_number DESC").
		Preload("Question").
		Find(&answers).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch attempts: %w", err)
	}

	responses := make([]dto.StudentAnswerResponse, len(answers))
	for i, answer := range answers {
		responses[i] = s.toAnswerResponse(&answer)
	}

	return responses, nil
}

// GetUnansweredQuestions retrieves all unanswered questions for a student in a lesson
func (s *StudentAnswerService) GetUnansweredQuestions(studentID uuid.UUID, lessonID string) ([]models.ObjectiveQuestion, error) {
	// Parse lesson ID
	lID, err := uuid.Parse(lessonID)
	if err != nil {
		return nil, errors.New("invalid lesson ID")
	}

	// Get all questions for the lesson
	var allQuestions []models.ObjectiveQuestion
	if err := s.db.Where("lesson_id = ? AND status = ? AND deleted_at IS NULL", lID, "active").
		Find(&allQuestions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch questions: %w", err)
	}

	// Get all answered question IDs for this student
	var answeredQuestionIDs []uuid.UUID
	if err := s.db.Model(&models.ObjectiveQuestionAnswer{}).
		Where("student_id = ? AND deleted_at IS NULL", studentID).
		Pluck("question_id", &answeredQuestionIDs).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch answered questions: %w", err)
	}

	// Create a map for quick lookup
	answeredMap := make(map[uuid.UUID]bool)
	for _, id := range answeredQuestionIDs {
		answeredMap[id] = true
	}

	// Filter unanswered questions
	var unansweredQuestions []models.ObjectiveQuestion
	for _, q := range allQuestions {
		if !answeredMap[q.ID] {
			unansweredQuestions = append(unansweredQuestions, q)
		}
	}

	return unansweredQuestions, nil
}

// SubmitAllAnswers submits all answers for a student in a lesson
func (s *StudentAnswerService) SubmitAllAnswers(studentID uuid.UUID, lessonID string, answers map[string]string) ([]dto.SubmitAnswerResponse, error) {
	var results []dto.SubmitAnswerResponse

	// Get all questions for the lesson
	var questions []models.ObjectiveQuestion
	if err := s.db.Where("lesson_id = ? AND status = ? AND deleted_at IS NULL", lessonID, "active").
		Preload("Options").
		Find(&questions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch questions: %w", err)
	}

	for _, question := range questions {
		// Check if answer exists for this question
		selectedOptionID, exists := answers[question.ID.String()]
		
		var req dto.SubmitAnswerRequest
		req.QuestionID = question.ID.String()
		req.TimeSpent = 0 // This could be calculated per question
		
		if exists && selectedOptionID != "" {
			// Student selected an option
			req.SelectedOptionID = selectedOptionID
		} else {
			// Student did not select an option - mark as wrong
			req.SelectedOptionID = ""
		}

		// Submit the answer
		response, err := s.SubmitAnswer(studentID, &req)
		if err != nil {
			// Log error but continue with other questions
			continue
		}
		results = append(results, *response)
	}

	return results, nil
}

// toAnswerResponse converts model to response DTO
func (s *StudentAnswerService) toAnswerResponse(answer *models.ObjectiveQuestionAnswer) dto.StudentAnswerResponse {
	return dto.StudentAnswerResponse{
		ID:                answer.ID.String(),
		StudentID:         answer.StudentID.String(),
		QuestionID:        answer.QuestionID.String(),
		SelectedOptionIDs: answer.SelectedOptionIDs,
		SelectedOptionID:  answer.SelectedOptionID.String(),
		IsCorrect:         answer.IsCorrect,
		Score:             answer.Score,
		TimeSpent:         answer.TimeSpent,
		AttemptNumber:     answer.AttemptNumber,
		Status:            answer.Status,
		CreatedAt:         answer.CreatedAt,
		UpdatedAt:         answer.UpdatedAt,
	}
}