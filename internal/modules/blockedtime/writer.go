package blockedtime

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type BlockedTimesWriter interface {
	GetBlockedTimes(ctx context.Context, tx *gorm.DB, businessID uuid.UUID, start, end time.Time) ([]models.BlockedTime, error)
}
