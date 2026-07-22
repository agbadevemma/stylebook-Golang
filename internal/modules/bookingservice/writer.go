package bookingservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type BookingServiceWriter interface {
	CreateBookingService(ctx context.Context, tx *gorm.DB, bookingService *models.BookingService) error
	GetBookingServicesByBookingID(ctx context.Context, tx *gorm.DB, bookingID uuid.UUID) ([]models.BookingService, error)
}
