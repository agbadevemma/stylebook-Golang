package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/utils"
)

type Booking struct {
	Base
	ClientID      uuid.UUID           `json:"client_id" gorm:"type:uuid;not null;index"`
	Client        Client              `json:"client" gorm:"foreignKey:ClientID"`
	BusinessID    uuid.UUID           `json:"business_id" gorm:"type:uuid;not null;index"`
	Business      Business            `json:"business" gorm:"foreignKey:BusinessID"`
	BarberID      *uuid.UUID          `json:"barber_id" gorm:"type:uuid"`
	Barber        *User               `json:"user" gorm:"foreignKey:BarberID"`
	ServiceID     uuid.UUID           `json:"service_id" gorm:"type:uuid;not null;index"`
	Service       Service             `json:"service" gorm:"foreignKey:ServiceID"`
	BookingTime   time.Time           `json:"booking_time" gorm:"not null"`
	BookingStatus utils.BookingStatus `json:"booking_status" gorm:"type:varchar(20);default:pending;not null"`
	PaymentStatus utils.PaymentStatus `json:"payment_status" gorm:"type:varchar(20);default:pending;not null"`
}
