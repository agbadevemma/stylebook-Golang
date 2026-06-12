package business

import (
	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateBusiness(tx *gorm.DB, business *models.Business) error {
	return tx.Create(business).Error
}

func (r *Repository) GetBusinessByID(id uuid.UUID) (*models.Business, error) {
	var business models.Business
	err := r.db.Where("id = ?", id).First(&business).Error
	if err != nil {
		return nil, err
	}
	return &business, nil
}

func (r *Repository) GetBusinessByUserID(id uuid.UUID) (*models.Business, error) {
	var business models.Business
	err := r.db.Preload("Owner", func(db *gorm.DB) *gorm.DB {
		return db.Omit("password")
	}).Where("owner_id = ?", id).First(&business).Error
	if err != nil {
		return nil, err
	}
	return &business, nil
}

func (r *Repository) UpdateBusinessBankDetails(
	businessID uuid.UUID,
	accountNumber string,
	transferRecipient string,
	bankName string,
	bankCode string,
) error {
	return r.db.Model(&models.Business{}).
		Where("id = ?", businessID).
		Updates(map[string]any{
			"account_number":     accountNumber,
			"bank_name":          bankName,
			"transfer_recipient": transferRecipient,
			"bank_code":          bankCode,
		}).Error
}
