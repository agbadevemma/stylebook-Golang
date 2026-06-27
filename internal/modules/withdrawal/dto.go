package withdrawal

import "github.com/stylebookgo/agbadevemma/internal/utils"

type TransferRequest struct {
	Source    string `json:"source"`
	Amount    int64  `json:"amount"`
	Recipient string `json:"recipient"`
	Reference string `json:"reference"`
	Reason    string `json:"reason"`
}

type TransferResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		TransferCode string `json:"transfer_code"`
		Status       string `json:"status"`
	} `json:"data"`
}

type WithdrawalHistoryResponse struct {
	ID        string                 `json:"id"`
	Amount    int64                  `json:"amount"`
	Status    utils.WithdrawalStatus `json:"status"`
	CreatedAt string                 `json:"createdAt"`
	Reference string                 `json:"reference"`
}
