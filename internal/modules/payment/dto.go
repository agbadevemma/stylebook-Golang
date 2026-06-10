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
