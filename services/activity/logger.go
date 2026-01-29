package activity

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"crm-go/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Logger struct {
	db *gorm.DB
}

func NewLogger(db *gorm.DB) *Logger {
	return &Logger{db: db}
}

func (l *Logger) log(ctx context.Context, data Event) error {
	log := models.ActivityLog{
		ID:         uuid.New(),
		UserID:     data.UserID,
		Action:     data.Action,
		EntityID:   data.EntityID,
		EntityType: data.EntityType,
		Details:    data.Details,
		IPAddress:  data.IPAddress,
		UserAgent:  data.UserAgent,
		CreatedAt:  time.Now(),
	}

	if len(data.Metadata) > 0 {
		meta, err := json.Marshal(data.Metadata)
		if err != nil {
			return fmt.Errorf("metadata marshal failed: %w", err)
		}
		log.Metadata = datatypes.JSON(meta)
	}

	if data.Tx != nil {
		return data.Tx.Create(&log).Error
	}

	return l.db.Create(&log).Error
}

func (l *Logger) LogWithTx(ctx context.Context, tx *gorm.DB, event Event) error {
	event.Tx = tx
	return l.log(ctx, event)
}

func (l *Logger) LogFromRequest(c *gin.Context, event Event) error {
	ip := c.ClientIP()
	if fwd := c.GetHeader("X-Forwarded-For"); fwd != "" {
		ip = strings.Split(fwd, ",")[0]
	}

	event.IPAddress = ip
	event.UserAgent = c.GetHeader("User-Agent")

	return l.log(c.Request.Context(), event)
}
