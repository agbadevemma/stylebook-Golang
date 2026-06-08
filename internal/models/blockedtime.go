package models

import (
	"time"

	"github.com/google/uuid"
)

type BlockedTime struct {
	Base
	StartTime  time.Time `json:"startTime" gorm:"not null;index"`
	EndTime    time.Time `json:"endTime" gorm:"not null;index"`
	Reason     string    `json:"reason" gorm:"type:text"`
	BusinessID uuid.UUID `json:"businessId" gorm:"not null;index"`
	Business   Business  `json:"business" gorm:"foreignKey:BusinessID;references:ID"`
}
