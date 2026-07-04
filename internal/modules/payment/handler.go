package payment

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) PaymentWebhook(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("webhook read error: %v", err)
		c.JSON(http.StatusOK, utils.SuccessResponse("acknowledged", nil))
		return
	}

	signature := c.GetHeader("x-paystack-signature")

	if !h.service.VerifySignature(body, signature) {
		log.Printf("invalid webhook signature")
		c.JSON(http.StatusOK, utils.SuccessResponse("acknowledged", nil))
		return
	}

	var event utils.PaystackWebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("invalid webhook payload: %v", err)
		c.JSON(http.StatusOK, utils.SuccessResponse("acknowledged", nil))
		return
	}

	// err = h.service.HandlePaymentWebhook(
	// 	event.Event,
	// 	event.Data.Reference,
	// )

	// if err != nil {

	// 	log.Printf("payment webhook processing error: %v", err)
	// }
	if err := h.service.HandleEvent(event); err != nil {
		log.Printf("webhook processing error: %v", err)
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("acknowledged", nil))
}

func (h *Handler) GetBusinessTransactions(c *gin.Context) {
	busID, ok := c.Get("businessId")
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("missing business id"))
		return
	}

	businessID, ok := busID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid business id"))
		return
	}

	page := 1
	pageSize := 10
	search := c.Query("search")

	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.Query("pageSize"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}

	params := ListTransactionsParams{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	data, total, err := h.service.GetBusinessTransactions(businessID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("failed to fetch transactions"))
		return
	}

	result := make([]BusinessTransactionResponse, 0, len(data))
	for _, payment := range data {
		var paidAt *string
		if payment.PaidAt != nil {
			formatted := payment.PaidAt.Format(time.RFC3339)
			paidAt = &formatted
		}

		result = append(result, BusinessTransactionResponse{
			ID:            payment.ID.String(),
			Reference:     payment.Reference,
			CustomerName:  payment.Name,
			CustomerEmail: payment.Email,
			Amount:        payment.Amount,
			ServiceAmount: payment.ServiceAmount,
			PlatformFee:   payment.PlatformFee,
			NetAmount:     payment.NetAmount,
			Currency:      payment.Currency,
			Status:        string(payment.Status),
			PaidAt:        paidAt,
			CreatedAt:     payment.CreatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Fetched transaction successfully", PaginatedTransactionsResponse{
		Data:     result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}))
}

func (h *Handler) ResolveBankAccount(c *gin.Context) {
	accountNumber := c.Query("accountNumber")
	bankCode := c.Query("bankCode")

	if accountNumber == "" || bankCode == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("missing query params"))
		return
	}

	res, err := h.service.ResolveBankAccount(accountNumber, bankCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("success", gin.H{
		"accountNumber": res.Data.AccountNumber,
		"accountName":   res.Data.AccountName,
	}))
}

func (h *Handler) GetBanks(c *gin.Context) {
	banks, err := h.service.GetBanks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("failed to fetch banks: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("banks fetched successfully", banks))
}

func (h *Handler) CreateSubAccount(c *gin.Context) {
	busID, ok := c.Get("businessId")
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("missing business id"))
		return
	}

	businessID, ok := busID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid business id"))
		return
	}
	var req CreateSubaccount

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	subaccountCode, err := h.service.CreateBusinessSubaccount(businessID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("subaccount created successfully", gin.H{
		"subaccountCode": subaccountCode,
	}))
}
