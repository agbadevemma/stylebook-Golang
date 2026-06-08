package models

import "github.com/google/uuid"

type Business struct {
	Base
	Name     string    `json:"name" gorm:"type:varchar(255);not null"`
	Timezone string    `gorm:"type:varchar(50);default:'Africa/Lagos'"`
	OwnerID  uuid.UUID `json:"ownerId" gorm:"type:varchar(255);not null;index"`
	Owner    User      `json:"owner" gorm:"foreignKey:OwnerID"`
}
