package businesshour

import (
	"context"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type BusinessHourWriter interface {
	CreateBusinessHourWithTx(ctx context.Context, tx *gorm.DB, hour *models.BusinessHour) error
	GetBusinessHour(ctx context.Context, tx *gorm.DB, businessID uuid.UUID, dayOfWeek int) (*models.BusinessHour, error)
}
