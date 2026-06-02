package models

import "github.com/google/uuid"

type BusinessHour struct {
	Base
	DayOfWeek int `json:"dayOfWeek" gorm:"not null"`
	// 0 = Sunday
	// 1 = Monday
	// 2 = Tuesday
	// 3 = Wednesday
	// 4 = Thursday
	// 5 = Friday
	// 6 = Saturday
	IsOpen     bool      `json:"isOpen" gorm:"not null;default:true"`
	StartTime  string    `json:"startTime" gorm:"type:varchar(10);not null"`
	EndTime    string    `json:"endTime" gorm:"type:varchar(10);not null"`
	BusinessID uuid.UUID `json:"businessId" gorm:"not null;index"`
	Business   Business  `json:"business" gorm:"foreignKey:BusinessID;references:ID"`
}
