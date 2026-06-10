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

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid start date format, expected YYYY-MM-DD"))
		return
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	if startDate.Before(today) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("start date cannot be in the past"))
		return
	}


	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid end date format, expected YYYY-MM-DD"))
		return
	}

	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("end date must be after start date"))
		return
	}

	// if endDate.Sub(startDate) > 31*24*time.Hour {
	// 	c.JSON(http.StatusBadRequest, utils.ErrorResponse("Date range cannot exceed 31 days"))
	// 	return
	// }

	result, err := h.service.GetAvailabilityRange(parsedBusinessId, startDate, endDate)
	if err != nil {
		log.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("An error occured"))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("availability fetched", result))
}
