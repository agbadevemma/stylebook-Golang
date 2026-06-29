package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetBusinessSettings(c *gin.Context) {
	businessID, ok := getBusinessIDFromContext(c)
	if !ok {
		return
	}

	business, err := h.service.GetBusinessSettings(businessID)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("business not found"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("business fetched successfully", toBusinessSettingsResponse(business)))
}

func (h *Handler) UpdateBusinessDetails(c *gin.Context) {
	businessID, ok := getBusinessIDFromContext(c)
	if !ok {
		return
	}

	var req UpdateBusinessDetailsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	business, err := h.service.UpdateBusinessDetails(businessID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("business updated successfully", toBusinessSettingsResponse(business)))
}

func getBusinessIDFromContext(c *gin.Context) (uuid.UUID, bool) {
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

func toBusinessSettingsResponse(business *models.Business) BusinessSettingsResponse {
	return BusinessSettingsResponse{
		ID:            business.ID.String(),
		Name:          business.Name,
		Timezone:      business.Timezone,
		AccountName:   business.AccountName,
		AccountNumber: business.AccountNumber,
		BankName:      business.BankName,
		BankCode:      business.BankCode,
	}
}
