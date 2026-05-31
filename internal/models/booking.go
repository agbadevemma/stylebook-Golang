package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/utils"
)

type Booking struct {
	Base
	ClientID      uuid.UUID           `json:"clientId" gorm:"type:uuid;not null;index"`
	Client        Client              `json:"client" gorm:"foreignKey:ClientID"`
	BusinessID    uuid.UUID           `json:"businessId" gorm:"type:uuid;not null;index"`
	Business      Business            `json:"business" gorm:"foreignKey:BusinessID"`
	BarberID      *uuid.UUID          `json:"barberId" gorm:"type:uuid"`
	Barber        *User               `json:"user" gorm:"foreignKey:BarberID"`
	ServiceID     uuid.UUID           `json:"serviceId" gorm:"type:uuid;not null;index"`
	Service       Service             `json:"service" gorm:"foreignKey:ServiceID"`
	BookingTime   time.Time           `json:"bookingTime" gorm:"not null"`
	BookingStatus utils.BookingStatus `json:"bookingStatus" gorm:"type:varchar(20);default:pending;not null"`
	PaymentStatus utils.PaymentStatus `json:"paymentStatus" gorm:"type:varchar(20);default:pending;not null"`
}
