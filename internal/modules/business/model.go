package business

import (
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type BusinessWriter interface {
	CreateBusiness(tx *gorm.DB, business *models.Business) error
}
