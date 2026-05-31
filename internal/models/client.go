package models

import (
	"time"
)

type Client struct {
	Base
	FullName      string    `json:"fullName" gorm:"type:varchar(255);not null"`
	Phone         string    `json:"phone" gorm:"type:varchar(255);unique;not null"`
	Email         string   `json:"email" gorm:"type:varchar(255)"`
	CreatedAt     time.Time `json:"createdAt" gorm:"autoCreateTime"`
	LastBookingAt *time.Time `json:"lastBookingAt" gorm:"autoUpdateTime"`
}
