package assignments

import (
	"errors"
	"log"
	"net/http"
	"slices"
	"sort"

	"github.com/anuragrao04/superlit-backend/database"
	"github.com/anuragrao04/superlit-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// this function takes an array of answers
// returns the answer with the given questionID
func getAnswer(answers []models.Answer, questionID uint) (models.Answer, error) {
	for _, answer := range answers {
		if answer.QuestionID == questionID {
			return answer, nil
		}
	}
	return models.Answer{}, errors.New("No Such Answer")
}

// The below function is used to strip the input, expected output and student output from the test cases
func stripCases(testCases []models.VerifiedTestCase) []models.VerifiedTestCase {
	for _, tc := range testCases {
		tc.Input = ""
		tc.ExpectedOutput = ""
		tc.ProducedOutput = ""
	}
	return testCases
}

// the below function is used to get the submission of a particular student
func GetStudentSubmission(c *gin.Context) {
	var request models.GetStudentSubmissionRequest
	err := c.BindJSON(&request)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	value, ok := c.Get("claims")
	claims, ok := value.(jwt.MapClaims)
	userIDFloat, ok := claims["userID"].(float64)
	userID := uint(userIDFloat)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized"})
		return
	}

	submission, questionIDs, questions, err := database.GetAssignmentSubmissionPerStudent(request.AssignmentID, userID)

	if err != nil {
		if err.Error() == "No Such AssignmentSubmission" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Wrong Assignment ID or User ID. No such assignment exists"})
			return
		}

		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// now we must do some formatting.
	// we need to find the min questionID to get the question numbers

	minQuestionID := slices.Min(questionIDs)

	type returnFormatArrayElement struct {
		QuestionNumber      uint                      `json:"questionNumber"`
		QuestionTitle       string                    `json:"questionTitle"`
		QuestionDescription string                    `json:"questionDescription"`
		Attempted           bool                      `json:"attempted"`
		Code                string                    `json:"code"`
		AIVerified          bool                      `json:"AIVerified"`
		AIVerdict           bool                      `json:"AIVerdict"`           // if AI has verified the code, this is the verdict. If true, it means it's aproved. else something is fishy
		AIVerdictFailReason string                    `json:"AIVerdictFailReason"` // if AI has flagged, why?
		AIVivaTaken         bool                      `json:"AIVivaTaken"`         // if AI Viva was taken
		AIVivaScore         int                       `json:"AIVivaScore"`         // how many viva questions did the student answer correctly
		Score               int                       `json:"score"`
		TestCases           []models.VerifiedTestCase `json:"testCases"`
		// we will remove the input, expected output and student output from the above.
		// It will only contain information about if test case is passed or not.
	}

	var returnArray = make([]returnFormatArrayElement, len(questionIDs))

	for _, question := range questions {
		questionNumberIndex := question.ID - minQuestionID
		returnArray[questionNumberIndex].QuestionTitle = question.Title
		returnArray[questionNumberIndex].QuestionDescription = question.Question
		returnArray[questionNumberIndex].QuestionNumber = questionNumberIndex + 1
		answer, err := getAnswer(submission.Answers, question.ID)
		if err != nil {
			returnArray[questionNumberIndex].Attempted = false
			continue
		} else {
			returnArray[questionNumberIndex].Attempted = true
		}
		returnArray[questionNumberIndex].Code = answer.Code
		returnArray[questionNumberIndex].AIVerified = answer.AIVerified
		returnArray[questionNumberIndex].AIVerdict = answer.AIVerdict
		returnArray[questionNumberIndex].AIVerdictFailReason = answer.AIVerdictFailReason
		returnArray[questionNumberIndex].AIVivaTaken = answer.AIVivaTaken
		returnArray[questionNumberIndex].AIVivaScore = answer.AIVivaScore
		returnArray[questionNumberIndex].Score = answer.Score
		returnArray[questionNumberIndex].TestCases = stripCases(answer.TestCases)
	}

	c.JSON(http.StatusOK, gin.H{"answers": returnArray, "totalScore": submission.TotalScore})
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

	// now we fetch all the submissions from the instant test with the given assigmment ID
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
		QuestionNumber      uint   `json:"questionNumber"`
		Score               uint   `json:"score"`
		AIVerified          bool   `json:"AIVerified"`
		AIVerdict           bool   `json:"AIVerdict"`
		AIVerdictFailReason string `json:"AIVerdictFailReason"`
		AIVivaTaken         bool   `json:"AIVivaTaken"`
		AIVivaScore         int    `json:"AIVivaScore"`
		StudentsCode        string `json:"studentsCode"`
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
			formattedSubmission.Submissions[questionNumber].AIVerdictFailReason = answer.AIVerdictFailReason
			formattedSubmission.Submissions[questionNumber].AIVivaTaken = answer.AIVivaTaken
			formattedSubmission.Submissions[questionNumber].AIVivaScore = answer.AIVivaScore
			formattedSubmission.Submissions[questionNumber].StudentsCode = answer.Code
		}
		formattedReturn = append(formattedReturn, formattedSubmission)
	}

	// sort formattedReturn by totalScore
	sort.Slice(formattedReturn, func(i, j int) bool {
		return formattedReturn[i].TotalScore > formattedReturn[j].TotalScore
	})

	blacklist, err := database.GetAssignmentBlacklist(getSubmissionsRequest.AssignmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"submissions": formattedReturn, "maxNumberOfQuestions": len(questionIDs), "blacklistedStudents": blacklist})
}
