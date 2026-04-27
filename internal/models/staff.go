package models

import (
	"github.com/google/uuid"
)

type Staff struct {
	Base
	BusinessID uuid.UUID `gorm:"uniqueIndex:idx_business_phone;uniqueIndex:idx_business_email" json:"businessId"`
	Business   Business  `gorm:"foreignKey:BusinessID" json:"business"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	Phone      string    `gorm:"uniqueIndex:idx_business_phone" json:"phone"`
	Email      string    `gorm:"uniqueIndex:idx_business_email" json:"email"`
}
