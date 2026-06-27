package business

type UpdateBusinessDetailsRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=100"`
	Timezone string `json:"timezone" binding:"required"`
}

type BusinessSettingsResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Timezone      string `json:"timezone"`
	AccountNumber string `json:"accountNumber"`
	BankName      string `json:"bankName"`
	BankCode      string `json:"bankCode"`
}
