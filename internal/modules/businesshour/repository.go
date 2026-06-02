package businesshour

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

func (r *Repository) CreateBusinessHour(hour *models.BusinessHour) error {
	return r.db.Create(hour).Error
}

func (r *Repository) CreateBusinessHourWithTx(tx *gorm.DB, hour *models.BusinessHour) error {
	return tx.Create(hour).Error
}

func (r *Repository) UpdateBusinessHour(tx *gorm.DB, hour *models.BusinessHour) error {
	return tx.Save(hour).Error
}

func (r *Repository) GetByBusinessAndDay(
	tx *gorm.DB,
	businessID uuid.UUID,
	day int,
) (*models.BusinessHour, error) {

	var hour models.BusinessHour

	err := tx.
		Where("businessId = ? AND dayOfWeek = ?", businessID, day).
		First(&hour).Error

	if err != nil {
		return nil, err
	}

	return &hour, nil
}

func (r *Repository) ListByBusiness(
	businessID uuid.UUID,
) ([]models.BusinessHour, error) {

	var hours []models.BusinessHour

	err := r.db.
		Where("businessId = ?", businessID).
		Order("dayOfWeek ASC").
		Find(&hours).Error

	return hours, err
}
