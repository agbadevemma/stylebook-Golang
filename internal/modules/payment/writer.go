package payment

import (
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
	"context"
)

type PaymentWriter interface {
	InitializePayment(ctx context.Context, tx *gorm.DB, booking *models.Booking) (string, error)
}
