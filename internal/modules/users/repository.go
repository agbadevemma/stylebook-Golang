package users

import (
	"context"

	"github.com/stylebookgo/agbadevemma/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, tx *gorm.DB, user *models.User) error {
	return tx.WithContext(ctx).Create(user).Error
}

func (r *Repository) GetUserByPhone(ctx context.Context, tx *gorm.DB, phone string) (*models.User, error) {
	var user models.User
	err := tx.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserbyEmail(ctx context.Context, tx *gorm.DB, email string) (*models.User, error) {
	var user models.User
	err := tx.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, tx *gorm.DB, user *models.User) error {

	return tx.WithContext(ctx).Save(user).Error
}
