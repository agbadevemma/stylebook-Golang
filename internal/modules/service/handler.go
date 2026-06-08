package service

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

func (h *Handler) CreateService(c *gin.Context) {
	var req CreateServiceRequest

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
	err := h.service.CreateService(&req, businessID)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse("Service created successfully", nil))
}

func (h *Handler) GetServiceByID(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Service ID is required"))
		return
	}
	parsedID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid UUID"))
		return
	}

	service, err := h.service.GetServiceByID(parsedID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Service not found"))
		return
	}

	res := ServiceResponse{
		ID:          service.ID,
		Name:        service.Name,
		Duration:    service.Duration,
		Price:       service.Price,
		Description: service.Description,
		IsActive:    service.IsActive,
		BusinessID:  service.BusinessID,
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Service found", res))
}

func (h *Handler) ListServices(c *gin.Context) {
	businessIDValue, exists := c.Get("businessId")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("business id not found"))
		return
	}

	businessID, ok := businessIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("invalid business id"))
		return
	}

	services, err := h.service.ListServices(businessID)
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

	res := make([]ServiceResponse, 0, len(services))

	for _, s := range services {
		res = append(res, ServiceResponse{
			ID:          s.ID,
			Name:        s.Name,
			Duration:    s.Duration,
			Description: s.Description,
			IsActive:    s.IsActive,
			Price:       s.Price,
			BusinessID:  s.BusinessID,
		})
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Services retrieved successfully", res))
}

func (h *Handler) UpdateService(c *gin.Context) {
	var req UpdateServiceRequest

	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Service ID is required"))
		return
	}
	parsedID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid UUID"))
		return
	}

	_, err = h.service.GetServiceByID(parsedID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Service not found"))
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	err = h.service.UpdateService(&req, parsedID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Service updated successfully", nil))
}

func (h *Handler) DeleteService(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Service ID is required"))
		return
	}
	parsedID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid UUID"))
		return
	}

	_, err = h.service.GetServiceByID(parsedID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Service not found"))
		return
	}

	err = h.service.DeleteService(parsedID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Service deleted successfully", nil))
}
