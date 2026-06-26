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

func (r *Repository) GetBusinessBookingStats(businessID uuid.UUID) (total int64, today int64, pending int64, clients int64, err error) {
	if err = r.db.Model(&models.Booking{}).
		Where("business_id = ?", businessID).
		Count(&total).Error; err != nil {
		return
	}

	todayStart := time.Now().UTC().Truncate(24 * time.Hour)
	todayEnd := todayStart.Add(24 * time.Hour)
	if err = r.db.Model(&models.Booking{}).
		Where("business_id = ? AND start_time >= ? AND start_time < ?", businessID, todayStart, todayEnd).
		Count(&today).Error; err != nil {
		return
	}

	if err = r.db.Model(&models.Booking{}).
		Where("business_id = ? AND booking_status = ?", businessID, "pending").
		Count(&pending).Error; err != nil {
		return
	}

	if err = r.db.Model(&models.Booking{}).
		Where("business_id = ?", businessID).
		Distinct("customer_email").
		Count(&clients).Error; err != nil {
		return
	}

	return
}

func (r *Repository) ListBusinessBookings(businessID uuid.UUID, params ListBookingsParams) ([]models.Booking, int64, error) {
	var bookings []models.Booking
	var total int64

	query := r.db.Model(&models.Booking{}).Where("business_id = ?", businessID)

	if params.Search != "" {
		like := "%" + params.Search + "%"
		query = query.Where(
			"customer_name ILIKE ? OR customer_email ILIKE ? OR customer_phone ILIKE ?",
			like, like, like,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Order("start_time DESC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&bookings).Error

	return bookings, total, err
}
