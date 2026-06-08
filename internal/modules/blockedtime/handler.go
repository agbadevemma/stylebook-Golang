package blockedtime

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

func (h *Handler) CreateBlockedTime(c *gin.Context) {
	var req CreateBlockedTimeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("An error occured"))
		return
	}

	businessID := c.MustGet("businessId").(uuid.UUID)

	if err := h.service.CreateBlockedTime(&req, businessID); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("An error occured"))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("blocked time created", nil))
}

func (h *Handler) ListBlockedTime(c *gin.Context) {

	busID, _ := c.Get("businessId")

	businessID, ok := busID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid Business Id"))
		return
	}

	blockedTimes, err := h.service.ListBlockedTime(businessID)

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, utils.ErrorResponse("blocked time not found"))
			return

		default:
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("failed to fetch "))
			return
		}

	}
	res := make([]ResponseBlockedTime, 0, len(blockedTimes))
	for _, blockTime := range blockedTimes {
		res = append(res, ResponseBlockedTime{
			ID:        blockTime.ID,
			StartTime: blockTime.StartTime,
			EndTime:   blockTime.EndTime,
			Reason:    blockTime.Reason,
		})
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Fetched blocked time succesfully", res))
}

func (h *Handler) GetBlockedTimeByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid blocked time id"))
		return
	}

	blockedTime, err := h.service.GetBlockedTimeByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("blocked time not found"))
			return
		}

		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("failed to get blocked time"))
		return
	}
	res := ResponseBlockedTime{
		ID:        blockedTime.ID,
		StartTime: blockedTime.StartTime,
		EndTime:   blockedTime.EndTime,
		Reason:    blockedTime.Reason,
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Fetched Blocked Time", res))
}

func (h *Handler) UpdateBlockedTime(c *gin.Context) {
	busID, _ := c.Get("businessId")

	businessID, ok := busID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid Business Id"))
		return
	}

	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid blocked time id"))
		return
	}

	var req UpdateBlockedTimeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	err = h.service.UpdateBlockedTime(id, businessID, &req)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, utils.ErrorResponse("blocked time not found"))
			return

		default:
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("failed to update blocked time"))
			return
		}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("blocked time updated", nil))
}

func (h *Handler) DeleteBlockedTime(c *gin.Context) {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid blocked time id"))
		return
	}

	if err := h.service.DeleteBlockedTime(id); err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, utils.ErrorResponse("blocked time not found"))
			return

		default:
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("failed to delete blocked time"))
			return
		}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("blocked time deleted", nil))

}
