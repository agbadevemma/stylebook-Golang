package businesshour

import (
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type BusinessHourWriter interface {
	CreateBusinessHourWithTx(tx *gorm.DB, hour *models.BusinessHour) error
}
