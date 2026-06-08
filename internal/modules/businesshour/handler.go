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
	res := make([]BusinessHoursResponse, 0, len(hours))

	for _, h := range hours {
		res = append(res, BusinessHoursResponse{
			ID:        h.ID,
			CreatedAt: h.CreatedAt,
			UpdatedAt: h.UpdatedAt,
			DayOfWeek: h.DayOfWeek,
			IsOpen:    h.IsOpen,
			StartTime: h.StartTime,
			EndTime:   h.EndTime,
		})

	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("An error Occurred"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Business hours fetched successfully", res))
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

	var req UpdateBusinessHourRequestWrapper

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	if len(req.Hours) == 0 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("request body cannot be empty"))
		return
	}

	if err := h.service.UpdateBusinessHour(
		businessID,
		req.Hours,
	); err != nil {

		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Business hours updated successfully", nil))
}
