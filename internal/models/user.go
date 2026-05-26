package models

import (
	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/utils"
)

type User struct {
	Base
	Name       string         `json:"name" gorm:"type:varchar(255);not null"`
	Email      string         `json:"email" gorm:"type:varchar(255);unique;not null"`
	Phone      string         `json:"phone" gorm:"type:varchar(255);unique;not null"`
	Password   string         `json:"password" gorm:"type:text;not null"`
	Role       utils.UserRole `json:"role" gorm:"type:varchar(30);default:'barber';not null"`
	BusinessID *uuid.UUID     `json:"business_id" gorm:"type:uuid;"` //only when registering barbers
}
