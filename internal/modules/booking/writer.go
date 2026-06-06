package booking

import (
	"time"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
)

type BookingWriter interface {
	GetBookings(businessID uuid.UUID, start, end time.Time) ([]models.Booking, error)
	GetBookingByID(id uuid.UUID) (*models.Booking, error)
}
