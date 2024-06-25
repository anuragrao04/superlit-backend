package googleSheets

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/anuragrao04/superlit-backend/database"
	"github.com/anuragrao04/superlit-backend/models"
	"github.com/anuragrao04/superlit-backend/prettyPrint"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/sheets/v4"
)

func PopulateGoogleSheetInstantTest(c *gin.Context) {
	var request models.PopulateGoogleSheetInstantTestSubmissionsRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	// now we extract the sheets ID from the provided link
	spreadSheetID, err := getSpreadSheetID(request.GoogleSheetLink)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the current spreadsheet properties
	// TODO: Implement a mutex lock for SRV and CTX. Idk if this supports concurrent operations but if it doesn't, we need a mutex.

	resp, err := SRV.Spreadsheets.Get(spreadSheetID).Context(CTX).Do()
	if err != nil {
		log.Printf("Unable to retrieve spreadsheet properties: %v", err)
		c.JSON(500, gin.H{"error": "Failed to retrieve spreadsheet properties"})
		return
	}

	// Count the number of sheets and generate a new sheet title based on the count
	numSheets := len(resp.Sheets)
	newSheetTitle := fmt.Sprintf("Sheet%d", numSheets+1)

	// Create a new sheet with the auto-assigned title
	requests := []*sheets.Request{
		{
			AddSheet: &sheets.AddSheetRequest{
				Properties: &sheets.SheetProperties{
					Title: newSheetTitle,
				},
			},
		},
	}

	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{Requests: requests}
	_, err = SRV.Spreadsheets.BatchUpdate(spreadSheetID, batchUpdateRequest).Context(CTX).Do()
	if err != nil {
		log.Printf("Unable to add new sheet: %v", err)
		c.JSON(500, gin.H{"error": "Failed to add new sheet"})
		return
	}

	log.Printf("Created new sheet '%s'\n", newSheetTitle)

	submissions, questionIDs, _, err := database.GetInstantTestSubmissions(request.PrivateCode)

	var valuesToInsert = make([][]interface{}, 0)
	minQuestionIndex := slices.Min(questionIDs)

	// Create the header row
	headRow := make([]interface{}, 0)
	headRow = append(headRow, "University ID")
	for _, questionID := range questionIDs {
		headRow = append(headRow, fmt.Sprintf("Question %d Score", questionID-minQuestionIndex+1))
		headRow = append(headRow, fmt.Sprintf("Question %d AI Verification", questionID-minQuestionIndex+1))
	}
	headRow = append(headRow, "Total Score")

	valuesToInsert = append(valuesToInsert, headRow)

	// Create the data rows
	for _, submission := range submissions {
		dataRow := make([]interface{}, len(headRow))
		dataRow[0] = submission.UniversityID
		for _, answer := range submission.Answers {
			dataRow[(answer.QuestionID-minQuestionIndex)*2+1] = answer.Score
			if answer.AIVerified {
				if answer.AIVerdict {
					dataRow[(answer.QuestionID-minQuestionIndex)*2+2] = "Verified Genuine"
				} else {
					dataRow[(answer.QuestionID-minQuestionIndex)*2+2] = "Check Student's Code!"
				}
			} else {
				dataRow[(answer.QuestionID-minQuestionIndex)*2+2] = "Not Verified"
			}
		}
		dataRow[2*len(questionIDs)+1] = submission.TotalScore
		valuesToInsert = append(valuesToInsert, dataRow)
	}

	valueRange := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         valuesToInsert,
	}

	lastColumnLetter := columnIndexToLetter(len(headRow))
	prettyPrint.PrettyPrint(valuesToInsert)

	updateRange := fmt.Sprintf("%s!A1:%s%d", newSheetTitle, lastColumnLetter, len(valuesToInsert))

	_, err = SRV.Spreadsheets.Values.Update(spreadSheetID, updateRange, valueRange).ValueInputOption("RAW").Context(CTX).Do()
	if err != nil {
		log.Printf("Unable to update sheet values: %v", err)
		c.JSON(500, gin.H{"error": "Failed to update sheet values"})
		return
	}

	c.JSON(200, gin.H{"message": "Sheet populated successfully"})
}

func columnIndexToLetter(index int) string {
	var result []rune
	for index > 0 {
		index--
		result = append(result, 'A'+rune(index%26))
		index /= 26
	}
	return string(result)
}

func getSpreadSheetID(link string) (spreadSheetID string, err error) {
	const prefix = "/d/"
	const suffix = "/edit"

	if !strings.Contains(link, prefix) || !strings.Contains(link, suffix) {
		return "", fmt.Errorf("invalid Google Sheets link")
	}

	start := strings.Index(link, prefix)
	if start == -1 {
		return "", fmt.Errorf("spreadsheet ID not found")
	}
	start += len(prefix)

	end := strings.Index(link[start:], suffix)
	if end == -1 {
		return "", fmt.Errorf("spreadsheet ID not found")
	}

	spreadsheetID := link[start : start+end]
	return spreadsheetID, nil
}
