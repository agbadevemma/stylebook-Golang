package businesshour

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
	tx   *gorm.DB
}

func NewService(
	repo *Repository,

) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListByBusiness(
	businessID uuid.UUID,
) ([]models.BusinessHour, error) {

	return s.repo.ListByBusiness(businessID)
}

func (s *Service) UpdateBusinessHour(
	businessID uuid.UUID,
	req []UpdateBusinessHourRequest,
) error {

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
	for _, r := range req {
		hour, err := s.repo.GetByBusinessAndDay(
			tx,
			businessID,
			r.DayOfWeek,
		)
		if err != nil {
			tx.Rollback()
			log.Printf("Error finding business hour: %v", err.Error())
			return errors.New("business hour not found")
		}
		if r.DayOfWeek != 0 {
			hour.DayOfWeek = r.DayOfWeek
		}
		if r.StartTime != "" {
			hour.StartTime = r.StartTime
		}
		if r.EndTime != "" {
			hour.EndTime = r.EndTime
		}
		if r.IsOpen != hour.IsOpen {
			hour.IsOpen = r.IsOpen
		}
		if err := s.repo.UpdateBusinessHour(tx, hour); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (s *Service) CreateBusinessHour(
	businessID uuid.UUID,
	req []CreateBusinessHourRequest,
) error {

	for _, r := range req {
		newHour := &models.BusinessHour{
			BusinessID: businessID,
			DayOfWeek:  r.DayOfWeek,
			IsOpen:     r.IsOpen,
			StartTime:  r.StartTime,
			EndTime:    r.EndTime,
		}

		err := s.repo.CreateBusinessHour(newHour)

		if err != nil {
			log.Printf("Error creating business hour: %v", err.Error())
			return errors.New("failed to create business hour")
		}
	}
	return nil
}
