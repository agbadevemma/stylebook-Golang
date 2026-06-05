package service

import (
	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
)

type ServiceWriter interface {
	GetServiceByID(id uuid.UUID) (*models.Service, error)
}
