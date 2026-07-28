package service

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

func (r *Repository) CreateService(ctx context.Context, service *models.Service) error {
	return r.db.WithContext(ctx).Create(service).Error
}

func (r *Repository) GetServiceByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*models.Service, error) {
	var service models.Service
	db := tx
	if db == nil {
		db = r.db
	}

	if err := db.WithContext(ctx).First(&service, id).Error; err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *Repository) UpdateService(ctx context.Context, service *models.Service) error {
	return r.db.WithContext(ctx).Save(service).Error
}

func (r *Repository) DeleteService(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Service{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *Repository) ListServices(ctx context.Context, businessID uuid.UUID) ([]models.Service, error) {
	var services []models.Service

	err := r.db.WithContext(ctx).Where("business_id = ?", businessID).Find(&services).Error

	return services, err
}
