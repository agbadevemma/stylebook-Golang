package models

import "github.com/google/uuid"

type Business struct {
	Base
	Name    string    `json:"name" gorm:"type:varchar(255);not null"`
	OwnerID uuid.UUID `json:"owner_id" gorm:"type:varchar(255);not null;index"`
	Owner   User      `json:"owner" gorm:"foreignKey:OwnerID"`
}
