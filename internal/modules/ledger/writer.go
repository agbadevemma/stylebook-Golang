package ledger

import (
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type LedgerWriter interface {
	CreateLedger(tx *gorm.DB, ledger *models.Ledger) error
}
