package business

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/utils"
	"gorm.io/gorm"
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

	c.JSON(http.StatusOK, utils.SuccessResponse("business fetched successfully", BusinessSettingsResponse{
		ID:             business.ID.String(),
		Name:           business.Name,
		Timezone:       business.Timezone,
		AccountName:    business.AccountName,
		AccountNumber:  business.AccountNumber,
		BankName:       business.BankName,
		BankCode:       business.BankCode,
		SubaccountCode: business.SubaccountCode,
	}))
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

	c.JSON(http.StatusOK, utils.SuccessResponse("business updated successfully", BusinessSettingsResponse{
		ID:             business.ID.String(),
		Name:           business.Name,
		Timezone:       business.Timezone,
		AccountName:    business.AccountName,
		AccountNumber:  business.AccountNumber,
		BankName:       business.BankName,
		BankCode:       business.BankCode,
		SubaccountCode: business.SubaccountCode,
	}))
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


func (h *Handler) GetPublicBusiness(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("business id is required"))
		return
	}

	businessID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid business id"))
		return
	}

	businessData, err := h.service.GetBusinessByID(businessID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("business not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("failed to fetch business"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Business fetched successfully", PublicBusinessResponse{
		ID:           businessData.ID,
		Name:         businessData.Name,
		Timezone:     businessData.Timezone,
		IsSubAccount: businessData.SubaccountCode != "",
	}))
}