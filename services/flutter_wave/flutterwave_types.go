// services/flutterwave_types.go
package services

// AccountResolutionRequest for verifying bank accounts
type AccountResolutionRequest struct {
	AccountNumber string `json:"account_number"`
	AccountBank   string `json:"account_bank"`
	Country       string `json:"country,omitempty"`
}

// AccountResolutionResponse from Flutterwave
type AccountResolutionResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AccountNumber string `json:"account_number"`
		AccountName   string `json:"account_name"`
	} `json:"data"`
}

// TransferRequest for initiating payouts
type TransferRequest struct {
	AccountBank   string `json:"account_bank"`
	AccountNumber string `json:"account_number"`
	Amount        int64  `json:"amount"`
	Narration     string `json:"narration"`
	Currency      string `json:"currency"`
	Reference     string `json:"reference"`
	CallbackURL   string `json:"callback_url,omitempty"`
	DebitCurrency string `json:"debit_currency,omitempty"`
}

// TransferResponse from Flutterwave
type TransferResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID                int64   `json:"id"`
		AccountNumber     string  `json:"account_number"`
		BankCode          string  `json:"bank_code"`
		FullName          string  `json:"full_name"`
		CreatedAt         string  `json:"created_at"`
		Currency          string  `json:"currency"`
		DebitCurrency     string  `json:"debit_currency"`
		Amount            int64   `json:"amount"`
		Fee               float64 `json:"fee"`
		Status            string  `json:"status"`
		Reference         string  `json:"reference"`
		Meta              any     `json:"meta"`
		Narration         string  `json:"narration"`
		CompleteMessage   string  `json:"complete_message"`
		RequiresApproval  int     `json:"requires_approval"`
		IsApproved        int     `json:"is_approved"`
		BankName          string  `json:"bank_name"`
	} `json:"data"`
}

// GetTransferResponse for checking transfer status
type GetTransferResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID              int64   `json:"id"`
		Amount          int64   `json:"amount"`
		Currency        string  `json:"currency"`
		Status          string  `json:"status"`
		Reference       string  `json:"reference"`
		CompleteMessage string  `json:"complete_message"`
	} `json:"data"`
}