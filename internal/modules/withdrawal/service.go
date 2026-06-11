package withdrawal

import (
	"errors"
	"fmt"

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

func (s *Service) CreateWithdrawal(businessID uuid.UUID, amount int64) (*TransferResponse, error) {

	balance, err := s.ledgerWriter.GetBusinessBalance(businessID)
	if err != nil {
		return nil, err
	}

	if amount > balance {
		return nil, errors.New("insufficient balance")
	}

	business, err := s.businessWriter.GetBusinessByID(businessID)
	if err != nil {
		return nil, err
	}

	if business.TransferRecipient == "" {
		return nil, errors.New("bank account not set")
	}

	transferRef := fmt.Sprintf("acv_%s", uuid.NewString())

	withdrawal := &models.Withdrawal{
		BusinessID:        businessID,
		Amount:            amount,
		Status:            utils.WithdrawalPending,
		RecipientCode:     business.TransferRecipient,
		TransferReference: transferRef,
	}

	if err := s.repo.CreateWithdrawal(withdrawal); err != nil {
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
		return nil, err
	}

	// withdrawal.TransferReference = res.Data.TransferCode
	// withdrawal.Status = utils.WithdrawalProcessing

	if err := s.repo.UpdateWithdrawal(businessID, res.Data.TransferCode, utils.WithdrawalProcessing); err != nil {
		return nil, err
	}

	return res, nil
}
