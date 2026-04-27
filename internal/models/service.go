package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	Base
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name       string    `gorm:"uniqueIndex" json:"name"`
	BusinessID uuid.UUID `gorm:"type:uuid;index" json:"businessId"`
	Business   Business  `gorm:"foreignKey:BusinessID" json:"business"`
	Price      int64     `json:"price"`
	Duration   int       `json:"duration"`
}

func (s *Service) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	return nil
}
