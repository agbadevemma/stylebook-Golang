package payment

import (
	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/utils"
)

type TransactionResponse struct {
	Result     []models.Payment `json:"result"`
	Pagination utils.Pagination `json:"pagination"`
}

type BusinessTransactionResponse struct {
	ID            string  `json:"id"`
	Reference     string  `json:"reference"`
	CustomerName  string  `json:"customerName"`
	CustomerEmail string  `json:"customerEmail"`
	Amount        int64   `json:"amount"`
	ServiceAmount int64   `json:"serviceAmount"`
	PlatformFee   int64   `json:"platformFee"`
	NetAmount     int64   `json:"netAmount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
	PaidAt        *string `json:"paidAt"`
	CreatedAt     string  `json:"createdAt"`
}

type PaginatedTransactionsResponse struct {
	Data     []BusinessTransactionResponse `json:"data"`
	Total    int64                         `json:"total"`
	Page     int                           `json:"page"`
	PageSize int                           `json:"pageSize"`
}

type ListTransactionsParams struct {
	Page     int
	PageSize int
	Search   string
	StartDate string // YYYY-MM-DD
	EndDate   string // YYYY-MM-DD
}

type InitRequest struct {
	Email     string `json:"email"`
	Amount    int64  `json:"amount"`
	Reference string `json:"reference"`
}

type PaystackResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
		Reference        string `json:"reference"`
	} `json:"data"`
}

type BankResolveResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AccountNumber string `json:"account_number"`
		AccountName   string `json:"account_name"`
		BankID        int    `json:"bank_id"`
	} `json:"data"`
}

type CreateSubaccount struct {
	AccountName   string `json:"accountName" binding:"required"`
	BankName      string `json:"bankName" binding:"required"`
	AccountNumber string `json:"accountNumber" binding:"required"`
	BankCode      string `json:"bankCode" binding:"required"`
}

type Bank struct {
	ID       int    `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Longname string `json:"longname"`
	Ussd     string `json:"ussd"`
	Active   bool   `json:"active"`
	Country  string `json:"country"`
	Currency string `json:"currency"`
	Type     string `json:"type"`
	Slug     string `json:"slug"`
}

type BanksResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    []Bank `json:"data"`
}

type CreateSubaccountRequest struct {
	BusinessName     string  `json:"business_name"`
	SettlementBank   string  `json:"settlement_bank"`
	AccountNumber    string  `json:"account_number"`
	PercentageCharge float32 `json:"percentage_charge"`
	Description      string  `json:"description,omitempty"`
}

type CreateSubaccountResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		SubaccountCode   string  `json:"subaccount_code"`
		BusinessName     string  `json:"business_name"`
		SettlementBank   string  `json:"settlement_bank"`
		AccountNumber    string  `json:"account_number"`
		PercentageCharge float32 `json:"percentage_charge"`
	} `json:"data"`
}

// Split Payment DTOs
type SplitConfig struct {
	Type           string `json:"type"` // "percentage" or "flat"
	SubaccountCode string `json:"subaccount_code"`
	Share          int64  `json:"share"` // percentage or amount in kobo
}

type InitSplitPaymentRequest struct {
	Email     string `json:"email"`
	Amount    int64  `json:"amount"`
	Reference string `json:"reference"`
	Split     struct {
		Type   string        `json:"type"`
		Splits []SplitConfig `json:"splits"`
	} `json:"split"`
}

type UpdateSubaccountRequest struct {
	BusinessName     string  `json:"business_name"`
	SettlementBank   string  `json:"settlement_bank"`
	AccountNumber    string  `json:"account_number"`
	PercentageCharge float32 `json:"percentage_charge"`
	Active           bool    `json:"active"`
	Description      string  `json:"description,omitempty"`
}