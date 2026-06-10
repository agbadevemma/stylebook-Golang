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

func CalculateFee(serviceAmount int64) (platformFee, netAmount, totalAmount int64) {

	fee := int64(float64(serviceAmount) * 0.10)

	if fee < 10000 {
		fee = 10000
	}

	if fee > 50000 {
		fee = 50000
	}

	total := serviceAmount + fee

	return fee, serviceAmount, total
}
