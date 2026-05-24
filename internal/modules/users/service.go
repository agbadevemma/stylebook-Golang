package users

import (
	"errors"

	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Signup(req *SignupRequest) error {
	_, err := s.repo.GetUserbyEmail(req.Email)
	if err == nil {
		return errors.New("user with this email already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)

	if err != nil {
		return errors.New("failed to hash password")
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     req.Role,
		Password: string(hashedPassword),
	}

	return s.repo.CreateUser(user)
}

func (s *Service) Login(req *LoginRequest) (string, error) {
	user, err := s.repo.GetUserbyEmail(req.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
