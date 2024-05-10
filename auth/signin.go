package auth

import (
	"log"
	"net/http"

	"github.com/anuragrao04/superlit-backend/database"
	"github.com/anuragrao04/superlit-backend/models"
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

	c.JSON(http.StatusOK, &user)
}
