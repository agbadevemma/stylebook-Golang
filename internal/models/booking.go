package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/utils"
)

type Booking struct {
	Base
	CustomerName  string              `json:"customerName" gorm:"type:varchar(100);not null"`
	CustomerPhone string              `json:"customerPhone" gorm:"type:varchar(20);not null"`
	CustomerEmail string              `json:"customerEmail" gorm:"type:varchar(100);not null"`
	BusinessID    uuid.UUID           `json:"businessId" gorm:"type:uuid;not null;index"`
	TotalDuration int                 `json:"totalDuration" gorm:"type:int;not null"`
	StartTime     time.Time           `json:"startTime" gorm:"type:timestamp;not null;index"`
	EndTime       time.Time           `json:"endTime" gorm:"type:timestamp;not null;index"`
	BookingStatus utils.BookingStatus `json:"bookingStatus" gorm:"type:varchar(20);default:pending;not null"`
	PaymentStatus utils.PaymentStatus `json:"paymentStatus" gorm:"type:varchar(20);default:pending;not null"`
}
