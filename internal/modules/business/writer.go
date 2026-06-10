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
	UpdateBusinessBankDetails(
		businessID uuid.UUID,
		accountNumber string,
		transferCode string,
		bankName string,
		bankCode string,
	) error
}
