package business

import (
	"context"

	"github.com/stylebookgo/agbadevemma/internal/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetBusiness(ctx context.Context) (models.Business, error) {
	business, err := s.repo.GetBusiness(ctx)
	return business, err
}
