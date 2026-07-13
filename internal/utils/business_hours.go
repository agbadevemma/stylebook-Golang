package utils

import "time"

// IsTimeWithinBusinessHours checks if a booking time is within business hours
// Handles business hours that cross midnight (e.g., 16:00 to 02:00)
func IsTimeWithinBusinessHours(
	bookingStartLocal, bookingEndLocal time.Time,
	openTime, closeTime time.Time,
) bool {
	// Extract just the hour and minute for comparison
	bookingStartHour := bookingStartLocal.Hour()
	bookingStartMin := bookingStartLocal.Minute()
	bookingEndHour := bookingEndLocal.Hour()
	bookingEndMin := bookingEndLocal.Minute()

	businessOpenHour := openTime.Hour()
	businessOpenMin := openTime.Minute()
	businessCloseHour := closeTime.Hour()
	businessCloseMin := closeTime.Minute()

	// Convert times to minutes for easier comparison
	bookingStart := bookingStartHour*60 + bookingStartMin
	bookingEnd := bookingEndHour*60 + bookingEndMin
	businessOpen := businessOpenHour*60 + businessOpenMin
	businessClose := businessCloseHour*60 + businessCloseMin

	// Case 1: Business hours do NOT cross midnight (e.g., 09:00 to 17:00)
	if businessOpen <= businessClose {
		// Both booking start and end must be between open and close
		return bookingStart >= businessOpen && bookingEnd <= businessClose
	}

	// Case 2: Business hours CROSS midnight (e.g., 16:00 to 02:00)
	// Booking is valid if:
	// - it starts after open time (regardless of the time being after midnight)
	// - AND it ends before close time (regardless of the time being before midnight)
	// - AND the booking end is on the same calendar day as the start OR crosses to next day
	if bookingEnd < bookingStart {
		// Booking crosses midnight (started today, ends tomorrow)
		// Valid if: booking starts after business open
		return bookingStart >= businessOpen
	} else {
		// Booking doesn't cross midnight
		// Valid if either:
		// 1. booking is after open (and business hours cross midnight)
		// 2. OR booking is before close (and business hours cross midnight)
		return bookingStart >= businessOpen || bookingEnd <= businessClose
	}
}
