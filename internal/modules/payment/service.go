package payment

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/funct"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/utils"
	"gorm.io/gorm"
)

type Service struct {
	repo     *Repository
	provider *Provider
}

func NewService(repo *Repository, provider *Provider) *Service {
	return &Service{repo: repo, provider: provider}
}

func (s *Service) InitializePayment(tx *gorm.DB, booking *models.Booking) (string, error) {

	// 2. Prevent duplicate payment
	existingPayment, err := s.repo.GetPaymentByBookingID(booking.ID)

	if err == nil && existingPayment.Status == utils.PaymentSuccess {
		return "", errors.New("payment already completed for this booking")
	}

	// 3. Calculate total amount from booking services
	totalAmount, err := s.repo.CalculateBookingAmount(booking.ID)
	if err != nil {
		return "", errors.New("failed to calculate booking amount")
	}

	// 4. Apply platform fee
	fee, net := funct.CalculateFee(totalAmount)

	// 5. Generate reference (idempotency key)
	ref := uuid.NewString()

	// 6. Create payment record (pending)
	payment := &models.Payment{
		BookingID:     booking.ID,
		BusinessID:    booking.BusinessID,
		Reference:     ref,
		Email:         booking.CustomerEmail, // replace with email if you add it
		Amount:        totalAmount,
		ServiceAmount: totalAmount,
		PlatformFee:   fee,
		NetAmount:     net,
		Status:        utils.PaymentPending,
		Provider:      "paystack",
	}

	err = s.repo.CreatePayment(tx, payment)
	if err != nil {
		tx.Rollback()
		return "", errors.New("failed to create payment record")
	}

	// 7. Initialize Paystack payment
	res, err := s.provider.InitializePayment(
		payment.Email,
		payment.Amount,
		payment.Reference,
	)

	if err != nil {
		tx.Rollback()
		return "", errors.New("failed to initialize paystack payment")
	}

	// 8. Update payment with Paystack response
	payment.AuthorizationURL = res.Data.AuthorizationURL
	payment.AccessCode = res.Data.AccessCode

	if err := s.repo.UpdatePayment(tx, payment); err != nil {
		tx.Rollback()
		return "", errors.New("failed to update payment initialization")
	}

	if err := tx.Commit().Error; err != nil {
		return "", errors.New("transaction commit failed")
	}
	// 9. Return redirect URL to frontend
	return res.Data.AuthorizationURL, nil
}

func (s *Service) VerifySignature(payload []byte, signature string) bool {

	mac := hmac.New(sha512.New, []byte(os.Getenv("PAYSTACK_SECRET_KEY")))
	mac.Write(payload)

	expected := hex.EncodeToString(mac.Sum(nil))

	return expected == signature
}

func (s *Service) HandleSuccessfulPayment(reference string) error {

	tx := s.repo.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	// 1. Find payment by reference
	payment, err := s.repo.GetPaymentByReference(reference)
	if err != nil {
		return errors.New("payment not found")
	}

	// 2. Idempotency check (VERY IMPORTANT)
	if payment.Status == utils.PaymentSuccess {
		return nil // already processed
	}

	// 3. Update payment
	now := time.Now()

	payment.Status = utils.PaymentSuccess
	payment.IsVerified = true
	payment.PaidAt = &now

	if err := s.repo.UpdatePayment(tx, payment); err != nil {
		tx.Rollback()
		return err
	}

	// 4. Update booking
	return s.repo.UpdateBookingStatus(
		tx,
		payment.BookingID,
		utils.BookingConfirmed,
	)
}
