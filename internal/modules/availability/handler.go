package availability

import (
	"log"
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

func (h *Handler) GetAvailability(c *gin.Context) {

	businessID := c.Query("businessId")
	if businessID == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("No business Id"))
		return
	}

	parsedBusinessId, err := uuid.Parse(businessID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid business ID format"))
		return
	}

	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")
	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("startDate and endDate query params are required"))
		return
	}

	// Validate date format
	if err := utils.ValidateDateRange(startDateStr, endDateStr); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	startDate, _ := time.Parse("2006-01-02", startDateStr)
	endDate, _ := time.Parse("2006-01-02", endDateStr)

	// Validate start date is not in the past
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	if startDate.Before(today) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("start date cannot be in the past"))
		return
	}

	if endDate.Sub(startDate) > 31*24*time.Hour {
		log.Printf("Date range error: startDate=%s, endDate=%s", startDateStr, endDateStr)
		c.JSON(
			http.StatusBadRequest,
			utils.ErrorResponse("Date range cannot exceed 31 days"),
		)
		return
	}

	result, err := h.service.GetAvailabilityRange(c.Request.Context(), parsedBusinessId, startDate, endDate)
	if err != nil {
		log.Printf("Get availability api error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("An error occurred"))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("availability fetched", result))
}
