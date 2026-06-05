package service

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

func (r *Repository) CreateService(service *models.Service) error {
	return r.db.Create(service).Error
}

func (r *Repository) GetServiceByID(id uuid.UUID) (*models.Service, error) {
	var service models.Service

	if err := r.db.First(&service, id).Error; err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *Repository) UpdateService(service *models.Service) error {
	return r.db.Save(service).Error
}

func (r *Repository) DeleteService(id uuid.UUID) error {
	result := r.db.Delete(&models.Service{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *Repository) ListServices(businessID uuid.UUID) ([]models.Service, error) {
	var services []models.Service

	err := r.db.Where("business_id = ?", businessID).Find(&services).Error

	return services, err
}
