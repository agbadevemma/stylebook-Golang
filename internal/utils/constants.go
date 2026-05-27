package utils

type BookingStatus string

const (
	Pending   BookingStatus = "pending"
	Confirmed BookingStatus = "confirmed"
	Cancelled BookingStatus = "cancelled"
	Completed BookingStatus = "completed"
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

type UserRole string

const(
	Admin UserRole = "admin"
	Barber UserRole = "barber"
)