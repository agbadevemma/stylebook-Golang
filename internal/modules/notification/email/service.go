package email

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
)

type Service struct {
	client *resend.Client
	from   string
}

func NewEmailService() *Service {
	return &Service{
		client: resend.NewClient(os.Getenv("RESEND_API_KEY")),
		from:   os.Getenv("EMAIL_FROM"),
	}
}

func (s *Service) SendOTP(email string, otp string) error {
	subject := "Your Verification Code"

	html := fmt.Sprintf(`
		<div style="font-family: Arial;">
			<h2>Email Verification</h2>
			<p>Your OTP code is:</p>

			<h1 style="letter-spacing: 4px;">%s</h1>

			<p>This code will expire in 10 minutes.</p>

			<hr />
			<small>If you did not request this, ignore this email.</small>
		</div>
	`, otp)

	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.from,
		To:      []string{email},
		Subject: subject,
		Html:    html,
	})

	return err
}

func (s *Service) SendPasswordReset(email string, otp string) error {
	subject := "Password Reset Code"

	html := fmt.Sprintf(`
		<div style="font-family: Arial;">
			<h2>Password Reset Request</h2>

			<p>Use the OTP below to reset your password:</p>

			<h1 style="letter-spacing: 4px;">%s</h1>

			<p>This code expires in 10 minutes.</p>

			<p>If you didn’t request this, please ignore this email.</p>
		</div>
	`, otp)

	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.from,
		To:      []string{email},
		Subject: subject,
		Html:    html,
	})

	return err
}
