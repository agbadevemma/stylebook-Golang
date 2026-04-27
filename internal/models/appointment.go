package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type appointment struct {
	Base
	BusinessID   uuid.UUID  `gorm:"type:uuid;index" json:"businessId"`
	StaffID      *uuid.UUID `gorm:"type:uuid;index" json:"staffId"`
	ServiceID    uuid.UUID  `gorm:"type:uuid;index" json:"serviceId"`
	Business     Business   `gorm:"foreignKey:BusinessID" json:"business"`
	Service      Service    `gorm:"foreignKey:ServiceID" json:"service"`
	Staff        *Staff     `gorm:"foreignKey:StaffID" json:"staff,omitempty"`
	CustomerName string     `json:"customerName"`
	Phone        string     `json:"phone"`
	Time         time.Time  `json:"time"` // recommended over string
	Status       string     `gorm:"type:varchar(20);check:status IN ('booked','cancelled','done','missed')" json:"status"`
}

func (a *appointment) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New()
	return nil
}
