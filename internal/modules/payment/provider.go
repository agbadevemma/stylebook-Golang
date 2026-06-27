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



func (p *Provider) ResolveBankAccount(accountNumber, bankCode string) (*BankResolveResponse, error) {
	url := fmt.Sprintf("%s/bank/resolve?account_number=%s&bank_code=%s",
		os.Getenv("PAYSTACK_BASE_URL"),
		accountNumber,
		bankCode,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("PAYSTACK_SECRET_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result BankResolveResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.Status {
		return nil, fmt.Errorf("paystack error: %s", result.Message)
	}

	return &result, nil
}


func (p *Provider) CreateTransferRecipient(
	name,
	accountNumber,
	bankCode string,
) (*CreateRecipientResponse, error) {

	url := fmt.Sprintf("%s/transferrecipient", os.Getenv("PAYSTACK_BASE_URL"))

	payload := CreateRecipientRequest{
		Type:          "nuban",
		Name:          name,
		AccountNumber: accountNumber,
		BankCode:      bankCode,
		Currency:      "NGN",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

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

	var result CreateRecipientResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.Status {
		return nil, fmt.Errorf("paystack error: %s", result.Message)
	}

	return &result, nil
}

func (p *Provider) GetBanks() ([]Bank, error) {
	url := fmt.Sprintf("%s/bank", os.Getenv("PAYSTACK_BASE_URL"))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("PAYSTACK_SECRET_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result BanksResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.Status {
		return nil, fmt.Errorf("paystack error: %s", result.Message)
	}

	return result.Data, nil
}


