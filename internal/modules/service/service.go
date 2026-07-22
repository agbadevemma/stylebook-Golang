package service

import (
	"context"
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

func (s *Service) CreateService(ctx context.Context, req *CreateServiceRequest, businessID uuid.UUID) error {
	service := &models.Service{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Duration:    req.Duration,
		BusinessID:  businessID,
	}

	err := s.repo.CreateService(ctx, service)

	if err != nil {
		log.Printf("Error creating service: %v", err.Error())
		return errors.New("failed to create service")
	}
	return nil
}

func (s *Service) ListServices(ctx context.Context, businessID uuid.UUID) ([]models.Service, error) {
	fmt.Println("service nil:", s == nil)
	fmt.Println("businessWriter nil:", s.businessWriter == nil)
	_, err := s.businessWriter.GetBusinessByID(ctx, s.repo.db, businessID)
	if err != nil {
		return nil, err
	}
	return s.repo.ListServices(ctx, businessID)
}

func (s *Service) GetServiceByID(ctx context.Context, serviceId uuid.UUID) (*models.Service, error) {

	return s.repo.GetServiceByID(ctx, s.repo.db, serviceId)
}

func (s *Service) UpdateService(ctx context.Context, req *UpdateServiceRequest, serviceId uuid.UUID) error {
	service, err := s.repo.GetServiceByID(ctx, s.repo.db, serviceId)
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

	err = s.repo.UpdateService(ctx, service)

	if err != nil {
		log.Printf("Error updating service: %v", err.Error())
		return errors.New("failed to update service")
	}

	return nil
}

func (s *Service) DeleteService(ctx context.Context, serviceId uuid.UUID) error {

	return s.repo.DeleteService(ctx, serviceId)
}
