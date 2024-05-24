package tokens

import (
	"crypto/ecdsa"
	"errors"
	"github.com/anuragrao04/superlit-backend/database"
	"github.com/anuragrao04/superlit-backend/models"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var JWT_SECRET *ecdsa.PrivateKey

func CreateForgotLink(universityID string) (link string, user *models.User, err error) {
	// first we see if there is a user with the given ID
	if database.DB == nil {
		log.Println("DB is nil. Not connected to database")
	}
	err = database.DB.Where("university_id = ?", universityID).First(&user).Error

	// guard clause to get out if user does not exist
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil, errors.New("User doesn't exist")
	}

	// user.Password is an encrypted version of the password.
	// we will verify the signature at the time of reset password using this key
	ttl := 15 * time.Minute
	token := jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"universityID": user.UniversityID,
			"userID":       user.ID,
			"exp":          time.Now().Add(ttl).Unix(),
		})

	if err != nil {
		log.Fatal(err)
	}
	signedToken, err := token.SignedString(JWT_SECRET)

	if err != nil {
		return "", nil, err
	}

	FRONTEND_FORGOT_PASS_URL := os.Getenv("FRONTEND_FORGOT_PASS_URL")
	// the ENV Variable must include the http(s):// prefix along with the path to the forgot password page
	// for example https://superlit.vercel.app/auth/forgotpassword/
	// INCLUDE TRAILING SLASH!

	link = FRONTEND_FORGOT_PASS_URL + signedToken

	return link, user, nil
}

func LoadECPrivateKey(pathToKeyFile string) error {
	keyData, err := os.ReadFile(pathToKeyFile)
	if err != nil {
		return errors.New("error reading key file")
	}

	privateKey, err := jwt.ParseECPrivateKeyFromPEM(keyData)
	if err != nil {
		return err
	}

	JWT_SECRET = privateKey
	return nil
}
