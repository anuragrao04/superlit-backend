package classroom

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/anuragrao04/superlit-backend/database"
	"github.com/anuragrao04/superlit-backend/models"
)

func AddUserToClassroom(c *gin.Context) {
	var addUserToClassroomRequest models.AddUserToClassroomRequest
	// see this structure in models/models.go for request structure
	err := c.BindJSON(&addUserToClassroomRequest)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	err = database.AddUserToClassroom(addUserToClassroomRequest.UserID, addUserToClassroomRequest.ClassroomCode, addUserToClassroomRequest.IsTeacher)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong in adding the user to the classroom: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User added to classroom successfully"})
}
