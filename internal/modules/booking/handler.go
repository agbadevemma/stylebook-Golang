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

	// Validate startTime is in UTC
	if err := utils.ValidateISO8601Timestamp(req.StartTime, "startTime"); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	res, err := h.service.CreateBooking(c.Request.Context(), &req, req.BusinessID)

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

	stats, err := h.service.GetBusinessBookingStats(c.Request.Context(), businessID)
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
	startDate := c.Query("startDate") // YYYY-MM-DD
	endDate := c.Query("endDate")     // YYYY-MM-DD

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

	params := ListBookingsParams{Page: page, PageSize: pageSize, Search: search, StartDate: startDate, EndDate: endDate}

	bookings, total, err := h.service.ListBusinessBookings(c.Request.Context(), businessID, params)
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

func (h *Handler) RescheduleBooking(c *gin.Context) {
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

	bookingID := c.Param("bookingId")
	parsedBookingID, err := uuid.Parse(bookingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid booking id"))
		return
	}

	var req RescheduleBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	// Validate newStartTime is in UTC
	if err := utils.ValidateISO8601Timestamp(req.NewStartTime, "newStartTime"); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	err = h.service.RescheduleBooking(c.Request.Context(), parsedBookingID, businessID, &req)
	if err != nil {
		log.Printf("reschedule error: %v", err.Error())
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Booking rescheduled successfully", nil))
}

func (h *Handler) UpdateBookingStatus(c *gin.Context) {
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

	bookingID := c.Param("bookingId")
	parsedBookingID, err := uuid.Parse(bookingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid booking id"))
		return
	}

	var req UpdateBookingStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	err = h.service.UpdateBookingStatus(c.Request.Context(), parsedBookingID, businessID, req.Status)
	if err != nil {
		log.Printf("update status error: %v", err.Error())
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Booking status updated successfully", nil))
}

func (h *Handler) GetBookingDetails(c *gin.Context) {
	bookingID := c.Param("bookingId")
	parsedBookingID, err := uuid.Parse(bookingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid booking id"))
		return
	}

	booking, services, business, err := h.service.GetBookingDetailsPublic(c.Request.Context(), parsedBookingID)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("booking not found"))
		return
	}

	type ServiceDetail struct {
		Name     string `json:"name"`
		Duration int    `json:"duration"`
		Price    int64  `json:"price"`
	}

	type BookingDetail struct {
		ID            string          `json:"id"`
		CustomerName  string          `json:"customerName"`
		CustomerEmail string          `json:"customerEmail"`
		CustomerPhone string          `json:"customerPhone"`
		BusinessName  string          `json:"businessName"`
		Timezone      string          `json:"timezone"`
		StartTime     string          `json:"startTime"`
		EndTime       string          `json:"endTime"`
		TotalDuration int             `json:"totalDuration"`
		BookingStatus string          `json:"bookingStatus"`
		PaymentStatus string          `json:"paymentStatus"`
		Services      []ServiceDetail `json:"services"`
	}

	serviceDetails := make([]ServiceDetail, len(services))
	for i, service := range services {
		serviceDetails[i] = ServiceDetail{
			Name:     service.Name,
			Duration: service.Duration,
			Price:    service.Price,
		}
	}

	response := BookingDetail{
		ID:            booking.ID.String(),
		CustomerName:  booking.CustomerName,
		CustomerEmail: booking.CustomerEmail,
		CustomerPhone: booking.CustomerPhone,
		BusinessName:  business.Name,
		Timezone:      business.Timezone,
		StartTime:     booking.StartTime.Format("2006-01-02T15:04:05Z07:00"),
		EndTime:       booking.EndTime.Format("2006-01-02T15:04:05Z07:00"),
		TotalDuration: booking.TotalDuration,
		BookingStatus: string(booking.BookingStatus),
		PaymentStatus: string(booking.PaymentStatus),
		Services:      serviceDetails,
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Booking details retrieved successfully", response))
}
