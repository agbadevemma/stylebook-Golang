package booking

import (
	"context"
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

func (r *Repository) CreateBooking(ctx context.Context, tx *gorm.DB, booking *models.Booking) error {
	return tx.WithContext(ctx).Create(booking).Error
}

func (r *Repository) GetBookingByID(ctx context.Context, id uuid.UUID) (*models.Booking, error) {
	var booking models.Booking
	if err := r.db.WithContext(ctx).First(&booking, id).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *Repository) GetAllBookings(ctx context.Context) ([]*models.Booking, error) {
	var bookings []*models.Booking
	if err := r.db.WithContext(ctx).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *Repository) GetBookings(ctx context.Context,
	businessID uuid.UUID,
	start, end time.Time,
) ([]models.Booking, error) {

	var bookings []models.Booking

	err := r.db.WithContext(ctx).
		Where("business_id = ? AND start_time < ? AND end_time > ?",
			businessID, end, start).
		Find(&bookings).Error

	return bookings, err
}

func (r *Repository) GetBusinessBookingStats(ctx context.Context, businessID uuid.UUID) (total int64, today int64, pending int64, clients int64, err error) {
	if err = r.db.WithContext(ctx).Model(&models.Booking{}).
		Where("business_id = ? AND booking_status != ?", businessID, "completed").
		Count(&total).Error; err != nil {
		return
	}

	now := time.Now().UTC()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	todayEnd := todayStart.Add(24 * time.Hour)
	if err = r.db.WithContext(ctx).Model(&models.Booking{}).
		Where("business_id = ? AND start_time >= ? AND start_time < ? AND booking_status != ?", businessID, todayStart, todayEnd, "completed").
		Count(&today).Error; err != nil {
		return
	}

	if err = r.db.WithContext(ctx).Model(&models.Booking{}).
		Where("business_id = ? AND booking_status = ?", businessID, "pending").
		Count(&pending).Error; err != nil {
		return
	}

	if err = r.db.WithContext(ctx).Model(&models.Booking{}).
		Where("business_id = ? AND booking_status != ?", businessID, "completed").
		Distinct("customer_email").
		Count(&clients).Error; err != nil {
		return
	}

	return
}

func (r *Repository) ListBusinessBookings(ctx context.Context, businessID uuid.UUID, params ListBookingsParams) ([]models.Booking, int64, error) {
	var bookings []models.Booking
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Booking{}).Where("business_id = ? AND booking_status != ?", businessID, "completed")

	if params.Search != "" {
		like := "%" + params.Search + "%"
		query = query.Where(
			"customer_name ILIKE ? OR customer_email ILIKE ? OR customer_phone ILIKE ?",
			like, like, like,
		)
	}

	// Add date range filter
	if params.StartDate != "" {
		startTime, err := time.Parse("2006-01-02", params.StartDate)
		if err == nil {
			query = query.Where("start_time >= ?", startTime)
		}
	}

	if params.EndDate != "" {
		endTime, err := time.Parse("2006-01-02", params.EndDate)
		if err == nil {
			// Add 24 hours to include the entire end date
			endTime = endTime.Add(24 * time.Hour)
			query = query.Where("start_time < ?", endTime)
		}
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

func (r *Repository) UpdateBooking(ctx context.Context, tx *gorm.DB, booking *models.Booking) error {
	return tx.WithContext(ctx).Model(booking).Updates(booking).Error
}
