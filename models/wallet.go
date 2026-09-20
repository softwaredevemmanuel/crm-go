package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletStatus string

const (
	WalletStatusActive   WalletStatus = "active"
	WalletStatusInactive WalletStatus = "inactive"
	WalletStatusFrozen   WalletStatus = "frozen"
)

type Wallet struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	Name string `gorm:"type:varchar(255);not null" json:"name"`

	Currency string `gorm:"type:varchar(10);not null;default:'NGN'" json:"currency"`

	// Store money in the smallest currency unit.
	// For NGN, ₦1 = 100 kobo.
	Balance int64 `gorm:"not null;default:0" json:"balance"`

	AvailableBalance int64 `gorm:"not null;default:0" json:"available_balance"`

	Status WalletStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Transactions []WalletTransaction `gorm:"foreignKey:WalletID" json:"transactions,omitempty"`
}

func (Wallet) TableName() string {
	return "wallets"
}

type WalletTransactionType string

const (
	WalletTransactionCredit WalletTransactionType = "credit"
	WalletTransactionDebit  WalletTransactionType = "debit"
)

type WalletTransactionStatus string

const (
	WalletTransactionPending   WalletTransactionStatus = "pending"
	WalletTransactionCompleted WalletTransactionStatus = "completed"
	WalletTransactionFailed    WalletTransactionStatus = "failed"
	WalletTransactionReversed  WalletTransactionStatus = "reversed"
)

type WalletTransaction struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	WalletID uuid.UUID `gorm:"type:uuid;not null;index" json:"wallet_id"`

	Type WalletTransactionType `gorm:"type:varchar(20);not null" json:"type"`

	Amount int64 `gorm:"not null;check:amount > 0" json:"amount"`

	BalanceBefore int64 `gorm:"not null" json:"balance_before"`

	BalanceAfter int64 `gorm:"not null" json:"balance_after"`

	Reference string `gorm:"type:varchar(100);uniqueIndex;not null" json:"reference"`

	Description string `gorm:"type:text" json:"description"`

	Status WalletTransactionStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`

	// Optional link to the person/admin who initiated it
	CreatedBy uuid.UUID `gorm:"type:uuid;index" json:"created_by"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Wallet Wallet `gorm:"foreignKey:WalletID" json:"wallet,omitempty"`
	User   User   `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
}

func (WalletTransaction) TableName() string {
	return "wallet_transactions"
}

type Beneficiary struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	AccountName string `gorm:"type:varchar(255);not null" json:"account_name"`

	AccountNumber string `gorm:"type:varchar(20);not null" json:"account_number"`

	AccountBank string `gorm:"type:varchar(20);not null" json:"account_bank"`

	BankName string `gorm:"type:varchar(255)" json:"bank_name"`

	// Optional provider identifier
	ProviderRecipientCode string `gorm:"type:varchar(255)" json:"provider_recipient_code"`

	IsVerified bool `gorm:"default:false" json:"is_verified"`

	CreatedBy uuid.UUID `gorm:"type:uuid;index" json:"created_by"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
}

func (Beneficiary) TableName() string {
	return "wallet_beneficiaries"
}

type DisbursementStatus string

const (
	DisbursementPending    DisbursementStatus = "pending"
	DisbursementProcessing DisbursementStatus = "processing"
	DisbursementSuccessful DisbursementStatus = "successful"
	DisbursementFailed     DisbursementStatus = "failed"
	DisbursementReversed   DisbursementStatus = "reversed"
)

type WalletDisbursement struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	WalletID uuid.UUID `gorm:"type:uuid;not null;index" json:"wallet_id"`

	BeneficiaryID uuid.UUID `gorm:"type:uuid;not null;index" json:"beneficiary_id"`

	Amount int64 `gorm:"not null;check:amount > 0" json:"amount"`

	Fee int64 `gorm:"not null;default:0" json:"fee"`

	TotalAmount int64 `gorm:"not null" json:"total_amount"`

	Reference string `gorm:"type:varchar(100);uniqueIndex;not null" json:"reference"`

	ProviderReference string `gorm:"type:varchar(255);index" json:"provider_reference"`

	Status DisbursementStatus `gorm:"type:varchar(30);not null;default:'pending'" json:"status"`

	Reason string `gorm:"type:text" json:"reason"`

	InitiatedBy uuid.UUID `gorm:"type:uuid;not null;index" json:"initiated_by"`

	ApprovedBy *uuid.UUID `gorm:"type:uuid;index" json:"approved_by,omitempty"`

	ApprovedAt *time.Time `json:"approved_at,omitempty"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Wallet      Wallet      `gorm:"foreignKey:WalletID" json:"wallet,omitempty"`
	Beneficiary Beneficiary `gorm:"foreignKey:BeneficiaryID" json:"beneficiary,omitempty"`
}

func (WalletDisbursement) TableName() string {
	return "wallet_disbursements"
}

type WalletFundingStatus string

const (
	FundingPending    WalletFundingStatus = "pending"
	FundingSuccessful WalletFundingStatus = "successful"
	FundingFailed     WalletFundingStatus = "failed"
)

type WalletFunding struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	WalletID uuid.UUID `gorm:"type:uuid;not null;index" json:"wallet_id"`

	Amount int64 `gorm:"not null;check:amount > 0" json:"amount"`

	Reference string `gorm:"type:varchar(100);uniqueIndex;not null" json:"reference"`

	ProviderReference string `gorm:"type:varchar(255);index" json:"provider_reference"`

	Method string `gorm:"type:varchar(50)" json:"method"`

	Status WalletFundingStatus `gorm:"type:varchar(30);not null;default:'pending'" json:"status"`

	Description string `gorm:"type:text" json:"description"`

	CreatedBy *uuid.UUID `gorm:"type:uuid;index" json:"created_by,omitempty"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Wallet Wallet `gorm:"foreignKey:WalletID" json:"wallet,omitempty"`
}

func (WalletFunding) TableName() string {
	return "wallet_fundings"
}
