package booking

import (
	"time"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/utils"
)

type CreateBookingRequest struct {
	FullName      string              `json:"name" binding:"required"`
	Email         string              `json:"email" binding:"required,email"`
	Phone         string              `json:"phone" binding:"required"`
	ClientID      uuid.UUID           `json:"clientId" binding:"required"`
	BusinessID    uuid.UUID           `json:"businessId" binding:"required"`
	ServiceID     uuid.UUID           `json:"serviceId" binding:"required"`
	BarberID      *uuid.UUID          `json:"barberId"`
	BookingTime   time.Time           `json:"bookingTime" binding:"required"`
	PaymentStatus utils.PaymentStatus `json:"paymentStatus" binding:"required,oneof=pending"`
}
