package service

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
)

type Service struct {
	repo           *Repository
	businessWriter business.BusinessWriter
}

func NewService(repo *Repository, bw business.BusinessWriter) *Service {
	return &Service{repo: repo, businessWriter: bw}
}

func (s *Service) CreateService(req *CreateServiceRequest, businessID uuid.UUID) error {
	service := &models.Service{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Duration:    req.Duration,
		BusinessID:  businessID,
	}

	err := s.repo.CreateService(service)

	if err != nil {
		log.Printf("Error creating service: %v", err.Error())
		return errors.New("failed to create service")
	}
	return nil
}

func (s *Service) ListServices(businessID uuid.UUID) ([]models.Service, error) {
	fmt.Println("service nil:", s == nil)
	fmt.Println("businessWriter nil:", s.businessWriter == nil)
	_, err := s.businessWriter.GetBusinessByID(businessID)
	if err != nil {
		return nil, err
	}
	return s.repo.ListServices(businessID)
}

func (s *Service) GetServiceByID(serviceId uuid.UUID) (*models.Service, error) {

	return s.repo.GetServiceByID(serviceId)
}

func (s *Service) UpdateService(req *UpdateServiceRequest, serviceId uuid.UUID) error {
	service, err := s.repo.GetServiceByID(serviceId)
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
	if req.IsActive != service.IsActive {
		service.IsActive = req.IsActive
	}

	err = s.repo.UpdateService(service)

	if err != nil {
		log.Printf("Error updating service: %v", err.Error())
		return errors.New("failed to update service")
	}

	return nil
}

func (s *Service) DeleteService(serviceId uuid.UUID) error {

	return s.repo.DeleteService(serviceId)
}
