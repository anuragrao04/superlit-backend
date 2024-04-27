package models

type RunRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"` // this will be the extension of the file. Example: "py" for python, "c" for C, "cpp" for C++, "java" for Java
	Input    string `json:"input"`
}
