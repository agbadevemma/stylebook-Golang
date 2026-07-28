package business

import (
	"context"
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

func (s *Service) GetBusinessSettings(ctx context.Context, businessID uuid.UUID) (*models.Business, error) {
	return s.repo.GetBusinessByID(ctx, s.repo.db, businessID)
}

func (s *Service) UpdateBusinessDetails(
	ctx context.Context,
	businessID uuid.UUID,
	req UpdateBusinessDetailsRequest,
) (*models.Business, error) {
	name := strings.TrimSpace(req.Name)
	timezone := strings.TrimSpace(req.Timezone)

	if err := s.repo.UpdateBusinessDetails(ctx, businessID, name, timezone); err != nil {
		return nil, err
	}

	return s.repo.GetBusinessByID(ctx, s.repo.db, businessID)
}

func (s *Service) GetBusinessByID(ctx context.Context, id uuid.UUID) (*models.Business, error) {
	return s.repo.GetBusinessByID(ctx, s.repo.db, id)
}
