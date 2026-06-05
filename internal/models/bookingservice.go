package models

import (
	"github.com/google/uuid"
)

type BookingService struct {
	Base
	BookingID uuid.UUID `json:"bookingId" gorm:"type:uuid;not null;index"`
	ServiceID uuid.UUID `json:"serviceId" gorm:"type:uuid;not null;index"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	Price     int64     `json:"price" gorm:"type:bigint;not null"`
	Duration  int       `json:"duration" gorm:"type:int;not null"`
}
