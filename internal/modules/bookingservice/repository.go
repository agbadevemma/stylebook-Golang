package bookingservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateBookingService(ctx context.Context, tx *gorm.DB, bookingService *models.BookingService) error {
	return tx.WithContext(ctx).Create(bookingService).Error
}

func (r *Repository) GetBookingServicesByBookingID(ctx context.Context, tx *gorm.DB, bookingID uuid.UUID) ([]models.BookingService, error) {
	var bookingServices []models.BookingService
	err := tx.WithContext(ctx).Where("booking_id = ?", bookingID).Find(&bookingServices).Error
	return bookingServices, err
}
