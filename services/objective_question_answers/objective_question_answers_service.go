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

	// Parse grade ID
	gradeID, err := uuid.Parse(req.GradeID)
	if err != nil {
		return nil, errors.New("invalid grade ID")
	}

	// Verify grade exists
	var grade models.ClassGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", gradeID).First(&grade).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("grade not found")
		}
		return nil, fmt.Errorf("failed to verify grade: %w", err)
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

	// Check if student has already answered this question in this grade
	var existingAnswers []models.ObjectiveQuestionAnswer
	if err := s.db.Where("student_id = ? AND question_id = ? AND grade_id = ? AND deleted_at IS NULL", studentID, questionID, gradeID).
		Order("created_at DESC").
		Find(&existingAnswers).Error; err != nil {
		return nil, fmt.Errorf("failed to check existing answers: %w", err)
	}

	previousAttempts := len(existingAnswers)

	// Validate answer based on question type
	var isCorrect bool
	var score int
	var selectedOptionIDsStr string
	var selectedOptionID uuid.UUID

	switch question.QuestionType {
	case "multiple_choice":
		if req.SelectedOptionID == "" {
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
		if req.SelectedOptionID == "" {
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
		if len(req.SelectedOptionIDs) == 0 {
			isCorrect = false
			score = 0
			selectedOptionIDsStr = ""
			break
		}

		var selectedOptionIDs []uuid.UUID
		for _, idStr := range req.SelectedOptionIDs {
			id, err := uuid.Parse(idStr)
			if err != nil {
				return nil, errors.New("invalid option ID format")
			}
			selectedOptionIDs = append(selectedOptionIDs, id)
		}

		correctCount := 0
		totalCorrect := 0
		allCorrect := true

		for _, opt := range question.Options {
			if opt.IsCorrect {
				totalCorrect++
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
				for _, selID := range selectedOptionIDs {
					if selID == opt.ID {
						allCorrect = false
						break
					}
				}
			}
		}

		isCorrect = allCorrect && correctCount == totalCorrect && len(selectedOptionIDs) == totalCorrect
		score = question.Points
		if !isCorrect {
			if correctCount > 0 {
				score = (correctCount * question.Points) / totalCorrect
			} else {
				score = 0
			}
		}

		var idStrs []string
		for _, id := range selectedOptionIDs {
			idStrs = append(idStrs, id.String())
		}
		selectedOptionIDsStr = strings.Join(idStrs, ",")

	case "matching":
		if len(req.SelectedOptionIDs) == 0 {
			isCorrect = false
			score = 0
			selectedOptionIDsStr = ""
			break
		}
		selectedOptionIDsStr = strings.Join(req.SelectedOptionIDs, ",")
		isCorrect = false
		score = 0

	case "ordering":
		if len(req.SelectedOptionIDs) == 0 {
			isCorrect = false
			score = 0
			selectedOptionIDsStr = ""
			break
		}
		selectedOptionIDsStr = strings.Join(req.SelectedOptionIDs, ",")
		isCorrect = false
		score = 0

	default:
		return nil, errors.New("unsupported question type")
	}

	// Create the answer record
	answer := &models.ObjectiveQuestionAnswer{
		ID:                uuid.New(),
		GradeID:           gradeID,
		StudentID:         studentID,
		QuestionID:        questionID,
		SelectedOptionIDs: selectedOptionIDsStr,
		SelectedOptionID:  selectedOptionID,
		IsCorrect:         isCorrect,
		Score:             score,
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

	if req.SelectedOptionID == "" && len(req.SelectedOptionIDs) == 0 {
		correctAnswer = "You did not select any option"
		explanation = "Please select an option to answer this question"
	}

	response := &dto.SubmitAnswerResponse{
		IsCorrect:        isCorrect,
		Score:            score,
		CorrectAnswer:    correctAnswer,
		Explanation:      explanation,
		TotalAttempts:    previousAttempts + 1,
		PreviousAttempts: previousAttempts,
		TimeSpent:        req.TimeSpent,
	}

	return response, nil
}

// GetStudentAnswers retrieves all answers with optional filters
func (s *StudentAnswerService) GetStudentAnswers(
	params *dto.StudentAnswerQueryParams,
) ([]dto.StudentAnswerResponse, error) {

	var answers []models.ObjectiveQuestionAnswer

	// Start building query
	query := s.db.Model(&models.ObjectiveQuestionAnswer{}).Where("objective_question_answers.deleted_at IS NULL")

	// If lesson_id is provided, join with questions to filter by lesson
	if params.LessonID != "" {
		lessonID, err := uuid.Parse(params.LessonID)
		if err == nil {
			query = query.Joins("JOIN objective_questions ON objective_question_answers.question_id = objective_questions.id").
				Where("objective_questions.lesson_id = ?", lessonID)
		}
	}

	// If grade_id is provided, filter by grade
	if params.GradeID != "" {
		gradeID, err := uuid.Parse(params.GradeID)
		if err == nil {
			query = query.Where("objective_question_answers.grade_id = ?", gradeID)
		}
	}

	// If student_id is provided, filter by student
	if params.StudentID != "" {
		sID, err := uuid.Parse(params.StudentID)
		if err == nil {
			query = query.Where("objective_question_answers.student_id = ?", sID)
		}
	}

	// Filter by question
	if params.QuestionID != "" {
		questionID, err := uuid.Parse(params.QuestionID)
		if err == nil {
			query = query.Where("objective_question_answers.question_id = ?", questionID)
		}
	}

	// Filter by correctness
	if params.IsCorrect != nil {
		query = query.Where("objective_question_answers.is_correct = ?", *params.IsCorrect)
	}

	// Filter by status
	if params.Status != "" {
		query = query.Where("objective_question_answers.status = ?", params.Status)
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

	query = query.Order(fmt.Sprintf("objective_question_answers.%s %s", sortBy, sortDirection))

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
		Preload("Student").
		Preload("Grade").
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

// GetLessonTestResults retrieves aggregated results for all students in a lesson
func (s *StudentAnswerService) GetLessonTestResults(lessonID string) (*dto.LessonTestResultsResponse, error) {
	lID, err := uuid.Parse(lessonID)
	if err != nil {
		return nil, errors.New("invalid lesson ID")
	}

	// Get all questions for this lesson
	var questions []models.ObjectiveQuestion
	if err := s.db.Where("lesson_id = ? AND deleted_at IS NULL", lID).
		Find(&questions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch questions: %w", err)
	}

	if len(questions) == 0 {
		return &dto.LessonTestResultsResponse{
			Results: []dto.LessonTestResult{},
			Total:   0,
			Stats: &dto.LessonTestStats{
				TotalStudents: 0,
			},
		}, nil
	}

	// Get all question IDs
	questionIDs := make([]uuid.UUID, len(questions))
	maxScore := 0
	for i, q := range questions {
		questionIDs[i] = q.ID
		maxScore += q.Points
	}

	// Get all answers for these questions
	var answers []models.ObjectiveQuestionAnswer
	if err := s.db.Where("question_id IN ? AND deleted_at IS NULL", questionIDs).
		Preload("Student").
		Preload("Question").
		Preload("Grade").
		Find(&answers).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch answers: %w", err)
	}

	// Group answers by student
	studentAnswers := make(map[uuid.UUID][]models.ObjectiveQuestionAnswer)
	for _, answer := range answers {
		studentAnswers[answer.StudentID] = append(studentAnswers[answer.StudentID], answer)
	}

	// Build results for each student
	var results []dto.LessonTestResult
	var totalScoreSum int
	var passedCount int
	var highestScore float64
	var lowestScore float64 = 100

	for studentID, studentAnswerList := range studentAnswers {
		var correctCount int
		var totalScore int

		for _, ans := range studentAnswerList {
			if ans.IsCorrect {
				correctCount++
			}
			totalScore += ans.Score
		}

		scorePercentage := 0.0
		if maxScore > 0 {
			scorePercentage = (float64(totalScore) / float64(maxScore)) * 100
		}

		// Get student info
		var student models.User
		if len(studentAnswerList) > 0 && studentAnswerList[0].Student.ID != uuid.Nil {
			student = studentAnswerList[0].Student
		} else {
			if err := s.db.Where("id = ? AND deleted_at IS NULL", studentID).First(&student).Error; err != nil {
				continue
			}
		}

		// Get grade info
		var gradeName, gradeID string
		if len(studentAnswerList) > 0 && studentAnswerList[0].Grade.ID != uuid.Nil {
			gradeName = studentAnswerList[0].Grade.Name
			gradeID = studentAnswerList[0].GradeID.String()
		}

		result := dto.LessonTestResult{
			StudentID:        studentID.String(),
			StudentFirstName: student.FirstName,
			StudentLastName:  student.LastName,
			StudentEmail:     student.Email,
			GradeID:          gradeID,
			GradeName:        gradeName,
			TotalQuestions:   len(questions),
			CorrectAnswers:   correctCount,
			IncorrectAnswers: len(questions) - correctCount,
			TotalScore:       totalScore,
			MaxScore:         maxScore,
			ScorePercentage:  scorePercentage,
			Status:           "completed",
			CreatedAt:        time.Now().Format(time.RFC3339),
			UpdatedAt:        time.Now().Format(time.RFC3339),
		}

		results = append(results, result)
		totalScoreSum += totalScore

		if scorePercentage >= 50 {
			passedCount++
		}

		if scorePercentage > highestScore {
			highestScore = scorePercentage
		}
		if scorePercentage < lowestScore {
			lowestScore = scorePercentage
		}
	}

	totalStudents := len(results)
	averageScore := 0.0
	if totalStudents > 0 {
		averageScore = float64(totalScoreSum) / float64(totalStudents)
	}

	stats := &dto.LessonTestStats{
		TotalStudents: totalStudents,
		AverageScore:  averageScore,
		HighestScore:  highestScore,
		LowestScore:   lowestScore,
		PassedCount:   passedCount,
		FailedCount:   totalStudents - passedCount,
	}

	return &dto.LessonTestResultsResponse{
		Results: results,
		Total:   int64(totalStudents),
		Stats:   stats,
	}, nil
}

// GetGradeTestResults retrieves aggregated results for all students in a grade
func (s *StudentAnswerService) GetGradeTestResults(gradeID string) (*dto.GradeTestResultsResponse, error) {
	gID, err := uuid.Parse(gradeID)
	if err != nil {
		return nil, errors.New("invalid grade ID")
	}

	// Get all answers for this grade
	var answers []models.ObjectiveQuestionAnswer
	if err := s.db.Where("grade_id = ? AND deleted_at IS NULL", gID).
		Preload("Student").
		Preload("Question").
		Preload("Grade").
		Find(&answers).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch answers: %w", err)
	}

	if len(answers) == 0 {
		return &dto.GradeTestResultsResponse{
			Results: []dto.GradeTestResult{},
			Total:   0,
			Stats:   &dto.GradeTestStats{TotalStudents: 0},
		}, nil
	}

	// Group answers by student
	studentAnswers := make(map[uuid.UUID][]models.ObjectiveQuestionAnswer)
	for _, answer := range answers {
		studentAnswers[answer.StudentID] = append(studentAnswers[answer.StudentID], answer)
	}

	// Calculate max possible score per question
	questionMaxScores := make(map[uuid.UUID]int)
	for _, ans := range answers {
		if ans.Question.ID != uuid.Nil {
			questionMaxScores[ans.QuestionID] = ans.Question.Points
		}
	}

	var results []dto.GradeTestResult
	var passedCount int
	var highestScore float64
	var lowestScore float64 = 100
	var totalScoreSum int

	for studentID, studentAnswerList := range studentAnswers {
		var correctCount int
		var totalScore int

		for _, ans := range studentAnswerList {
			if ans.IsCorrect {
				correctCount++
			}
			totalScore += ans.Score
		}

		// Get student info
		var student models.User
		if len(studentAnswerList) > 0 && studentAnswerList[0].Student.ID != uuid.Nil {
			student = studentAnswerList[0].Student
		} else {
			continue
		}

		// Get grade info
		var gradeName, gradeIDStr string
		if len(studentAnswerList) > 0 && studentAnswerList[0].Grade.ID != uuid.Nil {
			gradeName = studentAnswerList[0].Grade.Name
			gradeIDStr = studentAnswerList[0].GradeID.String()
		}

		scorePercentage := 0.0
		totalQuestions := len(studentAnswerList)
		if totalQuestions > 0 {
			scorePercentage = (float64(correctCount) / float64(totalQuestions)) * 100
		}

		result := dto.GradeTestResult{
			StudentID:        studentID.String(),
			StudentFirstName: student.FirstName,
			StudentLastName:  student.LastName,
			StudentEmail:     student.Email,
			GradeID:          gradeIDStr,
			GradeName:        gradeName,
			TotalQuestions:   totalQuestions,
			CorrectAnswers:   correctCount,
			IncorrectAnswers: totalQuestions - correctCount,
			TotalScore:       totalScore,
			MaxScore:         totalScore, // This should be calculated based on questions
			ScorePercentage:  scorePercentage,
			Status:           "completed",
			CreatedAt:        time.Now().Format(time.RFC3339),
			UpdatedAt:        time.Now().Format(time.RFC3339),
		}

		results = append(results, result)
		totalScoreSum += totalScore

		if scorePercentage >= 50 {
			passedCount++
		}

		if scorePercentage > highestScore {
			highestScore = scorePercentage
		}
		if scorePercentage < lowestScore {
			lowestScore = scorePercentage
		}
	}

	totalStudents := len(results)
	averageScore := 0.0
	if totalStudents > 0 {
		averageScore = float64(totalScoreSum) / float64(totalStudents)
	}

	stats := &dto.GradeTestStats{
		TotalStudents: totalStudents,
		AverageScore:  averageScore,
		HighestScore:  highestScore,
		LowestScore:   lowestScore,
		PassedCount:   passedCount,
		FailedCount:   totalStudents - passedCount,
	}

	return &dto.GradeTestResultsResponse{
		Results: results,
		Total:   int64(totalStudents),
		Stats:   stats,
	}, nil
}

// GetStudentStats retrieves statistics for a student
func (s *StudentAnswerService) GetStudentStats( lessonID string, gradeID string) (*dto.StudentAnswerStats, error) {
	// Build base query for answers
	baseQuery := s.db.Model(&models.ObjectiveQuestionAnswer{})
	var totalQuestions int64
	var correctAnswers int64
	var totalScore int

	// Filter by grade if provided
	if gradeID != "" {
		gID, err := uuid.Parse(gradeID)
		if err != nil {
			return nil, errors.New("invalid grade ID")
		}
		baseQuery = baseQuery.Where("grade_id = ?", gID)
	}

	if lessonID != "" {
		lID, err := uuid.Parse(lessonID)
		if err != nil {
			return nil, errors.New("invalid lesson ID")
		}

		// Use subquery to get question IDs for the lesson
		var questionIDs []uuid.UUID
		if err := s.db.Model(&models.ObjectiveQuestion{}).
			Where("lesson_id = ? AND deleted_at IS NULL", lID).
			Pluck("id", &questionIDs).Error; err != nil {
			return nil, fmt.Errorf("failed to get question IDs: %w", err)
		}

		if len(questionIDs) == 0 {
			return &dto.StudentAnswerStats{
				TotalQuestions:   0,
				CorrectAnswers:   0,
				IncorrectAnswers: 0,
				ScorePercentage:  0,
				TotalScore:       0,
			}, nil
		}

		query := baseQuery.Where("question_id IN ?", questionIDs)

		if err := query.Count(&totalQuestions).Error; err != nil {
			return nil, fmt.Errorf("failed to count questions: %w", err)
		}

		if err := query.Where("is_correct = ?", true).Count(&correctAnswers).Error; err != nil {
			return nil, fmt.Errorf("failed to count correct answers: %w", err)
		}

		if err := query.Select("COALESCE(SUM(score), 0)").Scan(&totalScore).Error; err != nil {
			return nil, fmt.Errorf("failed to get total score: %w", err)
		}
	} else {
		if err := baseQuery.Count(&totalQuestions).Error; err != nil {
			return nil, fmt.Errorf("failed to count questions: %w", err)
		}

		if err := baseQuery.Where("is_correct = ?", true).Count(&correctAnswers).Error; err != nil {
			return nil, fmt.Errorf("failed to count correct answers: %w", err)
		}

		if err := baseQuery.Select("COALESCE(SUM(score), 0)").Scan(&totalScore).Error; err != nil {
			return nil, fmt.Errorf("failed to get total score: %w", err)
		}
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
		TotalScore:       totalScore,
	}

	return stats, nil
}

// toAnswerResponse converts model to response DTO
func (s *StudentAnswerService) toAnswerResponse(answer *models.ObjectiveQuestionAnswer) dto.StudentAnswerResponse {
	response := dto.StudentAnswerResponse{
		ID:                answer.ID.String(),
		GradeID:           answer.GradeID.String(),
		StudentID:         answer.StudentID.String(),
		QuestionID:        answer.QuestionID.String(),
		SelectedOptionIDs: answer.SelectedOptionIDs,
		SelectedOptionID:  answer.SelectedOptionID.String(),
		IsCorrect:         answer.IsCorrect,
		Score:             answer.Score,
		Status:            answer.Status,
		CreatedAt:         answer.CreatedAt,
		UpdatedAt:         answer.UpdatedAt,
	}

	// Add student details if preloaded
	if answer.Student.ID != uuid.Nil {
		response.Student = &dto.UserResponse{
			ID:        answer.Student.ID.String(),
			FirstName: answer.Student.FirstName,
			LastName:  answer.Student.LastName,
			Email:     answer.Student.Email,
			Phone:     answer.Student.Phone,
			Role:      answer.Student.Role,
			Position:  answer.Student.Position,
		}
	}

	// Add question details if preloaded
	if answer.Question.ID != uuid.Nil {
		response.Question = &dto.ObjectiveQuestionResponse{
			ID:           answer.Question.ID.String(),
			QuestionText: answer.Question.QuestionText,
			QuestionType: answer.Question.QuestionType,
			Points:       answer.Question.Points,
			Status:       answer.Question.Status,
		}
	}

	// Add grade details if preloaded
	if answer.Grade.ID != uuid.Nil {
		response.Grade = &dto.ClassGradeResponse{
			ID:          answer.Grade.ID.String(),
			Name:        answer.Grade.Name,
			Code:        answer.Grade.Code,
			Level:       answer.Grade.Level,
			Description: answer.Grade.Description,
			Status:      answer.Grade.Status,
		}
	}

	return response
}