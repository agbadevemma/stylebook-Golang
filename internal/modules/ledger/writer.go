package ledger

import (
	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type LedgerWriter interface {
	CreateLedger(tx *gorm.DB, ledger *models.Ledger) error
	GetBusinessBalance(businessID uuid.UUID) (int64, error)
}
