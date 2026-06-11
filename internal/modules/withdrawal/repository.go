package withdrawal

import (
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateWithdrawal(tx *gorm.DB, withdrawal *models.Withdrawal) error {
	return tx.Create(withdrawal).Error
}

func (r *Repository) UpdateWithdrawal(tx *gorm.DB, withdrawal *models.Withdrawal) error {
	return tx.Save(withdrawal).Error
}

func (r *Repository) GetWithdrawalByReference(transferReference string) (*models.Withdrawal, error) {
	var withdrawal models.Withdrawal

	err := r.db.
		Where("transfer_reference = ?", transferReference).
		First(&withdrawal).Error

	if err != nil {
		return nil, err
	}

	return &withdrawal, nil
}
