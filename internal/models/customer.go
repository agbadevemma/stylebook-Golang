package models

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	Base
	BusinessID uuid.UUID  `gorm:"type:uuid;uniqueIndex:idx_business_phone" json:"businessId"`
	Business   Business   `gorm:"foreignKey:BusinessID" json:"business"`
	Name       *string    `json:"name"`
	Phone      string     `gorm:"uniqueIndex:idx_business_phone" json:"phone"`
	VisitCount int        `gorm:"default:0" json:"visitCount"`
	LastVisit  *time.Time `json:"lastVisit"`
}
