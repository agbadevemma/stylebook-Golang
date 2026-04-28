package models

import (
	"github.com/google/uuid"
)

type Business struct {
	Base
	Name           string        `json:"name"`
	Phone          string        `gorm:"uniqueIndex" json:"phone"`
	BusinessTypeID uuid.UUID     `gorm:"type:uuid;index" json:"businessTypeId"`
	BusinessType   BusinessType  `gorm:"foreignKey:BusinessTypeID" json:"businessType"`
	PaymentStatus  string        `gorm:"default:active;check:payment_status IN ('active','overdue','suspended')" json:"paymentStatus"`
	Services       []Service     `gorm:"foreignKey:BusinessID" json:"services"`
	Staff          []Staff       `gorm:"foreignKey:BusinessID" json:"staff"`
	Customers      []Customer    `gorm:"foreignKey:BusinessID" json:"customers"`
	Appointments   []Appointment `gorm:"foreignKey:BusinessID" json:"appointments"`
}
