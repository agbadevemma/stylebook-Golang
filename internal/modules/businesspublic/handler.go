package businesspublic

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
		ID:       businessData.ID,
		Name:     businessData.Name,
		Timezone: businessData.Timezone,
	}))
}
