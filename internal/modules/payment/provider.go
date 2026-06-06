package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Provider struct {
}

func NewProvider() *Provider {
	return &Provider{}
}

type InitRequest struct {
	Email  string `json:"email"`
	Amount int64  `json:"amount"`
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

func (p *Provider) InitializePayment(email string, amount int64, reference string) (*PaystackResponse, error) {
	url := fmt.Sprintf("%s/transaction/initialize", os.Getenv("PAYSTACK_BASE_URL"))
	payload := InitRequest{
		Email:  email,
		Amount: amount,
	}
	body, err := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("PAYSTACK_SECRET_KEY"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var paystackResp PaystackResponse
	err = json.NewDecoder(resp.Body).Decode(&paystackResp)
	if err != nil {
		return nil, err
	}

	if !paystackResp.Status {
		return nil, fmt.Errorf("paystack error: %s", paystackResp.Message)
	}

	return &paystackResp, nil
}
