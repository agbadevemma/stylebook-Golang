package withdrawal

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
	"github.com/stylebookgo/agbadevemma/internal/modules/ledger"
	"github.com/stylebookgo/agbadevemma/internal/utils"
)

type Service struct {
	repo           *Repository
	provider       *Provider
	ledgerWriter   ledger.LedgerWriter
	businessWriter business.BusinessWriter
}

func NewService(repo *Repository, provider *Provider, ledgerWriter ledger.LedgerWriter, businesWriter business.BusinessWriter) *Service {
	return &Service{repo: repo, provider: provider, ledgerWriter: ledgerWriter, businessWriter: businesWriter}
}

func (s *Service) VerifySignature(payload []byte, signature string) bool {

	mac := hmac.New(sha512.New, []byte(os.Getenv("PAYSTACK_SECRET_KEY")))
	mac.Write(payload)

	expected := hex.EncodeToString(mac.Sum(nil))

	return expected == signature
}

func (s *Service) GetBusinessWithdrawals(businessID uuid.UUID) ([]models.Withdrawal, error) {
	return s.repo.GetBusinessWithdrawals(businessID)
}

func (s *Service) CreateWithdrawal(businessID uuid.UUID, amount int64) (*TransferResponse, error) {
	tx := s.repo.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	balance, err := s.ledgerWriter.GetBusinessBalance(businessID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if amount <= 0 {
		tx.Rollback()
		return nil, errors.New("amount must be greater than zero")
	}

	if amount > balance {
		tx.Rollback()
		return nil, errors.New("insufficient balance")
	}

	business, err := s.businessWriter.GetBusinessByID(businessID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if business.TransferRecipient == "" {
		tx.Rollback()
		return nil, errors.New("bank account not set")
	}

	transferRef := fmt.Sprintf("wdr_%s", uuid.NewString())

	withdrawal := &models.Withdrawal{
		BusinessID:        businessID,
		Amount:            amount,
		Status:            utils.WithdrawalPending,
		RecipientCode:     business.TransferRecipient,
		TransferReference: transferRef,
	}

	if err := s.repo.CreateWithdrawal(tx, withdrawal); err != nil {
		tx.Rollback()
		return nil, err
	}

	withdrawal.Status = utils.WithdrawalProcessing

	if err := s.repo.UpdateWithdrawal(tx, withdrawal); err != nil {
		tx.Rollback()
		return nil, err
	}

	res, err := s.provider.CreateTransfer(TransferRequest{
		Source:    "balance",
		Amount:    amount,
		Reference: transferRef,
		Recipient: business.TransferRecipient,
		Reason:    "Business withdrawal",
	})

	if err != nil {
		withdrawal.Status = utils.WithdrawalFailed

		if saveErr := s.repo.UpdateWithdrawal(tx, withdrawal); saveErr != nil {
			tx.Rollback()
			return nil, saveErr
		}

		if commitErr := tx.Commit().Error; commitErr != nil {
			return nil, commitErr
		}

		return nil, err
	}

	if res.Data.TransferCode != "" {
	
		existingWithdrawal, err := s.repo.GetWithdrawalByPaystackTransferCode(res.Data.TransferCode)
		if err == nil && existingWithdrawal != nil && existingWithdrawal.ID != withdrawal.ID {
		
			tx.Rollback()
			return nil, fmt.Errorf("transfer code already exists: %s", res.Data.TransferCode)
		}


		withdrawal.PaystackTransferCode = res.Data.TransferCode
	}

	if err := s.repo.UpdateWithdrawal(tx, withdrawal); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) HandleTransferEvent(event utils.PaystackWebhookEvent) error {
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

	withdrawal, err := s.repo.GetWithdrawalByReference(event.Data.Reference)
	if err != nil {
		tx.Rollback()
		return err
	}
	previousStatus := withdrawal.Status
	idemKey := withdrawal.TransferReference + ":" + event.Event

	switch event.Event {

	case "transfer.success":
		if withdrawal.Status == utils.WithdrawalSuccess {
			tx.Rollback()
			return nil
		}

		exists, err := s.ledgerWriter.LedgerExistsByReference(
			tx,
			idemKey,
			utils.LedgerDebit,
		)
		if err != nil {
			tx.Rollback()
			return err
		}

		if !exists {

			ledger := &models.Ledger{
				ID:         uuid.New(),
				BusinessID: withdrawal.BusinessID,
				Amount:     withdrawal.Amount,
				Type:       utils.LedgerDebit,
				Reference:  idemKey,
				Reason:     "Withdrawal successful",
				CreatedAt:  time.Now(),
			}

			if err := s.ledgerWriter.CreateLedger(tx, ledger); err != nil {
				tx.Rollback()
				return err
			}

		}
		withdrawal.Status = utils.WithdrawalSuccess

	case "transfer.failed":
		if withdrawal.Status == utils.WithdrawalSuccess ||
			withdrawal.Status == utils.WithdrawalReversed {
			tx.Rollback()
			return nil
		}

		withdrawal.Status = utils.WithdrawalFailed

	case "transfer.reversed":

		if withdrawal.Status == utils.WithdrawalReversed {
			tx.Rollback()
			return nil
		}

		exists, err := s.ledgerWriter.LedgerExistsByReference(
			tx,
			idemKey,
			utils.LedgerCredit,
		)

		if err != nil {
			tx.Rollback()
			return err
		}

		if previousStatus == utils.WithdrawalSuccess && !exists {
			ledger := &models.Ledger{
				ID:         uuid.New(),
				BusinessID: withdrawal.BusinessID,
				Amount:     withdrawal.Amount,
				Type:       utils.LedgerCredit,
				Reference:  idemKey,
				Reason:     "Withdrawal reversed",
				CreatedAt:  time.Now(),
			}

			if err := s.ledgerWriter.CreateLedger(tx, ledger); err != nil {
				tx.Rollback()
				return err
			}
		}

		withdrawal.Status = utils.WithdrawalReversed

	default:
		tx.Rollback()
		return fmt.Errorf("unsupported transfer event: %v", event)
	}

	if err := s.repo.UpdateWithdrawal(tx, withdrawal); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
