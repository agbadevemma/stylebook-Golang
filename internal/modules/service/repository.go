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
	return r.db.Delete(&models.Service{}, id).Error
}

func (r *Repository) ListServices() ([]models.Service, error) {
	var services []models.Service
	if err := r.db.Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}
