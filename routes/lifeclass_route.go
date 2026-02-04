// routes/live_class_routes.go
package routes

import (
    "crm-go/controllers/liveclass"
    "crm-go/services/liveclass"
    "crm-go/services/activity"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func LiveClassRoutes(r *gin.Engine, db *gorm.DB) {
    liveClassService := services.NewLiveClassService(db)
    activityService := activity.NewService(db)
    liveClassController := controllers.NewLiveClassController(db, liveClassService, activityService)
    
    liveClassRoutes := r.Group("/api/live-classes")
    {
        liveClassRoutes.POST("", liveClassController.CreateLiveClass)
        // Add other routes: GET, PUT, DELETE, etc.
    }
    
}