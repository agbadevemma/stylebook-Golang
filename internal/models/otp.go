package models

import (
	"time"

	"github.com/stylebookgo/agbadevemma/internal/utils"
)

type OTP struct {
	Base
	Email       string        `gorm:"type:varchar(255);not null;index:idx_otp_email_type"`
	CodeHash    string        `gorm:"type:varchar(255);not null"`
	Type        utils.OTPType `gorm:"type:varchar(50);not null;index:idx_otp_email_type"`
	Attempts    uint8         `gorm:"default:0"`
	MaxAttempts uint8         `gorm:"default:5"`
	Used        bool          `gorm:"default:false"`
	VerifiedAt  *time.Time    `gorm:"default:null"`
	ExpiresAt   time.Time     `gorm:"not null;index"`
}
