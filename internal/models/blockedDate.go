package models

import "github.com/google/uuid"

type BlockedDate struct {
	Base
	StartTime  string    `json:"start_time" gorm:"type:varchar(10);not null"`
	EndTime    string    `json:"end_time" gorm:"type:varchar(10);not null"`
	Reason     string    `json:"reason" gorm:"type:text"`
	BusinessID uuid.UUID `json:"business_id" gorm:"not null;index"`
	Business   Business  `json:"business" gorm:"foreignKey:BusinessID;references:ID"`
}
