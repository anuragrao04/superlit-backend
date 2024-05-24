package mailers

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

// this function takes the signed resetLink string and email.
// It sends an email to the user with the given token with a nice message
// console and all
// be emotional
// users are dum but we're here to save the day!
func SendForgotPasswordEmail(email string, resetLink string) error {
	OUR_EMAIL := os.Getenv("EMAILID")
	OUR_PASSWORD := os.Getenv("EMAILPASSWORD")

	if OUR_EMAIL == "" || OUR_PASSWORD == "" {
		return fmt.Errorf("email credentials are not set in environment variables")
	}

	message := gomail.NewMessage()
	message.SetHeader("From", OUR_EMAIL)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Reset Your Password")
	message.SetBody("text/html", fmt.Sprintf(`
		<html>
		<body>
		<p>Dear User,</p>
		<p>To reset your password, please click on the link below:</p>
		<p><a href="%s">Reset Password</a></p>
		<p>This link will expire in 15 minutes, so be sure to use it soon.</p>
		<p>If you did not request a password reset, please ignore this email or contact our support team.</p>
		<p>Best regards,</p>
		<p>Superlit Team</p>
		</body>
		</html>
	`, resetLink))

	dialer := gomail.NewDialer("smtp.gmail.com", 587, OUR_EMAIL, OUR_PASSWORD)

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	fmt.Println("Reset password email sent successfully to", email)
	return nil
}
