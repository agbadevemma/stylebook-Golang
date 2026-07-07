package otp

import (
	"github.com/stylebookgo/agbadevemma/internal/utils"
	"gorm.io/gorm"
)

type OtpWriter interface {
	SendOTP(tx *gorm.DB, email string, otpType utils.OTPType) error
	VerifyOTP(tx *gorm.DB,email, code string, otpType utils.OTPType) error
}
