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

type SignUpRequest struct {
	UniversityID string `json:"universityID"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	IsTeacher    bool   `json:"isTeacher"` // true if the registering user is a teacher. Else, they are a student
}

type SignInRequestUniversityID struct {
	UniversityID string `json:"universityID"`
	Password     string `json:"password"`
	IsTeacher    bool   `json:"isTeacher"` // true if the user is a teacher. Else, they are a student
}

type SignInRequestEmail struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	IsTeacher bool   `json:"isTeacher"` // true if the user is a teacher. Else, they are a student
}

type ForgotPasswordRequest struct {
	UniversityID string `json:"universityID"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

type AddUserToClassroomRequest struct {
	ClassroomCode string `json:"classroomCode"`
	UserID        uint   `json:"userID"`
	IsTeacher     bool   `json:"isTeacher"`
}

type CreateClassroomRequest struct {
	Name      string `json:"name"`
	TeacherID uint   `json:"teacherID"`
}
