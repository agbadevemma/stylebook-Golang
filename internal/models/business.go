package models

import (
	"time"

	"github.com/google/uuid"
)

type Business struct {
	Base
	Name           string       `json:"name"`
	Phone          string       `gorm:"uniqueIndex" json:"phone"`
	BusinessTypeID uuid.UUID    `gorm:"type:uuid;index" json:"businessTypeId"`
	BusinessType   BusinessType `gorm:"foreignKey:BusinessTypeID" json:"businessType"`
	Plan           string       `json:"plan"`
	Status         string       `gorm:"index" json:"status"`
	PaidUntil      time.Time    `json:"paidUntil"`
}
