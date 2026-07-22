package payment

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/utils"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreatePayment(ctx context.Context, tx *gorm.DB, payment *models.Payment) error {
	return tx.WithContext(ctx).Create(payment).Error
}

func (r *Repository) UpdatePayment(ctx context.Context, tx *gorm.DB, payment *models.Payment) error {

	return tx.WithContext(ctx).Save(payment).Error
}

func (r *Repository) CalculateBookingAmount(ctx context.Context, tx *gorm.DB, bookingID uuid.UUID) (int64, error) {
	var totalAmount int64
	err := tx.WithContext(ctx).
		Model(&models.BookingService{}).
		Where("booking_id = ?", bookingID).
		Select("COALESCE(SUM(price), 0)").
		Scan(&totalAmount).Error

	if err != nil {
		return 0, err
	}
	return totalAmount, nil
}

func (r *Repository) GetPaymentByBookingID(ctx context.Context, tx *gorm.DB, bookingID uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Where("booking_id = ?", bookingID).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *Repository) GetPaymentByReference(ctx context.Context, tx *gorm.DB, reference string) (*models.Payment, error) {
	var payment models.Payment
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Where("reference = ?", reference).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *Repository) UpdateBookingStatus(ctx context.Context,
	tx *gorm.DB,
	bookingID uuid.UUID,
	bookingStatus utils.BookingStatus,
	paymentStatus utils.PaymentStatus,
) error {

	return tx.WithContext(ctx).
		Model(&models.Booking{}).
		Where("id = ?", bookingID).
		Updates(map[string]any{
			"booking_status": bookingStatus,
			"payment_status": paymentStatus,
		}).
		Error
}

func (r *Repository) GetPaymentsByBusiness(
	ctx context.Context,
	businessID uuid.UUID,
	params ListTransactionsParams,
) ([]models.Payment, int64, error) {

	var payments []models.Payment
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Payment{}). // This line is already correct, no change needed
									Where("business_id = ?", businessID)

	if params.Search != "" {
		like := "%" + params.Search + "%"
		query = query.Where(
			"name ILIKE ? OR email ILIKE ? OR reference ILIKE ? OR status ILIKE ?",
			like, like, like, like,
		)
	}

	// Add date range filter
	if params.StartDate != "" {
		startTime, err := time.Parse("2006-01-02", params.StartDate)
		if err == nil {
			query = query.Where("created_at >= ?", startTime)
		}
	}

	if params.EndDate != "" {
		endTime, err := time.Parse("2006-01-02", params.EndDate)
		if err == nil {
			// Add 24 hours to include the entire end date
			endTime = endTime.Add(24 * time.Hour)
			query = query.Where("created_at < ?", endTime)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Order("created_at DESC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&payments).Error

	if err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

func (r *Repository) GetBookingByID(ctx context.Context, bookingID uuid.UUID) (*models.Booking, error) {
	var booking models.Booking
	if err := r.db.WithContext(ctx).Where("id = ?", bookingID).First(&booking).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *Repository) GetBookingServices(ctx context.Context, bookingID uuid.UUID) ([]models.BookingService, error) {
	var services []models.BookingService
	if err := r.db.WithContext(ctx).Where("booking_id = ?", bookingID).Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}
