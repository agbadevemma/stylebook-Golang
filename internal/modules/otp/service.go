package otp

import (
	"context"
	"errors"
	"log"
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
	ctx context.Context,
	tx *gorm.DB,
	email string,
	otpType utils.OTPType,
) error {

	code, err := funct.GenerateOTP()
	if err != nil {
		return err
	}

	// Check if an existing OTP record exists for this email and type
	existingOtp, err := s.repo.FindLatestOTPByEmailAndType(ctx, tx, email, otpType)

	// If existing record found, update it
	if err == nil && existingOtp != nil {
		existingOtp.CodeHash = funct.HashOTP(code)
		existingOtp.ExpiresAt = time.Now().UTC().Add(10 * time.Minute)
		existingOtp.Attempts = 0
		existingOtp.Used = false
		existingOtp.VerifiedAt = nil

		if err := s.repo.Update(ctx, tx, existingOtp); err != nil {
			return err
		}

		emailCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		go func() {
			defer cancel()
			if err := s.emailService.SendOTP(emailCtx, email, code); err != nil {
				log.Printf("failed to send OTP email: %v", err)
			}
		}()
		return nil
	}

	// Otherwise create a new OTP record
	otp := &models.OTP{
		Email:     email,
		Type:      otpType,
		CodeHash:  funct.HashOTP(code),
		ExpiresAt: time.Now().UTC().Add(10 * time.Minute),
		Used:      false,
	}

	if err := s.repo.Create(ctx, tx, otp); err != nil {
		return err
	}

	emailCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	go func() {
		defer cancel()
		if err := s.emailService.SendOTP(emailCtx, email, code); err != nil {
			// log error instead of ignoring it
			log.Printf("failed to send OTP email: %v", err)
		}
	}()
	return nil
}

func (s *Service) VerifyOTP(
	ctx context.Context,
	tx *gorm.DB,
	email, code string,
	otpType utils.OTPType,
) error {

	otp, err := s.repo.FindValidOTP(ctx, tx, email, otpType)
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
		_ = s.repo.Update(ctx, tx, otp)
		return errors.New("invalid otp")
	}

	now := time.Now().UTC()
	otp.Used = true
	otp.VerifiedAt = &now

	return s.repo.Update(ctx, tx, otp)
}

func (s *Service) DeleteOTPRecords(ctx context.Context, tx *gorm.DB, email string, otpType utils.OTPType) error {
	return s.repo.DeleteByEmailAndType(ctx, tx, email, otpType)
}
