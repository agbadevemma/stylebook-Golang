package models

import (
	"time"

	"github.com/google/uuid"
)

type Appointment struct {
	Base

	BusinessID uuid.UUID `gorm:"type:uuid;index;not null" json:"businessId"`

	StaffID *uuid.UUID `gorm:"type:uuid;index" json:"staffId"`

	Business Business `gorm:"foreignKey:BusinessID" json:"business"`

	Staff *Staff `gorm:"foreignKey:StaffID" json:"staff,omitempty"`

	CustomerName string `gorm:"not null" json:"customerName"`

	PhoneNumber string `gorm:"not null" json:"phoneNumber"`

	ScheduledTime time.Time `gorm:"not null" json:"scheduledTime"`

	TotalPrice int64 `gorm:"default:0" json:"totalPrice"`

	CommissionAmount int64 `gorm:"default:0" json:"commissionAmount"`

	BusinessPayout int64 `gorm:"default:0" json:"businessPayout"`

	Status string `gorm:"type:varchar(20);default:'booked';check:status IN ('booked','cancelled','done','missed')" json:"status"`

	Services []AppointmentService `gorm:"foreignKey:AppointmentID" json:"services"`
}
