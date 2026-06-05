package bookingservice

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

func (r *Repository) CreateBookingService(tx *gorm.DB, bookingService *models.BookingService) error {
	return tx.Create(bookingService).Error
}
