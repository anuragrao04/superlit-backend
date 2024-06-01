package instantTest

import (
	"encoding/json"
	"fmt"
	"github.com/anuragrao04/superlit-backend/database"
	"github.com/anuragrao04/superlit-backend/models"
	"github.com/gin-gonic/gin"
)

func PrettyPrint(obj interface{}) {
	bytes, _ := json.MarshalIndent(obj, "\t", "\t")
	fmt.Println(string(bytes))
}

func CreateTest(c *gin.Context) {
	// create a new test
	// parse incoming JSON
	var testRequest models.CreateInstantTestRequest
	err := c.BindJSON(&testRequest)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	// PrettyPrint(testRequest)

	privateCode, publicCode, err := database.CreateInstantTest(testRequest.Questions)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create test"})
		return
	}

	c.JSON(200, gin.H{"privateCode": privateCode, "publicCode": publicCode})
}

func GetInstantTest(c *gin.Context) {
	var testRequest models.GetInstantTestRequest
	err := c.BindJSON(&testRequest)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}
	questions, err := database.GetInstantTest(testRequest.PublicCode)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get test"})
		return
	}
	c.JSON(200, gin.H{"questions": questions})
}
