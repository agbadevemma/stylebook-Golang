package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Business struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	Name           string       `json:"name"`
	Phone          string       `gorm:"uniqueIndex" json:"phone"`
	BusinessTypeID uuid.UUID    `gorm:"type:uuid;index" json:"businessTypeId"`
	BusinessType   BusinessType `gorm:"foreignKey:BusinessTypeID" json:"businessType"`
	Plan           string       `json:"plan"`
	Status         string       `gorm:"index" json:"status"`
	PaidUntil      time.Time    `json:"paidUntil"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
}

func (b *Business) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}
