package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) error {
	now := time.Now().UTC()

	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	b.CreatedAt = now
	b.UpdatedAt = now
	return nil
}

func (b *Base) BeforeUpdate(tx *gorm.DB) error {
	b.UpdatedAt = time.Now().UTC()
	return nil
}
