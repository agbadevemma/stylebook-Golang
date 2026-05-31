package models

import "github.com/google/uuid"

type Availability struct {
	Base
	DayOfWeek  string    `json:"dayOfWeek"  gorm:"type:varchar(20);not null"`
	StartTime  string    `json:"startTime" gorm:"type:varchar(10);not null"`
	EndTime    string    `json:"endTime" gorm:"type:varchar(10);not null"`
	BusinessID uuid.UUID `json:"businessId" gorm:"not null;index"`
	Business   Business  `json:"business" gorm:"foreignKey:BusinessID;references:ID"`
}
