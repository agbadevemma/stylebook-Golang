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

type WithdrawalStatus string

const (
	WithdrawalPending    WithdrawalStatus = "pending"
	WithdrawalProcessing WithdrawalStatus = "processing"
	WithdrawalSuccess    WithdrawalStatus = "success"
	WithdrawalFailed     WithdrawalStatus = "failed"
)

type LedgerType string

const (
	LedgerCredit LedgerType = "credit"
	LedgerDebit  LedgerType = "debit"
)
