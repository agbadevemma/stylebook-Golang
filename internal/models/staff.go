package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Staff struct {
	Base
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	BusinessID uuid.UUID `gorm:"type:uuid;index" json:"businessId"`
	Business   Business  `gorm:"foreignKey:BusinessID" json:"business"`
	Name       string    `json:"name"`
	Phone      string    `gorm:"uniqueIndex" json:"phone"`
	Email      string    `gorm:"uniqueIndex" json:"email"`
}

func (s *Staff) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	return nil
}
