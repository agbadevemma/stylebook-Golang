package payment

import (
	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/utils"
)

type PaystackWebhookEvent struct {
	Event string `json:"event"`

	Data struct {
		Reference string `json:"reference"`
		Amount    int64  `json:"amount"`
		Status    string `json:"status"`
	} `json:"data"`
}

type TransactionResponse struct {
	Result     []models.Payment `json:"result"`
	Pagination utils.Pagination `json:"pagination"`
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
