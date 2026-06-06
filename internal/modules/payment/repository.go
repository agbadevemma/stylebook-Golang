package payment

import (
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

func (r *Repository) CreatePayment(tx *gorm.DB, payment *models.Payment) error {
	tx = r.db.Create(payment)
	return tx.Error
}

func (r *Repository) UpdatePayment(tx *gorm.DB, payment *models.Payment) error {
	tx = r.db.Save(payment)
	return tx.Error
}

func (r *Repository) CalculateBookingAmount(bookingID uuid.UUID) (int64, error) {
	var totalAmount int64
	result := r.db.Model(&models.BookingService{}).Where("booking_id = ?", bookingID).Pluck("price * quantity", &totalAmount)
	if result.Error != nil {
		return 0, result.Error
	}
	return totalAmount, nil
}
