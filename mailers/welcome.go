package mailers

import (
	"log"

	"gopkg.in/gomail.v2"
)

func SendWelcomeEmail(email string) {
	message := gomail.NewMessage()
	message.SetHeader("From", OUR_EMAIL)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Welcome To Superlit")
	message.SetBody("text/html", `
		<html>
		<body>
		<p>Welcome To Superlit</p>
		<p>We are very excited to have you here. Thank you for signing up</p>
<br />
		<p>Best regards,</p>
		<p>Superlit Team</p>
		</body>
		</html>
	`)

	for i := 0; i < 3; i++ {
		// retry up to 3 times
		if err := DIALER.DialAndSend(message); err != nil {
			log.Println("attempt", i+1, "failed to send email: ", err)
			if i == 2 {
				// it failed on the last attempt. Try no more
				return
			}
			continue
		}
		break
	}

	log.Println("Instant test codes sent successfully to", email)
	return

}
