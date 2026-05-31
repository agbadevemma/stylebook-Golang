package blockedtime

import (
	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateBlockedTime(blockedTime *models.BlockedTime) error {
	return r.db.Create(blockedTime).Error
}

func (r *Repository) ListBlockedTime(businessID uuid.UUID) ([]models.BlockedTime, error) {
	var blockedTimes []models.BlockedTime

	err := r.db.
		Where("businessId = ?", businessID).
		Find(&blockedTimes).Error

	return blockedTimes, err
}

func (r *Repository) GetBlockedTimeByID(id uuid.UUID) (*models.BlockedTime, error) {
	var blockedTime models.BlockedTime

	err := r.db.
		First(&blockedTime, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return &blockedTime, nil
}

func (r *Repository) UpdateBlockedTime(blockedTime *models.BlockedTime) error {
	return r.db.Save(blockedTime).Error
}

func (r *Repository) DeleteBlockedTime(id uuid.UUID) error {
	result := r.db.Delete(&models.BlockedTime{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
