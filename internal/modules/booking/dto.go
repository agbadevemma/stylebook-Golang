package booking

import (
	"time"

	"github.com/google/uuid"
)

type CreateBookingRequest struct {
	CustomerName  string      `json:"customerName" binding:"required,min=2"`
	CustomerPhone string      `json:"customerPhone" binding:"required,min=7"`
	ServiceIDs    []uuid.UUID `json:"serviceIds" binding:"required,min=1"`
	StartTime     time.Time   `json:"startTime" binding:"required"`
	EndTime       time.Time   `json:"endTime" binding:"required,gtfield=StartTime"`
	TotalDuration time.Duration `json:"totalDuration" binding:"required"`
}

type BookingResponse struct {
	BookingID   uuid.UUID `json:"bookingId"`
	PaymentURL  string    `json:"paymentUrl"`
}