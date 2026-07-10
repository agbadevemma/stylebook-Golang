package bookingservice

import (
	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type BookingServiceWriter interface {
	CreateBookingService(tx *gorm.DB, bookingService *models.BookingService) error
	GetBookingServicesByBookingID(tx *gorm.DB, bookingID uuid.UUID) ([]models.BookingService, error)
}
