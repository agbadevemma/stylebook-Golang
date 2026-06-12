package payment

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

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

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid limit"))
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid offset"))
		return
	}

	data, total, err := h.service.GetBusinessTransactions(businessID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("failed to fetch transactions"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Fetched transaction successfully", TransactionResponse{Result: data, Pagination: utils.Pagination{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}}))
}

func (h *Handler) ResolveBankAccount(c *gin.Context) {
	accountNumber := c.Query("account_number")
	bankCode := c.Query("bank_code")

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

func (h *Handler) CreateTransferRecipient(c *gin.Context) {
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
	var req CreateTransferRecipientRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	res, err := h.service.CreateTransferRecipient(
		req.Name,
		req.AccountNumber,
		req.BankName,
		req.BankCode,
		businessID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("recipient created", gin.H{
		"recipientCode": res.Data.RecipientCode,
	}))
}
