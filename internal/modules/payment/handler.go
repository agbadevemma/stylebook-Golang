package payment

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
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
