package ledger

import (
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateLedger(tx *gorm.DB, ledger *models.Ledger) error {
	return tx.Create(ledger).Error
}
