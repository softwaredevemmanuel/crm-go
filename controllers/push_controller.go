package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/services"
)

type PushHandler struct {
	pushService *services.PushService
}

func NewPushHandler(pushService *services.PushService) *PushHandler {
	return &PushHandler{
		pushService: pushService,
	}
}

// Subscribe handles push subscription
// @Summary Subscribe to push notifications
// @Description Subscribe a user to push notifications
// @Tags Push Notifications
// @Accept json
// @Produce json
// @Param request body services.SubscriptionRequest true "Subscription request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/push/subscribe [post]
func (h *PushHandler) Subscribe(c *gin.Context) {
	var req services.SubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	err = h.pushService.Subscribe(
		userID,
		req.Subscription.Endpoint,
		req.Subscription.Keys.P256dh,
		req.Subscription.Keys.Auth,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save subscription",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Subscribed successfully",
	})
}

// Unsubscribe handles push unsubscription
// @Summary Unsubscribe from push notifications
// @Description Unsubscribe a user from push notifications
// @Tags Push Notifications
// @Accept json
// @Produce json
// @Param user_id body string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/push/unsubscribe [delete]
func (h *PushHandler) Unsubscribe(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	if err := h.pushService.Unsubscribe(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to unsubscribe",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Unsubscribed successfully",
	})
}

// SendTestNotification sends a test push notification
// @Summary Send test push notification
// @Description Send a test push notification to a user
// @Tags Push Notifications
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/push/test [post]
func (h *PushHandler) SendTestNotification(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id is required",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	err = h.pushService.SendPushNotification(
		userID,
		"🔔 Test Notification",
		"This is a test notification from Ehizua Hub!",
		"/notifications",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send notification: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Test notification sent successfully",
	})
}