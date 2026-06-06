package payment

type PaystackWebhookEvent struct {
	Event string `json:"event"`

	Data struct {
		Reference string `json:"reference"`
		Amount    int64  `json:"amount"`
		Status    string `json:"status"`
	} `json:"data"`
}
