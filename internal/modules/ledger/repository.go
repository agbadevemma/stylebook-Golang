package ledger

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateLedger(tx *gorm.DB, ledger *models.Ledger) error {
	var last models.Ledger
	err := tx.
		Where("business_id = ?", ledger.BusinessID).
		Order("created_at DESC, id DESC").
		Limit(1).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&last).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			last.Balance = 0
		} else {
			return err
		}
	}

	newBalance := last.Balance
	switch ledger.Type {
	case utils.LedgerCredit:
		newBalance += ledger.Amount
	case utils.LedgerDebit:
		newBalance -= ledger.Amount
	}

	ledger.Balance = newBalance

	return tx.Create(ledger).Error
}

func (r *Repository) GetBusinessBalance(businessID uuid.UUID) (int64, error) {
	var balance int64

	err := r.db.Model(&models.Ledger{}).
		Select("balance").
		Where("business_id = ?", businessID).
		Order("created_at DESC, id DESC").
		Limit(1).
		Scan(&balance).Error

	return balance, err
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

func (r *Repository) GetBusinessBalanceForUpdate(
	tx *gorm.DB,
	businessID uuid.UUID,
) (int64, error) {

	var balance int64

	err := tx.
		Where("business_id = ?", businessID).
		Order("created_at DESC, id DESC").
		Limit(1).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&balance).Error

	return balance, err
}
