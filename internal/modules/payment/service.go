package payment

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/funct"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/utils"
)

type Service struct {
	repo     *Repository
	provider *Provider

}

func NewService(repo *Repository, provider *Provider) *Service {
	return &Service{repo: repo, provider: provider}
}

func (s *Service) InitializePayment(bookingID uuid.UUID) (string, error) {

	// 1. Fetch booking
	booking, err := s.repo.GetBookingByID(bookingID)
	if err != nil {
		return "", errors.New("booking not found")
	}

	// 2. Prevent duplicate payment
	existingPayment, err := s.repo.GetByBookingID(bookingID)
	if err == nil && existingPayment.Status == utils.PaymentSuccess {
		return "", errors.New("payment already completed for this booking")
	}

	// 3. Calculate total amount from booking services
	totalAmount, err := s.repo.CalculateBookingAmount(bookingID)
	if err != nil {
		return "", errors.New("failed to calculate booking amount")
	}

	// 4. Apply platform fee
	fee, net := funct.CalculateFee(totalAmount)

	// 5. Generate reference (idempotency key)
	ref := "BOOK_" + booking.ID.String()

	// 6. Create payment record (pending)
	payment := &models.Payment{
		BookingID:     booking.ID,
		BusinessID:    booking.BusinessID,
		Reference:     ref,
		Email:         booking.CustomerPhone, // replace with email if you add it
		Amount:        totalAmount,
		ServiceAmount: totalAmount,
		PlatformFee:   fee,
		NetAmount:     net,
		Status:        utils.PaymentPending,
		Provider:      "paystack",
	}

	err = s.repo.CreatePayment(nil, payment)
	if err != nil {
		return "", errors.New("failed to create payment record")
	}

	// 7. Initialize Paystack payment
	res, err := s.provider.InitializePayment(
		payment.Email,
		payment.Amount,
		payment.Reference,
	)

	if err != nil {
		return "", errors.New("failed to initialize paystack payment")
	}

	// 8. Update payment with Paystack response
	payment.AuthorizationURL = res.Data.AuthorizationURL
	payment.AccessCode = res.Data.AccessCode

	if err := s.repo.UpdatePayment(nil, payment); err != nil {
		return "", errors.New("failed to update payment initialization")
	}

	// 9. Return redirect URL to frontend
	return res.Data.AuthorizationURL, nil
}
