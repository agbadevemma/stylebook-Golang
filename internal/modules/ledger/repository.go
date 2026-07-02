package ledger

import (
	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/utils"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateLedger(tx *gorm.DB, ledger *models.Ledger) error {

	return tx.Create(ledger).Error
}

func (r *Repository) LedgerExistsByReference(
	tx *gorm.DB,
	reference string,
	ledgerType utils.LedgerType,
) (bool, error) {
	var count int64

	err := tx.Model(&models.Ledger{}).
		Where("reference = ? AND type = ?", reference, ledgerType).
		Count(&count).Error

	return count > 0, err
}
