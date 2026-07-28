package otp

import (
	"context"
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

func (r *Repository) Create(
	ctx context.Context,
	tx *gorm.DB,
	otp *models.OTP,
) error {
	db := tx
	if db == nil {
		db = r.db
	}

	return db.WithContext(ctx).Create(otp).Error
}

func (r *Repository) FindValidOTP(
	ctx context.Context,
	tx *gorm.DB,
	email string,
	otpType utils.OTPType,
) (*models.OTP, error) {
	db := tx
	if db == nil {
		db = r.db
	}

	var otp models.OTP

	err := db.WithContext(ctx).
		Where(
			"email = ? AND type = ? AND used = ? AND expires_at > ?",
			email,
			otpType,
			false,
			time.Now().UTC(),
		).
		Order("created_at DESC").
		First(&otp).Error

	if err != nil {
		return nil, err
	}

	return &otp, nil
}

func (r *Repository) Update(
	ctx context.Context,
	tx *gorm.DB,
	otp *models.OTP,
) error {
	db := tx
	if db == nil {
		db = r.db
	}

	return db.WithContext(ctx).Save(otp).Error
}

func (r *Repository) DeleteByEmailAndType(
	ctx context.Context,
	tx *gorm.DB,
	email string,
	otpType utils.OTPType,
) error {
	db := tx
	if db == nil {
		db = r.db
	}

	return db.WithContext(ctx).
		Where("email = ? AND type = ?", email, otpType).
		Delete(&models.OTP{}).Error
}

func (r *Repository) FindLatestOTPByEmailAndType(
	ctx context.Context,
	tx *gorm.DB,
	email string,
	otpType utils.OTPType,
) (*models.OTP, error) {

	db := tx
	if db == nil {
		db = r.db
	}

	var otp models.OTP
	err := db.WithContext(ctx).
		Where("email = ? AND type = ?", email, otpType).
		Order("created_at DESC").
		First(&otp).Error

	if err != nil {
		return nil, err
	}

	return &otp, nil
}
