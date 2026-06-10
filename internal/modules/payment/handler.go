package payment

import (
	"encoding/json"
	"io"
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

func (h *Handler) PaystackWebhook(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid request body"))
		return
	}

	signature := c.GetHeader("x-paystack-signature")

	// 1. Verify signature
	if !h.service.VerifySignature(body, signature) {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("invalid signature"))
		return
	}

	// 2. Parse event
	var event PaystackWebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid payload"))
		return
	}

	// 3. Only handle successful payments
	if event.Event != "charge.success" {
		c.JSON(http.StatusOK, utils.SuccessResponse("ignored event", nil))
		return
	}

	// 4. Process payment
	err = h.service.HandleSuccessfulPayment(event.Data.Reference)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("payment processed successfully", nil))
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

	c.JSON(http.StatusOK, utils.SuccessResponse("Succe", map[string]any{
		"data":   data,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}))
}
