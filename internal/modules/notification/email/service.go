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

type BookingConfirmationData struct {
	CustomerName  string
	BusinessName  string
	AppointmentTime string
	Services      string
	TotalAmount   string
	BookingID     string
	CallbackURL   string
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
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
		</head>
		<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif; background-color: #f5f5f5;">
			<table width="100%%" cellpadding="0" cellspacing="0" style="background-color: #f5f5f5;">
				<tr>
					<td align="center" style="padding: 40px 20px;">
						<table width="100%%" max-width="600" cellpadding="0" cellspacing="0" style="background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1);">
							<!-- Header -->
							<tr>
								<td style="background-color: #000000; padding: 40px 30px; text-align: center; border-radius: 8px 8px 0 0;">
									<h1 style="margin: 0; color: #ffffff; font-size: 28px; font-weight: 700;">Email Verification</h1>
								</td>
							</tr>

							<!-- Content -->
							<tr>
								<td style="padding: 40px 30px;">
									<p style="margin: 0 0 30px 0; color: #333333; font-size: 16px; line-height: 1.6;">
										Thank you for signing up! Use the code below to verify your email address.
									</p>

									<!-- OTP Box -->
									<div style="background-color: #f9f9f9; border: 2px solid #000000; border-radius: 8px; padding: 30px; text-align: center; margin: 30px 0;">
										<p style="margin: 0 0 10px 0; color: #666666; font-size: 13px; text-transform: uppercase; letter-spacing: 2px;">Verification Code</p>
										<p style="margin: 0; color: #000000; font-size: 48px; font-weight: 700; letter-spacing: 8px; font-family: 'Courier New', monospace;">%s</p>
									</div>

									<p style="margin: 0 0 20px 0; color: #666666; font-size: 14px;">
										⏱️ <strong>This code expires in 10 minutes</strong>
									</p>

									<div style="background-color: #fef3cd; border-left: 4px solid #ffb400; padding: 15px 20px; border-radius: 4px; margin: 20px 0;">
										<p style="margin: 0; color: #664d03; font-size: 14px;">
											<strong>Security Tip:</strong> Never share this code with anyone. We'll never ask for it via email or phone.
										</p>
									</div>

									<p style="margin: 30px 0 0 0; color: #999999; font-size: 13px;">
										If you didn't request this code, please ignore this email or <a href="mailto:support@yourdomain.com" style="color: #000000; text-decoration: underline;">contact support</a>.
									</p>
								</td>
							</tr>

							<!-- Footer -->
							<tr>
								<td style="background-color: #f9f9f9; padding: 30px; text-align: center; border-radius: 0 0 8px 8px; border-top: 1px solid #eeeeee;">
									<p style="margin: 0; color: #999999; font-size: 12px;">
										This is an automated message. Please do not reply to this email.
									</p>
								</td>
							</tr>
						</table>
					</td>
				</tr>
			</table>
		</body>
		</html>
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
	subject := "Reset Your Password"

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
		</head>
		<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif; background-color: #f5f5f5;">
			<table width="100%%" cellpadding="0" cellspacing="0" style="background-color: #f5f5f5;">
				<tr>
					<td align="center" style="padding: 40px 20px;">
						<table width="100%%" max-width="600" cellpadding="0" cellspacing="0" style="background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1);">
							<!-- Header -->
							<tr>
								<td style="background-color: #000000; padding: 40px 30px; text-align: center; border-radius: 8px 8px 0 0;">
									<h1 style="margin: 0; color: #ffffff; font-size: 28px; font-weight: 700;">Password Reset Request</h1>
								</td>
							</tr>

							<!-- Content -->
							<tr>
								<td style="padding: 40px 30px;">
									<p style="margin: 0 0 30px 0; color: #333333; font-size: 16px; line-height: 1.6;">
										We received a request to reset the password for your account. Use the code below to create a new password.
									</p>

									<!-- Reset Code Box -->
									<div style="background-color: #f9f9f9; border: 2px solid #000000; border-radius: 8px; padding: 30px; text-align: center; margin: 30px 0;">
										<p style="margin: 0 0 10px 0; color: #666666; font-size: 13px; text-transform: uppercase; letter-spacing: 2px;">Reset Code</p>
										<p style="margin: 0; color: #000000; font-size: 48px; font-weight: 700; letter-spacing: 8px; font-family: 'Courier New', monospace;">%s</p>
									</div>

									<div style="background-color: #fff3cd; border-left: 4px solid #ffc107; padding: 15px 20px; border-radius: 4px; margin: 20px 0;">
										<p style="margin: 0; color: #856404; font-size: 14px;">
											⏱️ <strong>This code expires in 10 minutes</strong>
										</p>
									</div>

									<p style="margin: 30px 0 0 0; color: #333333; font-size: 14px;">
										👉 <strong>Next steps:</strong>
									</p>
									<ol style="color: #666666; font-size: 14px; line-height: 1.8;">
										<li>Enter the code above in the password reset form</li>
										<li>Create a new password (at least 8 characters)</li>
										<li>Save your new password somewhere safe</li>
									</ol>

									<div style="background-color: #f8d7da; border-left: 4px solid #dc3545; padding: 15px 20px; border-radius: 4px; margin: 30px 0;">
										<p style="margin: 0; color: #721c24; font-size: 14px;">
											<strong>Security Alert:</strong> If you didn't request this, someone may be trying to access your account. Change your password immediately and <a href="mailto:support@yourdomain.com" style="color: #721c24; text-decoration: underline; font-weight: bold;">contact support</a>.
										</p>
									</div>

									<p style="margin: 20px 0 0 0; color: #999999; font-size: 13px;">
										Questions? <a href="mailto:support@yourdomain.com" style="color: #000000; text-decoration: underline;">Get help from our support team</a>
									</p>
								</td>
							</tr>

							<!-- Footer -->
							<tr>
								<td style="background-color: #f9f9f9; padding: 30px; text-align: center; border-radius: 0 0 8px 8px; border-top: 1px solid #eeeeee;">
									<p style="margin: 0; color: #999999; font-size: 12px;">
										This is an automated message. Please do not reply to this email.
									</p>
								</td>
							</tr>
						</table>
					</td>
				</tr>
			</table>
		</body>
		</html>
	`, otp)

	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.from,
		To:      []string{email},
		Subject: subject,
		Html:    html,
	})

	return err
}

func (s *Service) SendBookingConfirmation(email string, data BookingConfirmationData) error {
	subject := fmt.Sprintf("✓ Booking Confirmed at %s", data.BusinessName)

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
		</head>
		<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif; background-color: #f5f5f5;">
			<table width="100%%" cellpadding="0" cellspacing="0" style="background-color: #f5f5f5;">
				<tr>
					<td align="center" style="padding: 40px 20px;">
						<table width="100%%" max-width="600" cellpadding="0" cellspacing="0" style="background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); overflow: hidden;">
							<!-- Success Header -->
							<tr>
								<td style="background-color: #000000; padding: 50px 30px; text-align: center;">
									<div style="font-size: 48px; margin-bottom: 20px;">✓</div>
									<h1 style="margin: 0; color: #ffffff; font-size: 32px; font-weight: 700;">Booking Confirmed!</h1>
									<p style="margin: 10px 0 0 0; color: #cccccc; font-size: 16px;">Your payment has been received</p>
								</td>
							</tr>

							<!-- Main Content -->
							<tr>
								<td style="padding: 40px 30px;">
									<p style="margin: 0 0 30px 0; color: #333333; font-size: 16px; line-height: 1.6;">
										Hi <strong>%s</strong>,<br>
										Thank you for booking with <strong>%s</strong>. Your appointment is confirmed and payment has been successfully processed!
									</p>

									<!-- Booking Details Card -->
									<div style="background-color: #f9f9f9; border: 1px solid #eeeeee; border-radius: 8px; padding: 0; overflow: hidden; margin: 30px 0;">
										<!-- Time Section -->
										<div style="border-bottom: 1px solid #eeeeee; padding: 25px 30px;">
											<p style="margin: 0 0 8px 0; color: #999999; font-size: 12px; text-transform: uppercase; letter-spacing: 0.5px;">📅 Appointment Date & Time</p>
											<p style="margin: 0; color: #000000; font-size: 20px; font-weight: 600;">%s</p>
										</div>

										<!-- Services Section -->
										<div style="border-bottom: 1px solid #eeeeee; padding: 25px 30px;">
											<p style="margin: 0 0 8px 0; color: #999999; font-size: 12px; text-transform: uppercase; letter-spacing: 0.5px;">✂️ Services</p>
											<p style="margin: 0; color: #000000; font-size: 14px; line-height: 1.8;">%s</p>
										</div>

										<!-- Amount Section -->
										<div style="padding: 25px 30px;">
											<p style="margin: 0 0 8px 0; color: #999999; font-size: 12px; text-transform: uppercase; letter-spacing: 0.5px;">💳 Amount Paid</p>
											<p style="margin: 0; color: #000000; font-size: 28px; font-weight: 700;">%s</p>
										</div>
									</div>

									<!-- Booking ID -->
									<div style="background-color: #f0f0f0; padding: 20px 30px; border-radius: 8px; text-align: center; margin: 30px 0;">
										<p style="margin: 0 0 5px 0; color: #999999; font-size: 12px; text-transform: uppercase;">Booking Reference</p>
										<p style="margin: 0; color: #000000; font-size: 18px; font-weight: 600; font-family: 'Courier New', monospace;">%s</p>
									</div>

									<!-- Important Info Box -->
									<div style="background-color: #fff3cd; border-left: 4px solid #ffc107; padding: 20px 25px; border-radius: 4px; margin: 30px 0;">
										<p style="margin: 0 0 10px 0; color: #856404; font-size: 14px; font-weight: 600;">📌 Important Reminders</p>
										<ul style="margin: 0; padding-left: 20px; color: #856404; font-size: 13px;">
											<li style="margin-bottom: 8px;">Arrive <strong>10 minutes early</strong> to your appointment</li>
											<li style="margin-bottom: 8px;">We're located at the address provided during booking</li>
											<li>If you need to reschedule, please do so at least 24 hours in advance</li>
										</ul>
									</div>

									<!-- CTA Button -->
									<div style="text-align: center; margin: 40px 0;">
										<a href="%s" style="background-color: #000000; color: #ffffff; padding: 14px 40px; text-decoration: none; border-radius: 6px; display: inline-block; font-weight: 600; font-size: 15px; transition: background-color 0.2s;">
											View Booking Details
										</a>
									</div>

									<!-- Support Info -->
									<div style="background-color: #f9f9f9; border: 1px solid #eeeeee; border-radius: 8px; padding: 25px 30px; margin: 30px 0;">
										<p style="margin: 0 0 15px 0; color: #000000; font-size: 14px; font-weight: 600;">Need Help?</p>
										<p style="margin: 0 0 10px 0; color: #666666; font-size: 13px;">
											💬 Questions about your booking?<br>
											📧 Reply to this email or contact us for support
										</p>
										<p style="margin: 10px 0 0 0; color: #666666; font-size: 13px;">
											⏱️ <strong>Cancellation Policy:</strong> Cancel at least 24 hours before your appointment to avoid forfeiture of payment.
										</p>
									</div>

									<p style="margin: 30px 0 0 0; color: #999999; font-size: 13px; line-height: 1.6;">
										Thank you for choosing us! We look forward to seeing you.
									</p>
								</td>
							</tr>

							<!-- Footer -->
							<tr>
								<td style="background-color: #f9f9f9; padding: 30px; text-align: center; border-top: 1px solid #eeeeee;">
									<p style="margin: 0 0 10px 0; color: #999999; font-size: 12px;">
										<strong>%s</strong>
									</p>
									<p style="margin: 0; color: #999999; font-size: 11px;">
										This is an automated confirmation email. Please do not reply to this message.
									</p>
								</td>
							</tr>
						</table>
					</td>
				</tr>
			</table>
		</body>
		</html>
	`,
	data.CustomerName,
	data.BusinessName,
	data.AppointmentTime,
	data.Services,
	data.TotalAmount,
	data.BookingID,
	data.CallbackURL,
	data.BusinessName,
	)

	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.from,
		To:      []string{email},
		Subject: subject,
		Html:    html,
	})

	return err
}
