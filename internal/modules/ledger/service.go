package ledger

import "github.com/google/uuid"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetBusinessBalance(businessID uuid.UUID) (int64, error) {
	return s.repo.GetBusinessBalance(businessID)
}
