package service

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateService(req *CreateServiceRequest) error {
	service := &models.Service{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Duration:    req.Duration,
	}

	err := s.repo.CreateService(service)

	if err != nil {
		log.Printf("Error creating service: %v", err.Error())
		return errors.New("failed to create service")
	}
	return nil
}

func (s *Service) GetServiceByID(id uuid.UUID) (*models.Service, error) {
	return s.repo.GetServiceByID(id)
}

func (s *Service) UpdateService(req *UpdateServiceRequest) error {
	service, err := s.repo.GetServiceByID(req.ID)

	if err != nil {
		log.Printf("Error finding service: %v", err.Error())
		return errors.New("service not found")
	}

	if req.Name != "" {
		service.Name = req.Name
	}
	if req.Description != "" {
		service.Description = req.Description
	}
	if req.Price != 0 {
		service.Price = req.Price
	}
	if req.Duration != 0 {
		service.Duration = req.Duration
	}
	if req.IsActive {
		service.IsActive = req.IsActive
	}
	err = s.repo.UpdateService(service)

	if err != nil {
		log.Printf("Error updating service: %v", err.Error())
		return errors.New("failed to update service")
	}

	return nil
}

func (s *Service) DeleteService(id uuid.UUID) error {
	_, err := s.repo.GetServiceByID(id)
	if err != nil {
		log.Printf("Error finding service: %v", err.Error())
		return errors.New("service not found")
	}

	
	err = s.repo.DeleteService(id)
	if err != nil {
		log.Printf("Error deleting service: %v", err.Error())
		return errors.New("failed to delete service")
	}
	return nil
}

func (s *Service) ListServices() ([]models.Service, error) {
	return s.repo.ListServices()
}
