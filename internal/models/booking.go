package models

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	Base
	GuestID       uuid.UUID `json:"guest_id" gorm:"type:uuid;not null"`
	BarberID      uuid.UUID `json:"barber_id" gorm:"type:uuid;not null"`
	ServiceID     uuid.UUID `json:"service_id" gorm:"type:uuid;not null"`
	BookingTime   time.Time `json:"booking_time" gorm:"not null"`
	BookingStatus string    `json:"booking_status" gorm:"type:enum('pending', 'confirmed', 'cancelled', 'completed');default:'pending';not null"`
	PaymentStatus string    `json:"payment_status" gorm:"type:enum('pending', 'paid', 'failed');default:'pending';not null"`
}
