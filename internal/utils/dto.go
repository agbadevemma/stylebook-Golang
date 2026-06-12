package utils

type Pagination struct {
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

type PaystackWebhookEvent struct {
	Event string `json:"event"`

	Data struct {
		Domain    string `json:"domain"`
		Amount    int64  `json:"amount"`
		Currency  string `json:"currency"`
		Reference string `json:"reference"`
		Status    string `json:"status"`

		// Payment fields
		Authorization struct {
			AuthorizationCode string `json:"authorization_code"`
			Bin               string `json:"bin"`
			Last4             string `json:"last4"`
			ExpMonth          string `json:"exp_month"`
			ExpYear           string `json:"exp_year"`
			Channel           string `json:"channel"`
			CardType          string `json:"card_type"`
			Bank              string `json:"bank"`
			CountryCode       string `json:"country_code"`
			Brand             string `json:"brand"`
		} `json:"authorization,omitempty"`

		Customer struct {
			ID           int64  `json:"id"`
			Email        string `json:"email"`
			CustomerCode string `json:"customer_code"`
		} `json:"customer,omitempty"`

		// Transfer fields (optional depending on event type)
		TransferCode string `json:"transfer_code,omitempty"`
		Reason       string `json:"reason,omitempty"`

		Recipient struct {
			Name          string `json:"name"`
			Type          string `json:"type"`
			AccountNumber string `json:"account_number"`
			BankCode      string `json:"bank_code"`
			BankName      string `json:"bank_name"`
		} `json:"recipient,omitempty"`

		// timestamps
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	} `json:"data"`
}
