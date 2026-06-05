package availability

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAvailability(c *gin.Context) {

	businessID := c.MustGet("businessId").(uuid.UUID)

	startDateStr := c.Query("start")
	endDateStr := c.Query("end")

	startDate, _ := time.Parse("2006-01-02", startDateStr)
	endDate, _ := time.Parse("2006-01-02", endDateStr)

	result, err := h.service.GetAvailabilityRange(businessID, startDate, endDate)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}
