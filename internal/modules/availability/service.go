package availability

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/funct"
	"github.com/stylebookgo/agbadevemma/internal/modules/blockedtime"
	"github.com/stylebookgo/agbadevemma/internal/modules/booking"
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
	"github.com/stylebookgo/agbadevemma/internal/modules/businesshour"
)

type Service struct {
	bookingWriter      booking.BookingWriter
	blockedTimesWriter blockedtime.BlockedTimesWriter
	businessHourWriter businesshour.BusinessHourWriter
	businessWriter     business.BusinessWriter
}

func NewService(bookingWriter booking.BookingWriter, blockedTimesWriter blockedtime.BlockedTimesWriter, businessHourWriter businesshour.BusinessHourWriter, businessWriter business.BusinessWriter) *Service {
	return &Service{
		bookingWriter:      bookingWriter,
		blockedTimesWriter: blockedTimesWriter,
		businessHourWriter: businessHourWriter,
		businessWriter:     businessWriter,
	}
}

func (s *Service) GetAvailabilityRange(
	businessID uuid.UUID,
	startDate, endDate time.Time,
) ([]DayAvailability, error) {

	business, err := s.businessWriter.GetBusinessByID(businessID)
	if err != nil {
		return nil, errors.New("business not found")
	}

	loc, err := time.LoadLocation(business.Timezone)
	if err != nil {
		return nil, err
	}

	startDate = time.Date(
		startDate.Year(),
		startDate.Month(),
		startDate.Day(),
		0, 0, 0, 0,
		loc,
	)

	endDate = time.Date(
		endDate.Year(),
		endDate.Month(),
		endDate.Day(),
		23, 59, 59, 0,
		loc,
	)

	var response []DayAvailability

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {

		slots, err := s.getDaySlots(businessID, d)
		if err != nil {
			return nil, err
		}

		response = append(response, DayAvailability{
			Date:  d.Format("2006-01-02"),
			Slots: slots,
		})
	}

	return response, nil
}

func (s *Service) getDaySlots(
	businessID uuid.UUID,
	date time.Time,
) ([]SlotResponse, error) {

	day := int(date.Weekday())

	hour, err := s.businessHourWriter.GetBusinessHour(businessID, day)

	if hour.StartTime == "" || hour.EndTime == "" {
		return []SlotResponse{}, nil
	}

	if err != nil || !hour.IsOpen {
		return []SlotResponse{}, nil
	}

	start, err := funct.Combine(date, hour.StartTime)
	if err != nil {
		return nil, err
	}
	end, err := funct.Combine(date, hour.EndTime)
	if err != nil {
		return nil, err
	}

	blocked, err := s.blockedTimesWriter.GetBlockedTimes(businessID, start, end)
	if err != nil {
		return nil, err
	}

	bookings, err := s.bookingWriter.GetBookings(businessID, start, end)
	if err != nil {
		return nil, err
	}

	slotDuration := 30 * time.Minute // or service-based later
	step := 15 * time.Minute

	var slots []SlotResponse
	cursor := start

	for !cursor.Add(slotDuration).After(end) {

		slotEnd := cursor.Add(slotDuration)

		if funct.IsBlocked(cursor, slotEnd, blocked) {
			cursor = cursor.Add(step)
			continue
		}

		if funct.IsBooked(cursor, slotEnd, bookings) {
			cursor = cursor.Add(step)
			continue
		}

		slots = append(slots, SlotResponse{
			Time: cursor.Format("15:04"),
		})

		cursor = cursor.Add(15 * time.Minute)
	}

	return slots, nil
}
