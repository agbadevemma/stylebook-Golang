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

func (p *Provider) CreateSubaccount(businessName, settlementBank, accountNumber string) (*CreateSubaccountResponse, error) {
	url := fmt.Sprintf("%s/subaccount", os.Getenv("PAYSTACK_BASE_URL"))

	payload := CreateSubaccountRequest{
		BusinessName:     businessName,
		SettlementBank:   settlementBank,
		AccountNumber:    accountNumber,
		PercentageCharge: 0,
		Description:      fmt.Sprintf("Subaccount for %s", businessName),
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

	var result CreateSubaccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.Status {
		return nil, fmt.Errorf("paystack error: %s", result.Message)
	}

	return &result, nil
}

func (p *Provider) UpdateSubaccount(subaccountCode, businessName, settlementBank, accountNumber string) (*CreateSubaccountResponse, error) {
	url := fmt.Sprintf("%s/subaccount/%s", os.Getenv("PAYSTACK_BASE_URL"), subaccountCode)
	payload := UpdateSubaccountRequest{
		BusinessName:     businessName,
		SettlementBank:   settlementBank,
		AccountNumber:    accountNumber,
		PercentageCharge: 0,
		Active:           true,
		Description:      fmt.Sprintf("Updated subaccount for %s", businessName),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
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

	var result CreateSubaccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.Status {
		return nil, fmt.Errorf("paystack error: %s", result.Message)
	}

	return &result, nil
}

func (p *Provider) InitializePaymentWithSplit(email string, totalAmount int64, reference string, businessSubaccountCode string, platformFee int64, bookingID string) (*PaystackResponse, error) {
	url := fmt.Sprintf("%s/transaction/initialize", os.Getenv("PAYSTACK_BASE_URL"))
	callbackURL := fmt.Sprintf("%s/book/success?bookingId=%s",
		os.Getenv("FRONTEND_URL"),
		bookingID,
	)

	payload := map[string]interface{}{
		"email":              email,
		"amount":             totalAmount,
		"reference":          reference,
		"subaccount":         businessSubaccountCode,
		"transaction_charge": platformFee,
		"callback_url":       callbackURL,
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

	var paystackResp PaystackResponse
	if err := json.NewDecoder(resp.Body).Decode(&paystackResp); err != nil {
		return nil, err
	}

	return &paystackResp, nil
}
