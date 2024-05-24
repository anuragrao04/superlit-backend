package tokens

import (
	"errors"
	"github.com/anuragrao04/superlit-backend/database"
	"github.com/anuragrao04/superlit-backend/models"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"os"
	"time"
)

func CreateForgotLink(universityID int) (link string, err error) {
	// first we see if there is a user with the given ID
	user := &models.User{}

	err = database.DB.Where("university_id = ?", universityID).First(&user).Error

	// guard clause to get out if user does not exist
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", errors.New("User doesn't exist")
	}

	// Making a JWT Token
	JWT_SECRET := os.Getenv("JWT_SECRET")

	// user.Password is an encrypted version of the password.
	// we will verify the signature at the time of reset password using this key
	uniqueSigningKey := JWT_SECRET + user.Password
	ttl := 15 * time.Minute
	token := jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"universityID": user.UniversityID,
			"userID":       user.ID,
			"exp":          time.Now().Add(ttl).Unix(),
		})

	signedToken, err := token.SignedString(uniqueSigningKey)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}
