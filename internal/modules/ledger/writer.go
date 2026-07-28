package ledger

import (
	"context"

	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/utils"
	"gorm.io/gorm"
)

type LedgerWriter interface {
	CreateLedger(ctx context.Context, tx *gorm.DB, ledger *models.Ledger) error
	LedgerExistsByReference(ctx context.Context, tx *gorm.DB, reference string, ledgerType utils.LedgerType) (bool, error)
}
