package models

import (
	"github.com/google/uuid"
)

type Service struct {
	Base
	Name       string    `gorm:"uniqueIndex:idx_business_service_name" json:"name"`
	BusinessID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_business_service_name" json:"businessId"`
	Business   Business  `gorm:"foreignKey:BusinessID" json:"business"`
	Price      int64     `json:"price"`
	Duration   int       `json:"duration"`
}
