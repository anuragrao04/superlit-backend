package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UniversityID string // this may be employee ID or student ID (EMP ID and SRN in terms of PES)
	Name         string
	Email        string
	Password     string
	IsTeacher    bool
	Classrooms   []Classroom `gorm:"many2many:user_classroom;"`
}

// aka class
type Classroom struct {
	gorm.Model
	Name        string
	Code        string       // this is the code that students will use to join the classroom
	Users       []User       `gorm:"many2many:user_classroom;"`
	Assignments []Assignment `gorm:"many2many:assignment_classroom;"`
}

// aka test
type Assignment struct {
	gorm.Model
	Name        string
	Description string
	StartTime   *time.Time
	EndTime     *time.Time
	Classrooms  []Classroom  `gorm:"many2many:assignment_classroom;"`
	Questions   []Question   `gorm:"foreignKey:AssignmentID"`
	Submissions []Submission `gorm:"foreignKey:AssignmentID"`

	// this is the classrooms in which the assignment is assigned.
	// Since every classroom can have multiple assignments, and one assignment may be assigned to multiple classrooms, we have a many to many relationship
}

// Instant tests are tests that are created on the fly by teachers.
// Private code is not shared. This code is used to edit the test
// Public code is shared with students. This code is used to attempt the test
type InstantTest struct {
	gorm.Model
	PrivateCode string
	PublicCode  string
	IsActive    bool
	Questions   []Question              `gorm:"foreignKey:InstantTestID"`
	Submissions []InstantTestSubmission `gorm:"foreignKey:InstantTestID"`
}

type InstantTestSubmission struct {
	gorm.Model
	InstantTestID uint
	UniversityID  string
	Answers       []Answer
	TotalScore    int
}

type Submission struct {
	gorm.Model
	AssignmentID uint
	UserID       uint
	User         User
	Answers      []Answer
	TotalScore   int
}

type Answer struct {
	gorm.Model
	SubmissionID            uint
	InstantTestSubmissionID uint
	QuestionID              uint
	Code                    string
	TestCases               []VerifiedTestCase
	Score                   int // total score for this particular question
}

type Question struct {
	gorm.Model
	AssignmentID   uint
	InstantTestID  uint
	Title          string            `json:"title"`
	Question       string            `json:"question"`
	ExampleCases   []ExampleTestCase `json:"exampleCases"`
	PreWrittenCode string            `json:"preWrittenCode"`
	TestCases      []TestCase        `json:"testCases"`
}

type ExampleTestCase struct {
	gorm.Model
	QuestionID     uint
	Input          string `'json:"input"`
	ExpectedOutput string `json:"expectedOutput"`
	Score          int    `json:"score"`
	Explanation    string `json:"explanation"`
}

type TestCase struct {
	gorm.Model
	QuestionID     uint
	Input          string `'json:"input"`
	ExpectedOutput string `json:"expectedOutput"`
	Score          int    `json:"score"`
}

type VerifiedTestCase struct {
	AnswerID       uint   `json:"answerID"`
	Passed         bool   `json:"passed"`
	Input          string `json:"input"`
	ExpectedOutput string `json:"expectedOutput"`
	ProducedOutput string `json:"producedOutput"`
}
