package mailers

import (
	"os"
)

// this function takes the signed resetLink string and email.
// It sends an email to the user with the given token with a nice message
// console and all
// be emotional
// users are dum

func SendForgotPasswordEmail(email string, resetLink string) error {

	OUR_EMAIL := os.Getenv("EMAILID")
	OUR_PASSWORD := os.Getenv("EMAILPASSWORD")

}
