package models

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	Base
	BusinessID  uuid.UUID  `gorm:"type:uuid;not null;index:idx_business_phone,uniquee" json:"businessId"`
	Business    Business   `gorm:"foreignKey:BusinessID" json:"business"`
	Name        *string    `json:"name"`
	PhoneNumber string     `gorm:"not null;index:idx_business_phone,uniquee" json:"phoneNumber"`
	VisitCount  int        `gorm:"default:0" json:"visitCount"`
	LastVisit   *time.Time `json:"lastVisit"`
}
