package payment

import (
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

func (r *Repository) CreatePayment(tx *gorm.DB, payment *models.Payment) error {
	return tx.Create(payment).Error
}

func (r *Repository) UpdatePayment(tx *gorm.DB, payment *models.Payment) error {

	return tx.Save(payment).Error
}

func (r *Repository) CalculateBookingAmount(tx *gorm.DB, bookingID uuid.UUID) (int64, error) {
	var totalAmount int64
	err := tx.
		Model(&models.BookingService{}).
		Where("booking_id = ?", bookingID).
		Select("COALESCE(SUM(price), 0)").
		Scan(&totalAmount).Error

	if err != nil {
		return 0, err
	}
	return totalAmount, nil
}

func (r *Repository) GetPaymentByBookingID(tx *gorm.DB, bookingID uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	if tx == nil {
		tx = r.db
	}
	if err := tx.Where("booking_id = ?", bookingID).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *Repository) GetPaymentByReference(tx *gorm.DB, reference string) (*models.Payment, error) {
	var payment models.Payment
	if tx == nil {
		tx = r.db
	}
	if err := tx.Where("reference = ?", reference).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *Repository) UpdateBookingStatus(
	tx *gorm.DB,
	bookingID uuid.UUID,
	bookingStatus utils.BookingStatus,
	paymentStatus utils.PaymentStatus,
) error {

	return tx.
		Model(&models.Booking{}).
		Where("id = ?", bookingID).
		Updates(map[string]any{
			"booking_status": bookingStatus,
			"payment_status": paymentStatus,
		}).
		Error
}

func (r *Repository) GetPaymentsByBusiness(
	businessID uuid.UUID,
	params ListTransactionsParams,
) ([]models.Payment, int64, error) {

	var payments []models.Payment
	var total int64

	query := r.db.Model(&models.Payment{}).
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
