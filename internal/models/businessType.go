package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusinessType struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name string    `gorm:"unique;not null" json:"name"`
}

func (bt *BusinessType) BeforeCreate(tx *gorm.DB) (err error) {
	bt.ID = uuid.New()
	return
}
