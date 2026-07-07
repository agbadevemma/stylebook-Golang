package otp

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

func (h *Handler) SendOtp(c *gin.Context) {
	var req SendOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	err := h.service.SendOTP(nil, req.Email, req.Type)
	if err != nil {
		log.Printf("Error -->> %v", err.Error())
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("An error  occured"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("OTP sent successfully", nil))
}

