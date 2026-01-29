package middleware

import (
	"crm-go/services/activity"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ActionHTTP = "http_request"
	EntityHTTP = "http_request"
)

func ActivityLogMiddleware(activitySvc *activity.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		if shouldSkip(c.Request.URL.Path) {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()

		userID, _ := c.Get("user_id")
		if userID == nil {
			userID = uuid.Nil
		}

		status := c.Writer.Status()

		if status < 400 && status != 200 && status != 201 {
			return
		}

		go func() {
			_ = activitySvc.Logger.LogFromRequest(
				c,
				activity.Event{
					UserID:     userID.(uuid.UUID),
					Action:     ActionHTTP,
					EntityID:   uuid.New(),
					EntityType: EntityHTTP,
					Details:    fmt.Sprintf("%s %s (%d)", c.Request.Method, c.Request.URL.Path, status),
					Metadata: map[string]interface{}{
						"method":  c.Request.Method,
						"path":    c.Request.URL.Path,
						"status":  status,
						"latency": time.Since(start).Milliseconds(),
					},
				},
			)
		}()
	}
}

func shouldSkip(path string) bool {
	skip := []string{"/health", "/metrics", "/favicon.ico", "/static", "/assets"}
	for _, s := range skip {
		if strings.HasPrefix(path, s) {
			return true
		}
	}
	return false
}
