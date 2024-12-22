package assignments

import (
	"log"
	"net/http"
	"sort"

	"github.com/anuragrao04/superlit-backend/database"
	"github.com/anuragrao04/superlit-backend/models"
	"github.com/gin-gonic/gin"
)

func GetAssignmentLeaderboard(c *gin.Context) {
	var getSubmissionsRequest models.GetAssignmentLeaderboardRequest
	err := c.BindJSON(&getSubmissionsRequest)

	// TODO: We do not perform authentication that the student has
	// access to this assignment. We must figure out a way to do this

	// now we fetch all the submissions from the test with the given assigmment ID
	submissions, _, err := database.GetAssignmentSubmissions(getSubmissionsRequest.AssignmentID)
	if err != nil {
		if err.Error() == "No Such Assignment" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Wrong Assignment ID. No such assignment exists"})
			return
		}

		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type studentSubmission struct {
		UniversityID string `json:"universityID"`
		TotalScore   uint   `json:"totalScore"`
	}

	var formattedReturn = make([]studentSubmission, 0)

	for _, submission := range submissions {
		var formattedSubmission studentSubmission
		formattedSubmission.UniversityID = submission.UniversityID
		formattedSubmission.TotalScore = uint(submission.TotalScore)
		formattedReturn = append(formattedReturn, formattedSubmission)
	}

	// sort formattedReturn by totalScore
	sort.Slice(formattedReturn, func(i, j int) bool {
		return formattedReturn[i].TotalScore > formattedReturn[j].TotalScore
	})

	c.JSON(http.StatusOK, gin.H{"leaderboard": formattedReturn})
}