package service

import (
	"fmt"
	"log"
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

	businessID, _ := c.Get("businessID")
	fmt.Print("businesId", businessID)
	idStr, ok := businessID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid Business Id"))
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}
	err := h.service.CreateService(&req, idStr)

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

	res := ServiceResponse{
		ID:          service.ID,
		Name:        service.Name,
		Duration:    service.Duration,
		Price:       service.Price,
		Description: service.Description,
		IsActive:    service.IsActive,
		BusinessID:  service.BusinessID,
	}

	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Service found", res))
}

func (h *Handler) ListServices(c *gin.Context) {
	services, err := h.service.ListServices()
	if err != nil {
		log.Printf("err %v", err.Error())
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to list services"))
		return
	}
	var res []ServiceResponse

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
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Service ID is required"))
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}
	err := h.service.UpdateService(&req, uuid.MustParse(id))
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
