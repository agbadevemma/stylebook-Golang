package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type ServiceWriter interface {
	GetServiceByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*models.Service, error)
}
