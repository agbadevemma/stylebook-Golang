package business

import (
	"strings"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetBusinessSettings(businessID uuid.UUID) (*models.Business, error) {
	return s.repo.GetBusinessByID(businessID)
}

func (s *Service) UpdateBusinessDetails(
	businessID uuid.UUID,
	req UpdateBusinessDetailsRequest,
) (*models.Business, error) {
	name := strings.TrimSpace(req.Name)
	timezone := strings.TrimSpace(req.Timezone)

	if err := s.repo.UpdateBusinessDetails(businessID, name, timezone); err != nil {
		return nil, err
	}

	return s.repo.GetBusinessByID(businessID)
}



func (s *Service) GetBusinessByID(id uuid.UUID) (*models.Business, error) {
	return s.repo.GetBusinessByID(id)
}
