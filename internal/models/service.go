package models

import "github.com/google/uuid"

type Service struct {
	Base

	Name        string `json:"name" gorm:"type:varchar(255);not null"`
	Description string `json:"description" gorm:"type:text"`

	// price in kobo/cents
	Price int64 `json:"price" gorm:"not null"`

	// duration in minutes
	Duration int `json:"duration" gorm:"not null"`

	IsActive bool `json:"isActive" gorm:"default:true"`

	BusinessID uuid.UUID `json:"businessId" gorm:"not null;index"`
	Business   Business  `json:"business" gorm:"foreignKey:BusinessID;references:ID"`
}
