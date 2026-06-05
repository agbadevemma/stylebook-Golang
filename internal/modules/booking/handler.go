package booking

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

func (h *Handler) CreateBooking(c *gin.Context) {

	var req CreateBookingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	businessID, ok := c.MustGet("businessId").(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid Business Id"))
		return
	}

	if err := h.service.CreateBooking(&req, businessID); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))

		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Booking created successfully", nil))
}
