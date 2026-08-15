// dto/address_dto.go
package dto

import (
	"time"
)

// CreateAddressRequest represents the request body for creating an address
type CreateAddressRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	Address     string `json:"address" binding:"required"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	PostalCode  string `json:"postal_code"`
	AddressType string `json:"address_type" binding:"omitempty,oneof=home school office other"`
	IsPrimary   bool   `json:"is_primary"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// UpdateAddressRequest represents the request body for updating an address
type UpdateAddressRequest struct {
	Address     string `json:"address"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	PostalCode  string `json:"postal_code"`
	AddressType string `json:"address_type" binding:"omitempty,oneof=home school office other"`
	IsPrimary   *bool  `json:"is_primary"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// AddressResponse represents the address response
type AddressResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	State       string    `json:"state"`
	Country     string    `json:"country"`
	PostalCode  string    `json:"postal_code"`
	AddressType string    `json:"address_type"`
	IsPrimary   bool      `json:"is_primary"`
	Status      string    `json:"status"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Nested user details
	User *UserResponse `json:"user,omitempty"`
}

// AddressListResponse represents paginated address list response
type AddressListResponse struct {
	Addresses  []AddressResponse `json:"addresses"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// AddressQueryParams represents query parameters for filtering addresses
type AddressQueryParams struct {
	Search       string `form:"search"`
	UserID       string `form:"user_id"`
	AddressType  string `form:"address_type" binding:"omitempty,oneof=home school office other"`
	Status       string `form:"status" binding:"omitempty,oneof=active inactive"`
	IsPrimary    *bool  `form:"is_primary"`
	Page         int    `form:"page" binding:"min=1"`
	Limit        int    `form:"limit" binding:"min=1,max=100"`
	SortBy       string `form:"sort_by" binding:"omitempty,oneof=address city state country address_type created_at"`
	SortOrder    string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}