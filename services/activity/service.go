package activity

import "gorm.io/gorm"

type Service struct {
	Logger             *Logger
	Users              *UserActivity
	Grades             *GradeActivity
	ObjectiveQuestions *ObjectiveActivity
}

func NewService(db *gorm.DB) *Service {
	logger := NewLogger(db)

	return &Service{
		Logger:             logger,
		Users:              &UserActivity{logger},
		Grades:             &GradeActivity{logger},
		ObjectiveQuestions: &ObjectiveActivity{logger},
	}
}
