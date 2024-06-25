package AI

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/anuragrao04/superlit-backend/database"
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

	// now we find the instant test with the given request.PrivateCode
	var instantTest models.InstantTest
	instantTest, err = database.GetInstantTestByPrivateCode(request.PrivateCode)
	if err != nil {
		c.JSON(400, gin.H{"error": "Something went wrong in fetching the test record from database: " + err.Error()})
		return
	}

	go verifyConstraintsInBackground(instantTest) // this will happen in the background
	c.JSON(200, gin.H{"message": "Verification started in the background. An email will be sent when finished"})

}

func verifyConstraintsInBackground(instantTest models.InstantTest) {
	for submissionIndex, submission := range instantTest.Submissions {
		for answerIndex, answer := range submission.Answers {
			// now we need to find the list of constraints for this answer
			question, err := getQuestion(instantTest.Questions, answer.QuestionID)
			constraints := question.Constraints
			if err != nil {
				// TODO: Send sad email
				return
			}

			// now we send a request to the server running ollama
			prompt := `Question Title:
			`

			prompt += question.Title + "\n"

			prompt += `Question Description:
			`
			prompt += question.Question + "\n"

			prompt += `Constraints: 
			` // the new line here is important
			for _, constraint := range constraints {
				prompt += "- " + constraint + "\n"
			}
			prompt += `Code:
			` // the new line here is important
			prompt += answer.Code

			postBody, _ := json.Marshal(map[string]interface{}{
				"model":  "superlit",
				"format": "json",
				"stream": false,
				"prompt": prompt,
			})
			responseBody := bytes.NewBuffer(postBody)
			resp, err := http.Post(os.Getenv("OLLAMA_URL")+"api/generate", "application/json", responseBody)
			if err != nil {
				// TODO: Send sad email
				log.Println("Failed to send request to ollama: ", err)
				continue
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			var response ollamaVerifyResponse
			err = json.Unmarshal(body, &response)

			if err != nil {
				// TODO: Send sad email
				log.Println("Failed to read response from ollama: ", err)
				continue
			}

			var ollamaVerdict ollamaVerifyAnswer
			err = json.Unmarshal([]byte(response.Response), &ollamaVerdict)
			log.Println("Verdict from LLM model: ", ollamaVerdict.Answer)
			if err != nil {
				// TODO: Send sad email, or write logic to retry.
				// but for now, I'm assuming that the model will return in JSON only.
				// that's what ollama's api docs say but just in case
				log.Println("LLM model returned some bullshit which could not be parsed: ", response.Response)
				continue
			}

			// now we have the verdict from ollama
			instantTest.Submissions[submissionIndex].Answers[answerIndex].AIVerified = true
			instantTest.Submissions[submissionIndex].Answers[answerIndex].AIVerdict = ollamaVerdict.Answer

		}
	}

	// after all of that is done, we save it to the database
	err := database.SaveInstantTest(instantTest)
	if err != nil {
		log.Println("Failed to save the instant test after AI verification: ", err)
		// TODO: try again or send sad email
	}

	// TODO: Send happy email!
}

type ollamaVerifyResponse struct {
	Response string `json:"response"`
}
type ollamaVerifyAnswer struct {
	Answer bool `json:"answer"`
}

// takes in the questions array, and a question ID and returns the question struct
func getQuestion(questions []models.Question, questionID uint) (models.Question, error) {
	for _, question := range questions {
		if question.ID == questionID {
			return question, nil
		}
	}
	log.Println("Someone is probably tryna raw dog the API. Question not found")
	return models.Question{}, errors.New("Question not found")
}
