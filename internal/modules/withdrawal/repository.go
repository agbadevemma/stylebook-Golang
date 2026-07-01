package withdrawal

import (
	"github.com/google/uuid"
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

func (r *Repository) GetWithdrawalByPaystackTransferCode(transferCode string) (*models.Withdrawal, error) {
	var withdrawal models.Withdrawal

	err := r.db.
		Where("paystack_transfer_code = ?", transferCode).
		First(&withdrawal).Error

	if err != nil {
		return nil, err
	}

	return &withdrawal, nil
}

func (r *Repository) GetBusinessWithdrawals(businessID uuid.UUID) ([]models.Withdrawal, error) {
	var withdrawals []models.Withdrawal

	err := r.db.
		Where("business_id = ?", businessID).
		Order("created_at DESC").
		Find(&withdrawals).Error

	if err != nil {
		return nil, err
	}

	return withdrawals, nil
}

func (r *Repository) UpdateWithdrawalTransferCode(
	withdrawalID uuid.UUID,
	transferCode string,
	status utils.WithdrawalStatus,
) error {
	return r.db.
		Model(&models.Withdrawal{}).
		Where("id = ?", withdrawalID).
		Updates(map[string]any{
			"paystack_transfer_code": transferCode,
			"status":                 status,
		}).
		Error
}

func (r *Repository) UpdateWithdrawalStatus(
	withdrawalID uuid.UUID,
	status utils.WithdrawalStatus,
) error {
	return r.db.
		Model(&models.Withdrawal{}).
		Where("id = ?", withdrawalID).
		Update("status", status).
		Error
}
