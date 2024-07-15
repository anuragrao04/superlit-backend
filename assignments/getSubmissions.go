package assignments

import (
	"log"
	"net/http"
	"slices"

	"github.com/anuragrao04/superlit-backend/database"
	"github.com/anuragrao04/superlit-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// the below function is used to get the submission of a particular student
// TODO: Implement this
func GetStudentSubmission(c *gin.Context) {

}

// the below function is used to get the submissions of a particular assignment
func GetAssignmentSubmissions(c *gin.Context) {
	var getSubmissionsRequest models.GetAssignmentSubmissionsRequest
	err := c.BindJSON(&getSubmissionsRequest)

	value, ok := c.Get("claims")
	claims, ok := value.(jwt.MapClaims)
	isTeacher := claims["isTeacher"].(bool)
	if !isTeacher || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to view this page"})
		return
	}

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	// TODO: We do not perform authentication that the teacher has
	// access to this assignment. We must figure out a way to do this
	// we only check that the user is a teacher

	// now we fetch all the submissions from the instant test with the given private code
	submissions, questionIDs, err := database.GetAssignmentSubmissions(getSubmissionsRequest.AssignmentID)
	if err != nil {
		if err.Error() == "No Such Assignment" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Wrong Assignment ID. No such assignment exists"})
			return
		}

		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// now we must do some formatting.

	type answerSubmission struct {
		QuestionNumber uint   `json:"questionNumber"`
		Score          uint   `json:"score"`
		AIVerified     bool   `json:"AIVerified"`
		AIVerdict      bool   `json:"AIVerdict"`
		StudentsCode   string `json:"studentsCode"`
	}

	type studentSubmission struct {
		UniversityID string             `json:"universityID"`
		Submissions  []answerSubmission `json:"submissionsAnswer"`
		TotalScore   uint               `json:"totalScore"`
	}

	var formattedReturn = make([]studentSubmission, 0)

	minQuestionID := slices.Min(questionIDs)

	for _, submission := range submissions {
		var formattedSubmission studentSubmission
		formattedSubmission.UniversityID = submission.UniversityID
		formattedSubmission.TotalScore = uint(submission.TotalScore)
		formattedSubmission.Submissions = make([]answerSubmission, len(questionIDs))
		for _, answer := range submission.Answers {
			questionNumber := answer.QuestionID - minQuestionID
			formattedSubmission.Submissions[questionNumber].Score = uint(answer.Score)
			formattedSubmission.Submissions[questionNumber].QuestionNumber = uint(questionNumber + 1)
			formattedSubmission.Submissions[questionNumber].AIVerified = answer.AIVerified
			formattedSubmission.Submissions[questionNumber].AIVerdict = answer.AIVerdict
			formattedSubmission.Submissions[questionNumber].StudentsCode = answer.Code
		}
		formattedReturn = append(formattedReturn, formattedSubmission)
	}

	c.JSON(http.StatusOK, gin.H{"submissions": formattedReturn, "maxNumberOfQuestions": len(questionIDs)})
}
