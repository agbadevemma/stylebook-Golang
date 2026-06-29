package business

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetBusinessSettings(businessID uuid.UUID) (*models.Business, error) {
	return s.repo.GetBusinessByID(businessID)
}

func (s *Service) UpdateBusinessDetails(
	businessID uuid.UUID,
	req UpdateBusinessDetailsRequest,
) (*models.Business, error) {
	name := strings.TrimSpace(req.Name)
	timezone := strings.TrimSpace(req.Timezone)

	if name == "" {
		return nil, errors.New("name is required")
	}

	if len(name) > 100 {
		return nil, errors.New("name is too long")
	}

	if timezone == "" {
		return nil, errors.New("timezone is required")
	}

	if err := s.repo.UpdateBusinessDetails(businessID, name, timezone); err != nil {
		return nil, err
	}

	return s.repo.GetBusinessByID(businessID)
}

func (s *Service) SetupWithdrawalPIN(businessID uuid.UUID, pin string) (*models.Business, error) {
	if len(pin) != 4 {
		return nil, errors.New("PIN must be exactly 4 digits")
	}

	// Verify it's actually numeric
	for _, r := range pin {
		if r < '0' || r > '9' {
			return nil, errors.New("PIN must contain only digits")
		}
	}

	// Hash the PIN using bcrypt
	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateWithdrawalPIN(businessID, string(hashedPIN)); err != nil {
		return nil, err
	}

	return s.repo.GetBusinessByID(businessID)
}

func (s *Service) UpdateWithdrawalPIN(businessID uuid.UUID, oldPin string, newPin string) (*models.Business, error) {
	if len(oldPin) != 4 || len(newPin) != 4 {
		return nil, errors.New("PINs must be exactly 4 digits")
	}

	// Verify both are numeric
	for _, r := range oldPin {
		if r < '0' || r > '9' {
			return nil, errors.New("old PIN must contain only digits")
		}
	}
	for _, r := range newPin {
		if r < '0' || r > '9' {
			return nil, errors.New("new PIN must contain only digits")
		}
	}

	business, err := s.repo.GetBusinessByID(businessID)
	if err != nil {
		return nil, err
	}

	if business.WithdrawalPIN == "" {
		return nil, errors.New("withdrawal PIN not set - please set up PIN first")
	}

	// Verify old PIN
	if err := bcrypt.CompareHashAndPassword([]byte(business.WithdrawalPIN), []byte(oldPin)); err != nil {
		return nil, errors.New("incorrect old PIN")
	}

	// Hash the new PIN
	hashedNewPIN, err := bcrypt.GenerateFromPassword([]byte(newPin), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateWithdrawalPIN(businessID, string(hashedNewPIN)); err != nil {
		return nil, err
	}

	return s.repo.GetBusinessByID(businessID)
}
