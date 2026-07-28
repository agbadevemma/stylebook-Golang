package otp

import (
	"context"

	"github.com/stylebookgo/agbadevemma/internal/utils"
	"gorm.io/gorm"
)

type OtpWriter interface {
	SendOTP(ctx context.Context, tx *gorm.DB, email string, otpType utils.OTPType) error
	VerifyOTP(ctx context.Context, tx *gorm.DB, email, code string, otpType utils.OTPType) error
	DeleteOTPRecords(ctx context.Context, tx *gorm.DB, email string, otpType utils.OTPType) error
}
