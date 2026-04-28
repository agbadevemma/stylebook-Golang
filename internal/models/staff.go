package models

import (
	"github.com/google/uuid"
)

type Staff struct {
	Base
	BusinessID uuid.UUID `gorm:"type:uuid;not null;index:idx_business_phone,unique;index:idx_business_email,unique" json:"businessId"`
	Business   Business  `gorm:"foreignKey:BusinessID" json:"business"`
	Name       string    `gorm:"not null" json:"name"`
	Role       string    `json:"role"`
	Phone      string    `gorm:"not null;index:idx_business_phone,unique" json:"phone"`
	Email      string    `gorm:"not null;index:idx_business_email,unique" json:"email"`
}
