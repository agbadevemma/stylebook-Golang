package blockedtime

import (
	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
)

type Service struct {
	repo           *Repository
	businessWriter business.BusinessWriter
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateBlockedTime(
	req *CreateBlockedTimeRequest,
	businessID uuid.UUID,
) error {

	blockedTime := &models.BlockedTime{
		BusinessID: businessID,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
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
	req *UpdateBlockedTimeRequest,
) error {

	blockedTime, err := s.repo.GetBlockedTimeByID(id)
	if err != nil {
		return err
	}

	if !req.StartTime.IsZero() {
		blockedTime.StartTime = req.StartTime
	}

	if !req.EndTime.IsZero() {
		blockedTime.EndTime = req.EndTime
	}

	if req.Reason != "" {
		blockedTime.Reason = req.Reason
	}

	return s.repo.UpdateBlockedTime(blockedTime)
}

func (s *Service) DeleteBlockedTime(id uuid.UUID) error {

	return s.repo.DeleteBlockedTime(id)
}
