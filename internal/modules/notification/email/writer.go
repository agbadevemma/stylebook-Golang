package email

type EmailWriter interface {
	SendOTP(email string, otp string) error
	SendPasswordReset(email string, otp string) error
	SendBookingConfirmation(email string, data BookingConfirmationData) error
}
