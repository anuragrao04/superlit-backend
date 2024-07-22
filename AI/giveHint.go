package AI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/anuragrao04/superlit-backend/database"
	"github.com/anuragrao04/superlit-backend/instantTest"
	"github.com/anuragrao04/superlit-backend/models"
	"github.com/gin-gonic/gin"
)

func GiveHint(c *gin.Context) {
	var request models.AIGiveHintRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}
	var question models.Question
	// TODO: Move this to database package
	database.DBLock.Lock()
	database.DB.First(&question, request.QuestionID)
	database.DBLock.Unlock()

	// although this function is part of instant test package
	// we can use it to calculate score of any test
	_, _, testCasesFailed, err := instantTest.CalculateScore(question, request.Code, request.Language)
	// TODO: change this entire format to fmt.Sprintf
	// The current way is confusing to new developers

	prompt := `Question Title:
	`
	prompt += question.Title + "\n"

	prompt += `Question Description:
	`
	prompt += question.Question + "\n"

	prompt += `Constraints: 
	` // the new line here is important
	for _, constraint := range question.Constraints {
		prompt += "- " + constraint + "\n"
	}
	prompt += `Code:
	` // the new line here is important
	prompt += request.Code

	prompt += `Test Cases Failed: 
	`
	for i, testCase := range testCasesFailed {
		prompt += `Case ` + fmt.Sprint(i+1) + `:
`
		prompt += `Input:
`
		prompt += testCase.Input + "\n"

		prompt += `Expected Output:
`
		prompt += testCase.ExpectedOutput + "\n"

		prompt += `Produced Output:
`
		prompt += testCase.ProducedOutput + "\n"
	}

	postBody, _ := json.Marshal(map[string]interface{}{
		"model":  "superlit-hint-lite",
		"format": "json",
		"stream": false,
		"prompt": prompt,
	})

	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(os.Getenv("OLLAMA_URL")+"api/generate", "application/json", responseBody)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	var response ollamaHintResponse
	err = json.Unmarshal(body, &response)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var ollamaHint ollamaHintAnswer
	err = json.Unmarshal([]byte(response.Response), &ollamaHint)
	if err != nil {
		// TODO: write logic to retry.
		// but for now, I'm assuming that the model will return in JSON only.
		// that's what ollama's api docs say but just in case
		log.Println("LLM model returned some bullshit which could not be parsed: ", response.Response)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hint": ollamaHint.Hint})
}

type ollamaHintResponse struct {
	Response string `json:"response"`
}

type ollamaHintAnswer struct {
	Hint string `json:"hint"`
}
