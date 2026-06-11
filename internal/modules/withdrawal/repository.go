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

func (r *Repository) CreateWithdrawal(withdrawal *models.Withdrawal) error {
	return r.db.Create(withdrawal).Error
}

func (r *Repository) UpdateWithdrawal(businessID uuid.UUID, transferReference string, status utils.WithdrawalStatus) error {
	return r.db.Model(&models.Withdrawal{}).Where("id = ?", businessID).
		Updates(map[string]any{
			"transfer_reference": transferReference,
			"status":             status,
		}).Error
}
