package booking

import (
	"log"
	"net/http"
	"strconv"

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

func (h *Handler) CreateBooking(c *gin.Context) {

	var req CreateBookingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	res, err := h.service.CreateBooking(&req, req.BusinessID)

	if err != nil {
		log.Printf("%v", err.Error())
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("An error occurred while creating the booking"))

		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Booking created successfully", res))

}

func (h *Handler) GetBusinessStats(c *gin.Context) {
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

	stats, err := h.service.GetBusinessBookingStats(businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("failed to fetch stats"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Stats retrieved successfully", stats))
}

func (h *Handler) ListBusinessBookings(c *gin.Context) {
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

	page := 1
	pageSize := 10
	search := c.Query("search")

	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.Query("pageSize"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}

	params := ListBookingsParams{Page: page, PageSize: pageSize, Search: search}

	bookings, total, err := h.service.ListBusinessBookings(businessID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("failed to fetch bookings"))
		return
	}

	res := make([]BusinessBookingResponse, 0, len(bookings))
	for _, booking := range bookings {
		res = append(res, BusinessBookingResponse{
			ID:            booking.ID,
			CustomerName:  booking.CustomerName,
			CustomerPhone: booking.CustomerPhone,
			CustomerEmail: booking.CustomerEmail,
			TotalDuration: booking.TotalDuration,
			StartTime:     booking.StartTime,
			EndTime:       booking.EndTime,
			BookingStatus: string(booking.BookingStatus),
			PaymentStatus: string(booking.PaymentStatus),
			CreatedAt:     booking.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Bookings retrieved successfully", PaginatedBookingsResponse{
		Data:     res,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}))
}
