package ledger

import (
	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/utils"
	"gorm.io/gorm"
)

type LedgerWriter interface {
	CreateLedger(tx *gorm.DB, ledger *models.Ledger) error
	GetBusinessBalance(businessID uuid.UUID) (int64, error)
	LedgerExistsByReference(tx *gorm.DB, reference string, ledgerType utils.LedgerType) (bool, error)
	GetBusinessBalanceForUpdate(tx *gorm.DB, businessID uuid.UUID) (int64, error)
}
