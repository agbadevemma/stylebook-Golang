package payment

import (
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type PaymentWriter interface {
	InitializePayment(tx *gorm.DB, booking *models.Booking) (string, error)
}
