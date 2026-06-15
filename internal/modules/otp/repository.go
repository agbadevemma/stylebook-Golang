package otp

import (
	"time"

	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/utils"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(tx *gorm.DB, otp *models.OTP) error {
	db := tx
	if db == nil {
		db = r.db
	}
	return db.Create(otp).Error
}

func (r *Repository) FindValidOTP(
	email string,
	otpType utils.OTPType,
) (*models.OTP, error) {

	var otp models.OTP

	err := r.db.
		Where("email = ? AND type = ? AND used = ? AND expires_at > ?",
			email, otpType, false, time.Now()).
		Order("created_at DESC").
		First(&otp).Error

	if err != nil {
		return nil, err
	}

	return &otp, nil
}

func (r *Repository) Update(tx *gorm.DB, otp *models.OTP) error {
	db := tx
	if db == nil {
		db = r.db
	}
	return db.Save(otp).Error
}
