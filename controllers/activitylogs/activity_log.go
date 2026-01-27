package activitylogs


import (
    "context"
    "encoding/json"
    "fmt"
    "strings"
    "time"
    "crm-go/models"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "gorm.io/datatypes"
)

type ActivityLogger struct {
    db *gorm.DB
}

func NewActivityLogger(db *gorm.DB) *ActivityLogger {
    return &ActivityLogger{db: db}
}

// Log creates a new activity log entry
func (al *ActivityLogger) Log(ctx context.Context, logData ActivityLogData) error {
    activityLog := models.ActivityLog{
        ID:         uuid.New(),
        UserID:     logData.UserID,
        Action:     logData.Action,
        EntityID:   logData.EntityID,
        EntityType: logData.EntityType,
        Details:    logData.Details,
        IPAddress:  logData.IPAddress,
        UserAgent:  logData.UserAgent,
        CreatedAt:  time.Now(),
    }
    
    // Add metadata if provided
    if len(logData.Metadata) > 0 {
        metadataJSON, err := json.Marshal(logData.Metadata)
        if err != nil {
            return fmt.Errorf("failed to marshal metadata: %w", err)
        }
        activityLog.Metadata = datatypes.JSON(metadataJSON)
    }
    
    // Create within transaction if provided
    if logData.Tx != nil {
        return logData.Tx.Create(&activityLog).Error
    }
    
    // Otherwise use the regular DB connection
    return al.db.Create(&activityLog).Error
}

// LogFromGinContext creates a log entry from Gin context
func (al *ActivityLogger) LogFromGinContext(c *gin.Context, logData ActivityLogData) error {
    // Get IP address
    ip := c.ClientIP()
    if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
        ip = strings.Split(forwarded, ",")[0]
    }
    
    // Get user agent
    userAgent := c.GetHeader("User-Agent")
    
    logData.IPAddress = ip
    logData.UserAgent = userAgent
    
    return al.Log(c.Request.Context(), logData)
}

// ActivityLogData contains all data needed to create an activity log
type ActivityLogData struct {
    UserID     uuid.UUID              `json:"user_id"`
    Action     string                 `json:"action"`
    EntityID   uuid.UUID              `json:"entity_id"`
    EntityType string                 `json:"entity_type"`
    Details    string                 `json:"details"`
    IPAddress  string                 `json:"ip_address,omitempty"`
    UserAgent  string                 `json:"user_agent,omitempty"`
    Metadata   map[string]interface{} `json:"metadata,omitempty"`
    Tx         *gorm.DB               `json:"-"` // Optional transaction
}

// Helper function to log assignment submission
func (al *ActivityLogger) LogAssignmentSubmission(tx *gorm.DB, userID, assignmentID, submissionID uuid.UUID, assignmentTitle string, metadata map[string]interface{}) error {
    if metadata == nil {
        metadata = make(map[string]interface{})
    }
    metadata["assignment_id"] = assignmentID
    metadata["submission_id"] = submissionID
    
    return al.Log(context.Background(), ActivityLogData{
        UserID:     userID,
        Action:     models.ActionAssignmentSubmit,
        EntityID:   submissionID,
        EntityType: "assignment_submission",
        Details:    fmt.Sprintf("Submitted assignment: %s", assignmentTitle),
        Metadata:   metadata,
        Tx:         tx,
    })
}

// LogUserLogin logs user login activity
func (al *ActivityLogger) LogUserLogin(c *gin.Context, userID uuid.UUID, success bool, metadata map[string]interface{}) error {
    if metadata == nil {
        metadata = make(map[string]interface{})
    }
    metadata["success"] = success
    
    action := models.ActionLogin
    details := "User logged in"
    if !success {
        details = "Failed login attempt"
    }
    
    logData := ActivityLogData{
        UserID:     userID,
        Action:     action,
        EntityID:   userID,
        EntityType: "user",
        Details:    details,
        Metadata:   metadata,
    }
    
    return al.LogFromGinContext(c, logData)
}