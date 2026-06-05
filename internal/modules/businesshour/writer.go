package businesshour

import (
	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type BusinessHourWriter interface {
	CreateBusinessHourWithTx(tx *gorm.DB, hour *models.BusinessHour) error
	GetBusinessHour(businessID uuid.UUID, dayOfWeek int) (*models.BusinessHour, error)
}
