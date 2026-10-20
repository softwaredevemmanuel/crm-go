package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/newsletter_campaign"
)

type NewsletterCampaignController struct {
	Service *services.NewsletterCampaignService
}

func NewNewsletterCampaignController(service *services.NewsletterCampaignService) *NewsletterCampaignController {
	return &NewsletterCampaignController{Service: service}
}

// CreateCampaign godoc
//
//	@Summary		Create a newsletter campaign
//	@Description	Admin creates a draft or scheduled newsletter campaign targeting a subscriber audience.
//	@Tags			Newsletter Campaigns
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.CreateCampaignRequest	true	"Campaign details"
//	@Success		201		{object}	dto.SingleCampaignResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Validation failed"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Security		BearerAuth
//	@Router			/api/newsletter/campaigns [post]
func (ctrl *NewsletterCampaignController) CreateCampaign(c *gin.Context) {
	var req dto.CreateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:  "validation failed",
			Fields: map[string]string{"body": err.Error()},
		})
		return
	}

	// Optional: read the authenticated user ID from context (set by auth middleware)
	var createdBy *uuid.UUID
	if uidVal, ok := c.Get("user_id"); ok {
		if s, ok := uidVal.(string); ok {
			if id, err := uuid.Parse(s); err == nil {
				createdBy = &id
			}
		}
	}

	campaign, err := ctrl.Service.Create(req, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.SingleCampaignResponse{Data: dto.FromCampaignModel(campaign)})
}

// ListCampaigns godoc
//
//	@Summary		List newsletter campaigns
//	@Description	Admin lists campaigns with optional search and status filter.
//	@Tags			Newsletter Campaigns
//	@Produce		json
//	@Param			page		query		int		false	"Page number (default 1)"			minimum(1)
//	@Param			limit		query		int		false	"Items per page (default 20, max 100)"
//	@Param			search		query		string	false	"Search by subject or title"
//	@Param			status		query		string	false	"Filter by status"	Enums(draft, scheduled, sending, sent, failed, cancelled)
//	@Param			sort_by		query		string	false	"Sort field"		Enums(created_at, updated_at, sent_at, subject)
//	@Param			sort_order	query		string	false	"Sort direction"	Enums(asc, desc)
//	@Success		200			{object}	dto.PaginatedCampaignsResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/fetch/newsletter/campaigns [get]
func (ctrl *NewsletterCampaignController) ListCampaigns(c *gin.Context) {
	var q dto.ListCampaignsQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid query", Message: err.Error()})
		return
	}

	campaigns, total, err := ctrl.Service.List(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	items := make([]dto.CampaignResponse, 0, len(campaigns))
	for i := range campaigns {
		items = append(items, dto.FromCampaignModel(&campaigns[i]))
	}
	totalPages := int((total + int64(q.Limit) - 1) / int64(q.Limit))

	c.JSON(http.StatusOK, dto.PaginatedCampaignsResponse{
		Data:       items,
		Total:      total,
		Page:       q.Page,
		Limit:      q.Limit,
		TotalPages: totalPages,
	})
}

// GetCampaign godoc
//
//	@Summary		Get a campaign by ID
//	@Tags			Newsletter Campaigns
//	@Produce		json
//	@Param			id	path		string	true	"Campaign UUID"	format(uuid)
//	@Success		200	{object}	dto.SingleCampaignResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/newsletter/campaigns/{id} [get]
func (ctrl *NewsletterCampaignController) GetCampaign(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid campaign id"})
		return
	}
	campaign, err := ctrl.Service.GetByID(id)
	if err != nil {
		if errors.Is(err, services.ErrCampaignNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.SingleCampaignResponse{Data: dto.FromCampaignModel(campaign)})
}

// UpdateCampaign godoc
//
//	@Summary		Update a campaign (partial)
//	@Tags			Newsletter Campaigns
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Campaign UUID"	format(uuid)
//	@Param			payload	body		dto.UpdateCampaignRequest	true	"Fields to update"
//	@Success		200		{object}	dto.SingleCampaignResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		404		{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/newsletter/campaigns/{id} [patch]
func (ctrl *NewsletterCampaignController) UpdateCampaign(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid campaign id"})
		return
	}
	var req dto.UpdateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:  "validation failed",
			Fields: map[string]string{"body": err.Error()},
		})
		return
	}
	campaign, err := ctrl.Service.Update(id, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrCampaignNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, services.ErrCampaignAlreadySent):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, dto.SingleCampaignResponse{Data: dto.FromCampaignModel(campaign)})
}

// DeleteCampaign godoc
//
//	@Summary		Delete a campaign
//	@Tags			Newsletter Campaigns
//	@Produce		json
//	@Param			id	path	string	true	"Campaign UUID"	format(uuid)
//	@Success		204	"No Content"
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/newsletter/campaigns/{id} [delete]
func (ctrl *NewsletterCampaignController) DeleteCampaign(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid campaign id"})
		return
	}
	if err := ctrl.Service.Delete(id); err != nil {
		if errors.Is(err, services.ErrCampaignNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// PreviewRecipients godoc
//
//	@Summary		Preview audience size
//	@Description	Returns how many subscribers match the given audience rules, without sending anything.
//	@Tags			Newsletter Campaigns
//	@Produce		json
//	@Param			audience	query		string	true	"Audience"	Enums(all, active, tagged, manual)
//	@Param			tag			query		string	false	"Tag (when audience=tagged)"
//	@Param			manual_ids	query		[]string false	"Manual subscriber UUIDs (when audience=manual)"
//	@Success		200	{object}	map[string]int64
//	@Failure		400	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/newsletter/campaigns/preview [get]
func (ctrl *NewsletterCampaignController) PreviewRecipients(c *gin.Context) {
	audience := c.Query("audience")
	tag := c.Query("tag")
	manualIDs := c.QueryArray("manual_ids")

	count, err := ctrl.Service.PreviewRecipients(audience, tag, manualIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"recipient_count": count})
}

// SendCampaign godoc
//
//	@Summary		Send a campaign
//	@Description	Queues the campaign for immediate delivery to its audience.
//	@Tags			Newsletter Campaigns
//	@Produce		json
//	@Param			id	path		string	true	"Campaign UUID"	format(uuid)
//	@Success		202	{object}	dto.SendCampaignResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		409	{object}	dto.ErrorResponse	"Already sent or no recipients"
//	@Security		BearerAuth
//	@Router			/api/newsletter/campaigns/{id}/send [post]
func (ctrl *NewsletterCampaignController) SendCampaign(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid campaign id"})
		return
	}
	count, err := ctrl.Service.Send(id)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrCampaignNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, services.ErrCampaignAlreadySent):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, services.ErrCampaignNoRecipients):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		}
		return
	}
	c.JSON(http.StatusAccepted, dto.SendCampaignResponse{
		Message:        "Campaign is being sent.",
		RecipientCount: count,
	})
}

// CancelCampaign godoc
//
//	@Summary		Cancel a scheduled campaign
//	@Tags			Newsletter Campaigns
//	@Produce		json
//	@Param			id	path		string	true	"Campaign UUID"	format(uuid)
//	@Success		200	{object}	dto.SingleCampaignResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/newsletter/campaigns/{id}/cancel [post]
func (ctrl *NewsletterCampaignController) CancelCampaign(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid campaign id"})
		return
	}
	campaign, err := ctrl.Service.Cancel(id)
	if err != nil {
		if errors.Is(err, services.ErrCampaignNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.SingleCampaignResponse{Data: dto.FromCampaignModel(campaign)})
}