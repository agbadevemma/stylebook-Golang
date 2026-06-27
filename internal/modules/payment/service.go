package payment

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/funct"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
	"github.com/stylebookgo/agbadevemma/internal/modules/ledger"
	"github.com/stylebookgo/agbadevemma/internal/modules/withdrawal"
	"github.com/stylebookgo/agbadevemma/internal/utils"
	"gorm.io/gorm"
)

type Service struct {
	repo             *Repository
	provider         *Provider
	business         business.BusinessWriter
	ledgerWriter     ledger.LedgerWriter
	withdrawalWriter withdrawal.WithdrawalWriter
}

func NewService(repo *Repository, provider *Provider, ledgerWriter ledger.LedgerWriter, business business.BusinessWriter, withdrawalWriter withdrawal.WithdrawalWriter) *Service {
	return &Service{repo: repo, provider: provider, ledgerWriter: ledgerWriter, business: business, withdrawalWriter: withdrawalWriter}
}

func (s *Service) InitializePayment(tx *gorm.DB, booking *models.Booking) (string, error) {

	existingPayment, err := s.repo.GetPaymentByBookingID(tx, booking.ID)
	if err == nil && existingPayment != nil {
		if existingPayment.Status == utils.PaymentSuccess {
			return "", errors.New("payment already completed for this booking")
		}

		if existingPayment.AuthorizationURL != "" {
			return existingPayment.AuthorizationURL, nil
		}
	}

	serviceAmount, err := s.repo.CalculateBookingAmount(tx, booking.ID)
	if err != nil {
		tx.Rollback()
		return "", errors.New("failed to calculate booking amount")
	}

	platformFee, netAmount, totalAmount := funct.CalculateFee(serviceAmount)
	// log.Printf("platformFee: %v, netAmount: %v, totalAmount: %v", platformFee, netAmount, totalAmount)

	ref := uuid.NewString()

	payment := &models.Payment{
		BookingID:     booking.ID,
		BusinessID:    booking.BusinessID,
		Reference:     ref,
		Name:          booking.CustomerName,
		Email:         booking.CustomerEmail,
		Amount:        totalAmount,
		ServiceAmount: serviceAmount,
		PlatformFee:   platformFee,
		NetAmount:     netAmount,
		Status:        utils.PaymentPending,
		Provider:      "paystack",
	}

	if err := s.repo.CreatePayment(tx, payment); err != nil {
		tx.Rollback()
		return "", errors.New("failed to create payment record")
	}

	res, err := s.provider.InitializePayment(
		payment.Email,
		payment.Amount,
		payment.Reference,
	)
	if err != nil {
		tx.Rollback()
		return "", errors.New("failed to initialize paystack payment")
	}

	payment.AuthorizationURL = res.Data.AuthorizationURL
	payment.AccessCode = res.Data.AccessCode

	if err := s.repo.UpdatePayment(tx, payment); err != nil {
		tx.Rollback()
		return "", errors.New("failed to update payment initialization")
	}

	return res.Data.AuthorizationURL, nil
}

func (s *Service) VerifySignature(payload []byte, signature string) bool {

	mac := hmac.New(sha512.New, []byte(os.Getenv("PAYSTACK_SECRET_KEY")))
	mac.Write(payload)

	expected := hex.EncodeToString(mac.Sum(nil))

	return expected == signature
}

func (s *Service) HandleEvent(event utils.PaystackWebhookEvent) error {

	switch event.Event {

	// 💳 PAYMENT EVENTS
	case "charge.success", "charge.failed", "charge.abandoned", "charge.refunded":
		return s.handlePaymentEvent(event)

	// 💸 TRANSFER EVENTS
	case "transfer.success", "transfer.failed", "transfer.reversed":
		return s.withdrawalWriter.HandleTransferEvent(event)

	default:
		// ignore unknown events safely
		log.Printf("ignored event: %s", event.Event)
		return nil
	}
}

func (s *Service) handlePaymentEvent(event utils.PaystackWebhookEvent) error {
	tx := s.repo.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	payment, err := s.repo.GetPaymentByReference(tx, event.Data.Reference)
	if err != nil {
		tx.Rollback()
		return err
	}

	switch event.Event {

	case "charge.success":
		return s.handlePaymentSuccess(tx, payment)

	case "charge.failed":
		return s.handlePaymentFailed(tx, payment)

	case "charge.abandoned":
		return s.handlePaymentAbandoned(tx, payment)

	case "refund.processed":
		return s.handlePaymentRefund(tx, payment)

	default:
		tx.Rollback()
		return nil
	}
}

func (s *Service) handlePaymentSuccess(tx *gorm.DB, payment *models.Payment) error {

	if payment.Status == utils.PaymentSuccess {
		tx.Rollback()
		return nil
	}

	now := time.Now()

	payment.Status = utils.PaymentSuccess
	payment.IsVerified = true
	payment.PaidAt = &now

	if err := s.repo.UpdatePayment(tx, payment); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.UpdateBookingStatus(
		tx,
		payment.BookingID,
		utils.BookingConfirmed,
		utils.PaymentSuccess,
	); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *Service) handlePaymentFailed(tx *gorm.DB, payment *models.Payment) error {

	if payment.Status == utils.PaymentSuccess {
		tx.Rollback()
		return nil
	}

	payment.Status = utils.PaymentFailed

	if err := s.repo.UpdatePayment(tx, payment); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
func (s *Service) handlePaymentAbandoned(tx *gorm.DB, payment *models.Payment) error {

	if payment.Status == utils.PaymentSuccess {
		tx.Rollback()
		return nil
	}

	payment.Status = utils.PaymentAbandoned

	if err := s.repo.UpdatePayment(tx, payment); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *Service) handlePaymentRefund(tx *gorm.DB, payment *models.Payment) error {

	if payment.Status == utils.PaymentRefunded {
		tx.Rollback()
		return nil
	}

	payment.Status = utils.PaymentRefunded

	if err := s.repo.UpdatePayment(tx, payment); err != nil {
		tx.Rollback()
		return err
	}

	// optional: reverse ledger entry
	return tx.Commit().Error
}

func (s *Service) GetBusinessTransactions(
	businessID uuid.UUID,
	params ListTransactionsParams,
) ([]models.Payment, int64, error) {

	return s.repo.GetPaymentsByBusiness(businessID, params)
}

func (s *Service) ResolveBankAccount(accountNumber, bankCode string) (*BankResolveResponse, error) {
	return s.provider.ResolveBankAccount(accountNumber, bankCode)
}

func (s *Service) CreateTransferRecipient(
	name,
	accountNumber,
	bankName,
	bankCode string,
	businessID uuid.UUID,
) (*CreateRecipientResponse, error) {

	_, err := s.business.GetBusinessByID(businessID)

	if err != nil {
		return nil, err
	}

	res, err := s.provider.CreateTransferRecipient(
		name,
		accountNumber,
		bankCode,
	)
	if err != nil {
		return nil, err
	}

	err = s.business.UpdateBusinessBankDetails(businessID, accountNumber, res.Data.RecipientCode, bankName, bankCode)
	if err != nil {
		return nil, err
	}

	return res, nil
}
