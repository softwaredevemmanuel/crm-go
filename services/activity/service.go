package activity

import "gorm.io/gorm"

type Service struct {
	Logger      *Logger
	Assignments *AssignmentActivity
	Users       *UserActivity
	Topics      *TopicActivity
	Grades      *GradeActivity
}

func NewService(db *gorm.DB) *Service {
	logger := NewLogger(db)

	return &Service{
		Logger:      logger,
		Assignments: &AssignmentActivity{logger},
		Users:       &UserActivity{logger},
		Topics:      &TopicActivity{logger},
		Grades:      &GradeActivity{logger},
	}
}


