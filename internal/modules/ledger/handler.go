package ledger

import (
	"net/http"

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

func (h *Handler) GetBusinessBalance(c *gin.Context) {
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

	balance, err := h.service.GetBusinessBalance(businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("failed to get balance"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("success", gin.H{
		"balance": balance,
	}))
}
