package bookingservice

import (
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type BookingServiceWriter interface {
	CreateBookingService(tx *gorm.DB, bookingService *models.BookingService) error
}
