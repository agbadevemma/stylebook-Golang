package availability

import (
	"context"
	"errors"
	"time"

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

func NewService(
	bookingWriter booking.BookingWriter,
	blockedTimesWriter blockedtime.BlockedTimesWriter,
	businessHourWriter businesshour.BusinessHourWriter,
	businessWriter business.BusinessWriter,
) *Service {

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

	// QUERY 1
	biz, err := s.businessWriter.GetBusinessByID(
		ctx,
		nil,
		businessID,
	)

	if err != nil {
		return nil, errors.New("business not found")
	}

	loc, err := time.LoadLocation(biz.Timezone)

	if err != nil {
		return nil, err
	}

	// QUERY 2
	hours, err := s.businessHourWriter.GetBusinessHours(
		ctx,
		nil,
		businessID,
	)

	if err != nil {
		return nil, err
	}

	hourMap := make(map[int]models.BusinessHour)

	for _, h := range hours {
		hourMap[h.DayOfWeek] = h
	}

	// QUERY 3
	bookings, err := s.bookingWriter.GetBookings(
		ctx,
		businessID,
		startDate.UTC(),
		endDate.AddDate(0, 0, 1).UTC(),
	)

	if err != nil {
		return nil, err
	}

	bookingMap := make(map[string][]models.Booking)

	for _, b := range bookings {

		key := b.StartTime.In(loc).
			Format("2006-01-02")

		bookingMap[key] = append(
			bookingMap[key],
			b,
		)
	}

	// QUERY 4
	blockedTimes, err := s.blockedTimesWriter.GetBlockedTimes(
		ctx,
		nil,
		businessID,
		startDate.UTC(),
		endDate.AddDate(0, 0, 1).UTC(),
	)

	if err != nil {
		return nil, err
	}

	blockedMap := make(map[string][]models.BlockedTime)

	for _, b := range blockedTimes {

		key := b.StartTime.In(loc).
			Format("2006-01-02")

		blockedMap[key] = append(
			blockedMap[key],
			b,
		)
	}

	var response []DayAvailability

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {

		localDate := d.In(loc)

		day := int(localDate.Weekday())

		hour, exists := hourMap[day]

		if !exists || !hour.IsOpen {

			response = append(response, DayAvailability{
				Date:  d.Format("2006-01-02"),
				Slots: []SlotResponse{},
			})

			continue
		}

		slots, err := generateSlots(
			localDate,
			hour,
			bookingMap[d.Format("2006-01-02")],
			blockedMap[d.Format("2006-01-02")],
		)

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

func generateSlots(
	dateLocal time.Time,
	hour models.BusinessHour,
	bookings []models.Booking,
	blocked []models.BlockedTime,
) ([]SlotResponse, error) {

	start, err := funct.Combine(
		dateLocal,
		hour.StartTime,
	)

	if err != nil {
		return nil, err
	}

	end, err := funct.Combine(
		dateLocal,
		hour.EndTime,
	)

	if err != nil {
		return nil, err
	}

	startUTC := start.UTC()
	endUTC := end.UTC()

	slotDuration := 30 * time.Minute
	step := 15 * time.Minute

	var slots []SlotResponse

	cursor := startUTC

	for !cursor.Add(slotDuration).After(endUTC) {

		slotEnd := cursor.Add(slotDuration)

		if funct.IsBlocked(
			cursor,
			slotEnd,
			blocked,
		) {

			cursor = cursor.Add(step)
			continue
		}

		if funct.IsBooked(
			cursor,
			slotEnd,
			bookings,
		) {

			cursor = cursor.Add(step)
			continue
		}

		slots = append(
			slots,
			SlotResponse{
				Time: cursor,
			},
		)

		cursor = cursor.Add(step)
	}

	return slots, nil
}
