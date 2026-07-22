package booking

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
)

type BookingWriter interface {
	GetBookings(ctx context.Context, businessID uuid.UUID, start, end time.Time) ([]models.Booking, error)
	GetBookingByID(ctx context.Context, id uuid.UUID) (*models.Booking, error)
}
