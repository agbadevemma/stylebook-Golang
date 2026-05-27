package client

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

func (r *Repository) CreateClient(tx *gorm.DB, client *models.Client) error {
	return tx.Create(client).Error
}

func (r *Repository) GetClientbyEmail(email string) (*models.Client, error) {
	var client models.Client
	if err := r.db.Where("email = ?", email).First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *Repository) GetClientByPhone(phone string) (*models.Client, error) {
	var client models.Client
	if err := r.db.Where("phone = ?", phone).First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}
