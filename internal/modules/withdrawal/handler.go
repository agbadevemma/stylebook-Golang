package withdrawal

import (
	"net/http"
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

func (h *Handler) GetWithdrawalHistory(c *gin.Context) {
	businessID, ok := getBusinessID(c)
	if !ok {
		return
	}

	withdrawals, err := h.service.GetBusinessWithdrawals(businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("failed to fetch withdrawals"))
		return
	}

	result := make([]WithdrawalHistoryResponse, 0, len(withdrawals))
	for _, withdrawal := range withdrawals {
		result = append(result, WithdrawalHistoryResponse{
			ID:        withdrawal.ID.String(),
			Amount:    withdrawal.Amount,
			Status:    withdrawal.Status,
			CreatedAt: withdrawal.CreatedAt.Format(time.RFC3339),
			Reference: withdrawal.TransferReference,
		})
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("withdrawals fetched successfully", result))
}

func (h *Handler) CreateWithdrawal(c *gin.Context) {
	businessID, ok := getBusinessID(c)
	if !ok {
		return
	}

	var req struct {
		Amount int64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("amount must be greater than zero"))
		return
	}

	res, err := h.service.CreateWithdrawal(businessID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("withdrawal initiated", gin.H{
		"transferCode": res.Data.TransferCode,
		"status":       res.Data.Status,
	}))
}

func getBusinessID(c *gin.Context) (uuid.UUID, bool) {
	busID, ok := c.Get("businessId")
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("missing business id"))
		return uuid.Nil, false
	}

	businessID, ok := busID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid business id"))
		return uuid.Nil, false
	}

	return businessID, true
}

// func (h *Handler) PaystackWebhook(c *gin.Context) {
// 	body, err := io.ReadAll(c.Request.Body)
// 	if err != nil {
// 		log.Printf("webhook read error: %v", err)
// 		c.JSON(http.StatusOK, utils.SuccessResponse("acknowledged", nil))
// 		return
// 	}

// 	signature := c.GetHeader("x-paystack-signature")

// 	if !h.service.VerifySignature(body, signature) {
// 		log.Printf("invalid webhook signature")
// 		c.JSON(http.StatusOK, utils.SuccessResponse("acknowledged", nil))
// 		return
// 	}

// 	var event TransferWebhookEvent
// 	if err := json.Unmarshal(body, &event); err != nil {
// 		log.Printf("invalid webhook payload: %v", err)
// 		c.JSON(http.StatusOK, utils.SuccessResponse("acknowledged", nil))
// 		return
// 	}

// 	err = h.service.HandleTransferWebhook(
// 		event.Data.Reference,
// 		event.Event,
// 	)

// 	if err != nil {
// 		log.Printf("transfer webhook processing failed: %v", err)
// 	}

// 	c.JSON(http.StatusOK, utils.SuccessResponse("acknowledged", nil))
// }
