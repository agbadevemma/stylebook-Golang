package email

import "context"

type EmailWriter interface {
	SendOTP(ctx context.Context, email string, otp string) error
	SendPasswordReset(ctx context.Context, email string, otp string) error
	SendBookingConfirmation(ctx context.Context, email string, data BookingConfirmationData) error
	SendAccountVerified(ctx context.Context, email string, data AccountVerifiedData) error
}
