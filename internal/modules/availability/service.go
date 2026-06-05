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

	if _, err := s.businessWriter.GetBusinessByID(businessID); err != nil {
		return nil, errors.New("business not found")
	}

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

	blocked, _ := s.blockedTimesWriter.GetBlockedTimes(businessID, start, end)
	bookings, _ := s.bookingWriter.GetBookings(businessID, start, end)

	slotDuration := 30 // or service-based later

	var slots []SlotResponse
	cursor := start

	for cursor.Add(time.Duration(slotDuration) * time.Minute).Before(end) {

		slotEnd := cursor.Add(time.Duration(slotDuration) * time.Minute)

		if funct.IsBlocked(cursor, slotEnd, blocked) {
			cursor = cursor.Add(15 * time.Minute)
			continue
		}

		if funct.IsBooked(cursor, slotEnd, bookings) {
			cursor = cursor.Add(15 * time.Minute)
			continue
		}

		slots = append(slots, SlotResponse{
			Time: cursor.Format("15:04"),
		})

		cursor = cursor.Add(15 * time.Minute)
	}

	return slots, nil
}
