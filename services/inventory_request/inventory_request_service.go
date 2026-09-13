// services/inventory_request_service.go
package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"crm-go/dto"
	"crm-go/models"
)

type InventoryRequestService struct {
	db *gorm.DB
}

func NewInventoryRequestService(db *gorm.DB) *InventoryRequestService {
	return &InventoryRequestService{db: db}
}

// generateRequestCode generates a unique request code
func (s *InventoryRequestService) generateRequestCode() (string, error) {
	prefix := "REQ-" + time.Now().Format("20060102") + "-"

	// Get the count of requests today
	var count int64
	today := time.Now().Format("2006-01-02")
	s.db.Model(&models.InventoryRequest{}).
		Where("DATE(created_at) = ?", today).
		Count(&count)

	return fmt.Sprintf("%s%04d", prefix, count+1), nil
}

// CreateInventoryRequest creates a new inventory request
func (s *InventoryRequestService) CreateInventoryRequest(userID uuid.UUID, req *dto.CreateInventoryRequestRequest) (*dto.InventoryRequestResponse, error) {
	// Parse section ID
	sectionID, err := uuid.Parse(req.InventorySectionID)
	if err != nil {
		return nil, errors.New("invalid inventory section ID")
	}

	// Parse item ID
	itemID, err := uuid.Parse(req.InventoryItemID)
	if err != nil {
		return nil, errors.New("invalid inventory item ID")
	}

	// Verify section exists
	var section models.InventorySection
	if err := s.db.Where("id = ? AND deleted_at IS NULL", sectionID).First(&section).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory section not found")
		}
		return nil, fmt.Errorf("failed to verify section: %w", err)
	}

	// Verify item exists
	var item models.InventoryItem
	if err := s.db.Where("id = ? AND deleted_at IS NULL", itemID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory item not found")
		}
		return nil, fmt.Errorf("failed to verify item: %w", err)
	}

	// Verify item belongs to section
	if item.InventorySectionID != sectionID {
		return nil, errors.New("inventory item does not belong to the specified section")
	}

	// Check if requested quantity is available
	if req.QuantityRequested > item.QuantityAvailable {
		return nil, fmt.Errorf("requested quantity (%d) exceeds available quantity (%d)", req.QuantityRequested, item.QuantityAvailable)
	}

	// Check for existing pending request for the same item by the same user
	var existingRequest models.InventoryRequest
	if err := s.db.Where(
		"requested_by_id = ? AND inventory_item_id = ? AND status = ? AND deleted_at IS NULL",
		userID, itemID, "pending",
	).First(&existingRequest).Error; err == nil {
		return nil, errors.New("you already have a pending request for this item")
	}

	// Generate request code
	requestCode, err := s.generateRequestCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate request code: %w", err)
	}

	// Parse needed by date
	var neededByDate *time.Time
	if req.NeededByDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.NeededByDate)
		if err != nil {
			return nil, errors.New("invalid needed by date format. Use YYYY-MM-DD")
		}
		neededByDate = &parsedDate
	}

	// Create request
	request := &models.InventoryRequest{
		ID:                 uuid.New(),
		RequestCode:        requestCode,
		RequestedByID:      userID,
		InventorySectionID: sectionID,
		InventoryItemID:    itemID,
		QuantityRequested:  req.QuantityRequested,
		QuantityApproved:   0,
		QuantityIssued:     0,
		Purpose:            req.Purpose,
		Status:             models.InventoryRequestStatusPending,
		NeededByDate:       neededByDate,
		Notes:              req.Notes,
		CreatedBy:          userID,
	}

	if err := s.db.Create(request).Error; err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Preload relationships
	if err := s.db.Preload("RequestedBy").
		Preload("InventorySection").
		Preload("InventoryItem").
		Preload("Creator").
		First(request, request.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load request details: %w", err)
	}

	return s.toRequestResponse(request), nil
}

// GetInventoryRequestByID retrieves a request by ID
func (s *InventoryRequestService) GetInventoryRequestByID(id string) (*dto.InventoryRequestResponse, error) {
	requestID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid request ID")
	}

	var request models.InventoryRequest
	if err := s.db.Where("id = ? AND deleted_at IS NULL", requestID).
		Preload("RequestedBy").
		Preload("InventorySection").
		Preload("InventoryItem").
		Preload("ApprovedBy").
		Preload("Creator").
		First(&request).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory request not found")
		}
		return nil, fmt.Errorf("failed to fetch request: %w", err)
	}

	return s.toRequestResponse(&request), nil
}

// GetInventoryRequests retrieves requests with filters and pagination
func (s *InventoryRequestService) GetInventoryRequests(params *dto.InventoryRequestQueryParams) (*dto.InventoryRequestListResponse, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}

	// Build query
	query := s.db.Model(&models.InventoryRequest{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.RequestedByID != "" {
		userID, err := uuid.Parse(params.RequestedByID)
		if err == nil {
			query = query.Where("requested_by_id = ?", userID)
		}
	}

	if params.InventorySectionID != "" {
		sectionID, err := uuid.Parse(params.InventorySectionID)
		if err == nil {
			query = query.Where("inventory_section_id = ?", sectionID)
		}
	}

	if params.InventoryItemID != "" {
		itemID, err := uuid.Parse(params.InventoryItemID)
		if err == nil {
			query = query.Where("inventory_item_id = ?", itemID)
		}
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("request_code ILIKE ? OR purpose ILIKE ? OR notes ILIKE ?", search, search, search)
	}

	if params.DateFrom != "" {
		dateFrom, err := time.Parse("2006-01-02", params.DateFrom)
		if err == nil {
			query = query.Where("created_at >= ?", dateFrom)
		}
	}

	if params.DateTo != "" {
		dateTo, err := time.Parse("2006-01-02", params.DateTo)
		if err == nil {
			query = query.Where("created_at <= ?", dateTo.Add(24*time.Hour))
		}
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count requests: %w", err)
	}

	// Apply sorting
	sortDirection := "DESC"
	if strings.ToLower(params.SortOrder) == "asc" {
		sortDirection = "ASC"
	}
	query = query.Order(fmt.Sprintf("%s %s", params.SortBy, sortDirection))

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	// Execute query
	var requests []models.InventoryRequest
	if err := query.
		Preload("RequestedBy").
		Preload("InventorySection").
		Preload("InventoryItem").
		Preload("ApprovedBy").
		Preload("Creator").
		Find(&requests).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch requests: %w", err)
	}

	// Convert to response
	responses := make([]dto.InventoryRequestResponse, len(requests))
	for i, request := range requests {
		responses[i] = *s.toRequestResponse(&request)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.InventoryRequestListResponse{
		Requests:   responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetMyRequests retrieves requests for the current user
func (s *InventoryRequestService) GetMyRequests(userID uuid.UUID, params *dto.InventoryRequestQueryParams) (*dto.InventoryRequestListResponse, error) {
	params.RequestedByID = userID.String()
	return s.GetInventoryRequests(params)
}

// UpdateInventoryRequest updates a pending request (only by requester)
func (s *InventoryRequestService) UpdateInventoryRequest(id string, userID uuid.UUID, req *dto.UpdateInventoryRequestRequest) (*dto.InventoryRequestResponse, error) {
	requestID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid request ID")
	}

	var request models.InventoryRequest
	if err := s.db.Where("id = ? AND deleted_at IS NULL", requestID).First(&request).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory request not found")
		}
		return nil, fmt.Errorf("failed to fetch request: %w", err)
	}

	// Only the requester can update
	if request.RequestedByID != userID {
		return nil, errors.New("you can only update your own requests")
	}

	// Only pending requests can be updated
	if request.Status != models.InventoryRequestStatusPending {
		return nil, errors.New("only pending requests can be updated")
	}

	if req.QuantityRequested != nil {
		// Verify item availability
		var item models.InventoryItem
		if err := s.db.Where("id = ? AND deleted_at IS NULL", request.InventoryItemID).First(&item).Error; err != nil {
			return nil, errors.New("inventory item not found")
		}
		if *req.QuantityRequested > item.QuantityAvailable {
			return nil, fmt.Errorf("requested quantity exceeds available quantity (%d)", item.QuantityAvailable)
		}
		request.QuantityRequested = *req.QuantityRequested
	}

	if req.Purpose != "" {
		request.Purpose = req.Purpose
	}

	if req.NeededByDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.NeededByDate)
		if err != nil {
			return nil, errors.New("invalid needed by date format")
		}
		request.NeededByDate = &parsedDate
	}

	if req.Notes != "" {
		request.Notes = req.Notes
	}

	if err := s.db.Save(&request).Error; err != nil {
		return nil, fmt.Errorf("failed to update request: %w", err)
	}

	// Reload with relationships
	if err := s.db.Preload("RequestedBy").
		Preload("InventorySection").
		Preload("InventoryItem").
		Preload("ApprovedBy").
		Preload("Creator").
		First(&request, request.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load request details: %w", err)
	}

	return s.toRequestResponse(&request), nil
}

// ApproveInventoryRequest approves a request (admin only)
func (s *InventoryRequestService) ApproveInventoryRequest(id string, approverID uuid.UUID, req *dto.ApproveInventoryRequestRequest) (*dto.InventoryRequestResponse, error) {
	requestID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid request ID")
	}

	var request models.InventoryRequest
	if err := s.db.Where("id = ? AND deleted_at IS NULL", requestID).First(&request).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory request not found")
		}
		return nil, fmt.Errorf("failed to fetch request: %w", err)
	}

	// Only pending requests can be approved
	if request.Status != models.InventoryRequestStatusPending {
		return nil, errors.New("only pending requests can be approved")
	}

	// Verify requested quantity doesn't exceed available stock
	var item models.InventoryItem
	if err := s.db.Where("id = ? AND deleted_at IS NULL", request.InventoryItemID).First(&item).Error; err != nil {
		return nil, errors.New("inventory item not found")
	}

	if req.QuantityApproved > item.QuantityAvailable {
		return nil, fmt.Errorf("approved quantity (%d) exceeds available quantity (%d)", req.QuantityApproved, item.QuantityAvailable)
	}

	if req.QuantityApproved > request.QuantityRequested {
		return nil, errors.New("approved quantity cannot exceed requested quantity")
	}

	now := time.Now()
	request.Status = models.InventoryRequestStatusApproved
	request.QuantityApproved = req.QuantityApproved
	request.ApprovedByID = &approverID
	request.ApprovedAt = &now

	if req.Notes != "" {
		request.Notes = req.Notes
	}

	if err := s.db.Save(&request).Error; err != nil {
		return nil, fmt.Errorf("failed to approve request: %w", err)
	}

	// Reload with relationships
	if err := s.db.Preload("RequestedBy").
		Preload("InventorySection").
		Preload("InventoryItem").
		Preload("ApprovedBy").
		Preload("Creator").
		First(&request, request.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load request details: %w", err)
	}

	return s.toRequestResponse(&request), nil
}

// RejectInventoryRequest rejects a request (admin only)
func (s *InventoryRequestService) RejectInventoryRequest(id string, rejectorID uuid.UUID, req *dto.RejectInventoryRequestRequest) (*dto.InventoryRequestResponse, error) {
	requestID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid request ID")
	}

	var request models.InventoryRequest
	if err := s.db.Where("id = ? AND deleted_at IS NULL", requestID).First(&request).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory request not found")
		}
		return nil, fmt.Errorf("failed to fetch request: %w", err)
	}

	// Only pending requests can be rejected
	if request.Status != models.InventoryRequestStatusPending {
		return nil, errors.New("only pending requests can be rejected")
	}

	now := time.Now()
	request.Status = models.InventoryRequestStatusRejected
	request.ApprovedByID = &rejectorID
	request.ApprovedAt = &now
	request.RejectionReason = req.RejectionReason

	if err := s.db.Save(&request).Error; err != nil {
		return nil, fmt.Errorf("failed to reject request: %w", err)
	}

	// Reload with relationships
	if err := s.db.Preload("RequestedBy").
		Preload("InventorySection").
		Preload("InventoryItem").
		Preload("ApprovedBy").
		Preload("Creator").
		First(&request, request.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load request details: %w", err)
	}

	return s.toRequestResponse(&request), nil
}

// FulfillInventoryRequest fulfills a request and reduces stock
func (s *InventoryRequestService) FulfillInventoryRequest(id string, fulfillerID uuid.UUID, req *dto.FulfillInventoryRequestRequest) (*dto.InventoryRequestResponse, error) {
	requestID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid request ID")
	}

	var request models.InventoryRequest
	if err := s.db.Where("id = ? AND deleted_at IS NULL", requestID).First(&request).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory request not found")
		}
		return nil, fmt.Errorf("failed to fetch request: %w", err)
	}

	// Only approved requests can be fulfilled
	if request.Status != models.InventoryRequestStatusApproved {
		return nil, errors.New("only approved requests can be fulfilled")
	}

	if req.QuantityIssued > request.QuantityApproved {
		return nil, errors.New("issued quantity cannot exceed approved quantity")
	}

	// Begin transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Verify item and check stock
	var item models.InventoryItem
	if err := tx.Where("id = ? AND deleted_at IS NULL", request.InventoryItemID).First(&item).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("inventory item not found")
	}

	if req.QuantityIssued > item.QuantityAvailable {
		tx.Rollback()
		return nil, fmt.Errorf("issued quantity (%d) exceeds available quantity (%d)", req.QuantityIssued, item.QuantityAvailable)
	}

	// Update item stock
	item.QuantityAvailable -= req.QuantityIssued
	item.QuantityInUse += req.QuantityIssued

	// Auto-update status if out of stock
	if item.QuantityAvailable == 0 {
		item.Status = models.InventoryItemStatusOutOfStock
	}

	if err := tx.Save(&item).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update item stock: %w", err)
	}

	// Update request
	now := time.Now()
	request.Status = models.InventoryRequestStatusFulfilled
	request.QuantityIssued = req.QuantityIssued
	request.FulfilledAt = &now

	if req.Notes != "" {
		request.Notes = req.Notes
	}

	if err := tx.Save(&request).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update request: %w", err)
	}

	tx.Commit()

	// Reload with relationships
	if err := s.db.Preload("RequestedBy").
		Preload("InventorySection").
		Preload("InventoryItem").
		Preload("ApprovedBy").
		Preload("Creator").
		First(&request, request.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load request details: %w", err)
	}

	return s.toRequestResponse(&request), nil
}

// CancelInventoryRequest cancels a request
func (s *InventoryRequestService) CancelInventoryRequest(id string, userID uuid.UUID) error {
	requestID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid request ID")
	}

	var request models.InventoryRequest
	if err := s.db.Where("id = ? AND deleted_at IS NULL", requestID).First(&request).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("inventory request not found")
		}
		return fmt.Errorf("failed to fetch request: %w", err)
	}

	// Only the requester can cancel their own request
	if request.RequestedByID != userID {
		return errors.New("you can only cancel your own requests")
	}

	// Only pending or approved requests can be cancelled
	if request.Status != models.InventoryRequestStatusPending && request.Status != models.InventoryRequestStatusApproved {
		return errors.New("only pending or approved requests can be cancelled")
	}

	request.Status = models.InventoryRequestStatusCancelled
	if err := s.db.Save(&request).Error; err != nil {
		return fmt.Errorf("failed to cancel request: %w", err)
	}

	return nil
}

// DeleteInventoryRequest soft deletes a request
func (s *InventoryRequestService) DeleteInventoryRequest(id string) error {
	requestID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid request ID")
	}

	var request models.InventoryRequest
	if err := s.db.Where("id = ? AND deleted_at IS NULL", requestID).First(&request).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("inventory request not found")
		}
		return fmt.Errorf("failed to fetch request: %w", err)
	}

	if err := s.db.Delete(&request).Error; err != nil {
		return fmt.Errorf("failed to delete request: %w", err)
	}

	return nil
}

// GetInventoryRequestStats retrieves statistics for inventory requests
func (s *InventoryRequestService) GetInventoryRequestStats(sectionID string) (*dto.InventoryRequestStats, error) {
	query := s.db.Model(&models.InventoryRequest{}).Where("deleted_at IS NULL")

	if sectionID != "" {
		sID, err := uuid.Parse(sectionID)
		if err != nil {
			return nil, errors.New("invalid section ID")
		}
		query = query.Where("inventory_section_id = ?", sID)
	}

	var stats dto.InventoryRequestStats

	// Total
	if err := query.Count(&stats.TotalRequests).Error; err != nil {
		return nil, fmt.Errorf("failed to count total requests: %w", err)
	}

	// By status
	if err := query.Where("status = ?", "pending").Count(&stats.PendingRequests).Error; err != nil {
		return nil, fmt.Errorf("failed to count pending requests: %w", err)
	}

	if err := query.Where("status = ?", "approved").Count(&stats.ApprovedRequests).Error; err != nil {
		return nil, fmt.Errorf("failed to count approved requests: %w", err)
	}

	if err := query.Where("status = ?", "rejected").Count(&stats.RejectedRequests).Error; err != nil {
		return nil, fmt.Errorf("failed to count rejected requests: %w", err)
	}

	if err := query.Where("status = ?", "fulfilled").Count(&stats.FulfilledRequests).Error; err != nil {
		return nil, fmt.Errorf("failed to count fulfilled requests: %w", err)
	}

	if err := query.Where("status = ?", "cancelled").Count(&stats.CancelledRequests).Error; err != nil {
		return nil, fmt.Errorf("failed to count cancelled requests: %w", err)
	}

	return &stats, nil
}

// toRequestResponse converts model to response DTO
func (s *InventoryRequestService) toRequestResponse(request *models.InventoryRequest) *dto.InventoryRequestResponse {
	response := &dto.InventoryRequestResponse{
		ID:                 request.ID.String(),
		RequestCode:        request.RequestCode,
		RequestedByID:      request.RequestedByID.String(),
		InventorySectionID: request.InventorySectionID.String(),
		InventoryItemID:    request.InventoryItemID.String(),
		QuantityRequested:  request.QuantityRequested,
		QuantityApproved:   request.QuantityApproved,
		QuantityIssued:     request.QuantityIssued,
		Purpose:            request.Purpose,
		Status:             string(request.Status),
		RejectionReason:    request.RejectionReason,
		NeededByDate:       request.NeededByDate,
		ApprovedAt:         request.ApprovedAt,
		FulfilledAt:        request.FulfilledAt,
		Notes:              request.Notes,
		CreatedBy:          request.CreatedBy.String(),
		CreatedAt:          request.CreatedAt,
		UpdatedAt:          request.UpdatedAt,
	}

	if request.ApprovedByID != nil {
		response.ApprovedByID = request.ApprovedByID.String()
	}

	// Add requested by details
	if request.RequestedBy.ID != uuid.Nil {
		response.RequestedBy = &dto.UserResponse{
			ID:        request.RequestedBy.ID.String(),
			FirstName: request.RequestedBy.FirstName,
			LastName:  request.RequestedBy.LastName,
			Email:     request.RequestedBy.Email,
			Phone:     request.RequestedBy.Phone,
			Role:      request.RequestedBy.Role,
			Position:  request.RequestedBy.Position,
		}
	}

	// Add section details
	if request.InventorySection.ID != uuid.Nil {
		response.InventorySection = &dto.InventorySectionResponse{
			ID:          request.InventorySection.ID.String(),
			Name:        request.InventorySection.Name,
			Code:        request.InventorySection.Code,
			Description: request.InventorySection.Description,
			Location:    request.InventorySection.Location,
			Status:      string(request.InventorySection.Status),
			Notes:       request.InventorySection.Notes,
			CreatedAt:   request.InventorySection.CreatedAt,
			UpdatedAt:   request.InventorySection.UpdatedAt,
		}
	}

	// Add item details
	if request.InventoryItem.ID != uuid.Nil {
		response.InventoryItem = &dto.InventoryItemResponse{
			ID:               request.InventoryItem.ID.String(),
			ItemCode:         request.InventoryItem.ItemCode,
			ItemName:         request.InventoryItem.ItemName,
			Description:      request.InventoryItem.Description,
			QuantityTotal:    request.InventoryItem.QuantityTotal,
			QuantityAvailable: request.InventoryItem.QuantityAvailable,
			Status:           string(request.InventoryItem.Status),
			Condition:        string(request.InventoryItem.Condition),
		}
	}

	// Add approved by details
	if request.ApprovedBy != nil && request.ApprovedBy.ID != uuid.Nil {
		response.ApprovedBy = &dto.UserResponse{
			ID:        request.ApprovedBy.ID.String(),
			FirstName: request.ApprovedBy.FirstName,
			LastName:  request.ApprovedBy.LastName,
			Email:     request.ApprovedBy.Email,
			Phone:     request.ApprovedBy.Phone,
			Role:      request.ApprovedBy.Role,
			Position:  request.ApprovedBy.Position,
		}
	}

	// Add creator details
	if request.Creator.ID != uuid.Nil {
		response.Creator = &dto.UserResponse{
			ID:        request.Creator.ID.String(),
			FirstName: request.Creator.FirstName,
			LastName:  request.Creator.LastName,
			Email:     request.Creator.Email,
			Phone:     request.Creator.Phone,
			Role:      request.Creator.Role,
			Position:  request.Creator.Position,
		}
	}

	return response
}