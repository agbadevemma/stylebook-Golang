package funct

import (
	"time"

	"github.com/stylebookgo/agbadevemma/internal/models"
)

func IsBlocked(start, end time.Time, blocks []models.BlockedTime) bool {
	for _, b := range blocks {
		if start.Before(b.EndTime) && end.After(b.StartTime) {
			return true
		}
	}
	return false
}

func IsBooked(start, end time.Time, bookings []models.Booking) bool {
	for _, b := range bookings {
		if start.Before(b.EndTime) && end.After(b.StartTime) {
			return true
		}
	}
	return false
}

func Combine(date time.Time, t string) (time.Time, error) {

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

func CalculateFee(amount int64) (platformFee int64, netAmount int64) {

	fee := int64(float64(amount) * 0.10)

	if fee < 100 {
		fee = 100
	}

	if fee > 500 {
		fee = 500
	}

	net := amount - fee

	return fee, net
}
