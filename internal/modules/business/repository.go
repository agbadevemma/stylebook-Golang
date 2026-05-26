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
