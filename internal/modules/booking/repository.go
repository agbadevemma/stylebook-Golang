package booking

import (
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateBooking(tx *gorm.DB, booking *models.Booking) error {
	return tx.Create(booking).Error
}
