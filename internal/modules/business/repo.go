package business

import (
	"context"
	"errors"

	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateNewBusiness(ctx context.Context, business models.Business) error {
	return r.db.WithContext(ctx).Create(&business).Error
}

func (r *Repository) GetAllBusinesses(ctx context.Context) ([]models.Business, error) {
	var business []models.Business
	return business, r.db.WithContext(ctx).Find(&business).Error
}

func (r *Repository) GetBusinessesByEmail(ctx context.Context, email string) (models.Business, error) {
	var business models.Business
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&business).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Business{}, nil
		}
		return models.Business{}, err
	}

	return business, r.db.WithContext(ctx).Where("email = ?", email).First(&business).Error
}
