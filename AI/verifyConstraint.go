package AI

import (
	"github.com/anuragrao04/superlit-backend/models"
	"github.com/gin-gonic/gin"
)

// this file contains the code to verify a particular constraint for
// all submissions from a particular test.

func AIVerifyConstraintsInstantTest(c *gin.Context) {
	// this function will verify the constraints for all submissions of a particular instant test
	var request models.AIVerifyConstraintsInstantTestRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}
}
