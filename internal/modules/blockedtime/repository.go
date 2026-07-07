package blockedtime

import (
	"time"

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
		Where("business_id = ?", businessID).
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

func (r *Repository) GetBlockedTimes(
	businessID uuid.UUID,
	start, end time.Time,
) ([]models.BlockedTime, error) {

	var blocks []models.BlockedTime

	err := r.db.
		Where("business_id = ? AND start_time < ? AND end_time > ?",
			businessID, end, start).
		Find(&blocks).Error

	return blocks, err
}

func (r *Repository) HasOverlappingBlockedTime(
	businessID uuid.UUID,
	start, end time.Time,
) (bool, error) {

	var count int64

	err := r.db.Model(&models.BlockedTime{}).
		Where("business_id = ?", businessID).
		Where("start_time < ? AND end_time > ?", end, start).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}