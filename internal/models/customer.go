package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	Base
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	BusinessID uuid.UUID  `gorm:"type:uuid;uniqueIndex:idx_business_phone" json:"businessId"`
	Business   Business   `gorm:"foreignKey:BusinessID" json:"business"`
	Name       *string    `json:"name"` // optional
	Phone      string     `gorm:"uniqueIndex:idx_business_phone" json:"phone"`
	VisitCount int        `gorm:"default:0" json:"visitCount"`
	LastVisit  *time.Time `json:"lastVisit"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

func (c *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return nil
}
