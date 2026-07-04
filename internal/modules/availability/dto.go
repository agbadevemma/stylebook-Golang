package availability

import "time"

type SlotResponse struct {
	Time time.Time `json:"time"` // ISO 8601 UTC timestamp
}

type DayAvailability struct {
	Date  string        `json:"date"`
	Slots []SlotResponse `json:"slots"`
}