// dto/inventory_item_dto.go
package dto

import (
	"time"
)

// ImageInput represents an image input
type ImageInput struct {
	ImageURL     string `json:"image_url" binding:"required,url"`
	Caption      string `json:"caption"`
	IsPrimary    bool   `json:"is_primary"`
	DisplayOrder int    `json:"display_order"`
}

// CreateInventoryItemRequest represents the request to create an inventory item
type CreateInventoryItemRequest struct {
	ItemCode            string `json:"item_code" binding:"required,min=2,max=100"`
	ItemName            string `json:"item_name" binding:"required,min=2,max=255"`
	Description         string `json:"description"`
	InventorySectionID  string `json:"inventory_section_id" binding:"required"`

	// Identification
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	SerialNumber string `json:"serial_number"`
	Barcode      string `json:"barcode"`

	// Stock
	QuantityTotal     int `json:"quantity_total" binding:"min=0"`
	QuantityAvailable int `json:"quantity_available" binding:"min=0"`
	QuantityInUse     int `json:"quantity_in_use" binding:"min=0"`
	QuantityDamaged   int `json:"quantity_damaged" binding:"min=0"`

	// Item Properties
	Condition         string `json:"condition" binding:"omitempty,oneof=new good fair poor damaged obsolete expired"`
	Status            string `json:"status" binding:"omitempty,oneof=active inactive archived disposed out_of_stock"`
	Location          string `json:"location"`
	StorageConditions string `json:"storage_conditions"`

	// Chemical-specific fields
	IsChemical         bool   `json:"is_chemical"`
	ChemicalFormula    string `json:"chemical_formula"`
	CASNumber          string `json:"cas_number"`
	HazardClass        string `json:"hazard_class" binding:"omitempty,oneof=none flammable toxic corrosive explosive radioactive biohazard oxidizer"`
	Concentration      string `json:"concentration"`
	Volume             string `json:"volume"`
	SafetyDataSheetURL string `json:"safety_data_sheet_url"`

	// Safety & Compliance
	RequiresSupervision  bool   `json:"requires_supervision"`
	SafetyInstructions   string `json:"safety_instructions"`
	DisposalInstructions string `json:"disposal_instructions"`
	RestrictedAccess     bool   `json:"restricted_access"`
	AllowedRoles         string `json:"allowed_roles"`

	// Media
	Images []ImageInput `json:"images"`
}

// UpdateInventoryItemRequest represents the request to update an inventory item
type UpdateInventoryItemRequest struct {
	ItemCode            string `json:"item_code" binding:"omitempty,min=2,max=100"`
	ItemName            string `json:"item_name" binding:"omitempty,min=2,max=255"`
	Description         string `json:"description"`
	InventorySectionID  string `json:"inventory_section_id"`

	// Identification
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	SerialNumber string `json:"serial_number"`
	Barcode      string `json:"barcode"`

	// Stock
	QuantityTotal     *int `json:"quantity_total"`
	QuantityAvailable *int `json:"quantity_available"`
	QuantityInUse     *int `json:"quantity_in_use"`
	QuantityDamaged   *int `json:"quantity_damaged"`

	// Item Properties
	Condition         string `json:"condition" binding:"omitempty,oneof=new good fair poor damaged obsolete expired"`
	Status            string `json:"status" binding:"omitempty,oneof=active inactive archived disposed out_of_stock"`
	Location          string `json:"location"`
	StorageConditions string `json:"storage_conditions"`

	// Chemical-specific fields
	IsChemical         *bool  `json:"is_chemical"`
	ChemicalFormula    string `json:"chemical_formula"`
	CASNumber          string `json:"cas_number"`
	HazardClass        string `json:"hazard_class" binding:"omitempty,oneof=none flammable toxic corrosive explosive radioactive biohazard oxidizer"`
	Concentration      string `json:"concentration"`
	Volume             string `json:"volume"`
	SafetyDataSheetURL string `json:"safety_data_sheet_url"`

	// Safety & Compliance
	RequiresSupervision  *bool  `json:"requires_supervision"`
	SafetyInstructions   string `json:"safety_instructions"`
	DisposalInstructions string `json:"disposal_instructions"`
	RestrictedAccess     *bool  `json:"restricted_access"`
	AllowedRoles         string `json:"allowed_roles"`

	// Media
	Images []ImageInput `json:"images"`
}

// ImageResponse represents an image in the response
type ImageResponse struct {
	ID           string `json:"id"`
	ImageURL     string `json:"image_url"`
	Caption      string `json:"caption"`
	IsPrimary    bool   `json:"is_primary"`
	DisplayOrder int    `json:"display_order"`
}

// InventoryItemResponse represents the response for an inventory item
type InventoryItemResponse struct {
	ID          string `json:"id"`
	ItemCode    string `json:"item_code"`
	ItemName    string `json:"item_name"`
	Description string `json:"description"`

	InventorySectionID string `json:"inventory_section_id"`

	// Identification
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	SerialNumber string `json:"serial_number"`
	Barcode      string `json:"barcode"`

	// Stock
	QuantityTotal     int `json:"quantity_total"`
	QuantityAvailable int `json:"quantity_available"`
	QuantityInUse     int `json:"quantity_in_use"`
	QuantityDamaged   int `json:"quantity_damaged"`

	// Item Properties
	Condition         string `json:"condition"`
	Status            string `json:"status"`
	Location          string `json:"location"`
	StorageConditions string `json:"storage_conditions"`

	// Chemical-specific fields
	IsChemical         bool   `json:"is_chemical"`
	ChemicalFormula    string `json:"chemical_formula"`
	CASNumber          string `json:"cas_number"`
	HazardClass        string `json:"hazard_class"`
	Concentration      string `json:"concentration"`
	Volume             string `json:"volume"`
	SafetyDataSheetURL string `json:"safety_data_sheet_url"`

	// Safety & Compliance
	RequiresSupervision  bool   `json:"requires_supervision"`
	SafetyInstructions   string `json:"safety_instructions"`
	DisposalInstructions string `json:"disposal_instructions"`
	RestrictedAccess     bool   `json:"restricted_access"`
	AllowedRoles         string `json:"allowed_roles"`

	// Media
	Images []ImageResponse `json:"images"`

	// Audit
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Nested relationships
	InventorySection *InventorySectionResponse `json:"inventory_section,omitempty"`
	Creator          *UserResponse             `json:"creator,omitempty"`
}

// InventoryItemListResponse represents a paginated list of inventory items
type InventoryItemListResponse struct {
	Items      []InventoryItemResponse `json:"items"`
	Total      int64                   `json:"total"`
	Page       int                     `json:"page"`
	Limit      int                     `json:"limit"`
	TotalPages int                     `json:"total_pages"`
}

// InventoryItemQueryParams represents query parameters for filtering items
type InventoryItemQueryParams struct {
	InventorySectionID string `form:"inventory_section_id"`
	Search             string `form:"search"`
	Category           string `form:"category"`
	Condition          string `form:"condition" binding:"omitempty,oneof=new good fair poor damaged obsolete expired"`
	Status             string `form:"status" binding:"omitempty,oneof=active inactive archived disposed out_of_stock"`
	IsChemical         *bool  `form:"is_chemical"`
	HazardClass        string `form:"hazard_class" binding:"omitempty,oneof=none flammable toxic corrosive explosive radioactive biohazard oxidizer"`
	Location           string `form:"location"`
	LowStock           *bool  `form:"low_stock"`
	Page               int    `form:"page" default:"1"`
	Limit              int    `form:"limit" default:"20"`
	SortBy             string `form:"sort_by" default:"created_at"`
	SortOrder          string `form:"sort_order" default:"desc" binding:"omitempty,oneof=asc desc"`
}

// InventoryItemStats represents statistics for inventory items
type InventoryItemStats struct {
	TotalItems          int64   `json:"total_items"`
	ActiveItems         int64   `json:"active_items"`
	InactiveItems       int64   `json:"inactive_items"`
	OutOfStockItems     int64   `json:"out_of_stock_items"`
	TotalQuantity       int64   `json:"total_quantity"`
	AvailableQuantity   int64   `json:"available_quantity"`
	DamagedQuantity     int64   `json:"damaged_quantity"`
	ChemicalItems       int64   `json:"chemical_items"`
	RestrictedItems     int64   `json:"restricted_items"`
}