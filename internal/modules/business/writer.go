package business

import (
	"context"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type BusinessWriter interface {
	GetBusinessByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*models.Business, error)
	GetBusinessByUserID(ctx context.Context, id uuid.UUID) (*models.Business, error)
	CreateBusiness(ctx context.Context, tx *gorm.DB, business *models.Business) error
	UpdateBusinessDetails(
		ctx context.Context,
		businessID uuid.UUID,
		name string,
		timezone string,
	) error
	UpdateBusinessBankDetails(
		ctx context.Context,
		businessID uuid.UUID,
		accountName,
		accountNumber,
		subaccountCode,
		bankName,
		bankCode string,
	) error
}
