package capstoneLogi

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LogFormat struct {
	UserID               string `json:"userID"`
	CurrentQuestionIndex int    `json:"currentQuestionIndex"`
	EditorContentBefore  string `json:"editorContentBefore"`
	EditorContentAfter   string `json:"editorContentAfter"`
	Timestamp            string `json:"timestamp"`
	IsPaste              bool   `json:"isPaste"`
	IsDeletion           bool   `json:"isDeletion"`
	IsCompilation        bool   `json:"isCompilation"`
	IsSubmission         bool   `json:"isSubmission"`
}

type LogiRequest struct {
	Logs []LogFormat `json:"logs" binding:"required"`
}

func HandleLogi(c *gin.Context) {
	var request LogiRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "yeverything galat",
		})
		return
	}

	filename := request.Logs[0].UserID // assuming this batch of logs comes from a single user

	f, err := os.OpenFile("./capstone-logi-logs/"+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error opening file",
		})
		return
	}

	for _, logLine := range request.Logs {
		text :=
			logLine.UserID +
				"," +
				fmt.Sprint(logLine.CurrentQuestionIndex) +
				"," +
				strconv.Quote(logLine.EditorContentBefore) +
				"," +
				strconv.Quote(logLine.EditorContentAfter) +
				"," +
				logLine.Timestamp +
				"," +
				fmt.Sprint(logLine.IsPaste) +
				"," +
				fmt.Sprint(logLine.IsDeletion) +
				"," +
				fmt.Sprint(logLine.IsCompilation) +
				"," +
				fmt.Sprint(logLine.IsSubmission) + "\n"

		_, err := f.WriteString(text)
		if err != nil {
			log.Println(err)
		}
	}

}
