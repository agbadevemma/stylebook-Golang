package blockedtime

import (
	"time"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
)

type BlockedTimesWriter interface {
	GetBlockedTimes(businessID uuid.UUID, start, end time.Time) ([]models.BlockedTime, error)
}
