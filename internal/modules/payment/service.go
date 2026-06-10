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
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
	"github.com/stylebookgo/agbadevemma/internal/modules/ledger"
	"github.com/stylebookgo/agbadevemma/internal/utils"
	"gorm.io/gorm"
)

type Service struct {
	repo         *Repository
	provider     *Provider
	business     business.BusinessWriter
	ledgerWriter ledger.LedgerWriter
}

func NewService(repo *Repository, provider *Provider, ledgerWriter ledger.LedgerWriter, business business.BusinessWriter) *Service {
	return &Service{repo: repo, provider: provider, ledgerWriter: ledgerWriter, business: business}
}

func (s *Service) InitializePayment(tx *gorm.DB, booking *models.Booking) (string, error) {

	// 2. Prevent duplicate payment
	existingPayment, err := s.repo.GetPaymentByBookingID(booking.ID)

	if err == nil && existingPayment.Status == utils.PaymentSuccess {
		return "", errors.New("payment already completed for this booking")
	}

	// 3. Calculate total amount from booking services
	amount, err := s.repo.CalculateBookingAmount(tx, booking.ID)
	if err != nil {
		tx.Rollback()
		return "", errors.New("failed to calculate booking amount")
	}

	platformFee, netAmount, totalAmount := funct.CalculateFee(amount)
	// log.Printf("platformFee: %v, netAmount: %v, totalAmount: %v", platformFee, netAmount, totalAmount)

	ref := uuid.NewString()

	payment := &models.Payment{
		BookingID:     booking.ID,
		BusinessID:    booking.BusinessID,
		Reference:     ref,
		Name:          booking.CustomerName,
		Email:         booking.CustomerEmail,
		Amount:        totalAmount,
		ServiceAmount: netAmount,
		PlatformFee:   platformFee,
		NetAmount:     netAmount,
		Status:        utils.PaymentPending,
		Provider:      "paystack",
	}

	err = s.repo.CreatePayment(tx, payment)
	if err != nil {

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

func (s *Service) HandleSuccessfulPayment(reference string) error {

	tx := s.repo.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	payment, err := s.repo.GetPaymentByReference(reference)
	if err != nil {
		return errors.New("payment not found")
	}

	if payment.IsVerified {
		return nil
	}

	if payment.Status == utils.PaymentSuccess {
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

	if err = s.repo.UpdateBookingStatus(
		tx,
		payment.BookingID,
		utils.BookingConfirmed,
		utils.PaymentSuccess,
	); err != nil {
		tx.Rollback()
		return err
	}

	ledger := &models.Ledger{
		ID:         uuid.New(),
		BusinessID: payment.BusinessID,
		Amount:     payment.ServiceAmount,
		Type:       utils.LedgerCredit,
		Reference:  "payment",
		Reason:     "payment for service",
		CreatedAt:  time.Now(),
	}

	if err := s.ledgerWriter.CreateLedger(tx, ledger); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *Service) GetBusinessTransactions(
	businessID uuid.UUID,
	limit int,
	offset int,
) ([]models.Payment, int64, error) {

	return s.repo.GetPaymentsByBusiness(businessID, limit, offset)
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
