package business

import (
	"context"

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

func (r *Repository) CreateBusiness(ctx context.Context, tx *gorm.DB, business *models.Business) error {
	return tx.WithContext(ctx).Create(business).Error
}

func (r *Repository) GetBusinessByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*models.Business, error) {
	var business models.Business
	if tx == nil {
		tx = r.db
	}
	err := tx.WithContext(ctx).Where("id = ?", id).First(&business).Error
	if err != nil {
		return nil, err
	}
	return &business, nil
}

func (r *Repository) GetBusinessByUserID(ctx context.Context, id uuid.UUID) (*models.Business, error) {
	var business models.Business
	err := r.db.Preload("Owner", func(db *gorm.DB) *gorm.DB {
		return db.Omit("password")
	}).WithContext(ctx).Where("owner_id = ?", id).First(&business).Error
	if err != nil {
		return nil, err
	}
	return &business, nil
}

func (r *Repository) UpdateBusinessDetails(
	ctx context.Context,
	businessID uuid.UUID,
	name string,
	timezone string,
) error {
	return r.db.WithContext(ctx).Model(&models.Business{}).
		Where("id = ?", businessID).
		Updates(map[string]any{
			"name":     name,
			"timezone": timezone,
		}).Error
}

func (r *Repository) UpdateBusinessBankDetails(
	ctx context.Context,
	businessID uuid.UUID,
	accountName,
	accountNumber,
	subaccountCode,
	bankName,
	bankCode string,
) error {
	return r.db.WithContext(ctx).Model(&models.Business{}).
		Where("id = ?", businessID).
		Updates(map[string]any{
			"account_name":    accountName,
			"account_number":  accountNumber,
			"bank_name":       bankName,
			"bank_code":       bankCode,
			"subaccount_code": subaccountCode,
		}).Error
}
