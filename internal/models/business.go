package models

import (
	"github.com/google/uuid"
)

type Business struct {
	Base
	Name                  string        `json:"name"`
	Phone                 string        `gorm:"uniqueIndex" json:"phone"`
	WhatsAppPhoneNumberID string        `gorm:"uniqueIndex" json:"whatsappPhoneNumberId"`
	Email                 string        `gorm:"uniqueIndex" json:"email"`
	Password              string        `gorm:"not null" json:"password"`
	WhatsAppToken         string        `gorm:"not null;index" json:"-"` //must be encrypted
	BusinessTypeID        uuid.UUID     `gorm:"type:uuid;index" json:"businessTypeId"`
	BusinessType          BusinessType  `gorm:"foreignKey:BusinessTypeID" json:"businessType"`
	Status                string        `gorm:"default:active;check:status IN ('active','inactive')" json:"status"`
	Services              []Service     `gorm:"foreignKey:BusinessID" json:"services"`
	Staff                 []Staff       `gorm:"foreignKey:BusinessID" json:"staff"`
	Customers             []Customer    `gorm:"foreignKey:BusinessID" json:"customers"`
	Appointments          []Appointment `gorm:"foreignKey:BusinessID" json:"appointments"`
}
