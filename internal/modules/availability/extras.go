package availability

import (
	"time"

	"github.com/stylebookgo/agbadevemma/internal/models"
)

func isBlocked(start, end time.Time, blocks []models.BlockedTime) bool {
	for _, b := range blocks {
		if start.Before(b.EndTime) && end.After(b.StartTime) {
			return true
		}
	}
	return false
}

func isBooked(start, end time.Time, bookings []models.Booking) bool {
	for _, b := range bookings {
		if start.Before(b.EndTime) && end.After(b.StartTime) {
			return true
		}
	}
	return false
}

func combine(date time.Time, t string) (time.Time, error) {

	parsed, err := time.Parse("15:04", t)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		parsed.Hour(),
		parsed.Minute(),
		0,
		0,
		date.Location(),
	), nil
}
