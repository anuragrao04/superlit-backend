package auth

import (
	"log"
	"net/http"

	"github.com/anuragrao04/superlit-backend/database"
	"github.com/anuragrao04/superlit-backend/models"
	"github.com/anuragrao04/superlit-backend/tokens"
	"github.com/gin-gonic/gin"
)

func SignInWithUniversityID(c *gin.Context) {
	var signInRequest models.SignInRequestUniversityID
	err := c.BindJSON(&signInRequest)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	// check if the user exists
	user, err := database.GetUserByUniversityIDPassword(signInRequest.UniversityID, signInRequest.Password)

	if err != nil {
		// means user does not exist
		c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
		return
	}

	token, err := tokens.CreateSignInToken(user.ID, user.UniversityID, user.IsTeacher, user.Name, user.Email)

	c.JSON(http.StatusOK, gin.H{"token": token, "isTeacher": user.IsTeacher})
}
