package utils

type BookingStatus string

const (
	BookingPending   BookingStatus = "pending"
	BookingConfirmed BookingStatus = "confirmed"
	BookingCancelled BookingStatus = "cancelled"
)

type PaymentStatus string
const (
	PaymentPending   PaymentStatus = "pending"
	PaymentSuccess   PaymentStatus = "success"
	PaymentFailed    PaymentStatus = "failed"
	PaymentAbandoned PaymentStatus = "abandoned"
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

const (
	Admin  UserRole = "admin"
	Barber UserRole = "barber"
)
