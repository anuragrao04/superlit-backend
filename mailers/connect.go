package mailers

import (
	"errors"
	"os"

	"github.com/anuragrao04/superlit-backend/events"
	"gopkg.in/gomail.v2"
)

// these variables are global to this package
var OUR_EMAIL string
var DIALER *gomail.Dialer

func Connect() error {

	OUR_EMAIL = os.Getenv("EMAILID")
	OUR_PASSWORD := os.Getenv("EMAILPASSWORD")

	if OUR_EMAIL == "" || OUR_PASSWORD == "" {
		return errors.New("email credentials are not set in environment variables")
	}

	DIALER = gomail.NewDialer("smtp.gmail.com", 587, OUR_EMAIL, OUR_PASSWORD)
	if DIALER == nil {
		return errors.New("failed to connect to email server")
	}

	events.SubscribeUserCreated(func(email string) {
		SendWelcomeEmail(email)
	})

	return nil
}
