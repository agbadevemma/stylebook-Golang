package models

import "github.com/google/uuid"

type AppointmentService struct {
	Base
	AppointmentID uuid.UUID   `gorm:"type:uuid;index;not null" json:"appointmentId"`
	ServiceID     uuid.UUID   `gorm:"type:uuid;index;not null" json:"serviceId"`
	Appointment   Appointment `gorm:"foreignKey:AppointmentID" json:"appointment"`
	Service       Service     `gorm:"foreignKey:ServiceID" json:"service"`
	ServiceName   string      `gorm:"not null" json:"serviceName"`
	Price         int64         `gorm:"not null" json:"price"`
	Duration      int         `json:"duration"`
}
