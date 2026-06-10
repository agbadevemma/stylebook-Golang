package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Provider struct {
}

func NewProvider() *Provider {
	return &Provider{}
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

func (p *Provider) InitializePayment(email string, amount int64, reference string) (*PaystackResponse, error) {
	url := fmt.Sprintf("%s/transaction/initialize", os.Getenv("PAYSTACK_BASE_URL"))

	payload := InitRequest{
		Email:     email,
		Amount:    amount,
		Reference: reference,
	}
	// fmt.Printf("sdsds %s", url)
	fmt.Printf("payload: %v", payload)

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Marshal error: %v", err)
		return nil, err
	}


	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("NewRequest error: %v", err)
		return nil, err
	}


	req.Header.Set("Authorization", "Bearer "+os.Getenv("PAYSTACK_SECRET_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}


	resp, err := client.Do(req)
	if err != nil {
		log.Printf("client.Do error: %v", err)
		return nil, err
	}


	defer resp.Body.Close()

	var paystackResp PaystackResponse
	err = json.NewDecoder(resp.Body).Decode(&paystackResp)
	if err != nil {
		log.Printf("Decode error: %v", err)
		return nil, err
	}

	return &paystackResp, nil
}
