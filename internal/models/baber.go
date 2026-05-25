package models

import (
	"github.com/google/uuid"
)

type Barber struct {
	Base
	UserID         uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	BusinessID     uuid.UUID `json:"business_id" gorm:"type:uuid;not null;index"`
	Business       Business  `json:"business" gorm:"foreignKey:BusinessID"`
	Slug           string    `json:"slug" gorm:"type:varchar(255);unique;not null"`
	Bio            string    `json:"bio" gorm:"type:text"`
	ImageURL       string    `json:"image_url" gorm:"type:text"`
	WhatsappNumber string    `json:"whatsapp_number" gorm:"type:varchar(255)"`
}
