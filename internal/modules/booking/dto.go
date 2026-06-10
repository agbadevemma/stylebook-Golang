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
