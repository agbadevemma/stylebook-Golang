package models

import (
	"github.com/google/uuid"
)

type Service struct {
	Base
	Name       string    `gorm:"not null;index:idx_business_service_name,unique" json:"name"`
	BusinessID uuid.UUID `gorm:"type:uuid;not null;index:idx_business_service_name,unique" json:"businessId"`
	Business   Business  `gorm:"foreignKey:BusinessID" json:"business"`
	Price      int64     `gorm:"not null" json:"price"`
	Duration   int       `gorm:"not null" json:"duration"`
}
