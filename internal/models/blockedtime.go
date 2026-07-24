package models

import (
	"time"

	"github.com/google/uuid"
)

type BlockedTime struct {
	Base
	StartTime  time.Time `json:"startTime" gorm:"not null;index:idx_blocked_availability,priority:2"`
	EndTime    time.Time `json:"endTime" gorm:"not null;index:idx_blocked_availability,priority:3"`
	Reason     string    `json:"reason" gorm:"type:text"`
	BusinessID uuid.UUID `json:"businessId" gorm:"not null;index:idx_blocked_availability,priority:1"`
	Business   Business  `json:"business" gorm:"foreignKey:BusinessID;references:ID"`
}
