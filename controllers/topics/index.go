package controllers

import(
	  "crm-go/services/topics"
    "crm-go/services/activity"
	    "gorm.io/gorm"

)
type TopicController struct {
    db                *gorm.DB
    createTopicService *services.CreateTopicService
    getTopicService    *services.GetTopicService
    updateTopicService *services.UpdateTopicService
    activity          *activity.Service
}

func NewCreateTopicController(
	db *gorm.DB, 
	createTopicService *services.CreateTopicService, 
	getTopicService *services.GetTopicService, 
	updateTopicService *services.UpdateTopicService, 
	activitySvc *activity.Service) *TopicController {
    return &TopicController{
        db:                db,
        createTopicService: createTopicService,
        getTopicService:    getTopicService,
        updateTopicService: updateTopicService,
        activity:          activitySvc,
    }
}