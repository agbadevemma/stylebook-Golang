package business

import (
	"context"

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

func (r *Repository) GetBusiness(ctx context.Context) (models.Business, error) {
	var business models.Business
	return business, r.db.WithContext(ctx).First(&business).Error
}

func (r *Repository) GetBusinessesByEmail(ctx context.Context, email string) (models.Business, error) {
	var business models.Business
	return business, r.db.WithContext(ctx).Where("email = ?", email).First(&business).Error
}
