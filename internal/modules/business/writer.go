package business

import (
	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type BusinessWriter interface {
	GetBusinessByID(id uuid.UUID) (*models.Business, error)
	GetBusinessByUserID(id uuid.UUID) (*models.Business, error)
	CreateBusiness(tx *gorm.DB, business *models.Business) error
	UpdateBusinessDetails(
		businessID uuid.UUID,
		name string,
		timezone string,
	) error
	UpdateBusinessBankDetails(
		businessID uuid.UUID,
		accountName,
		accountNumber,
		subaccountCode,
		bankName,
		bankCode string,
	) error
}
