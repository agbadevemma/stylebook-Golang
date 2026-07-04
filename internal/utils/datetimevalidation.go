package utils

import (
	"errors"
	"regexp"
	"time"
)

// ValidateDateFormat validates that a date string is in YYYY-MM-DD format
func ValidateDateFormat(dateStr string) error {
	if dateStr == "" {
		return errors.New("date cannot be empty")
	}

	// Check format with regex
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	if !dateRegex.MatchString(dateStr) {
		return errors.New("date must be in YYYY-MM-DD format")
	}

	// Try to parse to ensure it's a valid date
	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return errors.New("invalid date: " + err.Error())
	}

	return nil
}

// ValidateISO8601Timestamp validates that a timestamp string is in ISO 8601 UTC format
func ValidateISO8601Timestamp(timestamp time.Time, fieldName string) error {
	if timestamp.IsZero() {
		return errors.New(fieldName + " cannot be empty")
	}

	// Check that it's in UTC (offset should be +00:00 or Z)
	if timestamp.Location() != time.UTC && timestamp.Location().String() != "UTC" {
		// Allow local timezone but convert expectation
		// Actually, we should enforce UTC
		return errors.New(fieldName + " must be in UTC timezone")
	}

	return nil
}

// ValidateISO8601TimestampString validates that a timestamp string is valid ISO 8601 format
func ValidateISO8601TimestampString(timestampStr string, fieldName string) error {
	if timestampStr == "" {
		return errors.New(fieldName + " cannot be empty")
	}

	// Try to parse as RFC3339 (ISO 8601)
	parsedTime, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		return errors.New(fieldName + " must be in ISO 8601 format (e.g., 2026-07-04T08:00:00Z or 2026-07-04T08:00:00+00:00)")
	}

	// Verify it's in UTC
	if parsedTime.Location() != time.UTC {
		return errors.New(fieldName + " must be in UTC timezone (ends with Z)")
	}

	return nil
}

// ValidateDateRange validates that start date is before end date
func ValidateDateRange(startStr, endStr string) error {
	if err := ValidateDateFormat(startStr); err != nil {
		return errors.New("invalid startDate: " + err.Error())
	}
	if err := ValidateDateFormat(endStr); err != nil {
		return errors.New("invalid endDate: " + err.Error())
	}

	startDate, _ := time.Parse("2006-01-02", startStr)
	endDate, _ := time.Parse("2006-01-02", endStr)

	if endDate.Before(startDate) {
		return errors.New("endDate must be after or equal to startDate")
	}

	return nil
}

// ValidateTimeRange validates that start time is before end time
func ValidateTimeRange(startTime, endTime time.Time, startFieldName, endFieldName string) error {
	if err := ValidateISO8601Timestamp(startTime, startFieldName); err != nil {
		return err
	}
	if err := ValidateISO8601Timestamp(endTime, endFieldName); err != nil {
		return err
	}

	if endTime.Before(startTime) {
		return errors.New(endFieldName + " must be after " + startFieldName)
	}

	if startTime.Equal(endTime) {
		return errors.New(startFieldName + " and " + endFieldName + " cannot be the same")
	}

	return nil
}
