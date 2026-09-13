// models/inventory_item.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryItemCondition string

const (
	InventoryConditionNew      InventoryItemCondition = "new"
	InventoryConditionGood     InventoryItemCondition = "good"
	InventoryConditionFair     InventoryItemCondition = "fair"
	InventoryConditionPoor     InventoryItemCondition = "poor"
	InventoryConditionDamaged  InventoryItemCondition = "damaged"
	InventoryConditionObsolete InventoryItemCondition = "obsolete"
	InventoryConditionExpired  InventoryItemCondition = "expired"
)

type InventoryItemStatus string

const (
	InventoryItemStatusActive     InventoryItemStatus = "active"
	InventoryItemStatusInactive   InventoryItemStatus = "inactive"
	InventoryItemStatusArchived   InventoryItemStatus = "archived"
	InventoryItemStatusDisposed   InventoryItemStatus = "disposed"
	InventoryItemStatusOutOfStock InventoryItemStatus = "out_of_stock"
)

type InventoryItemHazardClass string

const (
	HazardNone        InventoryItemHazardClass = "none"
	HazardFlammable   InventoryItemHazardClass = "flammable"
	HazardToxic       InventoryItemHazardClass = "toxic"
	HazardCorrosive   InventoryItemHazardClass = "corrosive"
	HazardExplosive   InventoryItemHazardClass = "explosive"
	HazardRadioactive InventoryItemHazardClass = "radioactive"
	HazardBiohazard   InventoryItemHazardClass = "biohazard"
	HazardOxidizer    InventoryItemHazardClass = "oxidizer"
)

type InventoryItem struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ItemCode    string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"item_code"`
	ItemName    string    `gorm:"type:varchar(255);not null" json:"item_name"`
	Description string    `gorm:"type:text" json:"description"`

	InventorySectionID uuid.UUID `gorm:"type:uuid;not null;index" json:"inventory_section_id"`

	// Identification
	Brand        string `gorm:"type:varchar(255)" json:"brand"`
	Model        string `gorm:"type:varchar(255)" json:"model"`
	SerialNumber string `gorm:"type:varchar(255)" json:"serial_number"`
	Barcode      string `gorm:"type:varchar(255);index" json:"barcode"`

	// Stock
	QuantityTotal     int `gorm:"default:0;check:quantity_total >= 0" json:"quantity_total"`
	QuantityAvailable int `gorm:"default:0;check:quantity_available >= 0" json:"quantity_available"`
	QuantityInUse     int `gorm:"default:0;check:quantity_in_use >= 0" json:"quantity_in_use"`
	QuantityDamaged   int `gorm:"default:0;check:quantity_damaged >= 0" json:"quantity_damaged"`

	// Item Properties
	Condition         InventoryItemCondition `gorm:"type:varchar(20);default:'good';check:condition IN ('new','good','fair','poor','damaged','obsolete','expired')" json:"condition"`
	Status            InventoryItemStatus    `gorm:"type:varchar(20);default:'active';check:status IN ('active','inactive','archived','disposed','out_of_stock')" json:"status"`
	Location          string                 `gorm:"type:varchar(255)" json:"location"`
	StorageConditions string                 `gorm:"type:varchar(255)" json:"storage_conditions"`

	// Chemical-specific fields
	IsChemical         bool                     `gorm:"default:false" json:"is_chemical"`
	ChemicalFormula    string                   `gorm:"type:varchar(255)" json:"chemical_formula"`
	CASNumber          string                   `gorm:"type:varchar(50)" json:"cas_number"`
	HazardClass        InventoryItemHazardClass `gorm:"type:varchar(50);default:'none';check:hazard_class IN ('none','flammable','toxic','corrosive','explosive','radioactive','biohazard','oxidizer')" json:"hazard_class"`
	Concentration      string                   `gorm:"type:varchar(100)" json:"concentration"`
	Volume             string                   `gorm:"type:varchar(100)" json:"volume"`
	SafetyDataSheetURL string                   `gorm:"type:varchar(500)" json:"safety_data_sheet_url"`

	// Safety & Compliance
	RequiresSupervision  bool   `gorm:"default:false" json:"requires_supervision"`
	SafetyInstructions   string `gorm:"type:text" json:"safety_instructions"`
	DisposalInstructions string `gorm:"type:text" json:"disposal_instructions"`
	RestrictedAccess     bool   `gorm:"default:false" json:"restricted_access"`
	AllowedRoles         string `gorm:"type:text" json:"allowed_roles"`

	// Media
	Images []InventoryItemImage `gorm:"foreignKey:InventoryItemID;constraint:OnDelete:CASCADE" json:"images,omitempty"`

	// Audit
	CreatedBy uuid.UUID      `gorm:"type:uuid;not null;index" json:"created_by"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	InventorySection InventorySection        `gorm:"foreignKey:InventorySectionID" json:"inventory_section,omitempty"`
	Creator          User                    `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Transactions     []InventoryTransaction  `gorm:"foreignKey:ItemID" json:"transactions,omitempty"`
}

func (InventoryItem) TableName() string {
	return "inventory_items"
}

