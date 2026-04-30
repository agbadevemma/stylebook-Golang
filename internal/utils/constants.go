package utils

type BookingStatus string

const (
	Booked    BookingStatus = "booked"
	Cancelled BookingStatus = "cancelled"
	Done      BookingStatus = "done"
	Missed    BookingStatus = "missed"
)

type PaymentStatus string

const (
	Active    PaymentStatus = "paid"
	Overdue   PaymentStatus = "unpaid"
	Suspended PaymentStatus = "suspended"
)

type FlowState string

const (
	StateStart   FlowState = "start"
	StateService FlowState = "service"
	StateStaff   FlowState = "staff"
	StateTime    FlowState = "time"
	StateConfirm FlowState = "confirm"
)
