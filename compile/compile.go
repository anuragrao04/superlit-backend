// Description: This is the directory for compiling code. This doesn't do judging, this doesn't do scoring.
// The function for running the code through multiple test cases is also here.

package compile

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/anuragrao04/superlit-backend/models"
	"github.com/google/uuid"
)

// the below function is run on the /run route. See the main function for all other routes.
func RunCode(c *gin.Context) {
	var runRequest models.RunRequest
	// see this structure in models/models.go for request structure
	err := c.BindJSON(&runRequest)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	// write the code to a file
	codeFile, err := WriteCodeToFile(runRequest.Code, runRequest.Language)
	defer os.Remove(codeFile)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// now we need to handle the compilation different for different languages
	// for now, we'll handle C and python. Other languages can be added later

	// C
	if runRequest.Language == "c" {
		// compile the code
		compiledBinary, err := CompileBinary(codeFile, runRequest.Language)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Compilation Failed"})
			return
		}

		// run the binary
		output := RunBinary(runRequest.Input, compiledBinary)

		// we clean up the files
		os.Remove(compiledBinary)

		// send the output back
		c.JSON(http.StatusOK, gin.H{"output": output})

	} else if runRequest.Language == "py" {
		output := RunBinary(runRequest.Input, "python", codeFile)

		// we clean up the file
		os.Remove(codeFile)
		// send the output back
		c.JSON(http.StatusOK, gin.H{"output": output})
	}
}

// the below function is responsible for writing code to a file
// arguments:
// 1. code: the code to be written
// 2. language: the language of the code. This is used to determine the extension of the file
// returns:
// 1. the name of the file
// 2. error: any error that occurs during writing
func WriteCodeToFile(code, language string) (fileName string, err error) {
	// Write the code to a file
	// we'll write all code and compilations inside the playground directory

	// create the playground directory if it doesn't exist
	if _, err := os.Stat("./playground"); os.IsNotExist(err) {
		os.Mkdir("./playground", 0755)
	}
	// create a new file with a random name
	file, err := os.Create("./playground/" + uuid.New().String() + "." + language)
	if err != nil {
		log.Println(err)
		return "", errors.New("Error creating file")
	}
	// write the code to the file
	file.Write([]byte(code))

	if err = file.Close(); err != nil {
		log.Println(err)
		return "", errors.New("Error closing file")
	}

	return file.Name(), nil
}

// the below function is responsible for generating a compiled binary for a given file.
// it is only used in case of compiled languages
// arguments:
// 1. file: the file to be compiled
// 2. language: the language of the file. This is used to determine the compiler to be used
// the language argument would be the extension of a file. For example, "c" for C, "cpp" for C++, "java" for Java
// returns:
// 1. name of the compiled binary
// 2. error: any error that occurs during compilation
func CompileBinary(file string, language string) (compiledBinary string, err error) {
	// compile the code
	// C
	if language == "c" {
		cmd := exec.Command("gcc", file, "-o", file+".out")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		// log any error in the above compilation
		if err := cmd.Run(); err != nil {
			log.Println("Compilation Failed")
			log.Println(err)
			return "", err
		}
		return file + ".out", nil
	}
	return "", errors.New("Unsupported Language")
}

// the below function will take a compiled binary and run it.
// in case of an interpreted language, this will run the same
// arguments:
// 1. run command: the command to run the binary/program. This is sent as a slice of arguments
// for example, in order to run a python program, send: ["python", "<filename>"]
// 2. input: the input to be given to the program
// returns:
// output: the output of the binary including stdout and stderror
func RunBinary(input string, command ...string) string {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Run()
	return out.String()
}
