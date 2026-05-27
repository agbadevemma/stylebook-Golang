package client

import (
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type ClientWriter interface {
	CreateClient(tx *gorm.DB, client *models.Client) error
}
