package controllers

import (
	"crm-go/services/activity"
	services "crm-go/services/topics"

	"gorm.io/gorm"
)

type LessonController struct {
	db                  *gorm.DB
	createLessonService *services.CreateLessonService
	getLessonService    *services.GetLessonService
	updateLessonService *services.UpdateLessonService
	activity            *activity.Service
}

func NewCreateLessonController(
	db *gorm.DB,
	createLessonService *services.CreateLessonService,
	getLessonService *services.GetLessonService,
	updateLessonService *services.UpdateLessonService,
	activitySvc *activity.Service) *LessonController {
	return &LessonController{
		db:                  db,
		createLessonService: createLessonService,
		getLessonService:    getLessonService,
		updateLessonService: updateLessonService,
		activity:            activitySvc,
	}
}
