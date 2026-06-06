package booking

import (
	"time"

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

func (r *Repository) CreateBooking(tx *gorm.DB, booking *models.Booking) error {
	return tx.Create(booking).Error
}

func (r *Repository) GetBookingByID(id uuid.UUID) (*models.Booking, error) {
	var booking models.Booking
	if err := r.db.First(&booking, id).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *Repository) GetAllBookings() ([]*models.Booking, error) {
	var bookings []*models.Booking
	if err := r.db.Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *Repository) GetBookings(
	businessID uuid.UUID,
	start, end time.Time,
) ([]models.Booking, error) {

	var bookings []models.Booking

	err := r.db.
		Where("business_id = ? AND start_time < ? AND end_time > ?",
			businessID, end, start).
		Find(&bookings).Error

	return bookings, err
}
