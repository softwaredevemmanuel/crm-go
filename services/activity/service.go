package activity

import "gorm.io/gorm"

type Service struct {
	Logger             *Logger
	Assignments        *AssignmentActivity
	Users              *UserActivity
	Lessons            *LessonActivity
	Grades             *GradeActivity
	LiveClasses        *LiveClassActivity
	ObjectiveQuestions *ObjectiveActivity
}

func NewService(db *gorm.DB) *Service {
	logger := NewLogger(db)

	return &Service{
		Logger:             logger,
		Assignments:        &AssignmentActivity{logger},
		Users:              &UserActivity{logger},
		Lessons:            &LessonActivity{logger},
		Grades:             &GradeActivity{logger},
		LiveClasses:        &LiveClassActivity{logger},
		ObjectiveQuestions: &ObjectiveActivity{logger},
	}
}
