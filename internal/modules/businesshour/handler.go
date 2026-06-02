package businesshour

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

func (h *Handler) ListBusinessHours(c *gin.Context) {

	busId, _ := c.Get("businessId")

	businessID, ok := busId.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid Business Id"))
		return
	}

	hours, err := h.service.ListByBusiness(
		businessID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Business hours fetched successfully", hours))
}

func (h *Handler) CreateBusinessHour(c *gin.Context) {
	var req []CreateBusinessHourRequest

	busId, _ := c.Get("businessId")

	businessID, ok := busId.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid Business Id"))
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}
	err := h.service.CreateBusinessHour(businessID, req)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse("Business hours created successfully", nil))
}

func (h *Handler) UpdateBusinessHour(c *gin.Context) {

	businessID := c.MustGet("businessId").(uuid.UUID)

	var req []UpdateBusinessHourRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	if err := h.service.UpdateBusinessHour(
		businessID,
		req,
	); err != nil {

		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Business hours updated successfully", nil))
}
