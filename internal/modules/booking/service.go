package booking

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/funct"
	"github.com/stylebookgo/agbadevemma/internal/models"
	"github.com/stylebookgo/agbadevemma/internal/modules/blockedtime"
	"github.com/stylebookgo/agbadevemma/internal/modules/bookingservice"
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
	"github.com/stylebookgo/agbadevemma/internal/modules/businesshour"
	"github.com/stylebookgo/agbadevemma/internal/modules/payment"
	"github.com/stylebookgo/agbadevemma/internal/modules/service"
)

type Service struct {
	repo                 *Repository
	businessWriter       business.BusinessWriter
	bookingserviceWriter bookingservice.BookingServiceWriter
	businessHourWriter   businesshour.BusinessHourWriter
	blockedTimesWriter   blockedtime.BlockedTimesWriter
	serviceWriter        service.ServiceWriter
	paymentWriter        payment.PaymentWriter
}

func NewService(repo *Repository, bookingserviceWriter bookingservice.BookingServiceWriter, businessHourWriter businesshour.BusinessHourWriter, serviceWriter service.ServiceWriter, blockedTimesWriter blockedtime.BlockedTimesWriter, paymentWriter payment.PaymentWriter, businessWriter business.BusinessWriter) *Service {
	return &Service{repo: repo, bookingserviceWriter: bookingserviceWriter, businessHourWriter: businessHourWriter, blockedTimesWriter: blockedTimesWriter, serviceWriter: serviceWriter, paymentWriter: paymentWriter, businessWriter: businessWriter}
}

func (s *Service) CreateBooking(req *CreateBookingRequest, businessID uuid.UUID) (BookingResponse, error) {
	tx := s.repo.db.Begin()

	if tx.Error != nil {
		return BookingResponse{}, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// ================================
	// 1. LOAD SERVICES + TOTAL DURATION
	// ================================
	var (
		selectedServices []*models.Service
		totalDuration    int // minutes
	)

	for _, serviceID := range req.ServiceIDs {
		service, err := s.serviceWriter.GetServiceByID(serviceID)
		if err != nil {
			tx.Rollback()
			return BookingResponse{}, err
		}

		selectedServices = append(selectedServices, service)
		totalDuration += service.Duration
	}

	// ================================
	// 2. GET BUSINESS + TIMEZONE
	// ================================
	business, err := s.businessWriter.GetBusinessByID(businessID)
	if err != nil {
		tx.Rollback()
		return BookingResponse{}, err
	}

	// Check if business has set up a subaccount (required for payments)
	if business.SubaccountCode == "" || business.SubaccountCode == "null" {
		tx.Rollback()
		return BookingResponse{}, errors.New("business has not completed payment setup. Please add bank details")
	}


	startTimeUTC := req.StartTime
	endTimeUTC := startTimeUTC.Add(
		time.Duration(totalDuration) * time.Minute,
	)

	// ================================
	// 3. BUSINESS HOURS VALIDATION
	// ================================
	loc, err := time.LoadLocation(business.Timezone)
	if err != nil {
		tx.Rollback()
		return BookingResponse{}, errors.New("invalid business timezone")
	}

	// Convert UTC booking time to business timezone to determine calendar day
	bookingLocal := startTimeUTC.In(loc)
	dayOfWeek := int(bookingLocal.Weekday())

	// Get business hours for this day of week
	businessHour, err := s.businessHourWriter.GetBusinessHour(businessID, dayOfWeek)
	if err != nil {
		tx.Rollback()
		return BookingResponse{}, errors.New("failed to fetch business hours")
	}

	// Check if business is open on this day
	if !businessHour.IsOpen || businessHour.StartTime == "" || businessHour.EndTime == "" {
		tx.Rollback()
		return BookingResponse{}, errors.New("business is closed on this day")
	}

	// Parse business hours (local time strings like "09:00", "17:00")
	startHourLocal, err := funct.Combine(bookingLocal, businessHour.StartTime)
	if err != nil {
		tx.Rollback()
		return BookingResponse{}, errors.New("invalid business start time format")
	}

	endHourLocal, err := funct.Combine(bookingLocal, businessHour.EndTime)
	if err != nil {
		tx.Rollback()
		return BookingResponse{}, errors.New("invalid business end time format")
	}

	// Convert business hours back to UTC for comparison
	startHourUTC := startHourLocal.UTC()
	endHourUTC := endHourLocal.UTC()

	// Validate booking is within business hours
	if startTimeUTC.Before(startHourUTC) || startTimeUTC.After(endHourUTC) {
		tx.Rollback()
		return BookingResponse{}, errors.New("booking start time is outside business hours")
	}

	if endTimeUTC.After(endHourUTC) {
		tx.Rollback()
		return BookingResponse{}, errors.New("booking end time extends beyond business hours")
	}

	// ================================
	// 4. LOCK BUSINESS
	// ================================
	if err := tx.Exec(`
		SELECT id FROM businesses WHERE id = ? FOR UPDATE
	`, businessID).Error; err != nil {
		tx.Rollback()
		return BookingResponse{}, err
	}

	// ================================
	// 5. OVERLAP CHECK (UTC ONLY)
	// ================================
	var conflicts []models.Booking

	if err := tx.Raw(`
		SELECT *
		FROM bookings
		WHERE business_id = ?
		AND start_time < ?
		AND end_time > ?
		FOR UPDATE
	`, businessID, endTimeUTC, startTimeUTC).Scan(&conflicts).Error; err != nil {
		tx.Rollback()
		return BookingResponse{}, err
	}

	if len(conflicts) > 0 {
		tx.Rollback()
		return BookingResponse{}, errors.New("time slot already booked")
	}

	// ================================
	// 7. BLOCKED TIMES CHECK
	// ================================
	blocked, err := s.blockedTimesWriter.GetBlockedTimes(
		businessID,
		startTimeUTC,
		endTimeUTC,
	)
	log.Printf("%v", blocked)
	if err != nil {
		tx.Rollback()
		return BookingResponse{}, err
	}

	if funct.IsBlocked(startTimeUTC, endTimeUTC, blocked) {
		tx.Rollback()
		return BookingResponse{}, errors.New("time slot is blocked")
	}

	// ================================
	// 8. CREATE BOOKING
	// ================================
	booking := &models.Booking{
		CustomerName:  req.CustomerName,
		CustomerPhone: req.CustomerPhone,
		CustomerEmail: req.CustomerEmail,
		BusinessID:    businessID,
		TotalDuration: totalDuration,
		StartTime:     startTimeUTC,
		EndTime:       endTimeUTC,
	}

	if err := s.repo.CreateBooking(tx, booking); err != nil {
		tx.Rollback()
		return BookingResponse{}, err
	}

	// ================================
	// 9. ATTACH SERVICES
	// ================================
	for _, sID := range req.ServiceIDs {

		service, err := s.serviceWriter.GetServiceByID(sID)
		if err != nil {
			tx.Rollback()
			return BookingResponse{}, err
		}

		bookingService := &models.BookingService{
			BookingID: booking.ID,
			ServiceID: service.ID,
			Name:      service.Name,
			Duration:  service.Duration,
			Price:     service.Price,
		}

		if err := s.bookingserviceWriter.CreateBookingService(tx, bookingService); err != nil {
			tx.Rollback()
			return BookingResponse{}, err
		}
	}

	// ================================
	// 10. PAYMENT
	// ================================
	paymentURL, err := s.paymentWriter.InitializePayment(tx, booking)
	if err != nil {
		tx.Rollback()
		return BookingResponse{}, err
	}

	return BookingResponse{
		BookingID:  booking.ID,
		PaymentURL: paymentURL,
	}, tx.Commit().Error
}

func (s *Service) GetBusinessBookingStats(businessID uuid.UUID) (BookingStatsResponse, error) {
	_, err := s.businessWriter.GetBusinessByID(businessID)
	if err != nil {
		return BookingStatsResponse{}, err
	}

	total, today, pending, clients, err := s.repo.GetBusinessBookingStats(businessID)
	if err != nil {
		return BookingStatsResponse{}, err
	}

	return BookingStatsResponse{
		TotalBookings:   total,
		TodayBookings:   today,
		PendingBookings: pending,
		TotalClients:    clients,
	}, nil
}

func (s *Service) ListBusinessBookings(businessID uuid.UUID, params ListBookingsParams) ([]models.Booking, int64, error) {
	_, err := s.businessWriter.GetBusinessByID(businessID)
	if err != nil {
		return nil, 0, err
	}

	return s.repo.ListBusinessBookings(businessID, params)
}
