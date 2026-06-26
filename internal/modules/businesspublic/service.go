package businesspublic

import (
	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
)

type Service struct {
	businessWriter business.BusinessWriter
}

func NewService(businessWriter business.BusinessWriter) *Service {
	return &Service{businessWriter: businessWriter}
}

func (s *Service) GetBusinessByID(id uuid.UUID) (*models.Business, error) {
	return s.businessWriter.GetBusinessByID(id)
}
