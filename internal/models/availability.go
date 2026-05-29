package models

import "github.com/google/uuid"

type Availability struct {
	Base
	DayOfWeek  string    `json:"dayofweek"  gorm:"type:varchar(20);not null"`
	StartTime  string    `json:"start_time" gorm:"type:varchar(10);not null"`
	EndTime    string    `json:"end_time" gorm:"type:varchar(10);not null"`
	BusinessID uuid.UUID `json:"business_id" gorm:"not null;index"`
	Business   Business  `json:"business" gorm:"foreignKey:BusinessID;references:ID"`
}
