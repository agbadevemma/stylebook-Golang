package models

import (
	"time"

	"github.com/google/uuid"
)

type Barber struct {
	Base
	UserID         uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	ShopName       string    `json:"shop_name" gorm:"type:varchar(255);not null"`
	Slug           string    `json:"slug" gorm:"type:varchar(255);unique;not null"`
	Bio            string    `json:"bio" gorm:"type:text"`
	ImageURL       string    `json:"image_url" gorm:"type:text"`
	WhatsappNumber string    `json:"whatsapp_number" gorm:"type:varchar(255)"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
}
