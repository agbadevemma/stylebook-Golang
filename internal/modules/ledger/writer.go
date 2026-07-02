package ledger

import (
	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/utils"
	"gorm.io/gorm"
)

type LedgerWriter interface {
	CreateLedger(tx *gorm.DB, ledger *models.Ledger) error
	LedgerExistsByReference(tx *gorm.DB, reference string, ledgerType utils.LedgerType) (bool, error)
}
