package booking

import (
	"log"
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

func (h *Handler) CreateBooking(c *gin.Context) {

	var req CreateBookingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	res, err := h.service.CreateBooking(&req, req.BusinessID)

	if err != nil {
		log.Printf("%v", err.Error())
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("An error occurred while creating the booking"))

		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Booking created successfully", res))

}
