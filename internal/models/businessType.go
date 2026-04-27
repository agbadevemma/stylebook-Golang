package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusinessType struct {
	Base

	Name string    `gorm:"unique;not null" json:"name"`
}

func (bt *BusinessType) BeforeCreate(tx *gorm.DB) (err error) {
	bt.ID = uuid.New()
	return
}
