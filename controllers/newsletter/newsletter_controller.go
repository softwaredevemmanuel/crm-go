package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/newsletter"
)

type NewsletterController struct {
	Service *services.NewsletterService
}

func NewNewsletterController(service *services.NewsletterService) *NewsletterController {
	return &NewsletterController{Service: service}
}

// ============================================================
// PUBLIC endpoints
// ============================================================

// Subscribe godoc
//
//	@Summary		Subscribe to newsletter
//	@Description	Public endpoint. Adds the email as a pending subscriber and sends a confirmation email.
//	@Tags			Newsletter
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.SubscribeRequest	true	"Subscription details"
//	@Success		201		{object}	dto.SimpleMessage
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid payload"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/api/newsletter/subscribe [post]
func (ctrl *NewsletterController) Subscribe(c *gin.Context) {
	var req dto.SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:  "validation failed",
			Fields: map[string]string{"body": err.Error()},
		})
		return
	}

	ip := c.ClientIP()
	ua := c.Request.UserAgent()

	_, err := ctrl.Service.Subscribe(req, ip, ua)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// Don't reveal whether the email was already subscribed.
	c.JSON(http.StatusCreated, dto.SimpleMessage{
		Message: "Thanks for subscribing! Please check your email to confirm.",
	})
}

// ConfirmSubscription godoc
//
//	@Summary		Confirm subscription (link from email)
//	@Description	Confirms a pending subscriber using the token from the confirmation email.
//	@Tags			Newsletter
//	@Produce		json
//	@Param			token	query		string	true	"Confirmation token"
//	@Success		200		{object}	dto.SimpleMessage
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid token"
//	@Router			/api/newsletter/confirm [get]
func (ctrl *NewsletterController) ConfirmSubscription(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "missing token"})
		return
	}

	_, err := ctrl.Service.ConfirmByToken(token)
	if err != nil {
		if errors.Is(err, services.ErrInvalidToken) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SimpleMessage{
		Message: "Subscription confirmed. Welcome aboard!",
	})
}

// Unsubscribe godoc
//
//	@Summary		Unsubscribe from newsletter
//	@Description	Unsubscribes a subscriber using their unsubscribe token.
//	@Tags			Newsletter
//	@Produce		json
//	@Param			token	query		string	true	"Unsubscribe token"
//	@Param			reason	query		string	false	"Optional reason"
//	@Success		200		{object}	dto.SimpleMessage
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid token"
//	@Router			/api/newsletter/unsubscribe [get]
func (ctrl *NewsletterController) Unsubscribe(c *gin.Context) {
	token := c.Query("token")
	reason := c.Query("reason")
	if token == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "missing token"})
		return
	}

	_, err := ctrl.Service.UnsubscribeByToken(token, reason)
	if err != nil {
		if errors.Is(err, services.ErrInvalidToken) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SimpleMessage{
		Message: "You have been unsubscribed. We're sorry to see you go.",
	})
}

// ============================================================
// ADMIN endpoints
// ============================================================

// ListSubscribers godoc
//
//	@Summary		List subscribers (paginated)
//	@Description	Admin lists all subscribers with optional search, status, and source filters.
//	@Tags			Newsletter Admin
//	@Produce		json
//	@Param			page		query		int		false	"Page number (default 1)"			minimum(1)
//	@Param			limit		query		int		false	"Items per page (default 20, max 100)"
//	@Param			search		query		string	false	"Search by email or name"
//	@Param			status		query		string	false	"Filter by status"	Enums(pending, subscribed, unsubscribed, bounced, complained)
//	@Param			source		query		string	false	"Filter by source"
//	@Param			sort_by		query		string	false	"Sort field"		Enums(created_at, updated_at, email, status)
//	@Param			sort_order	query		string	false	"Sort direction"	Enums(asc, desc)
//	@Success		200			{object}	dto.PaginatedSubscribersResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/newsletter/subscribers [get]
func (ctrl *NewsletterController) ListSubscribers(c *gin.Context) {
	var q dto.ListSubscribersQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid query", Message: err.Error()})
		return
	}

	subs, total, err := ctrl.Service.List(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	items := make([]dto.SubscriberResponse, 0, len(subs))
	for i := range subs {
		items = append(items, dto.FromSubscriberModel(&subs[i]))
	}

	totalPages := int((total + int64(q.Limit) - 1) / int64(q.Limit))

	c.JSON(http.StatusOK, dto.PaginatedSubscribersResponse{
		Data:       items,
		Total:      total,
		Page:       q.Page,
		Limit:      q.Limit,
		TotalPages: totalPages,
	})
}

// GetSubscriber godoc
//
//	@Summary		Get a subscriber by ID
//	@Tags			Newsletter Admin
//	@Produce		json
//	@Param			id	path		string	true	"Subscriber UUID"	format(uuid)
//	@Success		200	{object}	dto.SingleSubscriberResponse
//	@Failure		400	{object}	dto.ErrorResponse	"Invalid UUID"
//	@Failure		404	{object}	dto.ErrorResponse	"Not found"
//	@Security		BearerAuth
//	@Router			/api/newsletter/subscribers/{id} [get]
func (ctrl *NewsletterController) GetSubscriber(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid subscriber id"})
		return
	}

	sub, err := ctrl.Service.GetByID(id)
	if err != nil {
		if errors.Is(err, services.ErrSubscriberNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SingleSubscriberResponse{Data: dto.FromSubscriberModel(sub)})
}

// CreateSubscriber godoc
//
//	@Summary		Admin creates a subscriber
//	@Description	Adds a subscriber directly, bypassing double opt-in.
//	@Tags			Newsletter Admin
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.AdminCreateSubscriberRequest	true	"Subscriber details"
//	@Success		201		{object}	dto.SingleSubscriberResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		409		{object}	dto.ErrorResponse	"Email already exists"
//	@Security		BearerAuth
//	@Router			/api/newsletter/subscribers [post]
func (ctrl *NewsletterController) CreateSubscriber(c *gin.Context) {
	var req dto.AdminCreateSubscriberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:  "validation failed",
			Fields: map[string]string{"body": err.Error()},
		})
		return
	}

	sub, err := ctrl.Service.AdminCreate(req)
	if err != nil {
		if errors.Is(err, services.ErrSubscriberAlreadyThere) {
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.SingleSubscriberResponse{Data: dto.FromSubscriberModel(sub)})
}

// UpdateSubscriber godoc
//
//	@Summary		Update a subscriber (partial)
//	@Tags			Newsletter Admin
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string							true	"Subscriber UUID"	format(uuid)
//	@Param			payload	body		dto.AdminUpdateSubscriberRequest	true	"Fields to update"
//	@Success		200		{object}	dto.SingleSubscriberResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		404		{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/newsletter/subscribers/{id} [patch]
func (ctrl *NewsletterController) UpdateSubscriber(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid subscriber id"})
		return
	}

	var req dto.AdminUpdateSubscriberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:  "validation failed",
			Fields: map[string]string{"body": err.Error()},
		})
		return
	}

	sub, err := ctrl.Service.AdminUpdate(id, req)
	if err != nil {
		if errors.Is(err, services.ErrSubscriberNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SingleSubscriberResponse{Data: dto.FromSubscriberModel(sub)})
}

// DeleteSubscriber godoc
//
//	@Summary		Delete a subscriber
//	@Tags			Newsletter Admin
//	@Produce		json
//	@Param			id	path	string	true	"Subscriber UUID"	format(uuid)
//	@Success		204	"No Content"
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/newsletter/subscribers/{id} [delete]
func (ctrl *NewsletterController) DeleteSubscriber(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid subscriber id"})
		return
	}

	if err := ctrl.Service.Delete(id); err != nil {
		if errors.Is(err, services.ErrSubscriberNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// NewsletterStats godoc
//
//	@Summary		Get newsletter stats
//	@Tags			Newsletter Admin
//	@Produce		json
//	@Success		200	{object}	dto.NewsletterStats
//	@Failure		500	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/newsletter/stats [get]
func (ctrl *NewsletterController) NewsletterStats(c *gin.Context) {
	stats, err := ctrl.Service.Stats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}