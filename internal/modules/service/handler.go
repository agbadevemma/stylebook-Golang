package service

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

func (h *Handler) CreateService(c *gin.Context) {
	var req CreateServiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}
	err := h.service.CreateService(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse("Service created successfully", nil))
}

func (h *Handler) GetServiceByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Service ID is required"))
		return
	}
	service, err := h.service.GetServiceByID(uuid.MustParse(id))

	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Service found", service))
}

func (h *Handler) ListServices(c *gin.Context) {
	services, err := h.service.ListServices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to list services"))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Services retrieved successfully", services))
}

func (h *Handler) UpdateService(c *gin.Context) {
	var req UpdateServiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}
	err := h.service.UpdateService(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Service updated successfully", nil))
}

func (h *Handler) DeleteService(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Service ID is required"))
		return
	}
	err := h.service.DeleteService(uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Service deleted successfully", nil))
}
