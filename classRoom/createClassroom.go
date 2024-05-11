package classroom

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/anuragrao04/superlit-backend/database"
	"github.com/anuragrao04/superlit-backend/models"
)

func CreateClassroom(c *gin.Context) {
	var createClassroomRequest models.CreateClassroomRequest
	err := c.BindJSON(&createClassroomRequest)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	classroom, err := database.CreateClassroom(createClassroomRequest.Name, createClassroomRequest.TeacherID)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong in creating the classroom: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, classroom)
}
