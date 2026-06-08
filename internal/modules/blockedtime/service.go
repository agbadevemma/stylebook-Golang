package blockedtime

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
)

type Service struct {
	repo           *Repository
	businessWriter business.BusinessWriter
}

func NewService(repo *Repository, businessWriter business.BusinessWriter) *Service {
	return &Service{repo: repo, businessWriter: businessWriter}
}

func (s *Service) CreateBlockedTime(
	req *CreateBlockedTimeRequest,
	businessID uuid.UUID,
) error {

	business, err := s.businessWriter.GetBusinessByID(businessID)
	if err != nil {
		return err
	}

	loc, err := time.LoadLocation(business.Timezone)
	if err != nil {
		return err
	}

	start := req.StartTime.In(loc)
	end := req.EndTime.In(loc)

	blockedTime := &models.BlockedTime{
		BusinessID: businessID,
		StartTime:  start.UTC(),
		EndTime:    end.UTC(),
		Reason:     req.Reason,
	}

	return s.repo.CreateBlockedTime(blockedTime)
}

func (s *Service) ListBlockedTime(businessID uuid.UUID) ([]models.BlockedTime, error) {
	_, err := s.businessWriter.GetBusinessByID(businessID)
	if err != nil {
		return nil, err
	}
	return s.repo.ListBlockedTime(businessID)
}

func (s *Service) GetBlockedTimeByID(id uuid.UUID) (*models.BlockedTime, error) {
	return s.repo.GetBlockedTimeByID(id)
}

func (s *Service) UpdateBlockedTime(
	id uuid.UUID,
	businessID uuid.UUID,
	req *UpdateBlockedTimeRequest,
) error {

	blockedTime, err := s.repo.GetBlockedTimeByID(id)
	if err != nil {
		return err
	}

	business, err := s.businessWriter.GetBusinessByID(businessID)
	if err != nil {
		return err
	}

	loc, err := time.LoadLocation(business.Timezone)
	if err != nil {
		return err
	}

	if !req.StartTime.IsZero() {
		start := req.StartTime.In(loc)
		blockedTime.StartTime = start.UTC()
	}

	if !req.EndTime.IsZero() {
		end := req.EndTime.In(loc)
		blockedTime.EndTime = end.UTC()
	}

	if !blockedTime.StartTime.Before(blockedTime.EndTime) {
		return errors.New("start time must be before end time")
	}

	if req.Reason != "" {
		blockedTime.Reason = req.Reason
	}

	return s.repo.UpdateBlockedTime(blockedTime)
}

func (s *Service) DeleteBlockedTime(id uuid.UUID) error {

	return s.repo.DeleteBlockedTime(id)
}
