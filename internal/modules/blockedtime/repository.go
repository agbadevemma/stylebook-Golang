package blockedtime

import (
	"context"
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

func (r *Repository) CreateBlockedTime(ctx context.Context, blockedTime *models.BlockedTime) error {
	return r.db.WithContext(ctx).Create(blockedTime).Error
}

func (r *Repository) ListBlockedTime(ctx context.Context, businessID uuid.UUID) ([]models.BlockedTime, error) {
	var blockedTimes []models.BlockedTime

	err := r.db.WithContext(ctx).
		Where("business_id = ?", businessID).
		Find(&blockedTimes).Error

	return blockedTimes, err
}

func (r *Repository) GetBlockedTimeByID(ctx context.Context, id uuid.UUID) (*models.BlockedTime, error) {
	var blockedTime models.BlockedTime

	err := r.db.WithContext(ctx).
		First(&blockedTime, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return &blockedTime, nil
}

func (r *Repository) UpdateBlockedTime(ctx context.Context, blockedTime *models.BlockedTime) error {
	return r.db.WithContext(ctx).Save(blockedTime).Error
}

func (r *Repository) DeleteBlockedTime(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.BlockedTime{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *Repository) GetBlockedTimes(ctx context.Context,
	tx *gorm.DB,
	businessID uuid.UUID,
	start, end time.Time,
) ([]models.BlockedTime, error) {


	var blocks []models.BlockedTime

	if tx == nil {
		tx = r.db
	}
	err := tx.WithContext(ctx).
		Where("business_id = ? AND start_time < ? AND end_time > ?",
			businessID, end, start).
		Find(&blocks).Error

	return blocks, err
}

func (r *Repository) HasOverlappingBlockedTime(ctx context.Context,
	businessID uuid.UUID,
	start, end time.Time,
) (bool, error) {

	var count int64

	err := r.db.WithContext(ctx).Model(&models.BlockedTime{}).
		Where("business_id = ?", businessID).
		Where("start_time < ? AND end_time > ?", end, start).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
