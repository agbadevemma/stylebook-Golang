package models

import (
	"github.com/google/uuid"
)

type WhatsAppSession struct {
	Base

	Phone string `gorm:"uniqueIndex" json:"phone"`

	BusinessID uuid.UUID `gorm:"type:uuid;index" json:"businessId"`
	Business   Business  `gorm:"foreignKey:BusinessID" json:"business"`
	State      string    `json:"state"`
	Payload    string    `json:"payload"` // store temporary booking data (JSON)
}
