package booking

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/funct"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/modules/blockedtime"
	"github.com/stylebookgo/agbadevemma/internal/modules/bookingservice"
	"github.com/stylebookgo/agbadevemma/internal/modules/businesshour"
	"github.com/stylebookgo/agbadevemma/internal/modules/service"
)

type Service struct {
	repo                 *Repository
	bookingserviceWriter bookingservice.BookingServiceWriter
	businessHourWriter   businesshour.BusinessHourWriter
	blockedTimesWriter   blockedtime.BlockedTimesWriter
	serviceWriter        service.ServiceWriter
}

func NewService(repo *Repository, bookingserviceWriter bookingservice.BookingServiceWriter, serviceWriter service.ServiceWriter) *Service {
	return &Service{repo: repo, bookingserviceWriter: bookingserviceWriter, serviceWriter: serviceWriter}
}

func (s *Service) CreateBooking(req *CreateBookingRequest, businessID uuid.UUID) error {
	tx := s.repo.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	startTime := req.StartTime
	endTime := startTime.Add(time.Duration(req.TotalDuration) * time.Minute)

	// ================================
	// 1. LOCK BUSINESS (light safety lock)
	// ================================
	if err := tx.Exec(`
		SELECT id FROM businesses WHERE id = ? FOR UPDATE
	`, businessID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// ================================
	// 2. CHECK BUSINESS HOURS
	// ================================
	day := int(startTime.Weekday())

	businessHour, err := s.businessHourWriter.GetBusinessHour(businessID, day)
	if err != nil {
		tx.Rollback()
		return errors.New("business hours not found")
	}

	if !businessHour.IsOpen {
		tx.Rollback()
		return errors.New("business is closed on this day")
	}

	openTime, err := funct.Combine(startTime, businessHour.StartTime)
	if err != nil {
		tx.Rollback()
		return err
	}

	closeTime, err := funct.Combine(startTime, businessHour.EndTime)
	if err != nil {
		tx.Rollback()
		return err
	}

	if startTime.Before(openTime) || endTime.After(closeTime) {
		tx.Rollback()
		return errors.New("booking outside business hours")
	}

	// ================================
	// 3. LOCK OVERLAPPING BOOKINGS (IMPORTANT PART)
	// ================================
	var conflicts []models.Booking

	if err := tx.Raw(`
		SELECT *
		FROM bookings
		WHERE business_id = ?
		AND start_time < ?
		AND end_time > ?
		FOR UPDATE
	`, businessID, endTime, startTime).Scan(&conflicts).Error; err != nil {
		tx.Rollback()
		return err
	}

	if len(conflicts) > 0 {
		tx.Rollback()
		return errors.New("time slot already booked")
	}

	// ================================
	// 4. CHECK BLOCKED TIMES
	// ================================
	blocked, _ := s.blockedTimesWriter.GetBlockedTimes(businessID, startTime, endTime)

	if funct.IsBlocked(startTime, endTime, blocked) {
		tx.Rollback()
		return errors.New("time slot is blocked")
	}

	// ================================
	// 5. CREATE BOOKING
	// ================================
	booking := &models.Booking{
		CustomerName:  req.CustomerName,
		CustomerPhone: req.CustomerPhone,
		BusinessID:    businessID,
		TotalDuration: req.TotalDuration,
		StartTime:     startTime,
		EndTime:       endTime,
	}

	if err := s.repo.CreateBooking(tx, booking); err != nil {
		tx.Rollback()
		return err
	}

	// ================================
	// 6. ATTACH SERVICES
	// ================================
	for _, serviceID := range req.ServiceIDs {

		service, err := s.serviceWriter.GetServiceByID(serviceID)
		if err != nil {
			tx.Rollback()
			return err
		}

		bookingService := &models.BookingService{
			BookingID: booking.ID,
			ServiceID: serviceID,
			Name:      service.Name,
			Duration:  service.Duration,
			Price:     service.Price,
		}

		if err := s.bookingserviceWriter.CreateBookingService(tx, bookingService); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
