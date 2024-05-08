package models

// INCOMING REQUEST MODELS

// this is the request for /run. It only runs the given code and sends the output back
type RunRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"` // this will be the extension of the file. Example: "py" for python, "c" for C, "cpp" for C++, "java" for Java
	Input    string `json:"input"`
}

type SubmitRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
	TestID   string `json:"testID"` // this will be the test ID of the test stored in the database
}

type RegisterRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	IsTeacher bool   `json:"isTeacher"` // true if the registering user is a teacher. Else, they are a student
}

type LoginRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	IsTeacher bool   `json:"isTeacher"` // true if the user is a teacher. Else, they are a student
}
