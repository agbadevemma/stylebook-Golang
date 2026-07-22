package booking

import (
	"time"

	"github.com/google/uuid"
)

type CreateBookingRequest struct {
	CustomerName  string      `json:"customerName" binding:"required,min=2"`
	CustomerPhone string      `json:"customerPhone" binding:"required,min=7"`
	CustomerEmail string      `json:"customerEmail" binding:"required"`
	ServiceIDs    []uuid.UUID `json:"serviceIds" binding:"required,min=1"`
	BusinessID    uuid.UUID   `json:"businessId" binding:"required" `
	StartTime     time.Time   `json:"startTime" binding:"required"`
}

type BookingResponse struct {
	BookingID  uuid.UUID `json:"bookingId"`
	PaymentURL string    `json:"paymentUrl"`
}

type BusinessBookingResponse struct {
	ID            uuid.UUID `json:"id"`
	CustomerName  string    `json:"customerName"`
	CustomerPhone string    `json:"customerPhone"`
	CustomerEmail string    `json:"customerEmail"`
	TotalDuration int       `json:"totalDuration"`
	StartTime     time.Time `json:"startTime"`
	EndTime       time.Time `json:"endTime"`
	BookingStatus string    `json:"bookingStatus"`
	PaymentStatus string    `json:"paymentStatus"`
	CreatedAt     time.Time `json:"createdAt"`
}

type PaginatedBookingsResponse struct {
	Data     []BusinessBookingResponse `json:"data"`
	Total    int64                     `json:"total"`
	Page     int                       `json:"page"`
	PageSize int                       `json:"pageSize"`
}

type BookingStatsResponse struct {
	TotalBookings   int64 `json:"totalBookings"`
	TodayBookings   int64 `json:"todayBookings"`
	PendingBookings int64 `json:"pendingBookings"`
	TotalClients    int64 `json:"totalClients"`
}

type ListBookingsParams struct {
	Page      int
	PageSize  int
	Search    string
	StartDate string // YYYY-MM-DD
	EndDate   string // YYYY-MM-DD
}

type RescheduleBookingRequest struct {
	NewStartTime time.Time `json:"newStartTime" binding:"required"`
}

type UpdateBookingStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
