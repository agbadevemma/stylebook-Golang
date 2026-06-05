package availability

type SlotResponse struct {
	Time  string `json:"time"`
}

type DayAvailability struct {
	Date  string        `json:"date"`
	Slots []SlotResponse `json:"slots"`
}