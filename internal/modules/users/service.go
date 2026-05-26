package users

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
	"github.com/stylebookgo/agbadevemma/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository

	businessWriter business.BusinessWriter
}

func NewService(repo *Repository, bus business.BusinessWriter) *Service {
	return &Service{repo: repo, businessWriter: bus}
}

func (s *Service) Signup(req *SignupRequest) error {
	
	if _, err := s.repo.GetUserbyEmail(req.Email); err == nil {
		return errors.New("user with this email already exists")
	}

	if _, err := s.repo.GetUserByPhone(req.Phone); err == nil {
		return errors.New("user with this phone number already exists")
	}

	var businessID *uuid.UUID

	if req.Role == utils.Barber {
		if req.BusinessID == uuid.Nil {
			return errors.New("business ID is required for barbers")
		}

		if _, err := s.businessWriter.GetBusinessByID(req.BusinessID); err != nil {
			return errors.New("business not found")
		}

		businessID = &req.BusinessID
	}


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

	if err := s.repo.CreateUser(tx, user); err != nil {
		tx.Rollback()
		return err
	}

	if req.Role == utils.Admin {
		business := &models.Business{
			Name:    req.BusinessName,
			OwnerID: user.ID,
		}

		if err := s.businessWriter.CreateBusiness(tx, business); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
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
