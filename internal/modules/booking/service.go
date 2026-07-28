package booking

import (
	"context"
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
	"github.com/stylebookgo/agbadevemma/internal/utils"
	"gorm.io/gorm/clause"
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

func (s *Service) CreateBooking(ctx context.Context, req *CreateBookingRequest, businessID uuid.UUID) (BookingResponse, error) {
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
		service, err := s.serviceWriter.GetServiceByID(ctx, tx, serviceID)
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
	business, err := s.businessWriter.GetBusinessByID(ctx, tx, businessID)
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
	businessHour, err := s.businessHourWriter.GetBusinessHour(ctx, tx, businessID, dayOfWeek)
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

	// Get booking times in business timezone
	startTimeLocal := startTimeUTC.In(loc)
	endTimeLocal := endTimeUTC.In(loc)

	// Check if booking is within business hours (handles midnight crossing)
	if !utils.IsTimeWithinBusinessHours(startTimeLocal, endTimeLocal, startHourLocal, endHourLocal) {
		tx.Rollback()
		return BookingResponse{}, errors.New("booking time is outside business hours")
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
		ctx,
		tx,
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

	if err := s.repo.CreateBooking(ctx, tx, booking); err != nil {
		tx.Rollback()
		return BookingResponse{}, err
	}

	// ================================
	// 9. ATTACH SERVICES
	// ================================
	for _, sID := range req.ServiceIDs {

		service, err := s.serviceWriter.GetServiceByID(ctx, tx, sID)
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

		if err := s.bookingserviceWriter.CreateBookingService(ctx, tx, bookingService); err != nil {
			tx.Rollback()
			return BookingResponse{}, err
		}
	}

	// ================================
	// 10. PAYMENT
	// ================================
	paymentURL, err := s.paymentWriter.InitializePayment(ctx, tx, booking)
	if err != nil {
		tx.Rollback()
		return BookingResponse{}, err
	}

	return BookingResponse{
		BookingID:  booking.ID,
		PaymentURL: paymentURL,
	}, tx.Commit().Error
}

func (s *Service) GetBusinessBookingStats(ctx context.Context, businessID uuid.UUID) (BookingStatsResponse, error) {
	_, err := s.businessWriter.GetBusinessByID(ctx, s.repo.db, businessID)
	if err != nil {
		return BookingStatsResponse{}, err
	}

	total, today, pending, clients, err := s.repo.GetBusinessBookingStats(ctx, businessID)
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

func (s *Service) ListBusinessBookings(ctx context.Context, businessID uuid.UUID, params ListBookingsParams) ([]models.Booking, int64, error) {
	_, err := s.businessWriter.GetBusinessByID(ctx, s.repo.db, businessID)
	if err != nil {
		return nil, 0, err
	}

	return s.repo.ListBusinessBookings(ctx, businessID, params)
}

func (s *Service) RescheduleBooking(ctx context.Context, bookingID uuid.UUID, businessID uuid.UUID, req *RescheduleBookingRequest) error {
	// Validate startTime is in UTC
	if err := utils.ValidateISO8601Timestamp(req.NewStartTime, "newStartTime"); err != nil {
		return errors.New("invalid start time format: " + err.Error())
	}

	// ================================
	// 1. START TRANSACTION WITH LOCKING
	// ================================
	tx := s.repo.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// ================================
	// 2. LOCK BOOKING ROW
	// ================================
	var booking models.Booking
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&booking, bookingID).Error; err != nil {
		tx.Rollback()
		return errors.New("booking not found or cannot be locked")
	}

	// Only paid bookings can be rescheduled
	if booking.PaymentStatus != utils.PaymentSuccess {
		tx.Rollback()
		return errors.New("only paid bookings can be rescheduled")
	}

	// ================================
	// 3. LOCK BUSINESS ROW
	// ================================
	var business models.Business

	if err := tx.Clauses(clause.Locking{
		Strength: "UPDATE",
	}).First(&business, businessID).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to lock business")
	}

	// ================================
	// 4. GET BUSINESS DETAILS
	// ================================
	// business, err := s.businessWriter.GetBusinessByID(ctx, tx, businessID)
	// if err != nil {
	// 	tx.Rollback()
	// 	return errors.New("business not found")
	// }

	// Calculate new end time based on total duration
	newEndTimeUTC := req.NewStartTime.Add(time.Duration(booking.TotalDuration) * time.Minute)

	// Load timezone
	loc, err := time.LoadLocation(business.Timezone)
	if err != nil {
		tx.Rollback()
		return errors.New("invalid business timezone")
	}

	// Convert UTC booking time to business timezone to determine calendar day
	bookingLocal := req.NewStartTime.In(loc)
	dayOfWeek := int(bookingLocal.Weekday())

	// Get business hours for this day of week
	businessHour, err := s.businessHourWriter.GetBusinessHour(ctx, tx, businessID, dayOfWeek)
	if err != nil {
		tx.Rollback()
		return errors.New("business hours not configured for this day")
	}

	// Parse business hours
	openTime, err := time.Parse("15:04", businessHour.StartTime)
	if err != nil {
		tx.Rollback()
		return errors.New("invalid business open time format")
	}

	closeTime, err := time.Parse("15:04", businessHour.EndTime)
	if err != nil {
		tx.Rollback()
		return errors.New("invalid business close time format")
	}

	// Get start and end times in business timezone
	bookingStart := req.NewStartTime.In(loc)
	bookingEnd := newEndTimeUTC.In(loc)

	// Check if booking is within business hours (handles midnight crossing)
	if !utils.IsTimeWithinBusinessHours(bookingStart, bookingEnd, openTime, closeTime) {
		tx.Rollback()
		return errors.New("booking time is outside business hours")
	}

	// ================================
	// 5. CHECK FOR BOOKING CONFLICTS WITH LOCKING
	// ================================
	var conflicts []models.Booking
	if err := tx.Raw(`
		SELECT *
		FROM bookings
		WHERE business_id = ?
		AND id != ?
		AND start_time < ?
		AND end_time > ?
		FOR UPDATE
	`, businessID, bookingID, newEndTimeUTC, req.NewStartTime).Scan(&conflicts).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to check for conflicts")
	}

	if len(conflicts) > 0 {
		tx.Rollback()
		return errors.New("new time slot conflicts with another booking")
	}

	// ================================
	// 6. CHECK FOR BLOCKED TIMES
	// ================================
	blockedTimes, err := s.blockedTimesWriter.GetBlockedTimes(ctx, tx, businessID, req.NewStartTime, newEndTimeUTC)
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(blockedTimes) > 0 {
		tx.Rollback()
		return errors.New("new time slot is blocked")
	}

	// ================================
	// 7. UPDATE BOOKING
	// ================================
	booking.StartTime = req.NewStartTime
	booking.EndTime = newEndTimeUTC

	if err := s.repo.UpdateBooking(ctx, tx, &booking); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *Service) UpdateBookingStatus(ctx context.Context, bookingID uuid.UUID, businessID uuid.UUID, status string) error {
	// Validate status is one of the allowed values
	validStatuses := map[string]bool{
		string(utils.BookingPending):   true,
		string(utils.BookingConfirmed): true,
		string(utils.BookingCancelled): true,
		string(utils.BookingCompleted): true,
	}

	if !validStatuses[status] {
		return errors.New("invalid booking status: " + status)
	}

	// Get the booking with locking
	tx := s.repo.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var booking models.Booking
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&booking, "id = ? AND business_id = ?", bookingID, businessID).Error; err != nil {
		tx.Rollback()
		return errors.New("booking not found")
	}

	// Verify booking belongs to the business
	if booking.BusinessID != businessID {
		tx.Rollback()
		return errors.New("booking does not belong to this business")
	}

	// Update the status
	booking.BookingStatus = utils.BookingStatus(status)

	if err := s.repo.UpdateBooking(ctx, tx, &booking); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *Service) GetBookingDetailsPublic(ctx context.Context, bookingID uuid.UUID) (*models.Booking, []models.BookingService, *models.Business, error) {
	booking, err := s.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, nil, nil, errors.New("booking not found")
	}

	// Get booking services
	services, err := s.bookingserviceWriter.GetBookingServicesByBookingID(ctx, s.repo.db, bookingID)
	if err != nil {
		return nil, nil, nil, errors.New("failed to fetch booking services")
	}

	// Get business details
	business, err := s.businessWriter.GetBusinessByID(ctx, s.repo.db, booking.BusinessID)
	if err != nil {
		return nil, nil, nil, errors.New("failed to fetch business details")
	}

	return booking, services, business, nil
}
