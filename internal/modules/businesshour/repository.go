package businesshour

import (
	"context"

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

func (r *Repository) CreateBusinessHour(ctx context.Context,tx *gorm.DB, hour *models.BusinessHour) error {
	if tx == nil {
		tx = r.db
	}
	return tx.WithContext(ctx).Create(hour).Error
}

func (r *Repository) CreateBusinessHourWithTx(ctx context.Context, tx *gorm.DB, hour *models.BusinessHour) error {
	return tx.WithContext(ctx).Create(hour).Error
}

func (r *Repository) UpdateBusinessHour(ctx context.Context, tx *gorm.DB, hour *models.BusinessHour) error {
	if tx == nil {
		tx = r.db
	}
	return tx.WithContext(ctx).Save(hour).Error
}

func (r *Repository) GetByBusinessAndDay(
	ctx context.Context,
	tx *gorm.DB,
	businessID uuid.UUID,
	day *int,
) (*models.BusinessHour, error) {

	var hour models.BusinessHour

	if tx == nil {
		tx = r.db
	}
	err := tx.WithContext(ctx).
		Where("business_id = ? AND day_of_week = ?", businessID, day).
		First(&hour).Error

	if err != nil {
		return nil, err
	}

	return &hour, nil
}

func (r *Repository) ListByBusiness(
	ctx context.Context,	
	businessID uuid.UUID,
) ([]models.BusinessHour, error) {

	var hours []models.BusinessHour

	err := r.db.WithContext(ctx).
		Where("business_id = ?", businessID).
		Order("day_of_week ASC").
		Find(&hours).Error

	return hours, err
}

func (r *Repository) GetBusinessHour(
	ctx context.Context,
	tx *gorm.DB,
	businessID uuid.UUID,
	dayOfWeek int,
) (*models.BusinessHour, error) {

	var hour models.BusinessHour

	if tx == nil {
		tx = r.db
	}
	err := tx.WithContext(ctx).
		Where("business_id = ? AND day_of_week = ?", businessID, dayOfWeek).
		First(&hour).Error

	if err != nil {

		return nil, err
	}

	return &hour, nil
}


func (r *Repository) GetBusinessHours(
	ctx context.Context,
	tx *gorm.DB,
	businessID uuid.UUID,
) ([]models.BusinessHour, error) {

	var hours []models.BusinessHour

	if tx == nil {
		tx = r.db
	}

	err := tx.WithContext(ctx).
		Where("business_id = ?", businessID).
		Find(&hours).Error

	return hours, err
}