package availability

import (
	"errors"
	"time"

	"context"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/funct"
	"github.com/stylebookgo/agbadevemma/internal/models"
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
	ctx context.Context,
	businessID uuid.UUID,
	startDate, endDate time.Time,
) ([]DayAvailability, error) {

	business, err := s.businessWriter.GetBusinessByID(ctx, nil, businessID)
	if err != nil {
		return nil, errors.New("business not found")
	}

	loc, err := time.LoadLocation(business.Timezone)
	if err != nil {
		return nil, err
	}

	// Fetch ALL bookings once
	bookings, err := s.bookingWriter.GetBookings(
		ctx,
		businessID,
		startDate.UTC(),
		endDate.AddDate(0, 0, 1).UTC(),
	)

	if err != nil {
		return nil, err
	}
	// Iterate through UTC dates, but determine day-of-week based on business timezone
	// This ensures we get business hours for the correct calendar day in business's local time
	var response []DayAvailability

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		// Convert UTC date to business timezone to determine which calendar day it represents
		dLocal := d.In(loc)

		slots, err := s.getDaySlots(ctx, businessID, dLocal, bookings)
		if err != nil {
			return nil, err
		}

		response = append(response, DayAvailability{
			Date:  d.Format("2006-01-02"), // Return UTC date as YYYY-MM-DD
			Slots: slots,
		})
	}

	return response, nil
}

func (s *Service) getDaySlots(
	ctx context.Context,
	businessID uuid.UUID,
	dateLocal time.Time,
	bookings []models.Booking,
) ([]SlotResponse, error) {

	// Use local time to determine which day of week (in business's timezone)
	day := int(dateLocal.Weekday())

	hour, err := s.businessHourWriter.GetBusinessHour(ctx, nil, businessID, day)

	if hour.StartTime == "" || hour.EndTime == "" {
		return []SlotResponse{}, nil
	}

	if err != nil || !hour.IsOpen {
		return []SlotResponse{}, nil
	}

	// Combine business hours (local time strings) with local date
	start, err := funct.Combine(dateLocal, hour.StartTime)
	if err != nil {
		return nil, err
	}
	startUTC := start.UTC()

	end, err := funct.Combine(dateLocal, hour.EndTime)
	if err != nil {
		return nil, err
	}
	endUTC := end.UTC()

	blocked, err := s.blockedTimesWriter.GetBlockedTimes(ctx, nil, businessID, startUTC, endUTC)
	if err != nil {
		return nil, err
	}

	slotDuration := 30 * time.Minute
	step := 15 * time.Minute

	var slots []SlotResponse
	cursor := startUTC

	for !cursor.Add(slotDuration).After(endUTC) {

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
			Time: cursor,
		})

		cursor = cursor.Add(15 * time.Minute)
	}

	return slots, nil
}
