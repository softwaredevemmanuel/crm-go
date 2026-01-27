import (
	controllers "crm-go/controllers"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"fmt"
	"strings"
)

// ActivityLogMiddleware creates middleware that automatically logs requests
func ActivityLogMiddleware(activityLogger *services.ActivityLogger) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Skip logging for certain paths
        if shouldSkipLogging(c.Request.URL.Path) {
            c.Next()
            return
        }
        
        // Start timer
        start := time.Now()
        
        // Process request
        c.Next()
        
        // Get user from context (if authenticated)
        userID, exists := c.Get("user_id")
        if !exists {
            userID = uuid.Nil
        }
        
        // Get response status
        status := c.Writer.Status()
        
        // Log the request
        go func() {
            metadata := map[string]interface{}{
                "method":     c.Request.Method,
                "path":       c.Request.URL.Path,
                "status":     status,
                "latency_ms": time.Since(start).Milliseconds(),
                "query":      c.Request.URL.RawQuery,
            }
            
            // Only log certain status codes
            if shouldLogStatus(status) {
                _ = activityLogger.LogFromGinContext(c, services.ActivityLogData{
                    UserID:     userID.(uuid.UUID),
                    Action:     "http_request",
                    EntityID:   uuid.New(), // Generate a unique ID for the request
                    EntityType: "http_request",
                    Details:    fmt.Sprintf("%s %s - %d", c.Request.Method, c.Request.URL.Path, status),
                    Metadata:   metadata,
                })
            }
        }()
    }
}

func shouldSkipLogging(path string) bool {
    skipPaths := []string{
        "/health",
        "/metrics",
        "/favicon.ico",
        "/static/",
        "/assets/",
    }
    
    for _, skipPath := range skipPaths {
        if strings.HasPrefix(path, skipPath) {
            return true
        }
    }
    return false
}

func shouldLogStatus(status int) bool {
    // Log all 4xx and 5xx errors, plus important 2xx responses
    return status >= 400 || status == 200 || status == 201
}