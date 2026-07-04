package otp

import (
	"errors"
	"time"

	"github.com/stylebookgo/agbadevemma/internal/funct"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/modules/notification/email"
	"github.com/stylebookgo/agbadevemma/internal/utils"
	"gorm.io/gorm"
)

type Service struct {
	repo         *Repository
	emailService email.EmailWriter
}

func NewService(repo *Repository, email email.EmailWriter) *Service {
	return &Service{repo: repo, emailService: email}
}

func (s *Service) SendOTP(
	tx *gorm.DB,
	email string,
	otpType utils.OTPType,
) error {

	code, err := funct.GenerateOTP()
	if err != nil {
		return err
	}

	otp := &models.OTP{

		Email:     email,
		Type:      otpType,
		CodeHash:  funct.HashOTP(code),
		ExpiresAt: time.Now().UTC().Add(10 * time.Minute),
		Used:      false,
	}

	if err := s.repo.Create(tx, otp); err != nil {
		return err
	}

	return s.emailService.SendOTP(email, code)
}

func (s *Service) VerifyOTP(
	tx *gorm.DB,
	email, code string,
	otpType utils.OTPType,
) error {

	otp, err := s.repo.FindValidOTP(email, otpType)
	if err != nil {
		return errors.New("invalid or expired otp")
	}

	if otp.Used {
		return errors.New("otp already used")
	}

	if otp.Attempts >= otp.MaxAttempts {
		return errors.New("too many attempts")
	}

	// increment attempts
	otp.Attempts++

	if funct.HashOTP(code) != otp.CodeHash {
		_ = s.repo.Update(tx, otp)
		return errors.New("invalid otp")
	}

	now := time.Now().UTC()
	otp.Used = true
	otp.VerifiedAt = &now

	return s.repo.Update(tx, otp)
}
