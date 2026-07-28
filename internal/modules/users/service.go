package users

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
	"github.com/stylebookgo/agbadevemma/internal/modules/businesshour"
	"github.com/stylebookgo/agbadevemma/internal/modules/notification/email"
	"github.com/stylebookgo/agbadevemma/internal/modules/otp"
	"github.com/stylebookgo/agbadevemma/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository

	businessWriter     business.BusinessWriter
	businessHourWriter businesshour.BusinessHourWriter
	otpWriter          otp.OtpWriter
	emailService       email.EmailWriter
}

func NewService(repo *Repository, bus business.BusinessWriter, busHour businesshour.BusinessHourWriter, otpWriter otp.OtpWriter, emailService email.EmailWriter) *Service {
	return &Service{repo: repo, businessWriter: bus, businessHourWriter: busHour, otpWriter: otpWriter, emailService: emailService}
}

func (s *Service) Signup(ctx context.Context, req *SignupRequest) error {
	tx := s.repo.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if _, err := s.repo.GetUserbyEmail(ctx, tx, req.Email); err == nil {
		tx.Rollback()
		return errors.New("user with this email already exists")
	}

	if _, err := s.repo.GetUserByPhone(ctx, tx, req.Phone); err == nil {
		tx.Rollback()
		return errors.New("user with this phone number already exists")
	}

	var businessID *uuid.UUID

	if req.Role == utils.Barber {
		if req.BusinessID == uuid.Nil {
			return errors.New("business ID is required for barbers")
		}

		if _, err := s.businessWriter.GetBusinessByID(ctx, s.repo.db, req.BusinessID); err != nil {
			return errors.New("business not found")
		}

		businessID = &req.BusinessID
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		tx.Rollback()
		return errors.New("failed to hash password")
	}

	user := &models.User{
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
		Role:       req.Role,
		Password:   string(hashedPassword),
		BusinessID: businessID,
	}

	if err := s.repo.CreateUser(ctx, tx, user); err != nil {
		tx.Rollback()
		return err
	}

	if req.Role == utils.Business {
		business := &models.Business{
			Name:    req.BusinessName,
			OwnerID: user.ID,
		}

		if err := s.businessWriter.CreateBusiness(ctx, tx, business); err != nil {
			tx.Rollback()
			return err
		}

		hours := []models.BusinessHour{
			{
				BusinessID: business.ID,
				DayOfWeek:  0,
				IsOpen:     false,
				StartTime:  "09:00",
				EndTime:    "18:00",
			},
			{
				BusinessID: business.ID,
				DayOfWeek:  1,
				IsOpen:     true,
				StartTime:  "09:00",
				EndTime:    "18:00",
			},
			{
				BusinessID: business.ID,
				DayOfWeek:  2,
				IsOpen:     true,
				StartTime:  "09:00",
				EndTime:    "18:00",
			},
			{
				BusinessID: business.ID,
				DayOfWeek:  3,
				IsOpen:     true,
				StartTime:  "09:00",
				EndTime:    "18:00",
			},
			{
				BusinessID: business.ID,
				DayOfWeek:  4,
				IsOpen:     true,
				StartTime:  "09:00",
				EndTime:    "18:00",
			},
			{
				BusinessID: business.ID,
				DayOfWeek:  5,
				IsOpen:     true,
				StartTime:  "09:00",
				EndTime:    "18:00",
			},
			{
				BusinessID: business.ID,
				DayOfWeek:  6,
				IsOpen:     true,
				StartTime:  "09:00",
				EndTime:    "18:00",
			},
		}
		for _, hour := range hours {
			if err := s.businessHourWriter.CreateBusinessHourWithTx(ctx, tx, &hour); err != nil {

				tx.Rollback()
				return err
			}
		}
	}
	if err := s.otpWriter.SendOTP(ctx, tx, req.Email, utils.OTPTypeEmailVerification); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (s *Service) VerifyOTP(ctx context.Context, req *VerifyOTPRequest) error {
	tx := s.repo.db.WithContext(ctx).Begin()

	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	err := s.otpWriter.VerifyOTP(ctx, tx, req.Email, req.Code, req.Type)
	if err != nil {
		tx.Rollback()
		return err
	}
	user, err := s.repo.GetUserbyEmail(ctx, tx, req.Email)
	if err != nil {
		tx.Rollback()
		return err
	}
	if user.IsEmailVerified {
		tx.Rollback()
		return errors.New("account already verified")
	}
	now := time.Now().UTC()
	user.IsEmailVerified = true
	user.EmailVerifiedAt = &now

	if err := s.repo.UpdateUser(ctx, tx, user); err != nil {
		tx.Rollback()
		return errors.New("failed to update password")
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete OTP records for this email after successful verification
	err = s.otpWriter.DeleteOTPRecords(ctx, tx, req.Email, req.Type)
	if err != nil {
		return errors.New("failed to delete OTP records")
	}

	// Send welcome email asynchronously with timeout
	emailCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	go func() {
		defer cancel()
		loginUrl := fmt.Sprintf("%s/login", os.Getenv("FRONTEND_URL"))
		accountVerifiedData := email.AccountVerifiedData{
			CustomerName:  user.Name,
			CustomerEmail: user.Email,
			LoginURL:      loginUrl,
		}
		if err := s.emailService.SendAccountVerified(emailCtx, user.Email, accountVerifiedData); err != nil {
			// Log error but don't fail the verification
			// In production, this should be logged to a proper logger
			log.Printf("failed to send verify otp email: %v", err)

		}
	}()

	return nil
}

func (s *Service) ChangePassword(ctx context.Context, req *ChangePasswordRequest) error {
	tx := s.repo.db.WithContext(ctx).Begin()

	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	user, err := s.repo.GetUserbyEmail(ctx, tx, req.Email)
	if err != nil {
		tx.Rollback()
		return errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err == nil {
		tx.Rollback()
		return errors.New("Create New Password")
	}
	if len(req.Password) < 8 {
		tx.Rollback()
		return errors.New("password too weak")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		tx.Rollback()
		return errors.New("error occured while hashing")
	}

	if err := s.otpWriter.VerifyOTP(ctx, tx, user.Email, req.Code, utils.OTPTypePasswordReset); err != nil {
		tx.Rollback()
		return errors.New(err.Error())
	}
	user.Password = string(hashedPassword)
	now := time.Now()
	user.PasswordChangedAt = &now
	if err := s.repo.UpdateUser(ctx, tx, user); err != nil {
		tx.Rollback()
		return errors.New("failed to update password")
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	err = s.otpWriter.DeleteOTPRecords(ctx, tx, user.Email, utils.OTPTypePasswordReset)
	if err != nil {
		return errors.New("failed to delete OTP records")
	}
	return nil
}

func (s *Service) Login(ctx context.Context, req *LoginRequest) (string, *models.Business, error) {
	user, err := s.repo.GetUserbyEmail(ctx, s.repo.db, req.Email)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}
	if !user.IsEmailVerified {
		return "", nil, errors.New("Email is not verified")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}
	business, err := s.businessWriter.GetBusinessByUserID(ctx, user.ID)
	if err != nil {
		return "", nil, errors.New("business not found")
	}
	token, err := utils.GenerateToken(user.ID.String(), user.Email, user.Role, business.ID)
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	return token, business, nil
}
