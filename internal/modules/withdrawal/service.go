package withdrawal

import (
	"github.com/stylebookgo/agbadevemma/internal/modules/ledger"
)

type Service struct {
	repo         *Repository
	provider     *Provider
	ledgerWriter ledger.LedgerWriter
}

func NewService(repo *Repository, provider *Provider, ledgerWriter ledger.LedgerWriter) *Service {
	return &Service{repo: repo, provider: provider, ledgerWriter: ledgerWriter}
}

// func (s *Service) CreateWithdrawal(businessID uuid.UUID, amount int64) (*TransferResponse, error) {

// 	balance, err := s.ledgerWriter.GetBusinessBalance(businessID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if amount > balance {
// 		return nil, errors.New("insufficient balance")
// 	}

// 	business, err := s.repo.GetBusinessByID(businessID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if business.RecipientCode == "" {
// 		return nil, errors.New("bank account not set")
// 	}

// 	withdrawal := &models.Withdrawal{
// 		BusinessID:    businessID,
// 		Amount:        amount,
// 		Status:        utils.WithdrawalPending,
// 		RecipientCode: business.RecipientCode,
// 	}

// 	if err := s.repo.CreateWithdrawal(withdrawal); err != nil {
// 		return nil, err
// 	}

// 	res, err := s.provider.CreateTransfer(TransferRequest{
// 		Source:    "balance",
// 		Amount:    amount,
// 		Recipient: business.RecipientCode,
// 		Reason:    "Business withdrawal",
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	withdrawal.TransferCode = res.Data.TransferCode
// 	withdrawal.Status = utils.WithdrawalProcessing

// 	if err := s.repo.UpdateWithdrawal(withdrawal); err != nil {
// 		return nil, err
// 	}

// 	return res, nil
// }
