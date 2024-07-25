package AI

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/anuragrao04/superlit-backend/database"
	"github.com/anuragrao04/superlit-backend/models"
	"github.com/gin-gonic/gin"
)

func GetVivaQuestions(c *gin.Context) {
	var request models.AIGetVivaQuestionsRequest
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

	postBody, _ := json.Marshal(map[string]interface{}{
		"model":  "superlit-viva-lite",
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

	var response ollamaVivaResponse
	err = json.Unmarshal(body, &response)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var ollamaViva ollamaVivaAnswer
	err = json.Unmarshal([]byte(response.Response), &ollamaViva)
	if err != nil {
		// TODO: write logic to retry.
		// but for now, I'm assuming that the model will return in JSON only.
		// that's what ollama's api docs say but just in case
		log.Println("LLM model returned some bullshit which could not be parsed: ", response.Response)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, ollamaViva)
}

type ollamaVivaResponse struct {
	Response string `json:"response"`
}

type ollamaVivaAnswer struct {
	Questions []vivaQuestions `json:"vivaQuestions"`
}

type vivaQuestions struct {
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	CorrectOption uint     `json:"correctOption"`
}
