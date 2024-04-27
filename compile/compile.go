// Description: This is the directory for compiling code. This doesn't do judging, this doesn't do scoring.
// Input: Code, language and input/test cases
// Output: Outputs for that code/test cases

package compile

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/anuragrao04/superlit-backend/models"
	"github.com/google/uuid"
)

func RunCode(c *gin.Context) {
	var runRequest models.RunRequest
	err := c.BindJSON(&runRequest)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	// Write the code to a file
	// we'll write all code and compilations inside the temp directory
	// we'll create a new file with a unique ID on each request

	// create the playground directory if it doesn't exist
	if _, err := os.Stat("./playground"); os.IsNotExist(err) {
		os.Mkdir("./playground", 0755)
	}

	// create a new file
	file, err := os.Create("./playground/" + uuid.New().String() + "." + runRequest.Language)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Writing to file"})
		return
	}
	// write the code to the file
	file.Write([]byte(runRequest.Code))

	if err := file.Close(); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error closing file"})
		return
	}

	// now we need to handle the compilation different for different languages
	// for now, we'll handle C and python. Further languages can be added later

	// C
	if runRequest.Language == "c" {
		// compile the code
		cmd := exec.Command("gcc", file.Name(), "-o", file.Name()+".out")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		// log any error in the above compilation
		if err := cmd.Run(); err != nil {
			log.Println("Compilation Failed")
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Compilation Failed"})
			return
		}

		// now we pipe the input to the compiled code, and obtain the output
		cmd = exec.Command(file.Name() + ".out")
		cmd.Stdin = strings.NewReader(runRequest.Input)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()

		// we clean up the files
		os.Remove(file.Name())
		os.Remove(file.Name() + ".out")

		// send the output back
		c.JSON(http.StatusOK, gin.H{"output": out.String()})

	} else if runRequest.Language == "python" {
		// we don't need to compiule python code, we can directly run it
		cmd := exec.Command("python", file.Name())
		cmd.Stdin = strings.NewReader(runRequest.Input)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()

		// we clean up the file
		os.Remove(file.Name())

		// send the output back
		c.JSON(http.StatusOK, gin.H{"output": out.String()})
	}

}
