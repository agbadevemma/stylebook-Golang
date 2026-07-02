package business

import "github.com/google/uuid"

type UpdateBusinessDetailsRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=100"`
	Timezone string `json:"timezone" binding:"required"`
}

type BusinessSettingsResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Timezone      string `json:"timezone"`
	AccountNumber string `json:"accountNumber"`
	AccountName   string `json:"accountName"`
	BankName      string `json:"bankName"`
	BankCode      string `json:"bankCode"`
	SubaccountCode string `json:"subaccountCode"`
}

type PublicBusinessResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Timezone     string    `json:"timezone"`
	IsSubAccount bool      `json:"isSubAccount"`
}
